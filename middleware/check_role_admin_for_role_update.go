package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

func CheckRoleAdminForRoleUpdate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role: not a string")
		}

		if c.FormValue("role") != "" {
			if role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid role: not authorized")
			}
		}

		return next(c)
	}
}
