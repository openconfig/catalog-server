// Copyright 2021 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program yang crawls all modules at certain commit of `https://github.com/openconfig/public`.
// Usage: yang [--path DIR] [--url URL]
//
// DIR is a comma separated list of paths that are appended as the search directory.
// If DIR appears as DIR/... then DIR and all direct and indirect subdirectories are checked.
//
// URL is github URL prefix of git commit that program is crawling.
// THIS PROGRAM IS STILL JUST A DEVELOPMENT TOOL.
package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openconfig/catalog-server/pkg/db"
	oc "github.com/openconfig/catalog-server/pkg/ygotgen"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ygot"
	"github.com/pborman/getopt"
)

// exitIfError writes errs to standard error and exits with an exit status of 1.
// If errs is empty then exitIfError does nothing and simply returns.
func exitIfError(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		stop(1)
	}
}

const (
	modelDir     = `release/models`   // modelDir is directory in openconfig/public github repo that contains modules in `models` directory.
	modelKeyword = `models`           // We check whether path of found module contains `models` to check whether it's under `models` directory.
	ietfDir      = `third_party/ietf` // ietfDir is directory in openconfig/public github repo that contains modules in `ietf` directory.
	ietfKeyword  = `ietf`             // We check whether path of found module contains `ietf` to check whether it's under `ietf` directory.
	orgName      = `openconfig`       // default orgName that is used when inserting modules into database.
)

// urlMap is map from model's name to its github URL.
var urlMap = map[string]string{}

// traverseDir traverses given directory *dir* to find all modules in this directory including its sub-directories.
// *url* is the url prefix of github repo at certain commit.
// It returns a slice of names of modules found.
func traverseDir(dir string, url string) ([]string, error) {
	var names []string
	walkDirFn := func(p string, d fs.DirEntry, err error) error {
		f, err := d.Info()
		if err != nil {
			return fmt.Errorf("traverseDir fails to call d.Info(): %v", err)
		}
		if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), ".yang") {
			// Only search for files with suffix of `.yang`
			currDir := path.Dir(p)
			// currDir is the name of current directory containing this module file.
			// For simplicity, here we directly choose the last level diretcory's name as currDir,
			// since openconfig repo only contains one level sub-directory inside *models* or *ietf* directories.
			currDir = currDir[strings.LastIndex(currDir, "/")+1:]
			file, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("traverseDir: failed to open file %s: %v", p, err)
			}
			// *name* is name of yang module/submodule, we get it by removing `.yang` from original file name.
			name := strings.TrimSuffix(f.Name(), ".yang")

			// Check whether found module is under `models` directory or `ietf` dirctory.
			if strings.Contains(p, modelKeyword) {
				// Modules/submodules under `models` dirctory are under subdirectory (*currDir*) of `models` directory.
				// Store url of found module/submodule in urlMap.
				urlMap[name] = url + modelDir + "/" + currDir + "/" + f.Name()
			} else {
				// Found modules must be either under `models` directory or `ietf` directory.
				if !(strings.Contains(p, ietfKeyword)) {
					return fmt.Errorf("traverseDir: model %s not in either models dir or ietd dir", f.Name())
				}
				// `ietf` directory does not have subdirectory.
				// Store url of found module/submodule in urlMap.
				urlMap[name] = url + ietfDir + "/" + f.Name()
			}
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()

			// Only append `module` and ignore `submodules`.
			// Note this is specific to openconfig repo only.
			// For other repos with comments before schema of modules, its content may not start with *module* keyword.
			if strings.HasPrefix(strings.TrimSpace(scanner.Text()), "module ") {
				names = append(names, name)
			} else {
				log.Println(p + " does not have modules")
			}
			file.Close()
		}
		return nil
	}

	err := filepath.WalkDir(dir, walkDirFn)
	return names, err
}

var stop = os.Exit

