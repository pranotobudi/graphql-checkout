package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pranotobudi/graphql-checkout/config"
	"github.com/pranotobudi/graphql-checkout/database"
	"github.com/pranotobudi/graphql-checkout/graph"
	"github.com/pranotobudi/graphql-checkout/graph/generated"
	"github.com/pranotobudi/graphql-checkout/store"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		// executed in development only,
		//for production set those on production environment settings

		// load local env variables to os
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("failed to load .env file")
		}
	}
	log.Println("bismillah")
	postgres := database.GetDB()
	postgres.MigrateDB("./database/init.sql")

	router := chi.NewRouter()
	// router.Use(auth.Middleware())
	storeObject := store.NewStore()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		StoreService: *store.NewStoreService(*storeObject, postgres),
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)
	config := config.AppConfig()
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
