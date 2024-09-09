package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/services"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/common"
)

type emailHandler struct {
	emailService *services.EmailService
}

func RegisterEmailHandler(appCtx *internal_http.AppContext, r chi.Router) {
	eHandler := &emailHandler{
		emailService: appCtx.EmailService,
	}
	common.RegisterEndpoint(r.Post, "/email", eHandler.SendEmail, "Send an email")
}

func (eh *emailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var createEmailRequest email.CreateNewEmailRequest
	err := json.NewDecoder(r.Body).Decode(&createEmailRequest)
	if err != nil {
		fmt.Printf("Parsing error: %s\n", err)
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	err = eh.emailService.SendEmail(r.Context(), &createEmailRequest)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkOperation(w)
}
