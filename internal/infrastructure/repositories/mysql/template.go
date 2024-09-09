package mysql

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/uptrace/bun"
)

type mysqlTemplateRepository struct {
	db *bun.DB
}

func NewMysqlTemplateRepository(connection *bun.DB) *mysqlTemplateRepository {
	return &mysqlTemplateRepository{
		db: connection,
	}
}

func (repo *mysqlTemplateRepository) GetById(ctx context.Context, id int64) (*email.Template, error) {
	tpl := &email.Template{Id: id}
	err := repo.db.NewSelect().Model(tpl).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, email.ErrTemplateNotFound(strconv.FormatInt(id, 10))
	} else if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (repo *mysqlTemplateRepository) Save(ctx context.Context, tpl *email.Template) error {
	err := upsert(tpl, ctx, repo.db)
	return err
}
