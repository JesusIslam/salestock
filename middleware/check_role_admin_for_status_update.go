package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

func CheckRoleAdminForStatusUpdate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role: not a string")
		}

		if c.FormValue("order_status") != "" || c.FormValue("shipment_status") != "" {
			if role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid role: not authorized")
			}
		}

		return next(c)
	}
}
