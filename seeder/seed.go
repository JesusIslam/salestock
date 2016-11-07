package main

import (
	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/logger"
	"github.com/JesusIslam/salestock/model"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	var err error

	db := database.New()
	defer db.Close()

	// add default users
	admin := &model.User{
		ID:       bson.NewObjectId(),
		Username: "administrator",
		Role:     "admin",
		Password: "administrator",
	}
	err = admin.Validate()
	if err != nil {
		logger.Fatal(err)
	}
	err = db.DB("salestock").C("users").Insert(admin)
	if err != nil {
		logger.Fatal(err)
	}

	customer := &model.User{
		ID:       bson.NewObjectId(),
		Username: "customer",
		Role:     "customer",
		Password: "customer",
	}
	err = customer.Validate()
	if err != nil {
		logger.Fatal(err)
	}
	err = db.DB("salestock").C("users").Insert(customer)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Finish seeding")
}
