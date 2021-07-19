package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openconfig/catalog-server/graph/generated"
	"github.com/openconfig/catalog-server/graph/model"
	"github.com/openconfig/catalog-server/pkg/db"
	"github.com/openconfig/catalog-server/pkg/dbtograph"
)

func (r *queryResolver) ModulesByOrgName(ctx context.Context, orgName *string) ([]*model.Module, error) {
	dbModules, err := db.QueryModulesByOrgName(orgName)
	if err != nil {
		return nil, err
	}
	return dbtograph.ModuleToGraphQL(dbModules), nil
}

func (r *queryResolver) ModulesByKey(ctx context.Context, name *string, version *string) ([]*model.Module, error) {
	dbModules, err := db.QueryModulesByKey(name, version)
	if err != nil {
		return nil, err
	}
	return dbtograph.ModuleToGraphQL(dbModules), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
