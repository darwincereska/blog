package main

import (
	"blog/internal/database"
	"blog/internal/config"
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
}
