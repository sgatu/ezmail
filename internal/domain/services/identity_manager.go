package services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/domain"
)

type IdentityManager interface {
	CreateIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error
	RefreshIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error
	DeleteIdentity(ctx context.Context, domainInfo *domain.DomainInfo) error
}
