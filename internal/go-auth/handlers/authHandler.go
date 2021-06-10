package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/kjunn2000/go-auth/config"
	"github.com/kjunn2000/go-auth/internal/go-auth/model"
	"github.com/kjunn2000/go-auth/internal/go-auth/postgresql"
	"go.uber.org/zap"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	username := user.Username
	password := user.AccPassword
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Log, _ := zap.NewDevelopment()
	as := postgresql.NewAuthStore(Log, config.NewConnString())

	u, err := as.FindUserByUsername(username)

	if err != nil || u.AccPassword != password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	atc := model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.Now(),
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 10)),
		},
	}
	rtc := model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.Now(),
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 45)),
		},
	}
	att := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)
	rtt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	ats, aerr := att.SignedString(model.SecretKey)
	rts, rerr := rtt.SignedString(model.SecretKey)

	if aerr != nil || rerr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accessTokenCookie := &http.Cookie{
		Name:    "access_token",
		Value:   ats,
		Expires: atc.ExpiresAt.Time,
	}
	refreshTokenCookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   rts,
		Expires: rtc.ExpiresAt.Time,
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

	if c.ExpiresAt.Time.Before(time.Now()) {
		w.WriteHeader(int(http.StatusBadRequest))
		return
	}

	atc := model.Claims{
		Username: c.Username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.Now(),
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 10)),
		},
	}

	att := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	ats, err := att.SignedString(model.SecretKey)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accessTokenCookie := &http.Cookie{
		Name:    "access_token",
		Value:   ats,
		Expires: c.ExpiresAt.Time,
	}

	http.SetCookie(w, accessTokenCookie)
}
