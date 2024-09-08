package template

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
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
