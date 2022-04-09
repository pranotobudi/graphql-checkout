package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pranotobudi/graphql-checkout/config"
	"github.com/pranotobudi/graphql-checkout/database"
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
	postgres := database.InitDB()
	postgres.MigrateDB("./database/init.sql")

	router := chi.NewRouter()
	// router.Use(auth.Middleware())

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", nil)

	config := config.AppConfig()
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
