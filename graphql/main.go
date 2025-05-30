package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountUrl string `envconfig: "ACCOUNT_SERVER_URL"`
	CatalogUrl string `envconfig: "CATALOG_SERVER_URL"`
	OrderUrl   string `envconfig: "ORDER_SERVER_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	// http.Handle("/graphql", handler.NewDefaultServer(s.ToExecutableSchema()))
	// http.Handle("/playground", Playground.Handler("sye", "/graphql"))

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", handler.Playground("sye", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
