package middleware

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/kjunn2000/go-auth/model"
)

func JwtTokenVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := r.Header.Get("Authorization")

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
		log.Printf("Successful authenticate user | Username => %s", c.Username)
		next.ServeHTTP(w, r)
	})
}
