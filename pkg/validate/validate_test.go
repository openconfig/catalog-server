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

package validate

import (
	"testing"
)

func TestValidateModule(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{
			input:   ``,
			wantErr: true,
		},
		{
			input:   `{}`,
			wantErr: true,
		},
		{
			input:   `{"openconfig-module-catalog:name": "", "openconfig-module-catalog:version": "version1"}`,
			wantErr: true,
		},
		{
			input:   `{"openconfig-module-catalog:name": "name1", "openconfig-module-catalog:version": "version1"}`,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		if _, err := ValidateModule(tc.input); (err != nil) != tc.wantErr {
			t.Errorf("ValidateModule test with string failed, input: %s, wantErr: %t, err: %v", tc.input, tc.wantErr, err)
		}
	}
}

func TestValidateFeatureBundle(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{
			input:   ``,
			wantErr: true,
		},
		{
			input:   `{}`,
			wantErr: true,
		},
		{
			input:   `{"openconfig-module-catalog:name": "", "openconfig-module-catalog:version": "version1"}`,
			wantErr: true,
		},
		{
			input:   `{"openconfig-module-catalog:name": "name1", "openconfig-module-catalog:version": "version1"}`,
			wantErr: false,
		},
	}
	for _, tc := range tests {
		if _, err := ValidateFeatureBundle(tc.input); (err != nil) != tc.wantErr {
			t.Errorf("ValidateFeatureBundle test with string failed, input: %s, wantErr: %t, err: %v", tc.input, tc.wantErr, err)
		}
	}
}
