package main

import (
	"context"
	"flag"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

const (
	projectID = `disco-idea-817`
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
