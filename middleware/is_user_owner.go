package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func IsUserOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, ok := c.Get("ID").(bson.ObjectId)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid ID: not an ObjectId")
		}

		if !bson.IsObjectIdHex(c.Param("id")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
		}

		sentID := bson.ObjectIdHex(c.Param("id"))
		if id != sentID {
			return echo.NewHTTPError(http.StatusForbidden, "Invalid ID: not data owner")
		}

		return next(c)
	}
}