// crawlModules is mainly from goyang program.
// It takes inputs of a list of *paths*, and a *url* prefix of the commit where the crawler is crawling.
// It first traverse all given paths to find all modules in these paths.
// Then it crawls data of each module using goyang library, and return a slice of points of yang.Module, and a sorted slice of names of modules.
func crawlModules(paths []string, url string) (map[string]*yang.Module, []string) {
	for _, path := range paths {
		expanded, err := yang.PathsWithModules(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		yang.AddPath(expanded...)
	}

	// files stores names of all modules to search for.
	files := getopt.Args()

	ms := yang.NewModules()

	// If names of modules to search for is not given, we traverse all given path to find
	// all modules in these directories and crawl them later.
	if len(files) == 0 {
		for _, path := range paths {
			names, err := traverseDir(path, url)
			if err != nil {
				log.Fatalf("traverse directory %s failed: %v", path, err)
			}
			// Append all found modules into *files*.
			files = append(files, names...)
		}
	}

	for _, name := range files {
		if err := ms.Read(name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}

	// Process the read files, exiting if any errors were found.
	exitIfError(ms.Process())

	// Keep track of the top level modules we read in.
	// Those are the only modules we want to print below.
	mods := map[string]*yang.Module{}
	var names []string

	for _, m := range ms.Modules {
		if mods[m.Name] == nil {
			mods[m.Name] = m
			names = append(names, m.Name)
		}
	}
	sort.Strings(names)
	return mods, names
}

// populateModule takes a pointer of yang.Module *mod*, and *name* of this module.
// It populates this yang.Module into a YANG catalog Module, and returns pointer to that module.
// If any error happens during the process, it just return a nil pointer.
func populateModule(mod *yang.Module, name string) *oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module {
	module := &oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module{
		Name:      &mod.Name,
		Namespace: &mod.Namespace.Name,
		Prefix:    &mod.Prefix.Name,
		Summary:   &mod.Description.Name,
	}

	version, err := yang.MatchingExtensions(mod, "openconfig-extensions", "openconfig-version")
	if err != nil || len(version) == 0 {
		log.Printf("%s do not have version\n", mod.Name)
		return nil
	}

	module.Version = &version[0].Argument

	// If there are multiple revisions, we directly use the lastest one.
	if len(mod.Revision) > 0 {
		module.Revision = &mod.Revision[0].Name
	}
	for i := 0; i < len(mod.Import); i++ {
		module.GetOrCreateDependencies().RequiredModule = append(module.GetOrCreateDependencies().RequiredModule, mod.Import[i].Name)
	}
	for i := 0; i < len(mod.Include); i++ {
		module.GetOrCreateSubmodules().GetOrCreateSubmodule(mod.Include[i].Name)
		submoduleURL := urlMap[mod.Include[i].Name]
		if submoduleURL == "" {
			log.Fatalf("cannot find url of submodule: %s", mod.Include[i].Name)
			return nil
		}
		module.GetOrCreateSubmodules().GetOrCreateSubmodule(mod.Include[i].Name).GetOrCreateAccess().Uri = &submoduleURL

	}
	moduleURL := urlMap[name]
	module.GetOrCreateAccess().Uri = &moduleURL
	return module
}

// insertModule marshalls the module into JSON string, and tries to insert it into database.
// It checks whether the key (name+version) of module already exists, if the module exists, then the insertion is skipped.
func insertModule(module *oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module) {
	// Serialize module struct into json for insertion.
	json, err := ygot.EmitJSON(module, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		log.Fatalf("Marshalling into json string failed\n")
	}

	// Query to check whether the key already exists before insertion.
	// As we crawl from the lastest version to the oldest one, we want to only insert the lastest data into database.
	queryRes, err := db.QueryModulesByKey(module.Name, module.Version)
	if err != nil {
		log.Printf("Query module, Name: %s, Version: %s failed: %v\n", module.GetName(), module.GetVersion(), err)
		return
	}

	// If the key already matches an existing module, we would not do an insertion.
	if len(queryRes) > 0 {
		log.Printf("module, Name: %s, Version: %s already exists in database, do not update it\n", module.GetName(), module.GetVersion())
		return
	}

	if err := db.InsertModule(orgName, module.GetName(), module.GetVersion(), json); err != nil {
		log.Printf("Insert module, Name: %s, Version: %s failed: %v\n", module.GetName(), module.GetVersion(), err)
		return
	}
	log.Printf("Inserting module succeeds, Name: %s, Version: %s\n", module.GetName(), module.GetVersion())
}

func main() {
	var help bool
	var paths []string
	var url string
	getopt.ListVarLong(&paths, "path", 'p', "comma separated list of directories to add to search path", "DIR[,DIR...]")
	getopt.BoolVarLong(&help, "help", 'h', "display help")
	// *url* is the url prefix of github repo at certain commit that we are crawling modules from.
	getopt.StringVarLong(&url, "url", 'u', "url prefix of git commit that we are crawling")
	getopt.SetParameters("")

	if err := getopt.Getopt(func(o getopt.Option) bool {
		return true
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	if help {
		getopt.CommandLine.PrintUsage(os.Stderr)
		os.Exit(0)
	}

	mods, names := crawlModules(paths, url)

	// Connect to DB
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Connect to db failed: %v", err)
		stop(1)
	}
	defer db.Close()

	// Convert all found modules into ygot go structure of Module and insert them into database.
	for _, n := range names {
		module := populateModule(mods[n], n)
		// if module is nil, it means that there is some issue when populating this module, thus we skip inserting it.
		if module == nil {
			continue
		}
		insertModule(module)
	}

}
