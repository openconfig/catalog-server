## How to support a new operation

#### Bascis about GraphQL 
+ Please go through the some bascis in GraphQL first. Tutorials in references may be helpful to get you familiar with GraphQL.
+ Important terms in GraphQL:
	+ Query: user-defined read operations.
	+ Mutation: user-defined write operations.

#### To add a new operation
+ Edit [schema.graphqls](../graph/schema.graphqls) to add a new operation you want. Depending on whether is a query or mutation, you need to add the new operation into `type Query {...}` or `type Mutation {...}`. 
	+ Formula for this new operation is `Operation-name(input parameters) [response type]`. 
	+ `Operation-name` is the name of the newly added operation.
	+ `input parameters` are list of input parameters for the new operation. The exclamation mark `!` after the input paramter implies that it is required and not optional. 
		+ For write operations (muation), we require input parameters contain a `token` parameter. A valid token with correct access is requred to perform such operation.
	+ `response type` defines response after new operation is executed.
+ Run `go run github.com/99designs/gqlgen generate` to generate codes which are required to support the newly added operation.
+ [schema.resolvers.go](../graph/schema.resolvers.go) file now contains a new resolver function with the same name as the newly added operation. You only need to implement all logic of the newly added operation inside the resolver function.

#### References
+ [How to GraphQL](https://www.howtographql.com/basics/0-introduction/)
+ [Introduction to GraphQL](https://graphql.org/learn/)
+ [Getting started on gqlgen](https://gqlgen.com/v0.13.0/getting-started/)