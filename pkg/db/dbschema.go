package db

// This file contains definition for structs in database schema

// Module is struct of Module table in db schema.
type Module struct {
	OrgName string // OrgName column refers to name of organization's name holding this Module.
	Name    string // Namme column refers to name of this Module.
	Version string // Version column refers to version of this Module.
	Data    string // Data column refers to json format string of this Module in YANG schema.
}
