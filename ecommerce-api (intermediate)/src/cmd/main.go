// E-Commerce API Server
// A simple e-commerce backend API built with Go
// Features: User authentication, product management, and categories
package main

import (
	"log"

	"github.com/Hakim-CS/go-ecommerce-api/api"
	"github.com/Hakim-CS/go-ecommerce-api/config"
	"github.com/Hakim-CS/go-ecommerce-api/constants"
	"github.com/Hakim-CS/go-ecommerce-api/infra/db"
	"github.com/Hakim-CS/go-ecommerce-api/infra/db/migrations"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()

	// Initialize the database connection
	err := db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Postgres, constants.Startup, err.Error())
	}

	// Run database migrations
	migrations.Up_1()

	// Start the HTTP server
	api.InitServer(cfg)
}
