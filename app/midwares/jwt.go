package midwares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"live-chat/app/configSecret"
	"net/http"
	"strings"
	"time"
)

const expiration = time.Hour * 4

type Claims struct {
	UserID int
	jwt.StandardClaims
}

func GenerateToken(userId int) (string, error) {
	//创建声明
	secret, err := configSecret.GetJwtKey()
	if err != nil {
		return "", err
	}
	var Secret = []byte(secret)
	a := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gin-jwt-demo",
			Id:        "",
			NotBefore: 0,
			Subject:   "",
		},
	}

	//哈希方法创建签名
	tt := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	tokenString, err := tt.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func parseToken(tokenStr string) (*Claims, error) {
	secret, err := configSecret.GetJwtKey()
	if err != nil {
		return nil, err
	}
	var Secret = []byte(secret)
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	//检验token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 200,
				"msg":  "auth is null",
			})
			c.Abort()
			return
		}

		parts := strings.Split(tokenStr, ".")
		if len(parts) != 3 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 200,
				"msg":  "auth is error",
			})
			c.Abort()
			return
		}
		mc, err := parseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 200,
				"msg":  "Token is not vaild",
			})
			c.Abort()
			return
		} else if time.Now().Unix() > mc.ExpiresAt {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 200,
				"msg":  "Token is overdue",
			})
			c.Abort()
			return
		}
		c.Set("ID", mc.UserID)
		c.Next()
	}
}
