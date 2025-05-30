package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"backend-github-trending/repository/repo_impl"
	"backend-github-trending/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
		log.Error(err.Error())
		return
	}

	defer sql.Close()
	e := echo.New()
	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepoImpl(sql),
	}
	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
	}
	api.SetupRoter() // thiết lập các route cho API
	e.Logger.Fatal(e.Start(":8080"))
}
