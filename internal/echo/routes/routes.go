package routes

import (
	"blog/internal/cache"
	"blog/internal/echo/handlers"
	"blog/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/labstack/echo-contrib/v5/echoprometheus"
)

// Helper method to get int param or default value
func GetIntQueryParam(c *echo.Context, name string, defaultValue int) int {
	value, err := strconv.Atoi(c.QueryParamOr(name, strconv.Itoa(defaultValue)))
	if err != nil {
		return defaultValue
	}

	return value
}

// Helper method to get bool param or default value
func GetBoolQueryParam(c *echo.Context, name string, defaultValue bool) bool {
	value, err := strconv.ParseBool(c.QueryParamOr(name, strconv.FormatBool(defaultValue)))
	if err != nil {
		return defaultValue
	}

	return value
}


// Souces is a struct that has the sources for any data that the routes need
type Sources struct {
	StrapiService *services.StrapiService
	Caches        []cache.Cache
}

func SetupRoutes(e *echo.Echo, sources Sources) {
	if sources.StrapiService == nil {
		log.Fatal("Error", "error", "strapi service is required")
	}

	// Global middleware
	
	// Remove trailing slash
	e.Pre(middleware.RemoveTrailingSlash())
	// Use GZIP compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{ Level: 5, }))
	// CORS origin
	e.Use(middleware.CORS("https://blog.darwincereska.dev", "http://localhost:3000", "http://0.0.0.0:3000"))
	// Request logger
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI: true,
		LogMethod: true,
		HandleError: true,
		LogLatency: true,
		LogRemoteIP: true,
		LogURIPath: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				e.Logger.Info("REQUEST", "method", v.Method, "uri", v.URI, "status", v.Status, "latency", v.Latency, "remote-ip", v.RemoteIP)
			} else {
				e.Logger.Error("REQUEST_ERROR", "error", v.Error.Error(), "method", v.Method, "uri", v.URIPath, "status", v.Status, "latency", v.Latency, "remote-ip", v.RemoteIP)
			}
			return nil
		},
	}))
	// Context timeout : DEFAULT 60s
	e.Use(middleware.ContextTimeout(time.Second * 60))
	// Prometheus metrics
	e.Use(echoprometheus.NewMiddleware("blog"))

	// Routes
	strapiRoutes(e, sources.StrapiService)

	// Static routes
	e.Static("/public", "web/static")

	// Special routes
	e.GET("/api*", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, e.Router().Routes())
	})

	e.GET("/metrics", echoprometheus.NewHandler())
}

// Setup Strapi routes
func strapiRoutes(e *echo.Echo, s *services.StrapiService) {
	// Post routes
	posts := e.Group("/api/posts") // Routing group

	// GET /api/posts/all
	posts.GET("/all", func(c *echo.Context) error {
		pageSize := GetIntQueryParam(c, "pageSize", 10)
		page := GetIntQueryParam(c, "page", 1)
	
		return handlers.GetPostSummaries(c, s, pageSize, page)
	})

	// GET /api/posts/featured
	posts.GET("/featured", func(c *echo.Context) error {
		pageSize := GetIntQueryParam(c, "pageSize", 10)
		page := GetIntQueryParam(c, "page", 1)

		return handlers.GetFeaturedPosts(c, s, pageSize, page)
	})

	// GET /api/posts/tag/:tag
	posts.GET("/tag/:tag", func(c *echo.Context) error {
		tag := c.Param("tag")
		pageSize := GetIntQueryParam(c, "pageSize", 10)
		page := GetIntQueryParam(c, "page", 1)

		return handlers.GetPostsByTag(c, s, tag, pageSize, page)
	})

	// GET /api/posts/post/:slug
	e.GET("/api/post/:slug", func(c *echo.Context) error {
		slug := c.Param("slug")

		return handlers.GetPost(c, s, slug)
	})
}
