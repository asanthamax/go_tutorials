package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthJwt(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")
	fmt.Print("Token string ", tokenString, "\n")
	if tokenString == "" {
		context.JSON(401, gin.H{"message": "Authorization header missing", "code": "UNAUTHORIZED"})
		context.Abort()
		return
	}

	err := ValidateToken(tokenString)
	if err != nil {
		context.JSON(401, gin.H{"message": err.Error(), "code": "UNAUTHORIZED"})
		context.Abort()
		return
	}
	context.Next()
}
