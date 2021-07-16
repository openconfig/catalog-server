/*
Package db contains functions related to database.
 * db.go includes conneting to db, query and insertion.
 * dbschema.go contains definitions of struct for db tables.
   Currently it only contains Module struct.
*/
package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	// Go postgres driver for Go's database/sql package
	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

// SQL stataments used in this package.
// Query statements can be appended based on its query parameters.
const (
	insertModule  = `INSERT INTO modules (orgName, name, version, data) VALUES($1, $2, $3, $4)`
	selectModules = `select * from modules`
)

var db *sql.DB

// Establish connection to database, *db* variable is assigned when opening database.
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
	var port int         // port number of target database
	var user string      // username of target database
	var password string  // password of target database
	var dbname string    // name of target database
	var host string      // host address of target database
	var socketDir string // (Cloud Run only) Directory of Unix socket
	var psqlconn string  // connection string used to connect to traget database

	// read db config from env

	port = 5432
	if portStr, ok := os.LookupEnv("DB_PORT"); !ok {
		glog.Infof("DB_PORT not set, setting port to %d", port)
	} else {
		var err error
		if port, err = strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("DB_PORT in incorrect format: %v", err)
		}
	}

	user, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return fmt.Errorf("DB_USERNAME not set")
	}

	if password, ok = os.LookupEnv("DB_PWD"); !ok {
		return fmt.Errorf("DB_PWD not set")
	}

	if dbname, ok = os.LookupEnv("DB_NAME"); !ok {
		return fmt.Errorf("DB_NAME not set")
	}

	if socketDir, ok = os.LookupEnv("DB_SOCKET_DIR"); !ok {
		socketDir = "/cloudsql"
	}

	host, ok = os.LookupEnv("DB_HOST")
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

// Insert module into database given values of four field of MODULE schema.
// Error is returned when insertion failed.
func InsertModule(orgName string, name string, version string, data string) error {
	if _, err := db.Exec(insertModule, orgName, name, version, data); err != nil {
		return fmt.Errorf("insert module into db failed: %v", err)
	}
	return nil
}

// Scan from queried modules from rows one by one, rows are *not* closed inside.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when scan rows failed.
func ReadModluesByRow(rows *sql.Rows) ([]Module, error) {
	var modules []Module
	for rows.Next() {
		var module Module
		if err := rows.Scan(&module.OrgName, &module.Name, &module.Version, &module.Data); err != nil {
			return nil, fmt.Errorf("scan db rows failure, %v", err)
		}
		modules = append(modules, module)
	}
	return modules, nil
}

// This function is used to generate query statement string based on query parameters.
// * parmNames is a list of names of all non-nil query parameters.
// * baseQuery is query statement without any query parameters.
func FormatQueryStr(parmNames []string, baseQuery string) string {
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

// Query modules of organization with *orgName* from database.
// If orgName is null then directly query all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByOrgName(orgName *string) ([]Module, error) {
	parms := []interface{}{} // parms is used to store value of non-nil query parameters
	parmNames := []string{}  // parmNames is used to store name of non-nil query parameters

	if orgName != nil {
		parms = append(parms, *orgName)
		parmNames = append(parmNames, "orgName")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := FormatQueryStr(parmNames, selectModules)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryModulesByOrgName failed: %v", err)
	}

	defer rows.Close()

	return ReadModluesByRow(rows)
}

// Query modules by its key (name, version), it is possible that parameters are null.
// If both parameters are null, this equals query for all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByKey(name *string, version *string) ([]Module, error) {
	parms := []interface{}{} // parms is used to store value of non-nil query parameters
	parmNames := []string{}  // parmNames is used to store name of non-nil query parameters

	if name != nil {
		parms = append(parms, *name)
		parmNames = append(parmNames, "name")
	}

	if version != nil {
		parms = append(parms, *version)
		parmNames = append(parmNames, "version")
	}

	// Format query statement string based on non-nil query parameters
	queryStmt := FormatQueryStr(parmNames, selectModules)

	rows, err := db.Query(queryStmt, parms...)
	if err != nil {
		return nil, fmt.Errorf("QueryModulesByOrgName failed: %v", err)
	}

	defer rows.Close()

	return ReadModluesByRow(rows)
}

// Close db connection
func Close() error {
	return db.Close()
}
