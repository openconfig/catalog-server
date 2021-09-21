### Overview

This directory contains scripts for administering the catalog system. It
currently provides the following functionalities:

+   Delete existing accounts. can alternatively be done in the Firebase UI. This
    may be deprecated in the future since it requires service account
    impersonation.

It does not provide functionalities for registering a new user as users can
register via the login page directly.

It also does not provide the functionality to grant write access to certain
organizations for an existing account. This must be done by directly inserting
and deleting entries within an administrative SQL table. Although Firebase
custom claims can also achieve this purpose, they require service accounts and
generating a private key for management, which may present more of a security
risk, and is more cumbersome than just running some SQL statements.

### Usage

+   To use these scripts the user must be admin of identity platform where the
    catalog system is deployed.
+   To `delete`, run `go run deleteaccount.go -email EMAIL-OF-ACCOUNT`.
