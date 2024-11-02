package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/services"
)

type EmailerStoreServiceMock struct {
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

func MockEmailStoreService() *EmailerStoreServiceMock {
	return &EmailerStoreServiceMock{
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

func (emssm *EmailerStoreServiceMock) SetCreateEmailResult(err error) {
	emssm.createEmailResult = err
}

func (emssm *EmailerStoreServiceMock) SetMarkEmailAsSentResult(err error) {
	emssm.markEmailAsSentResult = err
}

func (emssm *EmailerStoreServiceMock) SetGetByIdResult(em *email.Email, err error) {
	emssm.getByIdResult = struct {
		email *email.Email
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (emssm *EmailerStoreServiceMock) SetPrepareEmailResult(em *services.PreparedEmail, err error) {
	emssm.prepareEmailResult = struct {
		email *services.PreparedEmail
		err   error
	}{
		email: em,
		err:   err,
	}
}

func (emssm *EmailerStoreServiceMock) CreateEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error {
	emssm.CreateEmailCalls++
	return emssm.createEmailResult
}

func (emssm *EmailerStoreServiceMock) GetById(ctx context.Context, id int64) (*email.Email, error) {
	emssm.GetByIdCalls++
	return emssm.getByIdResult.email, emssm.getByIdResult.err
}

func (emssm *EmailerStoreServiceMock) PrepareEmail(ctx context.Context, id int64) (*services.PreparedEmail, error) {
	emssm.PrepareEmailCalls++
	return emssm.prepareEmailResult.email, emssm.prepareEmailResult.err
}

func (emssm *EmailerStoreServiceMock) MarkEmailAsSent(ctx context.Context, id int64) error {
	emssm.MarkEmailAsSentCalls++
	return emssm.markEmailAsSentResult
}
