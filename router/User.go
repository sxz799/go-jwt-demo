package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-jwt-study/middleware"
	"go-jwt-study/model"
	"log"
	"time"
)

func login(c *gin.Context) {
	var u model.User
	err := c.BindJSON(&u)
	if err != nil {
		c.String(200, "参数有误！")
		return
	}
	if u.Username == "a" && u.Password == "a" {

	} else {
		c.String(200, "用户名或密码有误！")
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &middleware.Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.String(200, "token异常！")
		return
	}
	c.SetCookie("token", tokenString, 3600, "", "", false, true)
	c.String(200, "登录成功！")

}
func logout(c *gin.Context) {

	c.SetCookie("token", "-", 0, "", "", false, true)
	c.String(200, "退出成功")
}

func index(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.String(200, "获取token失败")
		return
	}

	// Initialize a new instance of `Claims`
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})
	if err != nil {
		c.String(200, "token验证失败！")
		return
	}
	if !token.Valid {
		c.String(200, "token不合法或已过期！")
		return
	}
	log.Println(claims.RegisteredClaims.ExpiresAt)
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	c.String(200, "欢迎%s", claims.Username)
}

func User(e *gin.Engine) {
	e.GET("/login", login)
	e.GET("/index", index)
	e.GET("/logout", logout)

}
