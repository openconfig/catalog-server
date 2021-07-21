package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/openconfig/catalog-server/pkg/client"
)

// const values for demo purposes.
const (
	rawFieldName = `Data`                                             // Name of field in response containing json data.
	hostURL      = `https://helloworld-jpx33sh7ha-uc.a.run.app/query` // Address of deployed server.
	graphQuery   = `{ModulesByOrgName(OrgName:"1"){Data}}`            // GraphQL query.
	queryName    = `ModulesByOrgName`                                 // Name of graphQL query.
)

// Demo program to show how to use client library functions.
func main() {
	var tokenPathPtr = flag.String("token", "", "file path of auth token")
	flag.Parse()

	// Format whole query URL.
	queryURL := hostURL + `?query=` + graphQuery

	// Format HTTP query.
	req, err := client.FormatHTTPQuery(queryURL, *tokenPathPtr)
	if err != nil {
		glog.Fatalf("catalog client: format query failed: %v", err)
	}

	// Send query to server and receive response.
	resp, err := client.QueryServer(req)
	if err != nil {
		glog.Fatalf("catalog client: query server failed: %v", err)
	}

	// Parse query results.
	modules, err := client.ParseModule(resp, rawFieldName, queryName)
	if err != nil {
		glog.Fatalf("catalog client: parse response into go catalog modules failed: %v", err)
	}

	// Print out names of all matched modules just for testing parsing is correct.
	for i := 0; i < len(modules); i++ {
		fmt.Println("Name of matched module: ", modules[i].GetName())
	}

	// Users can further operate over YANG go structs (modules).

}
