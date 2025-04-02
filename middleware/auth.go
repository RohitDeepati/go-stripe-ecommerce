package middleware

import (
	"fmt"
	"go-ecommerce/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Missing authentication token"})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			fmt.Print("Invalid token format:", tokenParts)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			ctx.Abort()
			return
		}
		tokenString = tokenParts[1]
		claims, err := utils.VerifyToken(tokenString)
		if err != nil{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			ctx.Abort()
			return
		}
		user, userExists := claims["user"]
		role, roleExists := claims["role"]

		if !userExists || !roleExists{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid user or role"})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Set("role", role)
		ctx.Next()
	}
}

func UserRoleMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
		role, exists := ctx.Get("role")
		if !exists{
			ctx.JSON(http.StatusForbidden, gin.H{"error": "invalid user"})
			ctx.Abort()
			return
		}
		fmt.Println("role in token: ", role)
		ctx.Next()
	}
}