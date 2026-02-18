package routes

import (
	"blog/internal/cache"
	"blog/internal/echo/handlers"
	"blog/internal/echo/middleware"
	"blog/internal/services"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
)

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
	e.Use(middleware.ServerHandler)

	// Routes
	strapiRoutes(e, sources.StrapiService)

	// Special routes
	e.GET("/api", func(c *echo.Context) error {
		// Load all routes
		// routes, _ := json.MarshalIndent(e.Router().Routes(), "", "")
		return c.JSON(http.StatusOK, e.Router().Routes())
	})
}

// Setup Strapi routes
func strapiRoutes(e *echo.Echo, s *services.StrapiService) {
	// Post routes
	posts := e.Group("/api/posts") // Routing group

	// GET /api/posts/all
	posts.GET("/all", func(c *echo.Context) error {
		return handlers.GetAllPosts(c, s)
	})

	// GET /api/posts/featured
	posts.GET("/featured", func(c *echo.Context) error {
		return handlers.GetFeaturedPosts(c, s)
	})
}
