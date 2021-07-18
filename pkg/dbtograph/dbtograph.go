/*
Package dbtograph works as decomposition between struct of db schema
  and graphQL response type.

It contains functions to convert data in struct of db schema
  into graphql expected response struct.
*/
package dbtograph

import (
	"github.com/openconfig/catalog-server/graph/model"
	"github.com/openconfig/catalog-server/pkg/db"
)

// ModuleToGraphQL converts module schema in database to graphQL module response type.
// It return slice of graphQL module pointers.
func ModuleToGraphQL(dbModules []db.Module) []*model.Module {
	var models []*model.Module
	for i := 0; i < len(dbModules); i++ {
		models = append(models, &model.Module{
			OrgName: dbModules[i].OrgName,
			Name:    dbModules[i].Name,
			Version: dbModules[i].Version,
			Data:    dbModules[i].Data,
		})
	}
	return models
}
