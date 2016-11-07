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

func CouponCreate(c echo.Context) error {
	couponForm := &form.CouponCreate{}
	qtyInt64, err := strconv.ParseInt(c.FormValue("quantity"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity: "+err.Error())
	}
	couponForm.Quantity = int(qtyInt64)

	couponForm.ValidUntil, err = time.Parse(time.RFC3339, c.FormValue("valid_until"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid valid_until: "+err.Error())
	}

	couponForm.Discount, err = strconv.ParseFloat(c.FormValue("discount"), 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid discount: "+err.Error())
	}

	couponForm.DiscountType = c.FormValue("discount_type")
	err = couponForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	coupon := &model.Coupon{
		ID:           bson.NewObjectId(),
		Quantity:     couponForm.Quantity,
		ValidUntil:   couponForm.ValidUntil,
		Discount:     couponForm.Discount,
		DiscountType: couponForm.DiscountType,
	}
	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("coupons").Insert(coupon)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusCreated, response.Response{
		Message: coupon,
	})
	return err
}
