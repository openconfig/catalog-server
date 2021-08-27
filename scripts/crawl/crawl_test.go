// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"testing"

	oc "github.com/openconfig/catalog-server/pkg/ygotgen"
)

func TestCrawlModules(t *testing.T) {
	tests := []struct {
		name    string
		version string
		summary string
	}{
		{
			name:    "aug",
			version: "1.2.0",
			summary: "aug desc",
		},
		{
			name:    "base",
			version: "1.3.0",
			summary: "base desc",
		},
		{
			name:    "other",
			version: "1.4.0",
			summary: "other desc",
		},
		{
			name:    "subdir1",
			version: "1.5.0",
			summary: "subdir1 desc",
		},
	}

	paths := []string{"./models"}
	// Dummy url
	url := "https://github.com/openconfig/goyang/tree/master/testdata"
	mods, names := crawlModules(paths, url)
	var modules []*oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module
	for _, n := range names {
		module := populateModule(mods[n], n)
		if module == nil {
			continue
		}
		modules = append(modules, module)
	}

	for idx, tc := range tests {
		t.Run(fmt.Sprintf("TestCrawlModules name: %s, version: %s, summary: %s", tc.name, tc.version, tc.summary), func(t *testing.T) {
			if modules[idx].GetName() != tc.name || modules[idx].GetVersion() != tc.version || modules[idx].GetSummary() != tc.summary {
				t.Errorf("crawled modules mismatch, name: %s, version: %s summary: %s, expected name: %s, expected version: %s, expected summary: %s", modules[idx].GetName(), modules[idx].GetVersion(), modules[idx].GetSummary(), tc.name, tc.version, tc.summary)
			}
		})
	}
}
