// Copyright 2015 Google Inc.
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

// Program yang parses YANG files, displays errors, and possibly writes
// something related to the input on output.
//
// Usage: yang [--path DIR] [--format FORMAT] [FORMAT OPTIONS] [MODULE] [FILE ...]
//
// If MODULE is specified (an argument that does not end in .yang), it is taken
// as the name of the module to display.  Any FILEs specified are read, and the
// tree for MODULE is displayed.  If MODULE was not defined in FILEs (or no
// files were specified), then the file MODULES.yang is read as well.  An error
// is displayed if no definition for MODULE was found.
//
// If MODULE is missing, then all base modules read from the FILEs are
// displayed.  If there are no arguments then standard input is parsed.
//
// If DIR is specified, it is considered a comma separated list of paths
// to append to the search directory.  If DIR appears as DIR/... then
// DIR and all direct and indirect subdirectories are checked.
//
// FORMAT, which defaults to "tree", specifies the format of output to produce.
// Use "goyang --help" for a list of available formats.
//
// FORMAT OPTIONS are flags that apply to a specific format.  They must follow
// --format.
//
// THIS PROGRAM IS STILL JUST A DEVELOPMENT TOOL.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime/trace"
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/openconfig/catalog-server/pkg/db"
	oc "github.com/openconfig/catalog-server/pkg/ygotgen"
	"github.com/openconfig/goyang/pkg/indent"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ygot"
	"github.com/pborman/getopt"
)

// Each format must register a formatter with register.  The function f will
// be called once with the set of yang Entry trees generated.
type formatter struct {
	name  string
	f     func(io.Writer, []*yang.Entry)
	help  string
	flags *getopt.Set
}

var formatters = map[string]*formatter{}

func register(f *formatter) {
	formatters[f.name] = f
}

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
	modelURL     = `https://github.com/openconfig/public/tree/master/release/models`   // modelURL is prefix of github URL that would be appended for modules in `models` directory.
	modelKeyword = `models`                                                            // We check whether path of found module contains `models` to check whether it's under `models` directory.
	ietfURL      = `https://github.com/openconfig/public/tree/master/third_party/ietf` // ietfURL is prefix of github URL that would be appended for models in `ietf` directory.
	ietfKeyword  = `ietf`                                                              // We check whether path of found module contains `ietf` to check whether it's under `ietf` directory.
	orgName      = `openconfig`                                                        // default orgName that is used when inserting modules into database.
)

// urlMap is map from model's name to its github URL.
var urlMap = map[string]string{}

// traverseDir traverses given directory *dir* to find all modules in this directory including its sub-directories.
// It returns a slice of names of modules found.
func traverseDir(dir string) ([]string, error) {
	dirfiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("traverseDir: read files from directory %s failed: %v\n", dir, err)
	}
	var names []string
	var dirs []string
	for _, f := range dirfiles {
		if f.Mode().IsDir() {
			// Append subdirectories into *dirs*, and traverse them later.
			dirs = append(dirs, f.Name())
		} else if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), ".yang") {
			// Only search for files with suffix of `.yang`
			fullpath := dir + "/" + f.Name()
			currDir := path.Dir(fullpath)
			// currDir is the name of current directory containing this module file.
			currDir = currDir[strings.LastIndex(currDir, "/")+1:]
			file, err := os.Open(fullpath)
			if err != nil {
				return nil, fmt.Errorf("traverseDir: failed to open file %s: %v", fullpath, err)
			}
			// *name* is name of yang module/submodule, we get it by removing `.yang` from original file name.
			name := f.Name()[:len(f.Name())-5]

			// Check whether found module is under `models` directory or `ietf` dirctory.
			if strings.Contains(fullpath, modelKeyword) {
				// Modules/submodules under `models` dirctory are under subdirectory (*currDir*) of `models` directory.
				// Store url of found module/submodule in urlMap.
				urlMap[name] = modelURL + "/" + currDir + "/" + f.Name()
			} else {
				// Found modules must be either under `models` directory or `ietf` directory.
				if !(strings.Contains(fullpath, ietfKeyword)) {
					return nil, fmt.Errorf("traverseDir: model %s not in either models dir or ietd dir", f.Name())
				}
				// `ietf` directory does not have subdirectory.
				// Store url of found module/submodule in urlMap.
				urlMap[name] = ietfURL + "/" + f.Name()
			}
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()

			// Only append `module` and ignore `submodules`.
			if strings.HasPrefix(strings.TrimSpace(scanner.Text()), "module ") {
				names = append(names, name)
			} else {
				glog.Info(fullpath + " does not have modules")
			}
			file.Close()
		}
	}

	// Traverse found subdirectories.
	for _, dirName := range dirs {
		newnames, err := traverseDir(dir + "/" + dirName)
		if err != nil {
			return nil, err
		}
		names = append(names, newnames...)
	}
	return names, nil
}

var stop = os.Exit

