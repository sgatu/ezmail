package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
)

type IdentityManagerMock struct {
	createIdentityResult error
	deleteIdentityResult error
	CreateIdentityCalls  int
	DeleteIdentityCalls  int
}

func MockIdentityManager() *IdentityManagerMock {
	return &IdentityManagerMock{}
}

func (imm *IdentityManagerMock) CreateIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	imm.CreateIdentityCalls++
	return imm.createIdentityResult
}

func (imm *IdentityManagerMock) DeleteIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	imm.DeleteIdentityCalls++
	return imm.deleteIdentityResult
}

func (imm *IdentityManagerMock) SetCreateIdentityResult(r error) {
	imm.createIdentityResult = r
}

func (imm *IdentityManagerMock) SetDeleteIdentityResult(r error) {
	imm.deleteIdentityResult = r
}
