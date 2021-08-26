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

package db

// This file contains definition for structs in database schema

// Module is struct of Module table in db schema.
type Module struct {
	OrgName string // OrgName column refers to name of organization's name holding this Module.
	Name    string // Namme column refers to name of this Module.
	Version string // Version column refers to version of this Module.
	Summary string // Version column refers to summary of this Module.
	Data    string // Data column refers to json format string of this Module in YANG schema.
}

// FeatureBundle is struct of FeatureBundle table in db schema.
type FeatureBundle struct {
	OrgName string // OrgName column refers to name of organization's name holding this FeatureBundle.
	Name    string // Namme column refers to name of this FeatureBundle.
	Version string // Version column refers to version of this FeatureBundle.
	Data    string // Data column refers to json format string of this FeatureBundle in YANG schema.
}
