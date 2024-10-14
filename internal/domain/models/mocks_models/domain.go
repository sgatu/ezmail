package mocks_models

import "github.com/sgatu/ezmail/internal/domain/models/domain"

/*
GetAll(ctx context.Context) ([]DomainInfo, error)

	GetDomainInfoById(ctx context.Context, id int64) (*DomainInfo, error)
	GetDomainInfoByName(ctx context.Context, name string) (*DomainInfo, error)
	DeleteDomain(ctx context.Context, id int64) error
	Save(ctx context.Context, di *DomainInfo) error
*/
type domainRepositoryMock struct {
	returnGetById struct {
		domain *domain.DomainInfo
		err    error
	}
	returnSave   error
	GetByIdCalls int
	SaveCalls    int
}
