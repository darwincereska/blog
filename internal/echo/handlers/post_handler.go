package handlers

import (
	"blog/internal/services"
	"net/http"

	"github.com/labstack/echo/v5"
)

// GetAllPosts returns a list of all posts
func GetAllPosts(c *echo.Context, s *services.StrapiService) error {
	pageSize := GetIntParam(c, "pageSize", 10)
	page := GetIntParam(c, "page", 1)

	posts, err := s.GetAllPosts(c.Request().Context(), pageSize, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// GetFeaturedPosts returns a list of featured posts
func GetFeaturedPosts(c *echo.Context, s *services.StrapiService) error {
	pageSize := GetIntParam(c, "pageSize", 10)
	page := GetIntParam(c, "page", 1)

	posts, err := s.GetFeaturedPosts(c.Request().Context(), pageSize, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPost returns a specific post
func GetPost(c *echo.Context, s *services.StrapiService, slug string) error {
	post, err := s.GetPost(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, post)
}

// GetPostSummaries returns post summaries
func GetPostSummaries(c *echo.Context, s *services.StrapiService) error {
	pageSize := GetIntParam(c, "pageSize", 10)
	page := GetIntParam(c, "page", 1)

	posts, err := s.GetPostSummaries(c.Request().Context(), pageSize, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPostsByTag returns a list of posts by tag
func GetPostsByTag(c *echo.Context, s *services.StrapiService, tag string, pageSize, page int) error {
	posts, err := s.GetPostsByTag(c.Request().Context(), tag, pageSize, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}
