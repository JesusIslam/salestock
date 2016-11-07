package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

func IsRoleAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Invalid role: not a string")
		}

		if role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "Invalid role: forbidden")
		}

		return next(c)
	}
}
