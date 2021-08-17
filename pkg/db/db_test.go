package db

import (
	"fmt"
	"reflect"
	"testing"
)

// This file contains tests for db package by interacting with
// a real postgres database for testing purpose.
//
// To run these tests, you need to set up enviroment variables
// for connection to a testing and *clean* database.

// These are SQL statements for testing purpose to create Module table and drop Module table.
const (
	createModuleTable = `create table if not exists modules (
        orgName text NOT NULL, name text NOT NULL, version text NOT NULL,
        data jsonb NOT NULL, primary key (orgName, name, version)
		)`
	dropModuleTable = `drop table modules`
)

// CreateTestModuleTable is helper function to create module table in test database.
func CreateTestModuleTable() error {
	_, err := db.Exec(createModuleTable)
	if err != nil {
		return fmt.Errorf("testCreateTestModuleTable: failed to create testing Module table: %v", err)
	}
	return nil
}

// DropModuleTable is helper function to drop test table in test database.
func DropModuleTable() error {
	_, err := db.Exec(dropModuleTable)
	if err != nil {
		return fmt.Errorf("testDropModuleTable: failed to drop testing Module table: %v", err)
	}
	return nil
}

// TestInsertModule tests Insertion of Modules into database.
func TestInsertModule(t *testing.T) {
	tests := []struct {
		inOrgName string
		inName    string
		inVersion string
		inData    string
		wantErr   bool
		desc      string
	}{
		{
			inOrgName: "org_A",
			inName:    "module_A",
			inVersion: "version_A",
			inData:    "{}",
			wantErr:   false,
			desc:      "Test to insert one Module, expect to succeed",
		},
		{
			inOrgName: "org_B",
			inName:    "module_B",
			inVersion: "version_B",
			inData:    "",
			wantErr:   true,
			desc:      "Test to insert one Module with invalid json string, expect to fail",
		},
	}

	err := ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	if err := CreateTestModuleTable(); err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for _, tc := range tests {
		err = InsertModule(tc.inOrgName, tc.inName, tc.inVersion, tc.inData)
		if haserr := (err != nil); haserr != tc.wantErr {
			t.Errorf("insert module result mismatch, orgName: %s, name: %s, version: %s, data: %s, err: %v", tc.inOrgName, tc.inName, tc.inVersion, tc.inData, err)
		}
	}
	if err := DropModuleTable(); err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}

// TestQueryModulesByOrgName tests Query Module By orgName.
func TestQueryModulesByOrgName(t *testing.T) {
	inputs := struct {
		orgNames []string
		names    []string
		versions []string
		datas    []string
	}{
		orgNames: []string{"org1", "org1", "org2"},
		names:    []string{"name1_1", "name1_2", "name2"},
		versions: []string{"v1_1", "v1_2", "v2"},
		datas:    []string{"{}", "{}", "{}"},
	}
	tests := []struct {
		orgName string
		want    []Module
		desc    string
	}{
		{
			orgName: "org1",
			want: []Module{
				{
					OrgName: "org1",
					Name:    "name1_1",
					Version: "v1_1",
					Data:    "{}",
				},
				{
					OrgName: "org1",
					Name:    "name1_2",
					Version: "v1_2",
					Data:    "{}",
				},
			},
			desc: "Test to query with orgName with two matching Modules",
		},
		{
			orgName: "org2",
			want: []Module{
				{
					OrgName: "org2",
					Name:    "name2",
					Version: "v2",
					Data:    "{}",
				},
			},
			desc: "Test to query with orgName with only one matching Modules",
		},
		{
			orgName: "org3",
			want:    nil,
			desc:    "Test to query with orgName with no matching Module",
		},
		{
			orgName: "nil",
			want: []Module{
				{
					OrgName: "org1",
					Name:    "name1_1",
					Version: "v1_1",
					Data:    "{}",
				},
				{
					OrgName: "org1",
					Name:    "name1_2",
					Version: "v1_2",
					Data:    "{}",
				},
				{
					OrgName: "org2",
					Name:    "name2",
					Version: "v2",
					Data:    "{}",
				},
			},
			desc: "Test to query without orgName as query parameter, it equals to searching for all existing modules",
		},
	}

	err := ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	if err := CreateTestModuleTable(); err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for i := 0; i < len(inputs.names); i++ {
		err := InsertModule(inputs.orgNames[i], inputs.names[i], inputs.versions[i], inputs.datas[i])
		if err != nil {
			t.Errorf("pre insertion before query test failed: %v", err)
		}
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestQueryByOrgName %s, param orgName: %s", tc.desc, tc.orgName), func(t *testing.T) {
			var modules []Module
			if tc.orgName != "nil" {
				modules, err = QueryModulesByOrgName(&tc.orgName)
			} else {
				modules, err = QueryModulesByOrgName(nil)
			}
			if err != nil {
				t.Errorf("query by orgName failed, orgName: %s, err: %v", tc.orgName, err)
			}
			if !reflect.DeepEqual(modules, tc.want) {
				t.Errorf("query results mismatch, orgName: %s", tc.orgName)
			}
		})
	}

	if err := DropModuleTable(); err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}

