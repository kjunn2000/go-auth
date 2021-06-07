package model

import "github.com/dgrijalva/jwt-go"

var (
	SecretKey = []byte("secretkey")
)

type Users map[string]*User

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}
