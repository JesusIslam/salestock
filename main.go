package main

import (
	"github.com/JesusIslam/salestock/logger"
	"github.com/JesusIslam/salestock/router"
)

func main() {
	logger.Info("Transaction service is running")
	err := router.New().Run(router.NewEngine())
	if err != nil {
		logger.Fatal("Failed to run server: ", err.Error())
	}
}
