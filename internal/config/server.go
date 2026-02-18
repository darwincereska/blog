package config

import (
	"blog/internal/config/env"
	"github.com/charmbracelet/log"
)

type ServerConfig struct {
	Host           string
	Port           int
	StrapiEndpoint string
	RedisHost      string
	RedisPort      int
	StrapiToken    string
	CacheTTL       int
	EchoMode       string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) LoadConfig() {
	c.Host = env.GetString("HOST", "0.0.0.0")
	c.Port = env.GetInt("PORT", 3000)
	c.StrapiEndpoint = env.GetString("STRAPI_ENDPOINT", "https://strapi.darwincereska.dev")
	c.StrapiToken = env.GetString("STRAPI_TOKEN", "")
	c.RedisHost = env.GetString("REDIS_HOST", "localhost")
	c.RedisPort = env.GetInt("REDIS_PORT", 6379)
	c.CacheTTL = env.GetInt("CACHE_TTL", 3600)
	c.EchoMode = env.GetString("ECHO_MODE", "release")

	log.Info("Sucessfully loaded server config")
	log.Info("Host", "host", c.Host)
	log.Info("Port", "port", c.Port)
	log.Info("Redis Host", "host", c.RedisHost)
	log.Info("Strapi URL", "host", c.StrapiEndpoint)
	log.Info("Echo Mode", "mode", c.EchoMode)
	log.Info("Cache TTL", "ttl", c.CacheTTL)
}
