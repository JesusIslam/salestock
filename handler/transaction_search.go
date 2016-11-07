package handler

import (
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/labstack/echo"
)

func TransactionSearch(c echo.Context) (err error) {
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

	// support search by id only
	if len(params["id"]) > 0 {
		if !bson.IsObjectIdHex(params["id"][0]) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
		}
		search.ID = bson.ObjectIdHex(params["id"][0])
	}

	db := database.New()
	defer db.Close()

	transactions := []*model.Transaction{}
	err = db.DB("salestock").C("transactions").Find(bson.M{}).All(&transactions)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, transactions)
	return err
}
