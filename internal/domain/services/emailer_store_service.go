package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	scheduledEvRepo      events.ScheduledEventRepository
	snowFlakeNode        *snowflake.Node
}

func NewDefaultEmailStoreService(
	emailRepository email.EmailRepository,
	templateRepository email.TemplateRepository,
	domainInfoRepository domain.DomainInfoRepository,
	eventBus events.EventBus,
	snowFlakeNode *snowflake.Node,
	scheduledEvRepo events.ScheduledEventRepository,
) EmailStoreService {
	return &DefaultEmailStoreService{
		emailRepository:      emailRepository,
		templateRepository:   templateRepository,
		domainInfoRepository: domainInfoRepository,
		eventBus:             eventBus,
		snowFlakeNode:        snowFlakeNode,
		scheduledEvRepo:      scheduledEvRepo,
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
	addrParts := strings.SplitN(createEmail.From.ParsedEmail.Address, "@", 2)
	domainName := addrParts[1]

	if createEmail.To == nil || len(createEmail.To) == 0 {
		return fmt.Errorf("no email destinataries")
	}
	toEmailsS := make([]string, 0, len(createEmail.To))
	for _, to := range createEmail.To {
		toEmailsS = append(toEmailsS, to.StringEmail)
	}

	domain, err := dEmailer.domainInfoRepository.GetDomainInfoByName(ctx, domainName)
	if err != nil {
		return err
	}

	if !domain.Validated {
		return fmt.Errorf("domain not validated")
	}
	if createEmail.When != nil && dEmailer.scheduledEvRepo == nil {
		return fmt.Errorf("scheduled email without scheduling configuration")
	}

	toEmails, err := json.Marshal(toEmailsS)
	if err != nil {
		return err
	}
	replyToEmail := ""
	if createEmail.ReplyTo != nil {
		replyToEmail = createEmail.ReplyTo.StringEmail
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
		createEmail.From.StringEmail,
		string(toEmails),
		replyToEmail,
		string(bccEmails),
		contextInfo,
		createEmail.When,
	)
	if err != nil {
		return err
	}
	err = dEmailer.emailRepository.Save(ctx, emailEntity)
	if err != nil {
		return err
	}
	evt := events.CreateNewEmailEvent(emailEntity.Id)
	if createEmail.When != nil {
		dEmailer.scheduledEvRepo.Push(ctx, time.Time(*createEmail.When), evt)
	} else {
		dEmailer.eventBus.Push(ctx, evt)
	}
	return nil
}

func getBCC(createEmail *email.CreateNewEmailRequest) ([]byte, error) {
	if createEmail.BCC == nil || len(createEmail.BCC) == 0 {
		return []byte("[]"), nil
	}
	bccS := make([]string, 0, len(createEmail.BCC))
	for _, bcc := range createEmail.BCC {
		bccS = append(bccS, bcc.StringEmail)
	}
	bccEmails, err := json.Marshal(bccS)
	if err != nil {
		return nil, err
	}
	return bccEmails, nil
}
