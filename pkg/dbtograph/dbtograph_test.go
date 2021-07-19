package dbtograph

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/catalog-server/graph/model"
	"github.com/openconfig/catalog-server/pkg/db"
)

func TestModuleToGraphQL(t *testing.T) {
	tests := []struct {
		inputs []db.Module
		want   []model.Module
	}{
		{
			inputs: []db.Module{
				{
					OrgName: "org_A",
					Name:    "module_A",
					Version: "version_A",
					Data:    "{}",
				},
				{
					OrgName: "org_B",
					Name:    "module_B",
					Version: "version_B",
					Data:    "{}",
				},
			},
			want: []model.Module{
				{
					OrgName: "org_A",
					Name:    "module_A",
					Version: "version_A",
					Data:    "{}",
				},
				{
					OrgName: "org_B",
					Name:    "module_B",
					Version: "version_B",
					Data:    "{}",
				},
			},
		},
		{
			inputs: []db.Module{},
			want:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestModuleToGraphQL, number of inputs: %d", len(tc.inputs)), func(t *testing.T) {
			modules := ModuleToGraphQL(tc.inputs)
			for i := 0; i < len(modules); i++ {
				if diff := cmp.Diff(*modules[i], tc.want[i]); diff != "" {
					t.Errorf("module mismatch")
				}
			}
		})
	}

}
