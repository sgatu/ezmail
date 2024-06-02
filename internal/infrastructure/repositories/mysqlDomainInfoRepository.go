package repositories

import (
	"context"
	"database/sql"

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

func (repo *mysqlDomainInfoRepository) GetDomainInfoById(ctx context.Context, id string) (*domain.DomainInfo, error) {
	di := &domain.DomainInfo{Id: id}
	err := repo.db.NewSelect().Model(di).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDomainInfoNotFound
	} else if err != nil {
		return nil, err
	}
	return di, nil
}

func (repo *mysqlDomainInfoRepository) GetAllByUserId(ctx context.Context, userId string) ([]domain.DomainInfo, error) {
	var domainInfos []domain.DomainInfo
	err := repo.db.NewSelect().Model((*domain.DomainInfo)(nil)).Where("user_id = ?", userId).Scan(ctx, &domainInfos)
	if err != nil {
		return []domain.DomainInfo{}, err
	}
	return domainInfos, nil
}

func (repo *mysqlDomainInfoRepository) Save(ctx context.Context, di *domain.DomainInfo) error {
	_, err := repo.db.NewInsert().Model(di).Exec(ctx)
	return err
}
