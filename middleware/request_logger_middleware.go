package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLoggerMiddleware logs the request method, path, duration, and response status
func RequestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()          // Get the HTTP method (GET, POST, etc.)
		path := c.Path()              // Get the request path
		err := c.Next()               // Continue to the next middleware or handler
		duration := time.Since(start) // Measure the duration

		status := c.Response().StatusCode() // Get the HTTP status code of the response

		// Log the request information
		log.Printf("[%s] %s %s took %v, Status: %d", time.Now().Format(time.RFC3339), method, path, duration, status)

		return err
	}
}
