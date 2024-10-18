package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/services"
)

type EmailerMock struct {
	sendResult     error
	SendEmailCalls int
}

func MockEmailer() *EmailerMock {
	return &EmailerMock{}
}

func (em *EmailerMock) SendEmail(ctx context.Context, email *services.PreparedEmail) error {
	em.SendEmailCalls++
	return em.sendResult
}

func (em *EmailerMock) SetSendResult(e error) {
	em.sendResult = e
}
