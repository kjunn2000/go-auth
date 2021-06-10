package repository

import (
	"github.com/kjunn2000/go-auth/internal/go-auth/model"
)

type AuthenticationStore interface {
	SaveUser(*model.User) error
	FindUserByUsername(username string) (*model.User, error)
}
