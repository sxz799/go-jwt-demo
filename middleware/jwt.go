package middleware

import "github.com/golang-jwt/jwt/v4"

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
