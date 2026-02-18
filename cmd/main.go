package main

import (
	"blog/internal/cache"
	"blog/internal/config"
	"blog/internal/database"
	"blog/internal/services"
	"context"
	"os"

	"github.com/charmbracelet/log"
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
	// analytics_cache := cache.CreateCache(server_config.RedisHost, server_config.RedisPort, 1)

	// Create Strapi service
	strapi_service := services.NewStrapiService(server_config.StrapiEndpoint+"/graphql", server_config.StrapiToken, strapi_cache)

	// Strapi logger
	strapi_logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Prefix: "STRAPI",
	})

	// Test strapi get
	posts, err := strapi_service.GetFeaturedPosts(context.Background(), 10, 1) 
	if err != nil {
		strapi_logger.Error(err)
		os.Exit(0)
	}
	
	post := posts[0]

	strapi_logger.Info(post)
}
