package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/openconfig/catalog-server/pkg/client"
)

// Demo program to show how to use client library functions.
func main() {
	var queryPtr = flag.String("query", "", "query url")
	var tokenPathPtr = flag.String("token", "", "file path of auth token")
	var queryNamePtr = flag.String("queryName", "", "name of query")
	var fieldNamePtr = flag.String("filedName", "Data", "name of field in response containing json data")

	flag.Parse()

	// Format HTTP query.
	req, err := client.FormatHTTPQuery(*queryPtr, *tokenPathPtr)
	if err != nil {
		glog.Fatal("Format query failed: %v", err)
	}

	// Send query to server and receive response.
	resp, err := client.QueryServer(req)
	if err != nil {
		glog.Fatal("Query server failed: %v", err)
	}

	// Parse query results.
	modules, err := client.ParseModule(resp, *fieldNamePtr, *queryNamePtr)
	if err != nil {
		glog.Fatal("Parse into modules failed: %v", err)
	}

	// Print out names of all matched modules just for testing parsing is correct.
	for i := 0; i < len(modules); i++ {
		fmt.Println("Name of matched module: ", modules[i].GetName())
	}

	// Users can further operate over YANG go structs (modules).

}
