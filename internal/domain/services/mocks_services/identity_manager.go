package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
)

type identityManagerMock struct {
	createIdentityResult error
	deleteIdentityResult error
	CreateIdentityCalls  int
	DeleteIdentityCalls  int
}

func MockIdentityManager() *identityManagerMock {
	return &identityManagerMock{}
}

func (imm *identityManagerMock) CreateIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	imm.CreateIdentityCalls++
	return imm.createIdentityResult
}

func (imm *identityManagerMock) DeleteIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error {
	imm.DeleteIdentityCalls++
	return imm.deleteIdentityResult
}

func (imm *identityManagerMock) SetCreateIdentityResult(r error) {
	imm.createIdentityResult = r
}

func (imm *identityManagerMock) SetDeleteIdentityResult(r error) {
	imm.deleteIdentityResult = r
}
