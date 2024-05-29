package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sgatu/ezmail/internal/domain/models/user"
	"github.com/uptrace/bun"
)

type MysqlUserRepository struct {
	connection *bun.DB
}

func NewMysqlUserRepository(connection *bun.DB) *MysqlUserRepository {
	return &MysqlUserRepository{
		connection: connection,
	}
}

func (userRepo *MysqlUserRepository) GetById(ctx context.Context, id string) (*user.User, error) {
	usr := &user.User{Id: id}
	err := userRepo.connection.NewSelect().Model(usr).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFoundError
	} else if err != nil {
		return nil, err
	}
	return usr, nil
}

func (userRepo *MysqlUserRepository) FindByEmailAndPassword(ctx context.Context, email string, password string) (*user.User, error) {
	u := &user.User{}
	err := userRepo.connection.NewSelect().Model(u).Where("email = ?", email).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFoundError
	} else if err != nil {
		return nil, err
	}
	if !u.VerifyPassword(&user.BcryptPasswordHasher{}, password) {
		fmt.Println("Could not verify passsword")
		return nil, user.ErrUserNotFoundError
	}
	return u, nil
}

func (userRepo *MysqlUserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := userRepo.connection.NewInsert().Model(u).Exec(ctx)
	return err
}
