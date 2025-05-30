package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleSignin(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Welcome to Sign in",
	})
}
func HandleSignup(c echo.Context) error {
	type User struct {
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}
	user := User{
		Email:    "test@gmail.com",
		Fullname: "Test User",
		Age:      30,
	}
	return c.JSON(http.StatusOK, user)

}
