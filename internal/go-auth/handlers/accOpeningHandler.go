package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kjunn2000/go-auth/config"
	"github.com/kjunn2000/go-auth/internal/go-auth/model"
	"github.com/kjunn2000/go-auth/internal/go-auth/postgresql"
	"go.uber.org/zap"
)

func AccOpeningHandler(w http.ResponseWriter, r *http.Request) {

	Log, _ := zap.NewDevelopment()
	as := postgresql.NewAuthStore(Log, config.NewConnString())

	var user model.User
	json.NewDecoder(r.Body).Decode(&user)
	_, err := as.FindUserByUsername(user.Username)

	if err != sql.ErrNoRows {
		Log.Info("Username already exist.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user.UserId = uuid.New()

	err = as.SaveUser(&user)
	if err != nil {
		return
	}
}