func main() {
	var format string
	formats := make([]string, 0, len(formatters))
	for k := range formatters {
		formats = append(formats, k)
	}
	sort.Strings(formats)

	var traceP string
	var help bool
	var paths []string
	getopt.ListVarLong(&paths, "path", 'p', "comma separated list of directories to add to search path", "DIR[,DIR...]")
	getopt.StringVarLong(&format, "format", 'f', "format to display: "+strings.Join(formats, ", "), "FORMAT")
	getopt.StringVarLong(&traceP, "trace", 't', "write trace into to TRACEFILE", "TRACEFILE")
	getopt.BoolVarLong(&help, "help", 'h', "display help")
	getopt.BoolVarLong(&yang.ParseOptions.IgnoreSubmoduleCircularDependencies, "ignore-circdep", 'g', "ignore circular dependencies between submodules")
	getopt.SetParameters("[FORMAT OPTIONS] [SOURCE] [...]")

	if err := getopt.Getopt(func(o getopt.Option) bool {
		if o.Name() == "--format" {
			f, ok := formatters[format]
			if !ok {
				fmt.Fprintf(os.Stderr, "%s: invalid format.  Choices are %s\n", format, strings.Join(formats, ", "))
				stop(1)
			}
			if f.flags != nil {
				f.flags.VisitAll(func(o getopt.Option) {
					getopt.AddOption(o)
				})
			}
		}
		return true
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		getopt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	if traceP != "" {
		fp, err := os.Create(traceP)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		trace.Start(fp)
		stop = func(c int) { trace.Stop(); os.Exit(c) }
		defer func() { trace.Stop() }()
	}

	if help {
		getopt.CommandLine.PrintUsage(os.Stderr)
		fmt.Fprintf(os.Stderr, `
SOURCE may be a module name or a .yang file.

Formats:
`)
		for _, fn := range formats {
			f := formatters[fn]
			fmt.Fprintf(os.Stderr, "    %s - %s\n", f.name, f.help)
			if f.flags != nil {
				f.flags.PrintOptions(indent.NewWriter(os.Stderr, "   "))
			}
			fmt.Fprintln(os.Stderr)
		}
		stop(0)
	}

	for _, path := range paths {
		expanded, err := yang.PathsWithModules(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		yang.AddPath(expanded...)
	}

	if format == "" {
		format = "tree"
	}
	if _, ok := formatters[format]; !ok {
		fmt.Fprintf(os.Stderr, "%s: invalid format.  Choices are %s\n", format, strings.Join(formats, ", "))
		stop(1)
	}

	// files stores names of all modules to search for.
	files := getopt.Args()

	ms := yang.NewModules()

	// If names of modules to search for is not given, we traverse all given path to find
	// all modules in these directories and crawl them later.
	if len(files) == 0 {
		for _, path := range paths {
			names, err := traverseDir(path)
			if err != nil {
				glog.Fatalf("traverse directory %s failed: %v", path, err)
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

	// Connect to DB
	if err := db.ConnectDB(); err != nil {
		glog.Fatalf("Connect to db failed: %v", err)
		stop(1)
	}
	defer db.Close()

	// Convert all found modules into ygot go structure of Module and insert them into database.
	for _, n := range names {
		module := &oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module{
			Name:      &mods[n].Name,
			Namespace: &mods[n].Namespace.Name,
			Prefix:    &mods[n].Prefix.Name,
			Summary:   &mods[n].Description.Name,
		}

		version, err := yang.MatchingExtensions(mods[n], "openconfig-extensions", "openconfig-version")
		if err != nil || len(version) == 0 {
			glog.Infof("%s do not have version\n", mods[n].Name)
			continue
		}

		module.Version = &version[0].Argument

		// If there are multiple revisions, we directly use the lastest one.
		if len(mods[n].Revision) > 0 {
			module.Revision = &mods[n].Revision[0].Name
		}
		for i := 0; i < len(mods[n].Import); i++ {
			module.GetOrCreateDependencies().RequiredModule = append(module.GetOrCreateDependencies().RequiredModule, mods[n].Import[i].Name)
		}
		for i := 0; i < len(mods[n].Include); i++ {
			module.GetOrCreateSubmodules().GetOrCreateSubmodule(mods[n].Include[i].Name)
			submoduleURL := urlMap[mods[n].Include[i].Name]
			if submoduleURL == "" {
				glog.Fatalf("cannot find url of submodule: %s", mods[n].Include[i].Name)
				continue
			}
			module.GetOrCreateSubmodules().GetOrCreateSubmodule(mods[n].Include[i].Name).GetOrCreateAccess().Uri = &submoduleURL

		}
		moduleURL := urlMap[n]
		module.GetOrCreateAccess().Uri = &moduleURL

		// Serialize module struct into json for insertion.
		json, err := ygot.EmitJSON(module, &ygot.EmitJSONConfig{
			Format: ygot.RFC7951,
			Indent: "  ",
			RFC7951Config: &ygot.RFC7951JSONConfig{
				AppendModuleName: true,
			},
		})
		if err != nil {
			glog.Fatalf("Marshalling into json string failed\n")
		}

		if err := db.InsertModule(orgName, module.GetName(), module.GetVersion(), json); err != nil {
			glog.Fatalf("Insert module %s failed: %v\n", module.GetName(), err)
			break
		}
	}

}
