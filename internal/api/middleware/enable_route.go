package middleware

import (
	"github.com/labstack/echo/v4"
)

func EnableRoute(shouldEnable bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		if shouldEnable {
			return next
		}

		return echo.NotFoundHandler
	}
}
