package router

import (
	"github.com/JesusIslam/salestock/configuration"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
)

func New() *echo.Echo {
	app := echo.New()

	// add routes here

	return app
}

func NewEngine() *fasthttp.Server {
	return fasthttp.New(configuration.Server().Host)
}
