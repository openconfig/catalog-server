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
