package handler

import (
	"net/http"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func UserUpdate(c echo.Context) (err error) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
	}

	userForm := &form.UserUpdate{
		ID: bson.ObjectIdHex(id),
	}

	if c.FormValue("username") != "" {
		userForm.Username = c.FormValue("username")
	}

	if c.FormValue("password") != "" {
		userForm.Password = c.FormValue("password")
	}

	if c.FormValue("role") != "" {
		userForm.Role = c.FormValue("role")
	}

	err = userForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ID, data := userForm.ToUpdateData()

	db := database.New()
	defer db.Close()

	// check if username already taken
	n, err := db.DB("salestock").C("users").Find(bson.M{
		"username": userForm.Username,
	}).Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}
	if n > 0 {
		return echo.NewHTTPError(http.StatusConflict, "Invalid username: already taken")
	}

	err = db.DB("salestock").C("users").UpdateId(ID, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	// get updated user
	user := &model.User{}
	err = db.DB("salestock").C("users").FindId(ID).One(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}
	user.Password = ""

	err = c.JSON(http.StatusOK, response.Response{
		Message: user,
	})
	return err
}
