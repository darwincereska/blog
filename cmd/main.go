package main

import (
	"blog/internal/cache"
	"blog/internal/config"
	"blog/internal/database"
	"blog/internal/echo/logger"
	"blog/internal/echo/routes"
	"blog/internal/services"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
)

func main() {
	// Server config
	server_config := config.NewServerConfig()
	server_config.LoadConfig()

	// Database Config
	db_config := config.NewDatabaseConfig()
	db_config.LoadConfig()

	// Connect to database
	db, err := database.Connect(db_config.GetDSN())
	if err != nil {
		log.Error("Error", "error", err)
	}
	defer db.Close()

	// Create Redis caches
	strapi_cache := cache.CreateCache(server_config.RedisHost, server_config.RedisPort, 0)
	// analytics_cache := cache.CreateCache(server_config.RedisHost, server_config.RedisPort, 1)

	// Create Strapi service
	strapi_service := services.NewStrapiService(server_config.StrapiEndpoint+"/graphql", server_config.StrapiToken, strapi_cache)

	// Setup echo server
	e := echo.New()
	e.Logger = logger.NewCharmSlog()

	// Setup routes
	routes.SetupRoutes(e, routes.Sources{
		StrapiService: strapi_service,
	})

	// Start server
	host := fmt.Sprintf("%s:%d", server_config.Host, server_config.Port)
	if err := e.Start(host); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
