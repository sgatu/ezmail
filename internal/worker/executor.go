package worker

import (
	"github.com/sgatu/ezmail/internal/domain/services"
)

type executor struct {
	emailStoreService services.EmailStoreService
	emailerService    services.Emailer
}

func (e *executor) Run() {
}
