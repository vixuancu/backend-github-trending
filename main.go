package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"backend-github-trending/helper"
	"backend-github-trending/repository/repo_impl"
	"backend-github-trending/router"
	"backend-github-trending/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"time"
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
	e.Validator = utils.NewValidator()
	e.Use(utils.ValidationMiddleware())
	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepoImpl(sql),
	}

	// Khởi tạo RepoHandler
	githubRepo := repo_impl.NewGithubRepo(sql)
	repoHandler := handler.NewRepoHandler(githubRepo)

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		RepoHandler: repoHandler, // Truyền RepoHandler đã được khởi tạo vào API
	}
	api.SetupRouter() // thiết lập các route cho API
	go scheduleUpdateTrending(15*time.Second, *repoHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
func scheduleUpdateTrending(timeSchedule time.Duration, handler handler.RepoHandler) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Checking from github...")
				helper.CrawlRepo(handler.GithubRepo)
			}
		}
	}()
}
