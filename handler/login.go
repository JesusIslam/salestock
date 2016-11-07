package handler

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/JesusIslam/salestock/configuration"
	"github.com/JesusIslam/salestock/database"
	"github.com/JesusIslam/salestock/logger"
	"github.com/JesusIslam/salestock/model"
	"github.com/JesusIslam/salestock/response"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	jose "gopkg.in/square/go-jose.v2"
)

var signer jose.Signer

func init() {
	var err error
	signer, err = jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte(configuration.Server().JWTKey),
	}, nil)
	if err != nil {
		logger.Fatal("Failed to instantiate signer: ", err)
	}
}

func Login(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user := &model.User{}
	db := database.New()
	defer db.Close()
	err = db.DB("salestock").C("users").Find(bson.M{
		"username": username,
	}).One(user)
	if err == mgo.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "User not found: "+err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
	}

	// check password
	hashedPass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode password: "+err.Error())
	}
	err = bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password: "+err.Error())
	}

	// generate token
	claims := jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 7 * 24).Unix(),
		"nonce":    strconv.FormatInt(time.Now().UnixNano(), 10),
		"username": user.Username,
		"role":     user.Role,
		"ID":       user.ID.Hex(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(configuration.Server().JWTKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token: "+err.Error())
	}

	err = c.JSON(http.StatusOK, response.Response{
		Message: token,
	})
	return err
}
