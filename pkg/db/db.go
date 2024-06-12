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

/*
Package db contains functions related to database.
 * db.go includes conneting to db, query and insertion.
 * dbschema.go contains definitions of struct for db tables.
*/
package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/glog"

	// Go postgres driver for Go's database/sql package
	_ "github.com/lib/pq"
)

// These are SQL stataments used in this package.
// Query statements can be appended based on its query parameters.
const (
	// $4 and $5 should be assigned with the same value (the JSON data of module).
	insertModule  = `INSERT INTO modules (orgName, name, version, data) VALUES($1, $2, $3, $4) on conflict (orgName, name, version) do update set data=$5`
	selectModules = `select * from modules`
	// We want to ensure that user has to provide all three inputs,
	// instead of deleting too many modules by mistake with some fields missing.
	deleteModule = `delete from modules where orgName = $1 and name = $2 and version = $3`

	selectFeatureBundles = `select * from featureBundles`
	// $4 and $5 should be assigned with the same value (the JSON data of feature-bundle).
	insertFeatureBundle = `INSERT INTO featureBundles (orgName, name, version, data) VALUES($1, $2, $3, $4) on conflict (orgName, name, version) do update set data=$5`
	deleteFeatureBundle = `delete from featurebundles where orgName = $1 and name = $2 and version = $3`
)

// db is the global variable of connection to database.
// It would be assigned value when *ConnectDB* function is called.
//
// We choose to use this global variable due to that *gqlgen* automatically generates many server side codes.
// Resolver functions generated by *gqlgen* are handler functions for graphQL queries and are methods of *resolver* struct.
// If we define a struct with field of db connection instead of using the global variable
// and change *Query* functions to methods of that struct,
// we need to initialize db connection while initializing *resolver* object in server codes
// such that *resolver* functions can call these *Query* functions.
// However, initialization function of *resolver* struct is automatically generated and overwritten every time.
// Thus, we cannot initialize a db connection inside the *resolver* objects.
//
// Another option is to establish a new db connection for each *Query* function and close it after query finishes.
// However, that would be too expensive to connect to db every time server receives a new query.
var db *sql.DB

// ConnectDB establishes connection to database, *db* variable is assigned when opening database.
// This should only be called once before any other database function is called.
//
// Users need set environment variables for connection, including
//  * DB_HOST:          host address of target db instances, by default: localhost.
//  * DB_PORT:          port number of postgres db, by default: 5432.
//  * DB_USERNAME:      username of database, error would be returned if not set.
//  * DB_PWD:           password of target database, error would be returned if not set.
//  * DB_NAME:          name of database for connection, error would be returned if not set.
//  * DB_SOCKER_DIR:    directory of Unix socket in Cloud Run which serves as Cloud SQL
//                      Auth proxy to connect to postgres database.
//                      If service is deployed on Cloud Run, just use the default value.
//                      By default, it is set to `/cloudsql`.
func ConnectDB() error {
	// read db config from env

	// port number of target database
	port := 5432
	if portStr, ok := os.LookupEnv("DB_PORT"); !ok {
		glog.Infof("DB_PORT not set, setting port to %d", port)
	} else {
		var err error
		if port, err = strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("DB_PORT in incorrect format: %v", err)
		}
	}

	// username of target database
	user, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return fmt.Errorf("DB_USERNAME not set")
	}

	// password of target database
	password, ok := os.LookupEnv("DB_PWD")
	if !ok {
		return fmt.Errorf("DB_PWD not set")
	}

	// name of target database
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return fmt.Errorf("DB_NAME not set")
	}

	// (Cloud Run only) Directory of Unix socket
	socketDir, ok := os.LookupEnv("DB_SOCKET_DIR")
	if !ok {
		socketDir = "/cloudsql"
	}

	var psqlconn string // connection string used to connect to traget database

	// host address of target database
	host, ok := os.LookupEnv("DB_HOST")
	switch {
	case !ok:
		glog.Infoln("DB_HOST not set, setting host to localhost")
		host = "localhost"
		fallthrough
	case host == "localhost":
		// This connection string is used if service is not deployed on Cloud Run,
		// instead connection is made from localhost via Cloud SQL proxy.
		psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	default:
		psqlconn = fmt.Sprintf("host=%s/%s port=%d user=%s password=%s dbname=%s", socketDir, host, port, user, password, dbname)
	}

	// open database
	var err error
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("open database failed: %v", err)
	}

	// see if connection is established successfully
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database failed: %v", err)
	}

	return nil
}

