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
		desc    string
		inputs  []db.Module
		want    []model.Module
		wantErr bool
	}{
		{
			desc: "two modules with filled-out fields",
			inputs: []db.Module{
				{
					OrgName: "org_A",
					Name:    "module_A",
					Version: "version_A",
					Data:    `{"openconfig-module-catalog:name": "module_A", "openconfig-module-catalog:access": {"uri": "testlink_A"}, "openconfig-module-catalog:version": "version_A", "openconfig-module-catalog:summary": "foo"}`,
				},
				{
					OrgName: "org_B",
					Name:    "module_B",
					Version: "version_B",
					Data:    `{"openconfig-module-catalog:name": "module_B",  "openconfig-module-catalog:access": {"uri": "testlink_B"}, "openconfig-module-catalog:version": "version_B", "openconfig-module-catalog:summary": "bar"}`,
				},
			},
			want: []model.Module{
				{
					OrgName: "org_A",
					Name:    "module_A",
					Version: "version_A",
					URL:     "testlink_A",
					Summary: "foo",
					Data:    `{"openconfig-module-catalog:name": "module_A", "openconfig-module-catalog:access": {"uri": "testlink_A"}, "openconfig-module-catalog:version": "version_A", "openconfig-module-catalog:summary": "foo"}`,
				},
				{
					OrgName: "org_B",
					Name:    "module_B",
					Version: "version_B",
					URL:     "testlink_B",
					Summary: "bar",
					Data:    `{"openconfig-module-catalog:name": "module_B",  "openconfig-module-catalog:access": {"uri": "testlink_B"}, "openconfig-module-catalog:version": "version_B", "openconfig-module-catalog:summary": "bar"}`,
				},
			},
			wantErr: false,
		},
		{
			desc:   "empty input",
			inputs: []db.Module{},
			want:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			modules, err := ModuleToGraphQL(tc.inputs)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr mismatch, err: %v, wantErr: %t", err, tc.wantErr)
			}
			for i := 0; i < len(modules); i++ {
				if diff := cmp.Diff(*modules[i], tc.want[i]); diff != "" {
					t.Errorf("module mismatch:\n%s", diff)
				}
			}
		})
	}

}

func TestFeatureBundleToGraphQL(t *testing.T) {
	tests := []struct {
		inputs  []db.FeatureBundle
		want    []model.FeatureBundle
		wantErr bool
	}{
		{
			inputs: []db.FeatureBundle{
				{
					OrgName: "org_A",
					Name:    "feature_A",
					Version: "version_A",
					Data:    `{"openconfig-module-catalog:name": "feature_A","openconfig-module-catalog:version": "version_A"}`,
				},
				{
					OrgName: "org_B",
					Name:    "feature_B",
					Version: "version_B",
					Data:    `{"openconfig-module-catalog:name": "feature_B","openconfig-module-catalog:version": "version_A"}`,
				},
			},
			want: []model.FeatureBundle{
				{
					OrgName: "org_A",
					Name:    "feature_A",
					Version: "version_A",
					Data:    `{"openconfig-module-catalog:name": "feature_A","openconfig-module-catalog:version": "version_A"}`,
				},
				{
					OrgName: "org_B",
					Name:    "feature_B",
					Version: "version_B",
					Data:    `{"openconfig-module-catalog:name": "feature_B","openconfig-module-catalog:version": "version_A"}`,
				},
			},
			wantErr: false,
		},
		{
			inputs: []db.FeatureBundle{},
			want:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestModuleToGraphQL, number of inputs: %d", len(tc.inputs)), func(t *testing.T) {
			featureBundles, err := FeatureBundleToGraphQL(tc.inputs)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr mismatch, err: %v, wantErr: %t", err, tc.wantErr)
			}
			for i := 0; i < len(featureBundles); i++ {
				if diff := cmp.Diff(*featureBundles[i], tc.want[i]); diff != "" {
					t.Errorf("featureBundle mismatch")
				}
			}
		})
	}
}
