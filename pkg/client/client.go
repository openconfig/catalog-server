/*
Package client contains functions of client side library.
It provides helper functions to write a client program to:
 * Format query to catalog server.
 * Receive and parse responses into go structs.
*/
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	oc "github.com/openconfig/catalog-server/pkg/ygotgen"
)

// All responses from server are put inside a field called `data` of json string.
const dataField = "data"

// ReadAuthToken takes in *filepath* of token file, reads token and returns token string.
// This token is used when server is deployed on Google Cloud Run and only avaiable to permitted users.
// In this case, users need to include a header with identity token to get access to catalog server.
// This token can be obtained via the command `glcoud auth print-identity-token`.
// Referecne: https://cloud.google.com/sdk/gcloud/reference/auth/print-identity-token.
func ReadAuthToken(filepath string) (string, error) {
	token, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("Read Auth token from file failed: %v", err)
	}
	return string(token), nil
}

// FormatHTTPQuery generates a HTTP request that can be used to query catalog server.
// It takes in query string and token's filepath as parameters.
// Example query looks like: HOST_ADDR/query?query=GRAPHQL_QUERY.
// If filepath is "", then it means no filepath is given, do not append "Auth" header.
func FormatHTTPQuery(query string, filepath string) (*http.Request, error) {
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, fmt.Errorf("Format new HTTP request failed: %v", err)
	}
	if filepath != "" {
		token, err := ReadAuthToken(filepath)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

// QueryServer sends formatted query request (*req*) to catalog server and returns string of boday if no errors encountered.
func QueryServer(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Send formatted query to server failed: %v", err)
	}
	// Check whether status code is OK.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Query does not receive Status OK 200, status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read response body failed: %v", err)
	}
	return string(body), nil
}

// ParseModule parses query results into slice of ygot go structs of Module.
// It takes in three parameters:
// * resp: query response string in json format.
// * fieldName: name of filed of json string of Module in response.
// * queryName: name of query users want to extract response from as GraphQL supports composing multiple queries into one request.
func ParseModule(resp string, fieldName string, queryName string) ([]oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module, error) {
	var moduleMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &moduleMap); err != nil {
		return nil, fmt.Errorf("Unmarshal response into json failed: %v", err)
	}
	// Retrive *dataField* field from response json string.
	// Then we retrive responses from filed of *queryName*.
	dataJSON := moduleMap[dataField].(map[string]interface{})[queryName].([]interface{})
	var modules []oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module
	for i := 0; i < len(dataJSON); i++ {
		// Retrieve json string of Module from *filedName* field.
		moduleJson := dataJSON[i].(map[string]interface{})[fieldName].(string)
		module := &oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module{}
		if err := oc.Unmarshal([]byte(moduleJson), module); err != nil {
			return nil, fmt.Errorf("Cannot unmarshal JSON: %v", err)
		}
		modules = append(modules, *module)
	}
	return modules, nil
}
