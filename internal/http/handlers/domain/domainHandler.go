package domain

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/service/ses"
)

func DomainHandler(ctx *internal_http.AppContext, router chi.Router) {
	domHandler := &domainHandler{
		sesService: ctx.SESService,
	}
	router.Post("/domain", domHandler.createDomain)
}

type domainHandler struct {
	sesService *ses.SESService
}

type createDomainRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (dh *domainHandler) createDomain(w http.ResponseWriter, r *http.Request) {
	var createDomainReq createDomainRequest
	err := json.NewDecoder(r.Body).Decode(&createDomainReq)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	dom, err := dh.sesService.CreateDomain(r.Context(), currUser, createDomainReq.Name, createDomainReq.Region)
	if err != nil {
		common.ReturnReponse(err.Error(), 500, w)
		return
	}
	fmt.Printf("%+v\n", dom)
	common.ReturnReponse("", 201, w)
}
