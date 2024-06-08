package repositories

import (
	"context"
	"database/sql"

	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/uptrace/bun"
)

type mysqlAuthTokenRepository struct {
	db *bun.DB
}

func NewMysqlAuthTokenRepository(connection *bun.DB) *mysqlAuthTokenRepository {
	return &mysqlAuthTokenRepository{
		db: connection,
	}
}

func (repo *mysqlAuthTokenRepository) GetAuthTokenById(ctx context.Context, id string) (*auth.AuthToken, error) {
	authToken := &auth.AuthToken{Id: id}
	err := repo.db.NewSelect().Model(authToken).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, auth.ErrNoAuthTokenFound
	} else if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (repo *mysqlAuthTokenRepository) GetAuthTokenByToken(ctx context.Context, token string) (*auth.AuthToken, error) {
	authToken := &auth.AuthToken{}
	err := repo.db.NewSelect().Model(authToken).Where("token = ?", token).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, auth.ErrNoAuthTokenFound
	} else if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (repo *mysqlAuthTokenRepository) GetAuthTokensByUserId(ctx context.Context, userId string) ([]auth.AuthToken, error) {
	var authTokens []auth.AuthToken
	err := repo.db.NewSelect().Model((*auth.AuthToken)(nil)).Where("user_id = ?", userId).Scan(ctx, &authTokens)
	if err != nil {
		return []auth.AuthToken{}, err
	}
	return authTokens, nil
}

func (repo *mysqlAuthTokenRepository) Save(ctx context.Context, authToken *auth.AuthToken) error {
	err := upsert(authToken, ctx, repo.db)
	return err
}
