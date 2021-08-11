package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/openconfig/catalog-server/graph"
	"github.com/openconfig/catalog-server/graph/generated"
	"github.com/openconfig/catalog-server/pkg/db"
)

const defaultPort = "8080" // default port of launched catalog server

// Main function is automatically generated by *gqlgen*, except ConnectDB() to connect to database.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Establish connection with database, other codes are automatically generated.
	err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// Launch built-in graphQL frontend server.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// Set handler for all queries.
	http.Handle("/query", srv)

	// static file server to serve frontend webpages.
	fileServer := http.FileServer(http.Dir("frontend"))
	http.HandleFunc(
		"/static/",
		func(w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		},
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
