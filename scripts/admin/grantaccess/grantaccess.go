package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

const (
	projectID   = `disco-idea-817`
	accessField = `allow`
)

func main() {
	var keyPathPtr = flag.String("key", "", "file path of admin key")
	var emailPtr = flag.String("email", "", "email account that you want to change access for")
	var accessPtr = flag.String("access", "", "string of a list of organizations that account would be granted access to, seperated by delimiter")
	flag.Parse()

	if *keyPathPtr == "" || *emailPtr == "" || *accessPtr == "" {
		log.Fatalf("Please provide all required inputs\n")
	}

	opt := option.WithCredentialsFile(*keyPathPtr)

	// Set up firebase configuration.
	ctx := context.Background()
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
		log.Fatalln("error setting custom claims %v\n", err)
	}
	fmt.Println("Successfully set up new claims", currentClaims)

}
