/*
Package access contains function to validate token and parse access
  for write operations (i.e., create and update).
*/
package access

import (
	"context"
	"fmt"
	"strings"

	firebase "firebase.google.com/go/v4"
)

// const variables related to validate token.
// *projectID* is the id of project where the server is running in Google Cloud Platform.
// *delimiter* is delimiter for claim string of a list of names of organizations that the token owner has access to.
// *accessField* is field name of claim that contains the list of names of organizations that one has access to.
const (
	projectID   = `disco-idea-817`
	delimiter   = `,`
	accessField = `allow`
)

// ParseAcess takes input of a token string.
// It first valides whether the token is valid using firebase,
// then parse a list of names for organization that token owner has access to from the token.
// If token is invalid, an error is returned.
func ParseAccess(token string) ([]string, error) {
	// Set up firebase configuration to use correct token validation method.
	ctx := context.Background()
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ParseAccess: error initializing app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("ParseAccess: generate firebase authentication admin failed")
	}

	// Use firebase to valide token
	verifiedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("ParseAccess: error verifying ID token: %v\n", err)
	}

	// Retrieve *accessField* from claims, if the field does not exist, return an error.
	allowClaims, ok := verifiedToken.Claims[accessField]
	if !ok {
		return nil, fmt.Errorf("ParseAccess: verified token does not contain allow claims: %s", accessField)
	}

	// Split string into a slice of names of organizations.
	allowOrgs := strings.Split(allowClaims.(string), delimiter)
	return allowOrgs, nil
}
