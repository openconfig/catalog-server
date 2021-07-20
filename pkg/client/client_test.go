package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	// Mocked response template, server could set queryName inside this response.
	respTemplate = `{"data":{"%s":[{"Data":"{\"openconfig-module-catalog:name\": \"0\", \"openconfig-module-catalog:access\": {\"uri\": \"CPx\", \"md5-hash\": \"YC\"}, \"openconfig-module-catalog:prefix\": \"392Qv\", \"openconfig-module-catalog:summary\": \"e2W\", \"openconfig-module-catalog:version\": \"0\", \"openconfig-module-catalog:revision\": \"s\", \"openconfig-module-catalog:namespace\": \"V\", \"openconfig-module-catalog:submodules\": {\"submodule\": [{\"name\": \"1aX0v8\", \"access\": {\"uri\": \"63yc\", \"md5-hash\": \"wFJ3at\"}}, {\"name\": \"A94l\", \"access\": {\"uri\": \"LMSDG\", \"md5-hash\": \"fqF\"}}]}, \"openconfig-module-catalog:dependencies\": {\"required-module\": [\"rT\"]}, \"openconfig-module-catalog:classification\": {\"category\": \"openconfig-catalog-types:IETF_NETWORK_SERVICE\", \"subcategory\": \"openconfig-catalog-types:IETF_MODEL_TYPE\", \"deployment-status\": \"openconfig-catalog-types:EXPERIMENTAL\"}}"},{"Data":"{\"openconfig-module-catalog:name\": \"1\", \"openconfig-module-catalog:access\": {\"uri\": \"i\", \"md5-hash\": \"0DkF\"}, \"openconfig-module-catalog:prefix\": \"whAC\", \"openconfig-module-catalog:summary\": \"X5h\", \"openconfig-module-catalog:version\": \"1\", \"openconfig-module-catalog:revision\": \"OOq\", \"openconfig-module-catalog:namespace\": \"wnTaF\", \"openconfig-module-catalog:submodules\": {\"submodule\": [{\"name\": \"7R\", \"access\": {\"uri\": \"zV\", \"md5-hash\": \"DCdE\"}}, {\"name\": \"9E1\", \"access\": {\"uri\": \"lA0PP1\", \"md5-hash\": \"TLI\"}}, {\"name\": \"Cn94Q\", \"access\": {\"uri\": \"30nsE\", \"md5-hash\": \"zF5gm\"}}, {\"name\": \"de986\", \"access\": {\"uri\": \"KxGzK\", \"md5-hash\": \"JL\"}}]}, \"openconfig-module-catalog:dependencies\": {\"required-module\": [\"b\", \"iuaWU\"]}, \"openconfig-module-catalog:classification\": {\"category\": \"openconfig-catalog-types:IETF_MODEL_LAYER\", \"subcategory\": \"openconfig-catalog-types:IETF_TYPE_STANDARD\", \"deployment-status\": \"openconfig-catalog-types:PRODUCTION\"}}"},{"Data":"{\"openconfig-module-catalog:name\": \"2\", \"openconfig-module-catalog:access\": {\"uri\": \"xKbeLI\", \"md5-hash\": \"oep\"}, \"openconfig-module-catalog:prefix\": \"6gZ\", \"openconfig-module-catalog:summary\": \"3wb1lx\", \"openconfig-module-catalog:version\": \"2\", \"openconfig-module-catalog:revision\": \"zY82y\", \"openconfig-module-catalog:namespace\": \"E18EsV\", \"openconfig-module-catalog:submodules\": {\"submodule\": [{\"name\": \"26vi\", \"access\": {\"uri\": \"qMvvLw\", \"md5-hash\": \"TM7\"}}, {\"name\": \"Uf5Op\", \"access\": {\"uri\": \"mrXSm\", \"md5-hash\": \"5Gh\"}}]}, \"openconfig-module-catalog:dependencies\": {\"required-module\": [\"e5IN4\", \"b\", \"yD\"]}, \"openconfig-module-catalog:classification\": {\"category\": \"openconfig-catalog-types:IETF_NETWORK_ELEMENT\", \"subcategory\": \"openconfig-catalog-types:IETF_MODEL_TYPE\", \"deployment-status\": \"openconfig-catalog-types:EXPERIMENTAL\"}}"},{"Data":"{\"openconfig-module-catalog:name\": \"3\", \"openconfig-module-catalog:access\": {\"uri\": \"YZE\", \"md5-hash\": \"fE\"}, \"openconfig-module-catalog:prefix\": \"UWr\", \"openconfig-module-catalog:summary\": \"w3B\", \"openconfig-module-catalog:version\": \"3\", \"openconfig-module-catalog:revision\": \"zgCUn\", \"openconfig-module-catalog:namespace\": \"XAZ\", \"openconfig-module-catalog:submodules\": {\"submodule\": [{\"name\": \"pxN\", \"access\": {\"uri\": \"V65HM\", \"md5-hash\": \"Z\"}}]}, \"openconfig-module-catalog:dependencies\": {\"required-module\": [\"dDGxSc\", \"GM\", \"SllRk\", \"2\"]}, \"openconfig-module-catalog:classification\": {\"category\": \"openconfig-catalog-types:IETF_NETWORK_SERVICE\", \"subcategory\": \"openconfig-catalog-types:IETF_TYPE_VENDOR\", \"deployment-status\": \"openconfig-catalog-types:PRODUCTION\"}}"},{"Data":"{\"openconfig-module-catalog:name\": \"4\", \"openconfig-module-catalog:access\": {\"uri\": \"9c\", \"md5-hash\": \"ckdrUa\"}, \"openconfig-module-catalog:prefix\": \"I26\", \"openconfig-module-catalog:summary\": \"Axl\", \"openconfig-module-catalog:version\": \"4\", \"openconfig-module-catalog:revision\": \"0\", \"openconfig-module-catalog:namespace\": \"0D\", \"openconfig-module-catalog:submodules\": {\"submodule\": [{\"name\": \"nEa76V\", \"access\": {\"uri\": \"kv\", \"md5-hash\": \"vFF\"}}, {\"name\": \"x0O6Rb\", \"access\": {\"uri\": \"V1\", \"md5-hash\": \"XuUp\"}}]}, \"openconfig-module-catalog:dependencies\": {\"required-module\": [\"cKxRE\", \"WzvTc\"]}, \"openconfig-module-catalog:classification\": {\"category\": \"openconfig-catalog-types:IETF_NETWORK_SERVICE\", \"subcategory\": \"openconfig-catalog-types:IETF_TYPE_STANDARD\", \"deployment-status\": \"openconfig-catalog-types:EXPERIMENTAL\"}}"}]}}`
	fieldName    = `Data`
)

func setupHTTPTestServer() *httptest.Server {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)

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

	return ts
}

// TestQuery serves as a function to test all three functions inside client pkg,
// including: (1) FormatHTTPQuery, (2) QueryServer, and (3) ParseModule.
func TestQuery(t *testing.T) {
	ts := setupHTTPTestServer()
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
			wantModuleNames: []string{`0`, `1`, `2`, `3`, `4`},
		},
		{
			desc:            `Test query ModuleByKey`,
			query:           `{ModulesByOrgName(Name:"1"){Data}}`,
			queryName:       `ModulesByOrgName`,
			wantModuleNames: []string{`0`, `1`, `2`, `3`, `4`},
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestQueryServer %s", tc.desc), func(t *testing.T) {
			req, err := FormatHTTPQuery(ts.URL+"/query?query="+tc.query, "")
			if err != nil {
				t.Errorf("Format http query failed: %v", err)
			}
			resp, err := QueryServer(req)
			if err != nil {
				t.Errorf("Send query to server failed: %v", err)
			}
			modules, err := ParseModule(resp, fieldName, tc.queryName)
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