// Close function closes db connection
func Close() error {
	return db.Close()
}

// InsertModule inserts module into database given values of four field of MODULE schema.
// Or if there is existing module with existing key (orgName, name, version), update data field.
// Error is returned when insertion failed.
func InsertModule(orgName string, name string, version string, data string) error {
	if _, err := db.Exec(insertModule, orgName, name, version, data, data); err != nil {
		return fmt.Errorf("insert/update module into db failed: %v", err)
	}
	return nil
}

// readModulesByRow scans from queried modules from rows one by one, rows are closed inside.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when scan rows failed.
func readModulesByRow(rows *sql.Rows) ([]Module, error) {
	var modules []Module
	defer rows.Close()
	for rows.Next() {
		var module Module
		if err := rows.Scan(&module.OrgName, &module.Name, &module.Version, &module.Data); err != nil {
			return nil, fmt.Errorf("scan db rows failure, %v", err)
		}
		modules = append(modules, module)
	}
	return modules, nil
}

// formatQueryStr is used to generate query statement string based on query parameters.
// * parmNames is a list of names of all non-nil query parameters.
// * baseQuery is query statement without any query parameters.
func formatQueryStr(parmNames []string, baseQuery string) string {
	queryStmt := baseQuery
	for i := 0; i < len(parmNames); i++ {
		if i == 0 {
			queryStmt += " where"
		} else {
			queryStmt += " and"
		}
		queryStmt += fmt.Sprintf(" %s=$%d", parmNames[i], i+1)
	}
	return queryStmt
}

// QueryModulesByOrgName queries modules of organization with *orgName* from database.
// If orgName is null then directly query all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByOrgName(orgName *string) ([]Module, error) {
	var parms []interface{} // parms is used to store value of non-nil query parameters
	parmNames := []string{} // parmNames is used to store name of non-nil query parameters

	if orgName != nil {
		parms = append(parms, *orgName)
		parmNames = append(parmNames, "orgName")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := formatQueryStr(parmNames, selectModules)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryModulesByOrgName failed: %v", err)
	}

	defer rows.Close()

	return readModulesByRow(rows)
}

// QueryModulesByKey queries modules by its key (name, version), it is possible that parameters are null.
// If both parameters are null, this equals query for all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByKey(name *string, version *string) ([]Module, error) {
	var parms []interface{} // parms is used to store value of non-nil query parameters
	parmNames := []string{} // parmNames is used to store name of non-nil query parameters

	if name != nil {
		parms = append(parms, *name)
		parmNames = append(parmNames, "name")
	}

	if version != nil {
		parms = append(parms, *version)
		parmNames = append(parmNames, "version")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := formatQueryStr(parmNames, selectModules)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryModulesByOrgName failed: %v", err)
	}

	defer rows.Close()

	return readModulesByRow(rows)
}

// DeleteModule takes three string, orgName, name, version,
// whose combination is key of one Module in DB's Module table.
// If deletion fails, an non-nil error is returned.
// If the number of rows affected by this deletion is not 1, an error is also returned.
func DeleteModule(orgName string, name string, version string) error {
	result, err := db.Exec(deleteModule, orgName, name, version)
	if err != nil {
		return fmt.Errorf("DeleteModule failed: %v", err)
	}
	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteModule, access rows affected in result failed: %v", err)
	}
	// delete should only affect one row
	if num != 1 {
		return fmt.Errorf("DeleteModule: affected row is not one, it affects %d rows", num)
	}

	return nil
}

