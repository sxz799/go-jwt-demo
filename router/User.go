package router

import (
	"github.com/gin-gonic/gin"
	"go-jwt-study/middleware"
	"go-jwt-study/model"
	"log"
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

	accessTokenStr, refreshTokenStr, err := middleware.GenToken(u.Username)
	if err != nil {
		log.Println(err)
		c.String(200, "token生成失败")
	}
	c.SetCookie("access-token", accessTokenStr, 3600, "", "", false, true)
	c.SetCookie("refresh-token", refreshTokenStr, 3600, "", "", false, true)
	c.String(200, "登录成功！")

}
func logout(c *gin.Context) {

	c.SetCookie("access-token", "-", 0, "", "", false, true)
	c.SetCookie("refresh-token", "-", 0, "", "", false, true)
	c.String(200, "退出成功")
}

func index(c *gin.Context) {

	value, exists := c.Get("accessClaims")
	if exists {
		claims := value.(*middleware.Claims)
		c.String(200, "欢迎%s", claims.Username)
	} else {
		c.String(200, "获取claims失败")
	}

}

func User(e *gin.Engine) {
	e.GET("/login", login)
	e.GET("/index", middleware.JWTAuth(), index)
	e.GET("/logout", logout)

}
