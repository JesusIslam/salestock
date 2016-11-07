package handler

import (
	"net/http"
	"strconv"

	mgo "gopkg.in/mgo.v2"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func ProductSearch(c echo.Context) (err error) {
	search := &form.ProductSearch{}

	params := c.QueryParams()
	if len(params["page"]) > 0 {
		page64, err := strconv.ParseInt(params["page"][0], 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid page: "+err.Error())
		}
		search.Page = int(page64)
	}
	if len(params["per_age"]) > 0 {
		perPage64, err := strconv.ParseInt(params["per_page"][0], 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid per_page: "+err.Error())
		}
		search.PerPage = int(perPage64)
	}
	if len(params["order_by"]) > 0 {
		search.OrderBy = params["order_by"][0]
	}

	if len(params["id"]) > 0 {
		if !bson.IsObjectIdHex(params["id"][0]) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
		}
		search.ID = bson.ObjectIdHex(params["id"][0])
	}
	if len(params["name"]) > 0 {
		search.Name = params["name"][0]
	}
	if len(params["quantity_less_than_equal"]) > 0 {
		q64, err := strconv.ParseInt(params["quantity_less_than_equal"][0], 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity_less_than_equal: "+err.Error())
		}
		search.QuantityLessThanEqual = int(q64)
	}
	if len(params["quantity_more_than_equal"]) > 0 {
		q64, err := strconv.ParseInt(params["quantity_more_than_equal"][0], 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity_more_than_equal: "+err.Error())
		}
		search.QuantityMoreThanEqual = int(q64)
	}
	if len(params["price_less_than_equal"]) > 0 {
		search.PriceLessThanEqual, err = strconv.ParseFloat(params["price_less_than_equal"][0], 10)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid price_less_than_equal: "+err.Error())
		}
	}
	if len(params["price_more_than_equal"]) > 0 {
		search.PriceMoreThanEqual, err = strconv.ParseFloat(params["price_more_than_equal"][0], 10)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid price_more_than_equal: "+err.Error())
		}
	}

	err = search.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	query := search.ToSearchQuery()
	products := []*model.Product{}
	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("products").Find(query).All(&products)
	if err != mgo.ErrNotFound && err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: products,
	})
	return err
}
