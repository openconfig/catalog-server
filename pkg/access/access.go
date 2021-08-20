/*
Package access contains function to validate token and parse access from a valid token
  for write operations (i.e., create and update).
*/
package access

import (
	"context"
	"fmt"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
)

// const variables related to token validation.
// *projectID* is the id of project where the server is running in Google Cloud Platform.
// *delimiter* is delimiter for claim string of a list of names of organizations that the token owner has access to.
// *accessField* is field name of claim that contains the list of names of organizations that one has access to.
// Note that organization's names should not contain delimiter.
const (
	projectID       = `disco-idea-817`
	delimiter       = `,`
	baseAccessField = `allow`
)

func GetAccessField() (string, error) {
	// name of target database
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return "", fmt.Errorf("DB_NAME not set")
	}
	return (dbname + "-" + baseAccessField), nil
}

// ParseAcess takes input of a token string.
// It first validates whether the token is valid using firebase,
// then parses from the token's claims a list organization names to which that the token owner has write access.
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

	// Use firebase to validate token
	verifiedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("ParseAccess: error verifying ID token: %v\n", err)
	}

	accessField, err := GetAccessField()
	if err != nil {
		return nil, fmt.Errorf("ParseAccess: get access field failed: %v", err)
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

// This function takes input of a string of token and a string of organization's name.
// It checks whether the given token in valid and whether it contains access for write operation to *orgName*.
// If not, an error is returned.
func CheckAccess(token string, orgName string) error {
	// Validate token
	allowOrgs, err := ParseAccess(token)
	if err != nil {
		return fmt.Errorf("CheckAccess: user does not provide valid token: %v", err)
	}

	// Check whether this account has access to such orgnization.
	hasAccess := false
	for _, allowOrg := range allowOrgs {
		if allowOrg == orgName {
			hasAccess = true
			break
		}
	}

	// If the token does not contain access to input.OrgName, return an error.
	if !hasAccess {
		return fmt.Errorf("CheckAccess: user does not have access to organization %s", orgName)
	}

	return nil
}
