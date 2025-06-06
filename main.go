package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"backend-github-trending/helper"
	"backend-github-trending/log"
	"backend-github-trending/repository/repo_impl"
	"backend-github-trending/router"
	"backend-github-trending/utils"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"strconv"
	"time"
)

func init() {
	fmt.Println("DEV ENVIROMENT")
	fmt.Println("test Makefile")
	fmt.Println("test Makefile 4")
	//os.Setenv("APP_NAME", "github")
	log.InitLogger(false)
}
func main() {
	// connect database
	dbHost := getEnv("DB_HOST", "localhost")
	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		dbPort = 5432 // Mặc định nếu không thể chuyển đổi
	}

	sql := &db.Sql{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "123456"),
		Database: getEnv("DB_NAME", "github_trending"),
	}

	if err := sql.Connect(); err != nil {
		log.Error(err.Error())
		return
	}

	defer sql.Close()
	var email string
	err = sql.Db.GetContext(context.Background(), &email, "SELECT email FROM users WHERE email=$1", "abc@gmail.com")
	if err != nil {
		log.Error(err.Error())
	}
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
	go scheduleUpdateTrending(360*time.Second, *repoHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

// getEnv đọc biến môi trường, nếu không có thì sử dụng giá trị mặc định
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
