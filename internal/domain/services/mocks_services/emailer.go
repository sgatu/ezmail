package mocks_services

import (
	"context"

	"github.com/sgatu/ezmail/internal/domain/services"
)

type emailerMock struct {
	sendResult     error
	SendEmailCalls int
}

func MockEmailer() *emailerMock {
	return &emailerMock{}
}

func (em *emailerMock) SendEmail(ctx context.Context, email *services.PreparedEmail) error {
	em.SendEmailCalls++
	return em.sendResult
}

func (em *emailerMock) SetSendResult(e error) {
	em.sendResult = e
}
