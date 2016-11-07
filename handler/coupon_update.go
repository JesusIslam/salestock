package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func CouponUpdate(c echo.Context) (err error) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id: not a valid ObjectId")
	}

	couponForm := &form.CouponUpdate{
		ID:           bson.ObjectIdHex(id),
		DiscountType: c.FormValue("discount_type"),
	}
	if c.FormValue("quantity") != "" {
		qtyInt64, err := strconv.ParseInt(c.FormValue("quantity"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity: "+err.Error())
		}
		couponForm.Quantity = int(qtyInt64)
	}
	if c.FormValue("discount") != "" {
		couponForm.Discount, err = strconv.ParseFloat(c.FormValue("discount"), 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid discount: "+err.Error())
		}
	}
	if c.FormValue("valid_until") != "" {
		couponForm.ValidUntil, err = time.Parse(time.RFC3339, c.FormValue("valid_until"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid valid_until: "+err.Error())
		}
	}
	err = couponForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ID, data := couponForm.ToUpdateData()

	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("coupons").UpdateId(ID, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	// get updated coupon
	coupon := &model.Coupon{}
	err = db.DB("salestock").C("coupons").FindId(ID).One(coupon)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: coupon,
	})
	return err
}
