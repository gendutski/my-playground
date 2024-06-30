package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	secret string = "but it should be secret"
	// secret string = "hana maryam latifa"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		response := map[string]string{}

		claims := map[string]jwt.MapClaims{
			"validUser": {
				"id":      1,
				"name":    "Firman Darmawan",
				"email":   "mvp.firman.darmawan@gmail.com",
				"isValid": true,
				"exp":     time.Now().Add(time.Minute * 2).Unix(),
			},
			"invalidUser": {
				"id":      2,
				"name":    "Firman Darmawan",
				"email":   "firman@ruangguru.com",
				"isValid": false,
				"exp":     time.Now().Add(time.Minute * 2).Unix(),
			},
		}
		for key, cl := range claims {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
			tokenStr, err := token.SignedString([]byte(secret))
			if err != nil {
				return err
			}
			response[key] = tokenStr
		}

		return c.JSON(http.StatusOK, response)
	})

	usrGroup := e.Group("/user", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		ContextKey: "userToken",
		SuccessHandler: func(c echo.Context) {
			token, ok := c.Get("userToken").(*jwt.Token) // by default token is stored under `user` key
			if !ok {
				c.Error(errors.New("JWT token missing or invalid"))
			}
			claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
			if !ok {
				c.Error(errors.New("failed to cast claims as jwt.MapClaims"))
			}

			if claims["isValid"].(bool) {
				c.Set("user", User{
					ID:    int(claims["id"].(float64)),
					Name:  claims["name"].(string),
					Email: claims["email"].(string),
				})
			} else {
				c.Error(errors.New("user is invalid"))
			}
		},
	}))
	usrGroup.GET("", func(c echo.Context) error {
		user, ok := c.Get("user").(User)
		if !ok {
			return errors.New("invalid user")
		}
		return c.JSON(http.StatusOK, user)
	})

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
