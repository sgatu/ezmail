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

func (repo *mysqlDomainInfoRepository) GetDomainInfoByName(ctx context.Context, name string) (*domain.DomainInfo, error) {
	di := &domain.DomainInfo{}
	err := repo.db.NewSelect().Model(di).Where("domain_name = ?", name).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDomainInfoNotFound
	} else if err != nil {
		return nil, err
	}
	return di, nil
}

func (repo *mysqlDomainInfoRepository) GetAll(ctx context.Context) ([]domain.DomainInfo, error) {
	var domainInfos []domain.DomainInfo
	err := repo.db.NewSelect().Model((*domain.DomainInfo)(nil)).Scan(ctx, &domainInfos)
	if err != nil {
		return []domain.DomainInfo{}, err
	}
	return domainInfos, nil
}

func (repo *mysqlDomainInfoRepository) Save(ctx context.Context, di *domain.DomainInfo) error {
	err := upsert(di, ctx, repo.db)
	return err
}
