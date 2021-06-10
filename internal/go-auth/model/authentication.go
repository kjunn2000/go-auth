package model

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

var (
	SecretKey = []byte("secretkey")
)

type Users map[string]*User

type User struct {
	UserId      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	AccPassword string    `json:"password"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}
