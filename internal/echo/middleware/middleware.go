package middleware

import (
	// "time"

	// "github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
)

func ServerHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// start := time.Now()

		// Process the request
		err := next(c)
		if err != nil {
			c.Logger().Error(err.Error())
		}

		// stop := time.Now()
		// req := c.Request()

		// // Log using charmbracelet
		// log.Info("Request handlers",
		// 	"method", req.Method,
		// 	"path", req.URL.Path,
		// 	"latency", stop.Sub(start),
		// 	"ip", c.RealIP(),
		// )

		return nil
	}
}
