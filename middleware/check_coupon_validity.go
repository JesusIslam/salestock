package middleware

import (
	"net/http"
	"time"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/model"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CheckCouponValidity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.FormValue("coupon_id") == "" {
			return next(c)
		}

		if !bson.IsObjectIdHex(c.FormValue("coupon_id")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid coupon_id: not a valid ObjectId")
		}
		couponID := bson.ObjectIdHex(c.FormValue("coupon_id"))
		db := database.New()
		defer db.Close()

		coupon := &model.Coupon{}
		err := db.DB("salestock").C("coupons").Find(bson.M{
			"_id": couponID,
			"valid_until": bson.M{
				"$gt": time.Now(),
			},
			"quantity": bson.M{
				"$gt": 0,
			},
		}).One(coupon)
		if err == mgo.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Invalid coupon_id: not found or not valid anymore")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
		}

		return next(c)
	}
}
