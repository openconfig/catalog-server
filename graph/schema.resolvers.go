package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openconfig/catalog-server/graph/generated"
	"github.com/openconfig/catalog-server/graph/model"
	"github.com/openconfig/catalog-server/pkg/access"
	"github.com/openconfig/catalog-server/pkg/db"
	"github.com/openconfig/catalog-server/pkg/dbtograph"
	"github.com/openconfig/catalog-server/pkg/validate"
)

func (r *mutationResolver) CreateModule(ctx context.Context, input model.NewModule, token string) (string, error) {
	failMsg := `Fail`
	successMsg := `Success`

	// Validate the token and check whether it contains access to certain organization
	if err := access.CheckAccess(token, input.OrgName); err != nil {
		return failMsg, fmt.Errorf("CreateModule: validate token failed: %v", err)
	}

	// Validate module
	module, err := validate.ValidateModule(input.Data)
	if err != nil {
		return failMsg, fmt.Errorf("CreateModule: validate module failed: %v", err)
	}

	// Insert module if not exist, or update it.
	if err := db.InsertModule(input.OrgName, module.GetName(), module.GetVersion(), input.Data); err != nil {
		return failMsg, fmt.Errorf("CreateModule failed: %v", err)
	}

	return successMsg, nil
}

func (r *mutationResolver) DeleteModule(ctx context.Context, input model.ModuleKey, token string) (string, error) {
	failMsg := `Fail`
	successMsg := `Success`

	// Validate the token and check whether it contains access to certain organization
	if err := access.CheckAccess(token, input.OrgName); err != nil {
		return failMsg, fmt.Errorf("DeleteModule: validate token failed: %v", err)
	}

	// Delete a module
	if err := db.DeleteModule(input.OrgName, input.Name, input.Version); err != nil {
		return failMsg, fmt.Errorf("DeleteModule failed: %v", err)
	}

	return successMsg, nil
}

func (r *mutationResolver) CreateFeatureBundle(ctx context.Context, input model.NewFeatureBundle, token string) (string, error) {
	failMsg := `Fail`
	successMsg := `Success`

	// Validate the token and check whether it contains access to certain organization
	if err := access.CheckAccess(token, input.OrgName); err != nil {
		return failMsg, fmt.Errorf("CreateFeatureBundle: validate token failed: %v", err)
	}

	// Validate module
	featureBundle, err := validate.ValidateFeatureBundle(input.Data)
	if err != nil {
		return failMsg, fmt.Errorf("CreateFeatureBundle: validate featureBundle failed: %v", err)
	}

	// Insert module if not exist, or update it.
	if err := db.InsertFeatureBundle(input.OrgName, featureBundle.GetName(), featureBundle.GetVersion(), input.Data); err != nil {
		return failMsg, fmt.Errorf("CreateFeatureBundle failed: %v", err)
	}

	return successMsg, nil
}

func (r *mutationResolver) DeleteFeatureBundle(ctx context.Context, input model.FeatureBundleKey, token string) (string, error) {
	failMsg := `Fail`
	successMsg := `Success`

	// Validate the token and check whether it contains access to certain organization
	if err := access.CheckAccess(token, input.OrgName); err != nil {
		return failMsg, fmt.Errorf("DeleteFeatureBundle: validate token failed: %v", err)
	}

	// Delete a module
	if err := db.DeleteFeatureBundle(input.OrgName, input.Name, input.Version); err != nil {
		return failMsg, fmt.Errorf("DeleteFeatureBundle failed: %v", err)
	}

	return successMsg, nil
}

func (r *queryResolver) ModulesByOrgName(ctx context.Context, orgName *string) ([]*model.Module, error) {
	dbModules, err := db.QueryModulesByOrgName(orgName)
	if err != nil {
		return nil, err
	}
	return dbtograph.ModuleToGraphQL(dbModules)
}

func (r *queryResolver) ModulesByKey(ctx context.Context, name *string, version *string) ([]*model.Module, error) {
	dbModules, err := db.QueryModulesByKey(name, version)
	if err != nil {
		return nil, err
	}
	return dbtograph.ModuleToGraphQL(dbModules)
}

func (r *queryResolver) FeatureBundlesByOrgName(ctx context.Context, orgName *string) ([]*model.FeatureBundle, error) {
	dbFeatureBundles, err := db.QueryFeatureBundlesByOrgName(orgName)
	if err != nil {
		return nil, err
	}
	return dbtograph.FeatureBundleToGraphQL(dbFeatureBundles)
}

func (r *queryResolver) FeatureBundlesByKey(ctx context.Context, name *string, version *string) ([]*model.FeatureBundle, error) {
	dbFeatureBundles, err := db.QueryFeatureBundlesByKey(name, version)
	if err != nil {
		return nil, err
	}
	return dbtograph.FeatureBundleToGraphQL(dbFeatureBundles)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
