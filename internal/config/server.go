package config

import (
	"github.com/charmbracelet/log"
	"blog/internal/config/env"
)

type ServerConfig struct {
	Host string
	Port int
	StrapiHost string
	RedisHost string
	RedisPort int
	StrapiApiKey string
	CacheTTL int
	EchoMode string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) LoadConfig() {
	c.Host = env.GetString("HOST", "0.0.0.0")
	c.Port = env.GetInt("PORT", 3000)
	c.StrapiHost = env.GetString("STRAPI_HOST", "https://strapi.darwincereska.dev")
	c.StrapiApiKey = env.GetString("STRAPI_API_KEY", "")
	c.RedisHost = env.GetString("REDIS_HOST", "localhost")
	c.RedisPort = env.GetInt("REDIS_PORT", 6379)
	c.CacheTTL = env.GetInt("CACHE_TTL", 3600)
	c.EchoMode = env.GetString("ECHO_MODE", "release")

	log.Info("Sucessfully loaded server config")
	log.Info("Host", "host", c.Host)
	log.Info("Port", "port", c.Port)
	log.Info("Redis Host", "host", c.RedisHost)
	log.Info("Strapi URL", "host", c.StrapiHost)
	log.Info("Echo Mode", "mode", c.EchoMode)
	log.Info("Cache TTL", "ttl", c.CacheTTL)
}
