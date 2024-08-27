package email

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/service/ses"
)

type emailHandler struct {
	domainRepo domain.DomainInfoRepository
	sesService *ses.SESService
}

func RegisterEmailHandler(appCtx *internal_http.AppContext, r chi.Router) {
	eHandler := &emailHandler{
		domainRepo: appCtx.DomainInfoRepository,
		sesService: appCtx.SESService,
	}
	common.RegisterEndpoint(r.Post, "/email", eHandler.SendEmail, "Send an email")
}

type sendEmailRequest struct {
	From     string   `json:"from"`
	Subject  string   `json:"subject"`
	TextBody *string  `json:"text"`
	HtmlBody *string  `json:"html"`
	ReplyTo  *string  `json:"reply_to"`
	To       []string `json:"to"`
	Cco      []string `json:"cco"`
}

func (eh *emailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
}
