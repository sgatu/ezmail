package repositories

import (
	"context"

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
