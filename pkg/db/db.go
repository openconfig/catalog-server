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
	_ "github.com/lib/pq"
)

// Put all used SQL statments here for easy maintainess and tracking.
// To support more operations (query/update), just append this const list.
const (
	insertModule            = `INSERT INTO modules (orgName, name, version, data) VALUES($1, $2, $3, $4)`
	allModules              = `select * from modules`
	modulesByOrgName        = `select * from modules where orgName=$1`
	modulesByName           = `select * from modules where name=$1`
	modulesByVersion        = `select * from modules where version=$1`
	modulesByNameAndVersion = `select * from modules where name=$1 and version=$2`
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
	var ok bool
	var err error
	var port int         // port number of target database
	var user string      // username of target database
	var password string  // password of target database
	var dbname string    // name of target database
	var host string      // host address of target database
	var socketDir string // (Cloud Run only) Directory of Unix socket
	var psqlconn string  // connection string used to connect to traget database

	// read db config from env

	if portStr, ok := os.LookupEnv("DB_PORT"); !ok {
		fmt.Println("DB_PORT not set, set port to 5432")
		port = 5432
	} else {
		if port, err = strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("parse port failed: %v", err)
		}
	}

	if user, ok = os.LookupEnv("DB_USERNAME"); !ok {
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

	if host, ok = os.LookupEnv("DB_HOST"); !ok || host == "localhost" {
		if !ok {
			fmt.Println("DB_HOST not set, set host to localhost")
			host = "localhost"
		}
		// This connection string is used if service is not deployed on Cloud Run,
		// instead connection is made from localhost via Cloud SQL proxy.
		psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	} else {
		psqlconn = fmt.Sprintf("host=%s/%s port=%d user=%s password=%s dbname=%s", socketDir, host, port, user, password, dbname)
	}

	// open database
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("open database failed: %v", err)
	}

	// see if connection is established successfully
	if err = db.Ping(); err != nil {
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

// Query modules of organization with *orgName* from database.
// If orgName is null then directly query all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByOrgName(orgName *string) ([]Module, error) {
	var rows *sql.Rows
	var err error

	if orgName != nil {
		rows, err = db.Query(modulesByOrgName, *orgName)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByOrgName failed: %v, query param orgName: %s", err, *orgName)
		}
	} else {
		rows, err = db.Query(allModules)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByOrgName failed: %v", err)
		}
	}

	defer rows.Close()

	return ReadModluesByRow(rows)
}

// Query modules by its key (name, version), it is possible that parameters are null.
// If both parameters are null, this equals query for all modules.
// Return slice of db Module struct each field of which corresponds to one column in db.
// Error is returned when query or reading data failed.
func QueryModulesByKey(name *string, version *string) ([]Module, error) {
	var rows *sql.Rows
	var err error

	if name != nil && version != nil {
		rows, err = db.Query(modulesByNameAndVersion, *name, *version)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByKey failed: %v, query param name: %s, version: %s", err, *name, *version)
		}
	} else if name != nil {
		rows, err = db.Query(modulesByName, *name)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByKey failed: %v, query param name: %s", err, *name)
		}
	} else if version != nil {
		rows, err = db.Query(modulesByVersion, *version)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByKey failed: %v, query param version: %s", err, *version)
		}
	} else {
		rows, err = db.Query(allModules)
		if err != nil {
			return nil, fmt.Errorf("QueryModulesByKey failed: %v", err)
		}
	}

	defer rows.Close()

	return ReadModluesByRow(rows)
}

// Close db connection
func Close() error {
	return db.Close()
}
