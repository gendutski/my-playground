package main

import (
	"os"

	"github.com/gendutski/my-playground/dockerexample/config"
	"github.com/gendutski/my-playground/dockerexample/core/entity"
	"github.com/gendutski/my-playground/dockerexample/core/repository"
	"github.com/gendutski/my-playground/dockerexample/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// connect & auto migrate db
	db := config.Connect()
	db.AutoMigrate(&entity.User{})

	// set repository
	repo := repository.Init(db)
	// set handler
	controller := handler.Init(repo)

	// init echo
	e := echo.New()

	// routing
	e.GET("/", controller.Get)
	e.POST("/", controller.Post)

	// run
	e.Logger.Fatal(e.Start(":" + port))
}
