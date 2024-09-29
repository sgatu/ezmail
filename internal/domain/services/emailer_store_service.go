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
	"github.com/uptrace/bun"
)

type PreparedEmail struct {
	Html    string
	Text    string
	From    string
	To      string
	BCC     string
	ReplyTo string
	Subject string
	Domain  struct {
		Name string
		Id   int64 `json:",string"`
	}
	Id int64 `json:",string"`
}

type EmailStoreService interface {
	CreateEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error
	GetById(ctx context.Context, id int64) (*email.Email, error)
	PrepareEmail(ctx context.Context, id int64) (*PreparedEmail, error)
	MarkEmailAsSent(ctx context.Context, id int64) error
}

type DefaultEmailStoreService struct {
	emailRepository      email.EmailRepository
	templateRepository   email.TemplateRepository
	domainInfoRepository domain.DomainInfoRepository
	eventBus             events.EventBus
	snowFlakeNode        *snowflake.Node
}

func NewDefaultEmailStoreService(
	emailRepository email.EmailRepository,
	templateRepository email.TemplateRepository,
	domainInfoRepository domain.DomainInfoRepository,
	eventBus events.EventBus,
	snowFlakeNode *snowflake.Node,
) EmailStoreService {
	return &DefaultEmailStoreService{
		emailRepository:      emailRepository,
		templateRepository:   templateRepository,
		domainInfoRepository: domainInfoRepository,
		eventBus:             eventBus,
		snowFlakeNode:        snowFlakeNode,
	}
}

func (dEmailer *DefaultEmailStoreService) GetById(ctx context.Context, id int64) (*email.Email, error) {
	return dEmailer.emailRepository.GetById(ctx, id)
}

func (dEmailer *DefaultEmailStoreService) MarkEmailAsSent(ctx context.Context, id int64) error {
	email, err := dEmailer.emailRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	email.ProcessedDate = bun.NullTime{Time: time.Now().UTC()}
	return dEmailer.emailRepository.Save(ctx, email)
}

func (dEmailer *DefaultEmailStoreService) PrepareEmail(ctx context.Context, id int64) (*PreparedEmail, error) {
	email, err := dEmailer.emailRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	template, err := dEmailer.templateRepository.GetById(ctx, email.TemplateId)
	if err != nil {
		return nil, err
	}
	domain, err := dEmailer.domainInfoRepository.GetDomainInfoById(ctx, email.DomainId)
	if err != nil {
		return nil, err
	}
	htmlProcessed := template.Html
	textProcessed := template.Text
	subjectProcessed := template.Subject
	for k, v := range email.GetContext() {
		fKey := fmt.Sprintf("[[%s]]", k)
		htmlProcessed = strings.ReplaceAll(htmlProcessed, fKey, v)
		textProcessed = strings.ReplaceAll(textProcessed, fKey, v)
		subjectProcessed = strings.ReplaceAll(subjectProcessed, fKey, v)
	}
	return &PreparedEmail{
		Html:    htmlProcessed,
		Text:    textProcessed,
		Subject: subjectProcessed,
		To:      email.To,
		From:    email.From,
		ReplyTo: email.ReplyTo,
		BCC:     email.BCC,
		Id:      email.Id,
		Domain: struct {
			Name string
			Id   int64 "json:\",string\""
		}{
			Name: domain.DomainName,
			Id:   email.DomainId,
		},
	}, nil
}

func (dEmailer *DefaultEmailStoreService) CreateEmail(ctx context.Context, createEmail *email.CreateNewEmailRequest) error {
	if createEmail.TemplateId == 0 {
		return email.ErrTemplateNotFound("none")
	}
	_, err := dEmailer.templateRepository.GetById(ctx, createEmail.TemplateId)
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

	domain, err := dEmailer.domainInfoRepository.GetDomainInfoByName(ctx, domainName)
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
	bccEmails, err := getBCC(createEmail)
	if err != nil {
		return err
	}
	contextInfo := createEmail.Context
	if contextInfo == nil {
		contextInfo = make(map[string]string)
	}
	emailEntity, err := email.NewEmail(
		dEmailer.snowFlakeNode,
		domain.Id,
		createEmail.TemplateId,
		createEmail.From,
		string(toEmails),
		replyToEmail,
		string(bccEmails),
		contextInfo,
	)
	if err != nil {
		return err
	}
	err = dEmailer.emailRepository.Save(ctx, emailEntity)
	if err != nil {
		return err
	}
	dEmailer.eventBus.Push(ctx, events.CreateNewEmailEvent(emailEntity.Id))
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

func getBCC(createEmail *email.CreateNewEmailRequest) ([]byte, error) {
	if createEmail.BCC == nil {
		return []byte("[]"), nil
	}
	for _, bcc := range createEmail.BCC {
		_, err := mail.ParseAddress(bcc)
		if err != nil {
			return nil, err
		}
	}
	bccEmails, err := json.Marshal(createEmail.BCC)
	if err != nil {
		return nil, err
	}
	return bccEmails, nil
}
