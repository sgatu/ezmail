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
type DomainRepositoryMock struct {
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

func MockDomainRepository() *DomainRepositoryMock {
	return &DomainRepositoryMock{
		returnGetAll: struct {
			err    error
			diList []domain.DomainInfo
		}{
			err:    nil,
			diList: make([]domain.DomainInfo, 0),
		},
		returnGetById: struct {
			err    error
			domain *domain.DomainInfo
		}{},
		returnGetByName: struct {
			err    error
			domain *domain.DomainInfo
		}{},
	}
}

func (drm *DomainRepositoryMock) GetDomainInfoById(ctx context.Context, id int64) (*domain.DomainInfo, error) {
	drm.GetByIdCalls++
	return drm.returnGetById.domain, drm.returnGetById.err
}

func (drm *DomainRepositoryMock) GetDomainInfoByName(ctx context.Context, name string) (*domain.DomainInfo, error) {
	drm.GetByNameCalls++
	return drm.returnGetByName.domain, drm.returnGetByName.err
}

func (drm *DomainRepositoryMock) Save(ctx context.Context, dom *domain.DomainInfo) error {
	drm.SaveCalls++
	return drm.returnSave
}

func (drm *DomainRepositoryMock) DeleteDomain(ctx context.Context, id int64) error {
	drm.DeleteCalls++
	return drm.returnDelete
}

func (drm *DomainRepositoryMock) GetAll(ctx context.Context) ([]domain.DomainInfo, error) {
	drm.GetAllCalls++
	return drm.returnGetAll.diList, drm.returnGetAll.err
}

func (drm *DomainRepositoryMock) SetGetByIdReturn(dom *domain.DomainInfo, err error) {
	drm.returnGetById = struct {
		err    error
		domain *domain.DomainInfo
	}{
		domain: dom,
		err:    err,
	}
}

func (drm *DomainRepositoryMock) SetSaveReturn(err error) {
	drm.returnSave = err
}

func (drm *DomainRepositoryMock) SetDeleteReturn(err error) {
	drm.returnDelete = err
}

func (drm *DomainRepositoryMock) SetGetByNameReturn(dom *domain.DomainInfo, err error) {
	drm.returnGetByName = struct {
		err    error
		domain *domain.DomainInfo
	}{
		domain: dom,
		err:    err,
	}
}

func (drm *DomainRepositoryMock) SetGetAllReturn(diList []domain.DomainInfo, err error) {
	drm.returnGetAll = struct {
		err    error
		diList []domain.DomainInfo
	}{
		diList: diList,
		err:    err,
	}
}
