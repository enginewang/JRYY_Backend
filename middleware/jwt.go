package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	// Valid token
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware Middleware of JWT
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get token from Header.Authorization field.
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Authorization is null in Header",
			})
			fmt.Println("Authorization is null in Header")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Format of Authorization is wrong",
			})
			fmt.Println("Format of Authorization is wrong")
			c.Abort()
			return
		}
		// parts[0] is Bearer, parts is token.
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token 1",
			})
			fmt.Println("Invalid Token 1")
			c.Abort()
			return
		}
		//fmt.Println("token is: ", mc)
		// Store Account info into Context
		c.Set("token", mc)
		// After that, we can get Account info from c.Get("account")
		c.Next()
		return
	}
}

// 只有管理员才有权限
func AdminAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token, ok := c.Get("token")
		if ok != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token 2",
			})
			fmt.Println("Invalid Token 2")
			c.Abort()
			return
		}
		for key, value := range *token.(*jwt.MapClaims) {
			if key == "type" {
				if value.(float64) < 1 {
					c.Next()
					return
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": -1,
						"msg":  "抱歉，您没有权限访问该页面",
					})
					fmt.Println("抱歉，您没有权限访问该页面")
					c.Abort()
					return
				}
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Invalid Token 4",
		})
		fmt.Println("Invalid Token 4")
		c.Abort()
		return
	}
}

// 只有指定用户才有权限
func UserAuth(k string) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, ok := c.Get("token")
		sid := c.Param(k)
		if ok != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token 5",
			})
			fmt.Println("Invalid Token 5")
			c.Abort()
			return
		}
		for key, value := range *token.(*jwt.MapClaims) {
			if key == "type" && value.(float64) < 1 {
				c.Next()
				return
			}
		}
		for key, value := range *token.(*jwt.MapClaims) {
			if key == "sid" {
				fmt.Println("value:", value)
				fmt.Println("sid:", sid)
				if value == sid {
					c.Next()
					return
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": -1,
						"msg":  "抱歉，您无权查看其他参与者的页面",
					})
					fmt.Println("Invalid Token 6")
					c.Abort()
					return
				}
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Invalid Token 7",
		})
		fmt.Println("Invalid Token 7")
		c.Abort()
		return
	}
}


// 只有登录了才有权限
func LoginAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		_, ok := c.Get("token")
		if ok != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token 9",
			})
			fmt.Println("Invalid Token 9")
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
