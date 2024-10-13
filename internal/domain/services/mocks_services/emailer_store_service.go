package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/services"
)

type emailerStoreServiceMock struct {
	createEmailResult     error
	markEmailAsSentResult error
	getByIdResult         struct {
		email *email.Email
		err   error
	}
	prepareEmailResult struct {
		email *services.PreparedEmail
		err   error
	}
	CreateEmailCalls     int
	GetByIdCalls         int
	PrepareEmailCalls    int
	MarkEmailAsSentCalls int
}

func MockEmailStoreService() *emailerStoreServiceMock {
	return &emailerStoreServiceMock{
		getByIdResult: struct {
			email *email.Email
			err   error
		}{
			email: nil,
			err:   nil,
		},
		prepareEmailResult: struct {
			email *services.PreparedEmail
			err   error
		}{
			email: nil,
			err:   nil,
		},
	}
}

func (emssm *emailerStoreServiceMock) SetCreateEmailResult(err error) {
	emssm.createEmailResult = err
}

func (emssm *emailerStoreServiceMock) SetMarkEmailAsSentResult(err error) {
	emssm.markEmailAsSentResult = err
}

func (emssm *emailerStoreServiceMock) SetGetByIdResult(em *email.Email, err error) {
	emssm.getByIdResult = struct {
		email *email.Email
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (emssm *emailerStoreServiceMock) SetPrepareEmailResult(em *services.PreparedEmail, err error) {
	emssm.prepareEmailResult = struct {
		email *services.PreparedEmail
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (emssm *emailerStoreServiceMock) CreateEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error {
	emssm.CreateEmailCalls++
	return emssm.createEmailResult
}

func (emssm *emailerStoreServiceMock) GetById(ctx context.Context, id int64) (*email.Email, error) {
	emssm.GetByIdCalls++
	return emssm.getByIdResult.email, emssm.getByIdResult.err
}

func (emssm *emailerStoreServiceMock) PrepareEmail(ctx context.Context, id int64) (*services.PreparedEmail, error) {
	emssm.GetByIdCalls++
	return emssm.prepareEmailResult.email, emssm.prepareEmailResult.err
}

func (emssm *emailerStoreServiceMock) MarkEmailAsSent(ctx context.Context, id int64) error {
	emssm.MarkEmailAsSentCalls++
	return emssm.markEmailAsSentResult
}
