package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
	"time"

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

func NewEmailService(
	emailRepository email.EmailRepository,
	templateRepository email.TemplateRepository,
	domainInfoRepository domain.DomainInfoRepository,
	eventBus events.EventBus,
	snowFlakeNode *snowflake.Node,
) *EmailService {
	return &EmailService{
		emailRepository:      emailRepository,
		templateRepository:   templateRepository,
		domainInfoRepository: domainInfoRepository,
		eventBus:             eventBus,
		snowFlakeNode:        snowFlakeNode,
	}
}

func (emService *EmailService) SendEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error {
	if createEmail.TemplateId == 0 {
		return email.ErrTemplateNotFound("none")
	}
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

	if createEmail.To == nil {
		return fmt.Errorf("no email destinataries")
	}

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
	replyToEmail, err := getReplyTo(createEmail)
	if err != nil {
		return err
	}
	ccoEmails, err := getCCO(createEmail)
	if err != nil {
		return err
	}
	contextInfo := createEmail.Context
	if contextInfo == nil {
		contextInfo = make(map[string]string)
	}
	emailEntity, err := email.NewEmail(
		emService.snowFlakeNode,
		domain.Id,
		createEmail.TemplateId,
		createEmail.From,
		string(toEmails),
		replyToEmail,
		string(ccoEmails),
		contextInfo,
	)
	if err != nil {
		return err
	}
	err = emService.emailRepository.Save(ctx, emailEntity)
	if err != nil {
		return err
	}
	emService.eventBus.Push(ctx, &events.NewEmailEvent{
		When: time.Now().UTC(),
		Id:   emailEntity.Id,
		Type: "new_email",
	}, "emails")
	return nil
}

func getReplyTo(createEmail *email.CreateNewEmailRequest) (string, error) {
	replyToEmail := ""
	if createEmail.ReplyTo != nil && len(*createEmail.ReplyTo) > 0 {
		_, err := mail.ParseAddress(*createEmail.ReplyTo)
		if err != nil {
			return "", err
		}
		replyToEmail = *createEmail.ReplyTo
	}
	return replyToEmail, nil
}

func getCCO(createEmail *email.CreateNewEmailRequest) ([]byte, error) {
	if createEmail.CCO == nil {
		return []byte("[]"), nil
	}
	for _, cco := range createEmail.CCO {
		_, err := mail.ParseAddress(cco)
		if err != nil {
			return nil, err
		}
	}
	ccoEmails, err := json.Marshal(createEmail.CCO)
	if err != nil {
		return nil, err
	}
	return ccoEmails, nil
}
