package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/services"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/common"
)

type emailHandler struct {
	emailService services.EmailStoreService
}
type emailResponse struct {
	Created    time.Time         `json:"created"`
	Context    map[string]string `json:"context"`
	From       string            `json:"from"`
	ReplyTo    string            `json:"reply_to"`
	To         string            `json:"to"`
	BCC        string            `json:"bcc"`
	Processed  bool              `json:"processed"`
	TemplateId int64             `json:"template_id,string"`
	DomainId   int64             `json:"domain_id,string"`
	Id         int64             `json:"id,string"`
}

func createEmailResponse(email *email.Email) emailResponse {
	return emailResponse{
		Created:    email.Created,
		Context:    email.GetContext(),
		From:       email.From,
		ReplyTo:    email.ReplyTo,
		To:         email.To,
		BCC:        email.BCC,
		Processed:  email.Processed,
		TemplateId: email.TemplateId,
		DomainId:   email.DomainId,
		Id:         email.Id,
	}
}

func RegisterEmailHandler(appCtx *internal_http.AppContext, r chi.Router) {
	eHandler := &emailHandler{
		emailService: appCtx.EmailStoreService,
	}
	common.RegisterEndpoint(r.Post, "/email", eHandler.SendEmail, "Send an email")
	common.RegisterEndpoint(r.Get, "/email/{id}", eHandler.GetById, "Get email by id")
}

func (eh *emailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var createEmailRequest email.CreateNewEmailRequest
	err := json.NewDecoder(r.Body).Decode(&createEmailRequest)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	err = eh.emailService.CreateEmail(r.Context(), &createEmailRequest)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkOperation(w)
}

func (eh *emailHandler) GetById(w http.ResponseWriter, r *http.Request) {
	emailId := chi.URLParam(r, "id")
	if emailId == "" {
		common.ErrorResponse(common.EntityNotFoundError("email"), w)
		return
	}
	emailIdInt, err := strconv.ParseInt(emailId, 10, 64)
	if err != nil {
		common.ErrorResponse(common.EntityNotFoundError("email"), w)
		return
	}
	email, err := eh.emailService.GetById(r.Context(), emailIdInt)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkResponse(createEmailResponse(email), w)
}
