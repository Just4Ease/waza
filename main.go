package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"waza/config"
	"waza/graph"
	"waza/setup"
)

func main() {
	secrets := config.GetSecrets()
	logger := logrus.New()

	opts := setup.ConfigureServiceDependencies(logger)

	// GraphQL API ( using this because of the playground, so that you won't stress loading up postman. )
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(opts)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", secrets.Port)
	log.Fatal(http.ListenAndServe(":"+secrets.Port, nil))
}
