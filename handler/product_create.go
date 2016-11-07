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

func ProductCreate(c echo.Context) error {
	productForm := &form.ProductCreate{}
	productForm.Name = c.FormValue("name")
	qtyInt64, err := strconv.ParseInt(c.FormValue("quantity"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity: "+err.Error())
	}
	productForm.Quantity = int(qtyInt64)
	productForm.Price, err = strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price: "+err.Error())
	}
	err = productForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	product := &model.Product{
		ID:       bson.NewObjectId(),
		Name:     productForm.Name,
		Price:    productForm.Price,
		Quantity: productForm.Quantity,
	}
	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("products").Insert(product)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusCreated, response.Response{
		Message: product,
	})
	return err
}
