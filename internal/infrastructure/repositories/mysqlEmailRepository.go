package repositories

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/uptrace/bun"
)

type mysqlEmailRepository struct {
	db *bun.DB
}

func NewMysqlEmailRepository(connection *bun.DB) *mysqlEmailRepository {
	return &mysqlEmailRepository{
		db: connection,
	}
}

func (repo *mysqlEmailRepository) GetById(ctx context.Context, id int64) (*email.Email, error) {
	emailEntity := &email.Email{Id: id}
	err := repo.db.NewSelect().Model(emailEntity).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, email.ErrEmailNotFound(strconv.FormatInt(id, 10))
	} else if err != nil {
		return nil, err
	}
	return emailEntity, nil
}

func (repo *mysqlEmailRepository) Save(ctx context.Context, emailEntity *email.Email) error {
	err := upsert(emailEntity, ctx, repo.db)
	return err
}
