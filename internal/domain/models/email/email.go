package email

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models"
	"github.com/uptrace/bun"
)

type Template struct {
	bun.BaseModel `bun:"table:template,alias:tpl"`
	Created       time.Time `bun:",notnull"`
	Subject       string    `bun:",notnull"`
	Html          string    `bun:",notnull"`
	Text          string    `bun:",notnull"`
	Id            int64     `bun:",pk"`
}

type CompactTemplate struct {
	bun.BaseModel `bun:"table:template,alias:tpl"`
	Created       time.Time `bun:",notnull"`
	Subject       string    `bun:",notnull"`
	Id            int64     `bun:",pk" json:",string"`
}

type Email struct {
	bun.BaseModel `bun:"table:email,alias:em"`
	Created       time.Time `bun:",notnull"`
	ExpectedSent  time.Time `bun:"expected_sent,notnull"`
	context       map[string]string
	ProcessedDate bun.NullTime `bun:",nullzero"`
	From          string       `bun:",notnull"`
	ReplyTo       string       `bun:"reply_to"`
	To            string       `bun:",notnull"`
	BCC           string       `bun:"bcc"`
	ContextRaw    string       `bun:"context,notnull"`
	TemplateId    int64        `bun:",notnull"`
	DomainId      int64        `bun:",notnull"`
	Id            int64        `bun:",pk"`
}

type CreateNewEmailRequest struct {
	Context    map[string]string     `json:"context"`
	ReplyTo    *models.EmailAddress  `json:"reply_to"`
	When       *models.DateTime      `json:"when"`
	BCC        []models.EmailAddress `json:"bcc"`
	From       models.EmailAddress   `json:"from"`
	To         []models.EmailAddress `json:"to"`
	TemplateId int64                 `json:"template_id,string"`
}

type CreateTemplateRequest struct {
	Subject string
	Text    string
	Html    string
}

func (em *Email) GetContext() map[string]string {
	if em.context == nil {
		json.Unmarshal([]byte(em.ContextRaw), &em.context)
	}
	return em.context
}

func NewEmail(
	sNode *snowflake.Node,
	domainId int64,
	templateId int64,
	from string,
	to string,
	replyTo string,
	bcc string,
	context map[string]string,
	expectedSent *models.DateTime,
) (*Email, error) {
	em := &Email{
		Created:    time.Now().UTC(),
		From:       from,
		To:         to,
		BCC:        bcc,
		ReplyTo:    replyTo,
		context:    context,
		TemplateId: templateId,
		DomainId:   domainId,
		Id:         sNode.Generate().Int64(),
	}
	em.ExpectedSent = time.Now().UTC()
	if expectedSent != nil {
		em.ExpectedSent = time.Time(*expectedSent)
	}
	marshalResult, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}
	em.ContextRaw = string(marshalResult)
	return em, nil
}

func NewTemplate(
	sNode *snowflake.Node,
	text string,
	html string,
	subject string,
) *Template {
	return &Template{
		Id:      sNode.Generate().Int64(),
		Text:    text,
		Subject: subject,
		Html:    html,
		Created: time.Now().UTC(),
	}
}

func ErrTemplateNotFound(identifier string) error {
	return models.NewMissingEntityError("template not found", identifier)
}

func ErrEmailNotFound(identifer string) error {
	return models.NewMissingEntityError("email not found", identifer)
}

type TemplateRepository interface {
	GetById(ctx context.Context, id int64) (*Template, error)
	GetAll(ctx context.Context) ([]CompactTemplate, error)
	Save(ctx context.Context, tpl *Template) error
}
type EmailRepository interface {
	GetById(ctx context.Context, id int64) (*Email, error)
	Save(ctx context.Context, email *Email) error
}
