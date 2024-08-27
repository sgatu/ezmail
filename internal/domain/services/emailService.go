package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type EmailService struct {
	emailRepository      email.EmailRepository
	templateRepository   email.TemplateRepository
	domainInfoRepository domain.DomainInfoRepository
	eventBus             events.EventBus
	snowFlakeNode        *snowflake.Node
}

func (emService *EmailService) SendEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error {
	_, err := emService.templateRepository.GetById(ctx, createEmail.TemplateId)
	if err != nil {
		return err
	}
	addr, err := mail.ParseAddress(createEmail.From)
	if err != nil {
		return err
	}
	addrParts := strings.SplitN(addr.Address, "@", 2)
	domainName := addrParts[1]

	for _, to := range createEmail.To {
		_, err := mail.ParseAddress(to)
		if err != nil {
			return err
		}
	}
	domain, err := emService.domainInfoRepository.GetDomainInfoByName(ctx, domainName)
	if err != nil {
		return err
	}
	if !domain.Validated {
		return fmt.Errorf("domain not validated")
	}
	toEmails, err := json.Marshal(createEmail.To)
	if err != nil {
		return err
	}
	emailEntity, err := email.NewEmail(emService.snowFlakeNode, domain.Id, createEmail.TemplateId, createEmail.From, string(toEmails), *createEmail.ReplyTo, *createEmail.CCO, createEmail.Context)
	if err != nil {
		return err
	}
	err = emService.emailRepository.Save(ctx, emailEntity)
	if err != nil {
		return err
	}

	return nil
}
