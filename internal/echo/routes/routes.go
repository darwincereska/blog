package routes

import (
	"blog/internal/echo/handlers"
	"blog/internal/echo/middleware"
	"blog/internal/services"
	"strconv"

	"github.com/labstack/echo/v5"
)

// Helper method to get int param or default value
func getIntParam(c *echo.Context, name string, defaultValue int) int {
	value, err := strconv.Atoi(c.ParamOr(name, strconv.Itoa(defaultValue)))
	if err != nil {
		return defaultValue
	}

	return value
}

func SetupRoutes(e *echo.Echo, s *services.StrapiService) {
	// Global middleware
	e.Use(middleware.ServerHandler)

	// Post routes
	posts := e.Group("/posts") // Routing group

	// GET /posts/all
	posts.GET("/all", func(c *echo.Context) error {
		pageSize := getIntParam(c, "pageSize", 10)
		page := getIntParam(c, "page", 1)

		return handlers.GetAllPosts(c, s, pageSize, page)
	})
}
