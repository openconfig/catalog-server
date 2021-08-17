/*
Package validate contains functions to validate JSON strings
  which come with mutation operations (e.g., create or update)
  before such data are persisted into underlying storage.
*/
package validate

import (
	"fmt"

	oc "github.com/openconfig/catalog-server/pkg/ygotgen"
)

// ValidateModule is used to validate whether the input JSON data is in correct format.
// It takes a JSON string *data*, and returns a pointer to Module if *data* is in correct format.
// Otherwise, the function returns an error explaining why validation fails.
func ValidateModule(data string) (*oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module, error) {
	module := &oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module{}

	// First check whether the data can be correctly unmarshalled back to a Module struct.
	if err := oc.Unmarshal([]byte(data), module); err != nil {
		return nil, fmt.Errorf("ValidateModule: cannot unmarshal JSON: %v", err)
	}

	// Check whether Validate function for Module struct that comes from ygot package could pass.
	if err := module.Validate(); err != nil {
		return nil, fmt.Errorf("ValidateModule: Validate function failed: %v", err)
	}

	// Check whether this module contains non-empty key of Module (i.e, name and version).
	// If not, then return an error.
	if module.GetName() == "" || module.GetVersion() == "" {
		return nil, fmt.Errorf("ValidateModule: Module cannot have empty name or version")
	}
	return module, nil
}
