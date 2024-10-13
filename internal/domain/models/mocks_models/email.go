package mocks_models

import "context"

type emailRepositoryMock struct{}

func (erm *emailRepositoryMock) GetById(ctx context.Context, id int64) (*Email, error) {
}

func (erm *emailRepositoryMock) Save(ctx context.Context, email *Email) error {
}
