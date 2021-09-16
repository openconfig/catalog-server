### Overview
This directory contains scripts for helping admin of catalog system. It provides functionality to:
+ Delete existing account.
+ Grant write access of organizations to an existing account.

It does not provide functionalities for admin to register a new user as it can register 
via the login page by itself.

### Usage
+ To use these scripts the user must be admin of identity platform where the catalog system is deployed.
+ To `delete`, run `go run deleteaccount.go -email EMAIL-OF-ACCOUNT`.
+ To `grant access`, run `go run grantaccess.go -db NAME-OF-DB -email EMAIL-OF-ACCOUNT -access STRING-OF_CLAIMS`. `STRING-OF_CLAIMS` is a string of list of organizations seperated by comma. That is, we expect the name of organization do not contain comma. Or we can choose a different delimiter with no conflicts with names of organizations. If user does not provide `-access STRING-OF_CLAIMS` is equivalent to setting access of this account to `empty access` (i.e., no write access to any organization).
+ To `read all existing users' access`, run `go run grantacces.go -db NAME-OF-DB -all`.
