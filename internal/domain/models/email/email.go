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

type Email struct {
	bun.BaseModel `bun:"table:email,alias:em"`
	Created       time.Time `bun:",notnull"`
	context       map[string]string
	From          string `bun:",notnull"`
	ReplyTo       string `bun:""`
	To            string `bun:",notnull"`
	CCO           string `bun:"cco"`
	ContextRaw    string `bun:"context,notnull"`
	Processed     bool   `bun:",notnull"`
	TemplateId    int64  `bun:",notnull"`
	DomainId      int64  `bun:",notnull"`
	Id            int64  `bun:",pk"`
}

type CreateNewEmailRequest struct {
	Context    map[string]string
	ReplyTo    *string
	CCO        []string
	From       string
	To         []string
	TemplateId int64 `json:",string"`
}
type CreateTemplateRequest struct {
	Subject string
	Text    string
	Html    string
}

func (em *Email) getContext() map[string]string {
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
	cco string,
	context map[string]string,
) (*Email, error) {
	em := &Email{
		Created:    time.Now().UTC(),
		From:       from,
		To:         to,
		CCO:        cco,
		ReplyTo:    replyTo,
		context:    context,
		TemplateId: templateId,
		DomainId:   domainId,
		Id:         sNode.Generate().Int64(),
		Processed:  false,
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
	Save(ctx context.Context, tpl *Template) error
}
type EmailRepository interface {
	GetById(ctx context.Context, id int64) (*Email, error)
	Save(ctx context.Context, email *Email) error
}
