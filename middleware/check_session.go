package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JesusIslam/salestock/configuration"
	"github.com/JesusIslam/salestock/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := getToken(c.FormValue("token"), c.QueryParam("token"), c.Request().Header().Get("token"))
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: not found")
		}

		// validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(configuration.Server().JWTKey), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
		}

		var nonce, username, role string
		var ID bson.ObjectId
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: already expired")
			}
			nonce, ok = claims["nonce"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: nonce not a string")
			}
			username, ok = claims["username"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: username not a string")
			}
			role, ok = claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: role not a string")
			}
			id, ok := claims["ID"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: ID not a string")
			}
			if !bson.IsObjectIdHex(id) {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid token: ID not a valid ObjectId")
			}
			ID = bson.ObjectIdHex(id)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: invalid claims")
		}

		// check if nonce in blacklist
		db := database.New()
		defer db.Close()
		n, err := db.DB("salestock").C("blacklist").Find(bson.M{
			"nonce": nonce,
		}).Count()
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "Database error: "+err.Error())
		}
		if n > 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: nonce is blacklisted")
		}

		// set values to session
		c.Set("nonce", nonce)
		c.Set("username", username)
		c.Set("role", role)
		c.Set("ID", ID)

		return next(c)
	}
}

var bearer = "Bearer"

func getToken(fromForm, fromQuery, fromHeader string) (token string) {
	if fromForm != "" {
		token = fromForm
		return token
	}

	if fromQuery != "" {
		token = fromQuery
		return token
	}

	if fromHeader != "" {
		l := len(bearer)
		if len(fromHeader) > l+1 && fromHeader[:l] == bearer {
			token = fromHeader[l+1:]
		}
		return token
	}

	return token
}
