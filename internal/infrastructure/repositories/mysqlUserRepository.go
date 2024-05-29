package repositories

import (
	"context"
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
	user := &user.User{Id: id}
	err := userRepo.connection.NewSelect().Model(user).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *MysqlUserRepository) FindByEmailAndPassword(ctx context.Context, email string, password string) (*user.User, error) {
	u := &user.User{}
	err := userRepo.connection.NewSelect().Model(u).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if !u.VerifyPassword(&user.BcryptPasswordHasher{}, password) {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (userRepo *MysqlUserRepository) Save(ctx context.Context, u user.User) error {
	_, err := userRepo.connection.NewInsert().Model(u).Exec(ctx)
	return err
}
