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

const (
	createModuleTable = `create table if not exists modules (
        orgName text NOT NULL, name text not null, version text not null,
        data jsonb NOT NULL, primary key (orgName, name, version)
		)`
	dropModuleTable = `drop table modules`
)

// helper function to create module table in test database
func CreateModuleTable() error {
	_, err := db.Exec(createModuleTable)
	if err != nil {
		return fmt.Errorf("testCreateModuleTable: failed to create testing Module table: %v", err)
	}
	return nil
}

// helper function to drop test table in test database
func DropModuleTable() error {
	_, err := db.Exec(dropModuleTable)
	if err != nil {
		return fmt.Errorf("testDropModuleTable: failed to drop testing Module table: %v", err)
	}
	return nil
}

// Test Insertion of Modules into database
func TestInsertModule(t *testing.T) {
	tests := []struct {
		orgName string
		name    string
		version string
		data    string
		wanterr bool
	}{
		{
			orgName: "org_A",
			name:    "module_A",
			version: "version_A",
			data:    "{}",
			wanterr: false,
		},
		// Invalid json string, insertion should fail
		{
			orgName: "org_B",
			name:    "module_B",
			version: "version_B",
			data:    "",
			wanterr: true,
		},
	}
	var err error
	err = ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	err = CreateModuleTable()
	if err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for _, tc := range tests {
		err = InsertModule(tc.orgName, tc.name, tc.version, tc.data)
		haserr := (err != nil)
		if haserr != tc.wanterr {
			t.Errorf("insert module result mismatch, orgName: %s, name: %s, version: %s, data: %s, err: %v", tc.orgName, tc.name, tc.version, tc.data, err)
		}
	}
	err = DropModuleTable()
	if err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}

// Test Query Module By orgName
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
		},
		{
			orgName: "org3",
			want:    []Module{},
		},
		// special test case, set *string orgName to nil for this test
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
		},
	}

	var err error
	err = ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	err = CreateModuleTable()
	if err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for i := 0; i < len(inputs.names); i++ {
		err = InsertModule(inputs.orgNames[i], inputs.names[i], inputs.versions[i], inputs.datas[i])
		if err != nil {
			t.Errorf("pre insertion before query test failed: %v", err)
		}
	}
	for _, tc := range tests {
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
	}

	err = DropModuleTable()
	if err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}

// Test query Module by its key (name, version)
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
		},
		{
			name:    "name5",
			version: "v5",
			want:    []Module{},
		},
	}

	var err error
	err = ConnectDB()
	if err != nil {
		t.Errorf("connect to db failed: %v", err)
	}
	defer Close()
	err = CreateModuleTable()
	if err != nil {
		t.Errorf("create table failed: %v", err)
	}
	for i := 0; i < len(inputs.names); i++ {
		err = InsertModule(inputs.orgNames[i], inputs.names[i], inputs.versions[i], inputs.datas[i])
		if err != nil {
			t.Errorf("pre insertion before query test failed: %v", err)
		}
	}
	for _, tc := range tests {
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
	}

	err = DropModuleTable()
	if err != nil {
		t.Errorf("drop table failed, err: %v", err)
	}
}
