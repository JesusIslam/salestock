package handler

import (
	"net/http"
	"strconv"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func ProductUpdate(c echo.Context) (err error) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
	}

	productForm := &form.ProductUpdate{
		ID: bson.ObjectIdHex(id),
	}
	if c.FormValue("name") != "" {
		productForm.Name = c.FormValue("name")
	}
	if c.FormValue("quantity") != "" {
		qtyInt64, err := strconv.ParseInt(c.FormValue("quantity"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity: "+err.Error())
		}
		productForm.Quantity = int(qtyInt64)
	}
	if c.FormValue("price") != "" {
		productForm.Price, err = strconv.ParseFloat(c.FormValue("price"), 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid price: "+err.Error())
		}
	}
	err = productForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ID, data := productForm.ToUpdateData()

	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("products").UpdateId(ID, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	// get updated product
	product := &model.Product{}
	err = db.DB("salestock").C("products").FindId(ID).One(product)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: product,
	})
	return err
}
