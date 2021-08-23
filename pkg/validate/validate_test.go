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
