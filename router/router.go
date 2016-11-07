package router

import (
	"github.com/JesusIslam/salestock/configuration"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"

	"github.com/JesusIslam/salestock/handler"
	"github.com/JesusIslam/salestock/middleware"
)

func New() *echo.Echo {
	app := echo.New()

	app.Post("/login", handler.Login)
	app.Post("/logout", handler.Logout, middleware.CheckSession)

	app.Get("/product", handler.ProductSearch)
	app.Post("/product", handler.ProductCreate, middleware.CheckSession, middleware.IsRoleAdmin)
	app.Put("/product/:id", handler.ProductUpdate, middleware.CheckSession, middleware.IsRoleAdmin)
	app.Delete("/product/:id", handler.ProductDelete, middleware.CheckSession, middleware.IsRoleAdmin)

	app.Get("/coupon", handler.CouponSearch)
	app.Post("/coupon", handler.CouponCreate, middleware.CheckSession, middleware.IsRoleAdmin)
	app.Put("/coupon/:id", handler.CouponUpdate, middleware.CheckSession, middleware.IsRoleAdmin)
	app.Delete("/coupon/:id", handler.CouponDelete, middleware.CheckSession, middleware.IsRoleAdmin)

	app.Get("/user", handler.UserSearch)
	app.Post("/user", handler.UserCreate, middleware.CheckSession, middleware.IsRoleAdmin)
	app.Put("/user/:id", handler.UserUpdate, middleware.CheckSession, middleware.CheckRoleAdminForRoleUpdate, middleware.IsUserOwner)
	app.Delete("/user/:id", handler.UserDelete, middleware.CheckSession, middleware.IsRoleAdmin)

	app.Get("/transaction", handler.TransactionSearch)
	app.Post("/transaction", handler.TransactionCreate, middleware.CheckCouponValidity)
	app.Put("/transaction/:id", handler.TransactionUpdate, middleware.CheckCouponValidity, middleware.CheckRoleAdminForStatusUpdate)
	app.Delete("/transaction/:id", handler.TransactionDelete, middleware.CheckSession, middleware.IsRoleAdmin)

	return app
}

func NewEngine() *fasthttp.Server {
	return fasthttp.New(configuration.Server().Host)
}
