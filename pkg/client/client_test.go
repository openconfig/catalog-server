package client

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	// File name of testdata of mocked response template, server could set queryName inside this response.
	testdataFile = `testdata/testdata`
	// Name of filed containing raw JSON string of Module in response.
	rawFieldName = `Data`
)

// ReadTestData reads testdata from file.
func ReadTestData(filename string) (string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var testdata string

	// Append testdata line by line.
	for scanner.Scan() {
		// Strip \newline in each line, otherwise parsing string back to Module struct would fail.
		testdata += scanner.Text()
	}
	return testdata, nil
}

// SetupHTTPTestServer sets up mock server for testing purpose.
func SetupHTTPTestServer() (*httptest.Server, error) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	respTemplate, err := ReadTestData(testdataFile)
	if err != nil {
		return nil, fmt.Errorf("Read test data failed: %v", err)
	}
	// Setup query handler functions.
	mux.HandleFunc(`/query`, (func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()["query"]) == 1 {
			if strings.Contains(r.URL.Query()["query"][0], "ModulesByOrgName") {
				fmt.Fprintf(w, fmt.Sprintf(respTemplate, "ModulesByOrgName"))
			} else if strings.Contains(r.URL.Query()["query"][0], "ModulesByKey") {
				fmt.Fprintf(w, fmt.Sprintf(respTemplate, "ModulesByKey"))
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
		}

	}))

	return ts, nil
}

// TestQuery serves as a function to test all three functions inside client pkg,
// including: (1) FormatHTTPQuery, (2) QueryServer, and (3) ParseModule.
func TestQuery(t *testing.T) {
	ts, err := SetupHTTPTestServer()
	if err != nil {
		t.Errorf("Set up HTTP server failed: %v", err)
		// Stop this test if server failed to launch.
		return
	}
	defer ts.Close()

	tests := []struct {
		desc            string
		query           string
		queryName       string
		wantModuleNames []string
	}{
		{
			desc:            `Test query ModuleByOrgName`,
			query:           `{ModulesByOrgName(OrgName:"1"){Data}}`,
			queryName:       `ModulesByOrgName`,
			wantModuleNames: []string{`0`, `1`, `2`},
		},
		{
			desc:            `Test query ModuleByKey`,
			query:           `{ModulesByOrgName(Name:"1"){Data}}`,
			queryName:       `ModulesByOrgName`,
			wantModuleNames: []string{`0`, `1`, `2`},
		},
	}
	c, err := NewClient("")
	if err != nil {
		t.Errorf("NewClient failed: %v", err)
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestQueryServer %s", tc.desc), func(t *testing.T) {
			resp, err := c.QueryServer(ts.URL + "/query?query=" + tc.query)
			if err != nil {
				t.Errorf("Send query to server failed: %v", err)
			}
			modules, err := ParseModule(resp, rawFieldName, tc.queryName)
			if err != nil {
				t.Errorf("Parse Module failed: %v", err)
			}
			if len(modules) != len(tc.wantModuleNames) {
				t.Errorf("Parsed Modules length mismatch, get: %d, want: %d", len(modules), len(tc.wantModuleNames))
			}
			for i := 0; i < len(modules); i++ {
				if modules[i].GetName() != tc.wantModuleNames[i] {
					t.Errorf("Parsed Module name mismatch, get: %s, want: %s", modules[i].GetName(), tc.wantModuleNames[i])
				}
			}
		})
	}
}
