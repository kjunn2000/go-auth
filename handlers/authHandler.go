package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kjunn2000/go-auth/model"
)

var Users model.Users

func init() {

	Users = make(map[string]*model.User)
	user := &model.User{
		Username: "kaijun",
		Password: "pass",
	}
	Users["kaijun"] = user
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	username := user.Username
	password := user.Password
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, ok := Users[username]
	if !ok || u.Password != password {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("hello1")
		return
	}
	atc := model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	rtc := model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 45).Unix(),
		},
	}
	att := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)
	rtt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	ats, err := att.SignedString(model.SecretKey)
	rts, err := rtt.SignedString(model.SecretKey)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	accessTokenCookie := &http.Cookie{
		Name:    "access_token",
		Value:   ats,
		Expires: time.Unix(atc.ExpiresAt, 0),
	}
	refreshTokenCookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   rts,
		Expires: time.Unix(rtc.ExpiresAt, 0),
	}
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	rt, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v := rt.Value

	if v == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c := &model.Claims{}
	token, err := jwt.ParseWithClaims(v, c,
		func(t *jwt.Token) (interface{}, error) {
			return model.SecretKey, nil
		})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if time.Unix(c.ExpiresAt, 0).Before(time.Now()) {
		w.WriteHeader(int(http.StatusBadRequest))
		return
	}

	atc := model.Claims{
		Username: c.Username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}

	att := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	ats, err := att.SignedString(model.SecretKey)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	accessTokenCookie := &http.Cookie{
		Name:    "access_token",
		Value:   ats,
		Expires: time.Unix(atc.ExpiresAt, 0),
	}

	http.SetCookie(w, accessTokenCookie)
}
