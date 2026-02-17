package main

import (
	"blog/internal/database"
	"blog/internal/config"
	"github.com/charmbracelet/log"
	"blog/internal/cache"
	"blog/internal/services"
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
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Create Redis caches
	strapi_cache := cache.CreateCache(server_config.RedisHost, server_config.RedisPort, 0)
	analytics_cache := cache.CreateCache(server_config.RedisHost, server_config.RedisPort, 1)

	// Create Strapi service
	strapi_service := services.NewStrapiService(server_config.StrapiHost, server_config.StrapiApiKey, strapi_cache)
}
