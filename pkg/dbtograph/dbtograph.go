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
// It return slice of graphQL module pointers and error if there is any.
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
// It return slice of graphQL FeatureBundle pointers and error if there is any.
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