// readFeatureBundlesByRow scans from queried FeatureBundles from rows one by one, rows are closed inside.
// Return slice of db FeatureBundle struct each field of which corresponds to one column in db.
// Error is returned when scan rows failed.
func readFeatureBundlesByRow(rows *sql.Rows) ([]FeatureBundle, error) {
	var featureBundles []FeatureBundle
	defer rows.Close()
	for rows.Next() {
		var featureBundle FeatureBundle
		if err := rows.Scan(&featureBundle.OrgName, &featureBundle.Name, &featureBundle.Version, &featureBundle.Data); err != nil {
			return nil, fmt.Errorf("readFeatureBundlesByRow: scan db rows failure, %v", err)
		}
		featureBundles = append(featureBundles, featureBundle)
	}
	return featureBundles, nil
}

// QueryFeatureBundlesByOrgName queries feature-bundles of organization with *orgName* from database.
// If orgName is null then directly query all feature-bundles.
// Return slice of db FeatureBundle struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryFeatureBundlesByOrgName(orgName *string) ([]FeatureBundle, error) {
	var parms []interface{} // parms is used to store value of non-nil query parameters
	parmNames := []string{} // parmNames is used to store name of non-nil query parameters

	if orgName != nil {
		parms = append(parms, *orgName)
		parmNames = append(parmNames, "orgName")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := formatQueryStr(parmNames, selectFeatureBundles)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryFeatureBundlesByOrgName failed: %v", err)
	}

	return readFeatureBundlesByRow(rows)
}

// InsertFeatureBundle inserts FeatureBundle into database given values of four field of FeatureBundle schema.
// Or if there is existing FeatureBundle with existing key (orgName, name, version), update data field.
// Error is returned when insertion failed.
func InsertFeatureBundle(orgName string, name string, version string, data string) error {
	if _, err := db.Exec(insertFeatureBundle, orgName, name, version, data, data); err != nil {
		return fmt.Errorf("insert/update FeatureBundle into db failed: %v", err)
	}
	return nil
}

// DeleteFeatureBundle takes three pointer of string, orgName, name, version,
// whose combination is key of one FeatureBundle in DB's FeatureBundle table.
// If deletion fails, an non-nil error is returned.
// If the number of rows affected by this deletion is not 1, an error is also returned.
func DeleteFeatureBundle(orgName string, name string, version string) error {
	result, err := db.Exec(deleteFeatureBundle, orgName, name, version)
	if err != nil {
		return fmt.Errorf("DeleteFeatureBundle failed: %v", err)
	}
	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteFeatureBundle, access rows affected in result failed: %v", err)
	}
	// delete should only affect one row
	if num != 1 {
		return fmt.Errorf("DeleteFeatureBundle: affected row is not one, it affects %d rows", num)
	}

	return nil
}

// QueryFeatureBundlesByKey queries feature-bundles by its key (name, version), it is possible that
// If both parameters are null, this equals query for all feature-bundles.
// Return slice of db FeatureBundle struct each field of which corresponds to one column in
// Error is returned when query or reading data failed.
func QueryFeatureBundlesByKey(name *string, version *string) ([]FeatureBundle, error) {
	var parms []interface{} // parms is used to store value of non-nil query paramete
	parmNames := []string{} // parmNames is used to store name of non-nil query param

	if name != nil {
		parms = append(parms, *name)
		parmNames = append(parmNames, "name")
	}

	if version != nil {
		parms = append(parms, *version)
		parmNames = append(parmNames, "version")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := formatQueryStr(parmNames, selectFeatureBundles)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryFeatureBundlesByKey failed: %v", err)
	}

	return readFeatureBundlesByRow(rows)
}
