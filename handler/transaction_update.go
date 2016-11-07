package handler

import (
	"net/http"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TransactionUpdate(c echo.Context) (err error) {
	// only support update shipment id, shipment status, and order status
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
	}

	transactionForm := &form.TransactionUpdate{
		ID: bson.ObjectIdHex(id),
	}

	if c.FormValue("shipment_id") != "" {
		transactionForm.Shipment.ID = c.FormValue("shipment_id")
	}
	if c.FormValue("shipment_status") != "" {
		transactionForm.Shipment.Status = c.FormValue("shipment_status")
	}
	if c.FormValue("order_status") != "" {
		transactionForm.OrderStatus = c.FormValue("order_status")
	}
	err = transactionForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ID, data := transactionForm.ToUpdateData()

	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("transactions").UpdateId(ID, data)
	if err == mgo.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid id: not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	transaction := &model.Transaction{}
	err = db.DB("salestock").C("transactions").FindId(transactionForm.ID).One(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, transaction)
	return err
}
