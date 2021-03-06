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
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/openconfig/catalog-server/pkg/client"
)

// const values for demo purposes.
const (
	hostURL    = `https://helloworld-jpx33sh7ha-uc.a.run.app/query` // Address of deployed server.
	graphQuery = `{ModulesByOrgName(OrgName:"1"){Data}}`            // GraphQL query.
	queryName  = `ModulesByOrgName`                                 // Name of graphQL query.
)

// Demo program to show how to use client library functions.
func main() {
	var tokenPathPtr = flag.String("token", "token", "file path of auth token")
	flag.Parse()

	// Format whole query URL.
	queryURL := hostURL + `?query=` + graphQuery

	// Initialize a new Client.
	c, err := client.NewClient(*tokenPathPtr)
	if err != nil {
		glog.Fatalf("catalog client: new client failed: %v", err)
	}

	// Send query to server and receive response.
	resp, err := c.QueryServer(queryURL)
	if err != nil {
		glog.Fatalf("catalog client: query server failed: %v", err)
	}

	// Parse query results.
	modules, err := client.ParseModule(resp, queryName)
	if err != nil {
		glog.Fatalf("catalog client: parse response into go catalog modules failed: %v", err)
	}

	// Print out names of all matched modules just for testing parsing is correct.
	for i := 0; i < len(modules); i++ {
		fmt.Println("Name of matched module: ", modules[i].GetName())
	}

	// Users can further operate over YANG go structs (modules).

}
