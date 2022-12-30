package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenToken(u string) (accessTokenStr, refreshTokenStr string, err error) {
	// Create the JWT claims, which includes the username and expiry time
	accessClaims := &Claims{
		Username: u,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)),
		},
	}
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenStr, err = accessToken.SignedString(JwtKey)
	refreshTokenStr, err = refreshToken.SignedString(JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", "", err
	}
	return accessTokenStr, refreshTokenStr, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		accessTokenStr, err := c.Cookie("access-token")
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			c.String(200, "获取token失败")
			c.Abort()
			return
		}

		// Initialize a new instance of `Claims`
		accessClaims := &Claims{}
		refreshClaims := &Claims{}
		accessToken, err := jwt.ParseWithClaims(accessTokenStr, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		log.Println("accessToken过期时间:", accessClaims.ExpiresAt)
		if err != nil || !accessToken.Valid {
			refreshTokenStr, err1 := c.Cookie("refresh-token")
			refreshToken, err2 := jwt.ParseWithClaims(refreshTokenStr, refreshClaims, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})
			log.Println("refreshToken过期时间:", refreshClaims.ExpiresAt)
			if err1 != nil || err2 != nil || !refreshToken.Valid {
				c.String(200, "refreshToken不合法或已过期,需要重新登录")
				c.Abort()
				return
			}
			if refreshToken.Valid {
				accessTokenStr, refreshTokenStr, err := GenToken(accessClaims.Username)
				if err != nil {
					c.String(200, "更新token失败！")
					c.Abort()
					return
				}
				c.SetCookie("access-token", accessTokenStr, 3600, "", "", false, true)
				c.SetCookie("refresh-token", refreshTokenStr, 3600, "", "", false, true)
			}
		}
		c.Set("accessClaims", accessClaims)
		c.Next()
	}

}
