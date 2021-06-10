package postgresql

import (
	"fmt"

	squirrel "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/kjunn2000/go-auth/internal/go-auth/model"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AuthStore struct {
	log *zap.Logger
	db  *sqlx.DB
	sq  squirrel.StatementBuilderType
}

func NewAuthStore(log *zap.Logger, connStr string) *AuthStore {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Warn("Postgres db connection failed.")
	}
	return &AuthStore{
		log: log,
		db:  db,
		sq:  sq,
	}
}

func (s *AuthStore) SaveUser(u *model.User) error {
	sql, arg, err := s.sq.Insert("account_credential").
		Columns("userid", "username", "accpassword").
		Values(u.UserId, u.Username, u.AccPassword).ToSql()
	if err != nil {
		s.log.Warn("Failed to create insert statement.")
		return err
	}
	res, err := s.db.Exec(sql, arg...)
	if err != nil {
		fmt.Println(err)
		s.log.Warn("Failed to insert record to db.")
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		s.log.Warn("Failed to read result data.")
		return err
	}
	s.log.Info("Successful insert record to db.", zap.Int64("count", r))
	return nil
}

func (s *AuthStore) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	fmt.Println(user)
	sta, arg, err := s.sq.Select("*").From("account_credential").Where(squirrel.Eq{"username": username}).ToSql()
	if err != nil {
		s.log.Warn("Unable to create select sql.")
		return nil, err
	}

	err = s.db.Get(&user, sta, arg...)
	if err != nil {
		return nil, err
	}
	s.log.Info("Successful select username : ", username)
	return &user, nil
}
