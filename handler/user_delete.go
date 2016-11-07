package handler

import (
	"net/http"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func UserDelete(c echo.Context) (err error) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
	}

	db := database.New()
	defer db.Close()

	err = db.DB("salestock").C("users").RemoveId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id: "+err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: "User deleted: " + id,
	})
	return err
}
