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

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"

	// Go postgres driver for Go's database/sql package
	_ "github.com/lib/pq"
)

const (
	// getEmailWriteAccess returns a list of orgnames to which the given email has write access.
	getEmailWriteAccess = `SELECT orgname FROM email_write_access WHERE email = $1`
)

// readEmailOrgAccesses scans from a single-attribute relation of orgName strings.
func readEmailOrgAccesses(rows *sql.Rows) (map[string]bool, error) {
	orgNames := map[string]bool{}
	for rows.Next() {
		var orgName string
		if err := rows.Scan(&orgName); err != nil {
			return nil, fmt.Errorf("readEmailAccesses db scan: %v", err)
		}
		orgNames[orgName] = true
	}
	return orgNames, nil
}

// getEmailOrgAccesses gets the set of orgNames to which a particular email has
// write access.
func getEmailOrgAccesses(email string) (map[string]bool, error) {
	rows, err := db.Query(getEmailWriteAccess, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readEmailOrgAccesses(rows)
}

// CheckAccess takes input of a string of token and a string of organization's name.
// It checks whether the given token in valid and whether the associated user
// has write access for the given organization name.
// If not, an error is returned.
func CheckAccess(token string, orgName string) error {
	// Set up firebase configuration to use correct token validation method.
	ctx := context.Background()
	projectID, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		return fmt.Errorf("$PROJECT_ID not set")
	}
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return fmt.Errorf("access: error initializing app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("access: generate firebase authentication admin failed")
	}

	// Use firebase to validate token
	verifiedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return fmt.Errorf("error verifying ID token: %v\n", err)
	}

	userRecord, err := client.GetUser(ctx, verifiedToken.UID)
	if err != nil {
		return fmt.Errorf("access: GetUser failed: %v", err)
	}

	accesses, err := getEmailOrgAccesses(userRecord.Email)
	if err != nil {
		return fmt.Errorf("access: getEmailOrgAccesses failed: %v", err)
	}

	if !accesses[orgName] {
		// If the token does not contain access to input.OrgName, return an error.
		return fmt.Errorf("user does not have access to organization %s", orgName)
	}

	return nil
}
