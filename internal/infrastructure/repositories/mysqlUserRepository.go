package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sgatu/ezmail/internal/domain/models/user"
	"github.com/uptrace/bun"
)

type mysqlUserRepository struct {
	db *bun.DB
}

func NewMysqlUserRepository(connection *bun.DB) *mysqlUserRepository {
	return &mysqlUserRepository{
		db: connection,
	}
}

func (userRepo *mysqlUserRepository) GetById(ctx context.Context, id string) (*user.User, error) {
	usr := &user.User{Id: id}
	err := userRepo.db.NewSelect().Model(usr).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFoundError
	} else if err != nil {
		return nil, err
	}
	return usr, nil
}

func (userRepo *mysqlUserRepository) FindByEmailAndPassword(ctx context.Context, email string, password string) (*user.User, error) {
	u := &user.User{}
	err := userRepo.db.NewSelect().Model(u).Where("email = ?", email).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFoundError
	} else if err != nil {
		return nil, err
	}
	if !u.VerifyPassword(user.BcryptPasswordHasher, password) {
		fmt.Println("Could not verify passsword")
		return nil, user.ErrUserNotFoundError
	}
	return u, nil
}

func (userRepo *mysqlUserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := userRepo.db.NewInsert().Model(u).Exec(ctx)
	return err
}
