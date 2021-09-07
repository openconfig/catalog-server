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
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main() {
	var keyPathPtr = flag.String("key", "", "file path of admin key")
	var emailPtr = flag.String("email", "", "email account that you want to change access for")
	flag.Parse()

	if *keyPathPtr == "" || *emailPtr == "" {
		log.Fatalf("Please provide all required inputs\n")
	}

	opt := option.WithCredentialsFile(*keyPathPtr)

	// Set up firebase configuration.
	ctx := context.Background()
	projectID, ok := os.LookupEnv("CLOUDSDK_CORE_PROJECT")
	if !ok {
		log.Fatalf("$CLOUDSDK_CORE_PROJECT not set.")
	}
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config, opt)

	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Generate firebase authentication admin failed\n")
	}

	// Fetch user via email.
	user, err := client.GetUserByEmail(ctx, *emailPtr)
	// Check whether this user exists or not.
	if user == nil {
		log.Fatalf("Account %s does not exist\n", *emailPtr)
	}

	// Delete the account.
	if err := client.DeleteUser(ctx, user.UID); err != nil {
		log.Fatalf("error deleting user: %v\n", err)
	}
	log.Printf("Successfully deleted user: %s\n", *emailPtr)

}