// TestQueryModulesByKey tests query Module by its key (name, version).
func TestQueryModulesByKey(t *testing.T) {
	inputs := struct {
		orgNames []string
		names    []string
		versions []string
		datas    []string
	}{
		orgNames: []string{"org1_1", "org1_2", "org2", "org3", "org4"},
		names:    []string{"name1", "name1", "name2", "name2", "name3"},
		versions: []string{"v1", "v1", "v2", "v3", "v2"},
		datas:    []string{"{}", "{}", "{}", "{}", "{}"},
	}
	// all "nil" would be treated as null pointer in this test
	tests := []struct {
		name    string
		version string
		want    []Module
		desc    string
	}{
		{
			name:    "name1",
			version: "v1",
			want: []Module{
				{
					OrgName: "org1_1",
					Name:    "name1",
					Version: "v1",
					Data:    "{}",
				},
				{
					OrgName: "org1_2",
					Name:    "name1",
					Version: "v1",
					Data:    "{}",
				},
			},
			desc: "Test to query with both name and version",
		},
		{
			name:    "name2",
			version: "nil",
			want: []Module{
				{
					OrgName: "org2",
					Name:    "name2",
					Version: "v2",
					Data:    "{}",
				},
				{
					OrgName: "org3",
					Name:    "name2",
					Version: "v3",
					Data:    "{}",
				},
			},
			desc: "Test to query with only name as query parameter",
		},
		{
			name:    "nil",
			version: "v2",
			want: []Module{
				{
					OrgName: "org2",
					Name:    "name2",
					Version: "v2",
					Data:    "{}",
				},
				{
					OrgName: "org4",
					Name:    "name3",
					Version: "v2",
					Data:    "{}",
				},
			},
			desc: "Test to query with only version as query parameter",
		},
		{
			name:    "nil",
			version: "nil",
			want: []Module{
				{
					OrgName: "org1_1",
					Name:    "name1",
					Version: "v1",
					Data:    "{}",
				},
				{
					OrgName: "org1_2",
					Name:    "name1",
					Version: "v1",
					Data:    "{}",
				},
				{
					OrgName: "org2",
					Name:    "name2",
					Version: "v2",
					Data:    "{}",
				},
				{
					OrgName: "org3",
					Name:    "name2",
					Version: "v3",
					Data:    "{}",
				},
				{
					OrgName: "org4",
					Name:    "name3",
					Version: "v2",
					Data:    "{}",
				},
			},
			desc: "Test to query with no query parameters which equals query all models",
		},
		{
			name:    "name5",
			version: "v5",
			want:    nil,
			desc:    "Test to query with key of no matching modules",
		},
	}

	err := ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	if err := CreateTestModuleTable(); err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for i := 0; i < len(inputs.names); i++ {
		err = InsertModule(inputs.orgNames[i], inputs.names[i], inputs.versions[i], inputs.datas[i])
		if err != nil {
			t.Errorf("pre insertion before query test failed: %v", err)
		}
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestQueryByKey %s, param name: %s, version: %s", tc.desc, tc.name, tc.version), func(t *testing.T) {
			var modules []Module
			if tc.name != "nil" && tc.version != "nil" {
				modules, err = QueryModulesByKey(&tc.name, &tc.version)
			} else if tc.name != "nil" {
				modules, err = QueryModulesByKey(&tc.name, nil)
			} else if tc.version != "nil" {
				modules, err = QueryModulesByKey(nil, &tc.version)
			} else {
				modules, err = QueryModulesByKey(nil, nil)
			}
			if err != nil {
				t.Errorf("query by orgName failed, name: %s, version: %s, err: %v", tc.name, tc.version, err)
			}
			if !reflect.DeepEqual(modules, tc.want) {
				t.Errorf("query results mismatch, name: %s, version: %s", tc.name, tc.version)
			}
		})
	}

	if err := DropModuleTable(); err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}

func TestDeleteModule(t *testing.T) {
	inputs := []struct {
		orgName string
		name    string
		version string
		data    string
	}{
		{
			orgName: "test",
			name:    "1",
			version: "1",
			data:    "{}",
		},
		{
			orgName: "test",
			name:    "2",
			version: "2",
			data:    "{}",
		},
		{
			orgName: "newtest",
			name:    "3",
			version: "3",
			data:    "{}",
		},
	}

	tests := []struct {
		orgName string
		name    string
		version string
		wantErr bool
	}{
		{
			orgName: "test",
			name:    "1",
			version: "1",
			wantErr: false,
		},
		{
			orgName: "test",
			name:    "1",
			version: "1",
			wantErr: true,
		},
		{
			orgName: "test",
			name:    "1",
			version: "0",
			wantErr: true,
		},
		{
			orgName: "test",
			name:    "2",
			version: "2",
			wantErr: false,
		},
		{
			orgName: "newtest",
			name:    "3",
			version: "3",
			wantErr: false,
		},
	}

	err := ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	if err := CreateTestModuleTable(); err != nil {
		t.Errorf("create table failed: %v", err)
	}

	for i := 0; i < len(inputs); i++ {
		err := InsertModule(inputs[i].orgName, inputs[i].name, inputs[i].version, inputs[i].data)
		if err != nil {
			t.Errorf("pre insertion before query test failed: %v", err)
		}
	}

	for _, tc := range tests {
		if err := DeleteModule(tc.orgName, tc.name, tc.version); (err != nil) != tc.wantErr {
			t.Errorf("DeleteModule test failed: to delete, orgName: %s, name: %s, version: %s, wantErr: %t", tc.orgName, tc.name, tc.version, tc.wantErr)
		}
		// check whether after deletion, the module still exists or not.
		if !tc.wantErr {
			models, err := QueryModulesByKey(&tc.name, &tc.version)
			if err != nil {
				t.Errorf("DeleteModule, query after deleting encountered error: %v", err)
			}
			if len(models) != 0 {
				t.Errorf("DeleteModule, after deletion module still exists, orgName: %s, name: %s, version: %s", tc.orgName, tc.name, tc.version)
			}
		}
	}

	if err := DropModuleTable(); err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}
