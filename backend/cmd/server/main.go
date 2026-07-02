// @title           The_Scan API
// @version         1.0
// @description     Document scanning REST API with JWT authentication
// @host            localhost:8008
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth
package main

import (
	"log"

	"github.com/neathhatake/the_Scan/internal/api"
	"github.com/neathhatake/the_Scan/internal/config"
	"github.com/neathhatake/the_Scan/internal/database"
	"github.com/neathhatake/the_Scan/pkg/logger"
)

// @Summary Health check
// @Description Returns API health status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	logger.Init(config.App.AppEnv)

	// Connect to database
	db := database.ConnectDatabase(config.App)
	database.RunMigrations(db, "migrations", config.App.GetDSN())


	// Run API
	api.Run(config.App, db)
}
