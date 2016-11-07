package handler

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
)

func Logout(c echo.Context) (err error) {
	nonce, ok := c.Get("nonce").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session nonce is not a string")
	}

	// add the nonce to blacklist
	db := database.New()
	defer db.Close()
	_, err = db.DB("salestock").C("blacklists").Upsert(bson.M{
		"nonce": nonce,
	}, bson.M{
		"nonce": nonce,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: "Successfully logged out: nonce blacklisted",
	})

	return err
}
