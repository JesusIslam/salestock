package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/form"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TransactionCreate(c echo.Context) (err error) {
	transactionForm := &form.TransactionCreate{}
	if c.FormValue("coupon_id") != "" {
		if !bson.IsObjectIdHex(c.FormValue("coupon_id")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid coupon_id: not a valid ObjectId")
		}
		transactionForm.CouponID = bson.ObjectIdHex(c.FormValue("coupon_id"))
	}

	if !bson.IsObjectIdHex(c.FormValue("customer_id")) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer_id: not a valid ObjectId")
	}
	transactionForm.CustomerID = bson.ObjectIdHex(c.FormValue("customer_id"))
	transactionForm.OrderStatus = "submitted"
	transactionForm.Shipment.Address = c.FormValue("address")
	transactionForm.Shipment.Email = c.FormValue("email")
	transactionForm.Shipment.Name = c.FormValue("name")
	transactionForm.Shipment.PhoneNumber = c.FormValue("phone_number")

	transactionForm.Products = []bson.ObjectId{}
	productsStr := strings.Split(c.FormValue("products"), ",")
	for _, p := range productsStr {
		if !bson.IsObjectIdHex(p) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid products' element id: not a valid ObjectId")
		}
		transactionForm.Products = append(transactionForm.Products, bson.ObjectIdHex(p))
	}
	// no shipment status & shipment id, update only by admin

	err = transactionForm.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transaction := &model.Transaction{
		ID:         bson.NewObjectId(),
		CouponID:   transactionForm.CouponID,
		CustomerID: transactionForm.CustomerID,
		Shipment: &model.Shipment{
			Address:     transactionForm.Shipment.Address,
			Email:       transactionForm.Shipment.Email,
			Name:        transactionForm.Shipment.Name,
			PhoneNumber: transactionForm.Shipment.PhoneNumber,
		},
		OrderStatus: transactionForm.OrderStatus,
	}

	db := database.New()
	defer db.Close()

	// fill the transaction.products and reduce each product's quantity
	if len(transaction.Products) > 0 {
		for _, productID := range transactionForm.Products {
			// get product info
			product := &model.Product{}
			err = db.DB("salestock").C("products").Find(bson.M{
				"_id": productID,
				"quantity": bson.M{
					"$gt": 0,
				},
			}).One(product)
			if err == mgo.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Invalid products' id element or quantity depleted")
			}
			if err != nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
			}
			// reduce quantity
			err = db.DB("salestock").C("products").UpdateId(productID, bson.M{
				"$inc": bson.M{
					"quantity": -1,
				},
			})
			if err == mgo.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Invalid products' id element or quantity depleted")
			}
			if err != nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
			}
			transaction.Products = append(transaction.Products, &struct {
				ID        bson.ObjectId `json:"id,omitempty"`
				Name      string        `json:"name,omitempty"`
				PriceEach float64       `json:"price_each,omitempty"`
			}{
				ID:        productID,
				Name:      product.Name,
				PriceEach: product.Price,
			})
		}
	}
	// reduce coupon quantity and cut total products' price
	if transaction.CouponID != "" {
		coupon := &model.Coupon{}
		err = db.DB("salestock").C("coupons").Find(bson.M{
			"_id": transaction.CouponID,
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
		err = db.DB("salestock").C("coupons").UpdateId(transaction.CouponID, bson.M{
			"$inc": bson.M{
				"quantity": -1,
			},
		})
		if err == mgo.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Invalid coupon_id: not found or not valid anymore")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
		}
		for _, product := range transaction.Products {
			transaction.TotalPrice += product.PriceEach
		}
		if coupon.DiscountType == "nominal" {
			transaction.TotalPrice -= coupon.Discount
		} else {
			transaction.TotalPrice -= (transaction.TotalPrice * (coupon.Discount / 100))
		}
	}

	err = transaction.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = db.DB("salestock").C("transactions").Insert(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	err = c.JSON(http.StatusCreated, response.Response{
		Message: transaction,
	})
	return err
}
