package mocks_models

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
)

/*
GetAll(ctx context.Context) ([]DomainInfo, error)

	GetDomainInfoById(ctx context.Context, id int64) (*DomainInfo, error)
	GetDomainInfoByName(ctx context.Context, name string) (*DomainInfo, error)
	DeleteDomain(ctx context.Context, id int64) error
	Save(ctx context.Context, di *DomainInfo) error
*/
type domainRepositoryMock struct {
	returnSave    error
	returnDelete  error
	returnGetById struct {
		err    error
		domain *domain.DomainInfo
	}
	returnGetByName struct {
		err    error
		domain *domain.DomainInfo
	}
	returnGetAll struct {
		err    error
		diList []domain.DomainInfo
	}
	GetByIdCalls   int
	SaveCalls      int
	GetAllCalls    int
	DeleteCalls    int
	GetByNameCalls int
}

func (drm *domainRepositoryMock) GetDomainInfoById(ctx context.Context, id int64) (*domain.DomainInfo, error) {
	drm.GetByIdCalls++
	return drm.returnGetById.domain, drm.returnGetById.err
}

func (drm *domainRepositoryMock) GetDomainInfoByName(ctx context.Context, name string) (*domain.DomainInfo, error) {
	drm.GetByNameCalls++
	return drm.returnGetByName.domain, drm.returnGetByName.err
}

func (drm *domainRepositoryMock) Save(ctx context.Context, dom *domain.DomainInfo) error {
	drm.SaveCalls++
	return drm.returnSave
}

func (drm *domainRepositoryMock) DeleteDomain(ctx context.Context, id int64) error {
	drm.DeleteCalls++
	return drm.returnDelete
}

func (drm *domainRepositoryMock) GetAll(ctx context.Context) ([]domain.DomainInfo, error) {
	drm.GetAllCalls++
	return drm.returnGetAll.diList, drm.returnGetAll.err
}

func (drm *domainRepositoryMock) SetGetByIdReturn(dom *domain.DomainInfo, err error) {
	drm.returnGetById = struct {
		err    error
		domain *domain.DomainInfo
	}{
		domain: dom,
		err:    err,
	}
}

func (drm *domainRepositoryMock) SetSaveReturn(err error) {
	drm.returnSave = err
}

func (drm *domainRepositoryMock) SetDeleteReturn(err error) {
	drm.returnDelete = err
}

func (drm *domainRepositoryMock) SetGetByNameReturn(dom *domain.DomainInfo, err error) {
	drm.returnGetByName = struct {
		err    error
		domain *domain.DomainInfo
	}{
		domain: dom,
		err:    err,
	}
}

func (drm *domainRepositoryMock) SetGetAllReturn(diList []domain.DomainInfo, err error) {
	drm.returnGetAll = struct {
		err    error
		diList []domain.DomainInfo
	}{
		diList: diList,
		err:    err,
	}
}
