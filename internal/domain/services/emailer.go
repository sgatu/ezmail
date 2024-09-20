package services

import "context"

type Emailer interface {
	SendEmail(ctx context.Context, email *PreparedEmail) error
}
