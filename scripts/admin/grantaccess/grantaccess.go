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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	baseAccessField = `allow`
)

func main() {
	var emailPtr = flag.String("email", "", "email account that you want to change access for")
	var accessPtr = flag.String("access", "", "string of a list of organizations that account would be granted access to, seperated by delimiter. If not set, it means set empty access for this account")
	var listall = flag.Bool("all", false, "whether to list all current users' claims")
	var dbnamePtr = flag.String("db", "", "name of db that you want to grant user access to")

	flag.Parse()

	if *dbnamePtr == "" {
		flag.CommandLine.SetOutput(os.Stderr)
		fmt.Fprintf(os.Stderr, "Please provide db name\n")
		flag.Usage()
		os.Exit(1)
	}

	if !*listall && *emailPtr == "" {
		flag.CommandLine.SetOutput(os.Stderr)
		fmt.Fprintf(os.Stderr, "Please provide either provide email address or specify list `all`\n")
		flag.Usage()
		os.Exit(1)
	}

	projectID, ok := os.LookupEnv("CLOUDSDK_CORE_PROJECT")
	if !ok {
		log.Fatalf("$CLOUDSDK_CORE_PROJECT not set.")
	}

	ctx := context.Background()
	// Impersonate service account.
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: fmt.Sprintf("sa-claims@%s.iam.gserviceaccount.com", projectID),
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		log.Fatalf("impersonation failed: %v", err)
	}
	opt := option.WithTokenSource(ts)

	// Set up firebase configuration.
	config := &firebase.Config{
		ProjectID: projectID,
	}
	app, err := firebase.NewApp(ctx, config, opt)

	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Generate firebase authentication admin failed: %v\n", err)
	}

	accessField := *dbnamePtr + "-" + baseAccessField

	// list all existing users and their claims
	if *listall {
		iter := client.Users(ctx, "")
		for {
			user, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("error listing users: %s\n", err)
			}
			// Assume all users have email addresses.
			if user.CustomClaims == nil {
				log.Printf("%s does not have claims\n", user.Email)
			} else if access, ok := user.CustomClaims[accessField]; !ok {
				log.Printf("%s does not have access\n", user.Email)
			} else {
				log.Printf("%s current access is: %v\n", user.Email, access)
			}
		}
	} else {
		// Fetch user via email.
		user, err := client.GetUserByEmail(ctx, *emailPtr)
		if err != nil {
			log.Fatalf("Retrieve account by email failed: %v\n", err)
		}
		// Check whether this user exists or not.
		if user == nil {
			log.Fatalf("Account %s does not exist\n", *emailPtr)
		}

		// Obtain current claims from the user
		currentClaims := user.CustomClaims
		fmt.Println("Current claims", currentClaims)
		// If the account does not contain any claims, then its claim is nil,
		// we need to create a map of claims for it.
		if currentClaims == nil {
			currentClaims = make(map[string]interface{})
		}

		// Set claims to grant access.
		currentClaims[accessField] = *accessPtr
		if err := client.SetCustomUserClaims(ctx, user.UID, currentClaims); err != nil {
			log.Fatalf("error setting custom claims %v\n", err)
		}
		fmt.Println("Successfully set up new claims", currentClaims)
	}
}
