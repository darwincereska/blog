package config

import (
	"blog/internal/config/env"
	"fmt"
	"github.com/charmbracelet/log"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	TimeZone string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{}
}

func (c *DatabaseConfig) LoadConfig() {
	c.Host = env.GetString("DB_HOST", "localhost")
	c.Port = env.GetInt("DB_PORT", 5432)
	c.Name = env.GetString("DB_NAME", "blog")
	c.User = env.GetString("DB_USER", "blog")
	c.Password = env.GetString("DB_PASSWORD", "blog")
	c.SSLMode = env.GetString("DB_SSL_MODE", "disable")
	c.TimeZone = env.GetString("DB_TIME_ZONE", "America/New_York")

	log.Info("Successfully loaded database config")
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode, c.TimeZone)
}
