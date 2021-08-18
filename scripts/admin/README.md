### Overview
This directory contains scrips for helping admin of catalog system. It provides functionality to:
+ Delete existing account.
+ Grant write access of organizations to an existing account.

It does not provide functionalities for admin to register a new user as it can register 
via the login page by itself.

### Usage
+ Admin needs to first run `bash create-ca.sh` to download a secret `KEY` to operate functionalities above.
+ To `delete`, run `go run deleteaccount.go -email EMAIL-OF-ACCOUNT -key PATH-TO-KEY`.
+ TO `grant access`, run `go run grantaccess.go -email EMAIL-OF-ACCOUNT -key PATH-TO-KEY -claims STRING-OF_CLAIMS`. `STRING-OF_CLAIMS` is a string of list of organizations seperated by comma. That is, we expect the name of organization do not contain comma. Or we can choose a different delimiter with no conflicts with names of organizations.