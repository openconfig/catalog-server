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
Package dbtograph works as decomposition between struct of db schema
  and graphQL response type.

It contains functions to convert data in struct of db schema
  into graphql expected response struct.
*/
package dbtograph

import (
	"fmt"

	oc "github.com/openconfig/catalog-server/pkg/ygotgen"

	"github.com/openconfig/catalog-server/graph/model"
	"github.com/openconfig/catalog-server/pkg/db"
)

// ModuleToGraphQL converts module schema in database to graphQL module response type.
// It returns a slice of graphQL module pointers and an error if there is any.
func ModuleToGraphQL(dbModules []db.Module) ([]*model.Module, error) {
	var models []*model.Module
	for i := 0; i < len(dbModules); i++ {
		module := &oc.OpenconfigModuleCatalog_Organizations_Organization_Modules_Module{}
		// First check whether the data can be correctly unmarshalled back to a Module struct.
		if err := oc.Unmarshal([]byte(dbModules[i].Data), module); err != nil {
			return nil, fmt.Errorf("ModuleToGraphQL: cannot unmarshal JSON: %v", err)
		}

		model := model.Module{
			OrgName: dbModules[i].OrgName,
			Name:    dbModules[i].Name,
			Version: dbModules[i].Version,
			Summary: module.GetSummary(),
			Data:    dbModules[i].Data,
		}
		if module.GetAccess() != nil {
			model.URL = module.GetAccess().GetUri()
		}
		models = append(models, &model)
	}
	return models, nil
}

// FeatureBundleToGraphQL converts FeatureBundle schema in database to graphQL FeatureBundle response type.
// It returns a slice of graphQL FeatureBundle pointers and an error if there is any.
func FeatureBundleToGraphQL(dbFeatureBundles []db.FeatureBundle) ([]*model.FeatureBundle, error) {
	var featureBundles []*model.FeatureBundle
	for i := 0; i < len(dbFeatureBundles); i++ {
		featureBundles = append(featureBundles, &model.FeatureBundle{
			OrgName: dbFeatureBundles[i].OrgName,
			Name:    dbFeatureBundles[i].Name,
			Version: dbFeatureBundles[i].Version,
			Data:    dbFeatureBundles[i].Data,
		})
	}
	return featureBundles, nil
}
