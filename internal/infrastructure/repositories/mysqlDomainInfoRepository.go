package repositories

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/uptrace/bun"
)

type mysqlDomainInfoRepository struct {
	db *bun.DB
}

func NewMysqlDomainInfoRepository(connection *bun.DB) *mysqlDomainInfoRepository {
	return &mysqlDomainInfoRepository{
		db: connection,
	}
}

func (repo *mysqlDomainInfoRepository) Save(ctx context.Context, di *domain.DomainInfo) error {
	_, err := repo.db.NewInsert().Model(di).Exec(ctx)
	return err
}
