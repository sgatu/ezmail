package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/common"
)

type templateHandler struct {
	templateRepository email.TemplateRepository
	snowflakeNode      *snowflake.Node
}

func RegisterTemplateHandler(appCtx *internal_http.AppContext, r chi.Router) {
	tHandler := &templateHandler{
		templateRepository: appCtx.TemplateRepository,
		snowflakeNode:      appCtx.SnowflakeNode,
	}
	common.RegisterEndpoint(r.Post, "/template", tHandler.CreateTemplate, "Create a new template")
	common.RegisterEndpoint(r.Get, "/template/{id}", tHandler.GetTemplate, "Get template")
	common.RegisterEndpoint(r.Get, "/template", tHandler.GetAllTemplates, "Get compact version of all templates")
}

type templateResponse struct {
	Html      string `json:"html"`
	Text      string `json:"text"`
	Subject   string `json:"subject"`
	CreatedAt string `json:"created_at"`
	Id        int64  `json:"id,string"`
}

func getTemplateResponse(tpl *email.Template) *templateResponse {
	tr := templateResponse{
		Id:        tpl.Id,
		Html:      tpl.Html,
		CreatedAt: tpl.Created.Format(time.RFC3339),
		Text:      tpl.Text,
		Subject:   tpl.Subject,
	}
	return &tr
}

func (th *templateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var createTemplateRequest email.CreateTemplateRequest
	err := json.NewDecoder(r.Body).Decode(&createTemplateRequest)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	template := email.NewTemplate(th.snowflakeNode, createTemplateRequest.Text, createTemplateRequest.Html, createTemplateRequest.Subject)
	err = th.templateRepository.Save(r.Context(), template)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.EntityCreated(strconv.FormatInt(template.Id, 10), "template", w)
}

func (th *templateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := th.templateRepository.GetAll(r.Context())
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkResponse(templates, w)
}

func (th *templateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := chi.URLParam(r, "id")
	if templateId == "" {
		common.ErrorResponse(common.EntityNotFoundError("template"), w)
		return
	}
	templateIdInt, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		common.ErrorResponse(common.EntityNotFoundError("template"), w)
		return
	}
	template, err := th.templateRepository.GetById(r.Context(), templateIdInt)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkResponse(getTemplateResponse(template), w)
}
