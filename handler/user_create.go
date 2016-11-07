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

func UserCreate(c echo.Context) (err error) {
	userForm := &form.UserCreate{}
	userForm.Username = c.FormValue("username")
	userForm.Role = c.FormValue("role")
	userForm.Password = c.FormValue("password")
	err = userForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &model.User{
		ID:       bson.NewObjectId(),
		Password: userForm.Password,
		Username: userForm.Username,
		Role:     userForm.Role,
	}
	err = user.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := database.New()
	defer db.Close()

	// check if username already taken
	n, err := db.DB("salestock").C("users").Find(bson.M{
		"username": user.Username,
	}).Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}
	if n > 0 {
		return echo.NewHTTPError(http.StatusConflict, "Invalid username: already taken")
	}

	err = db.DB("salestock").C("users").Insert(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusCreated, response.Response{
		Message: user,
	})
	return err
}
