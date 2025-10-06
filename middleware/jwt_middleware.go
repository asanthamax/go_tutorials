package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"tododemo.com/m/models"
)

type loginst struct {
	Username string `json:"username,omitempty"`
}

func GetSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func Login(c *gin.Context) {

	loginParams := loginst{}
	c.ShouldBindJSON(&loginParams)
	fmt.Print("Login params ", loginParams, "\n")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": loginParams.Username,
		// let the token be valid for one year
		"nbf": time.Date(2022, 01, 01, 12, 0, 0, 0, time.UTC).Unix(), //nbf: not before
		"exp": time.Date(2023, 01, 01, 12, 0, 0, 0, time.UTC).Unix(), //exp: expire
	})

	tokenStr, err := token.SignedString([]byte(GetSecretKey()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SignedResponse{
		Token:   tokenStr,
		Message: "logged in",
	})
	return
}
