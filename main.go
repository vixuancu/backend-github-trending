package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	// connect database
	sql := &db.Sql{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "123456",
		Database: "github_trending",
	}
	if err := sql.Connect(); err != nil {
		return // error connecting to database
	}

	defer sql.Close()
	e := echo.New()
	e.GET("/", handler.HandlerWelcome)
	e.GET("/user/sign-in", handler.HandleSignin)
	e.GET("/user/sign-up", handler.HandleSignup)

	e.Logger.Fatal(e.Start(":8080"))
}
