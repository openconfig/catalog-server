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

// ValidateFeatureBundle is used to validate whether the input JSON data is in correct format.
// It takes a JSON string *data*, and returns a pointer to FeatureBundle if *data* is in correct format.
// Otherwise, the fucntion returns an error explaining why validation fails.
func ValidateFeatureBundle(data string) (*oc.OpenconfigModuleCatalog_Organizations_Organization_FeatureBundles_FeatureBundle, error) {
	featureBundle := &oc.OpenconfigModuleCatalog_Organizations_Organization_FeatureBundles_FeatureBundle{}

	// First check whether the data can be correctly unmarshalled back to a FeatureBundle struct.
	if err := oc.Unmarshal([]byte(data), featureBundle); err != nil {
		return nil, fmt.Errorf("ValidateFeatureBundle: cannot unmarshal JSON: %v", err)
	}

	// Check whether Validate function for FeatureBundle struct that comes from ygot package could pass.
	if err := featureBundle.Validate(); err != nil {
		return nil, fmt.Errorf("ValidateFeatureBundle: Validate function failed: %v", err)
	}

	// Check whether this featureBundle contains non-empty key of FeatureBundle (i.e, name and version).
	// If not, then return an error.
	if featureBundle.GetName() == "" || featureBundle.GetVersion() == "" {
		return nil, fmt.Errorf("ValidateFeatureBundle: FeatureBundle cannot have empty name or version")
	}
	return featureBundle, nil
}
