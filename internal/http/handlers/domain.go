package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/common"
	"github.com/sgatu/ezmail/internal/service/ses"
)

func RegisterDomainHandler(ctx *internal_http.AppContext, router chi.Router) {
	domHandler := &domainHandler{
		sesService:           ctx.SESService,
		domainInfoRepository: ctx.DomainInfoRepository,
		snowflakeNode:        ctx.SnowflakeNode,
	}
	common.RegisterEndpoint(router.Post, "/domain", domHandler.createDomain, "Register new domain in the system")
	common.RegisterEndpoint(router.Get, "/domain/{id}", domHandler.getDomain, "Get a domain identified by {id}")
	common.RegisterEndpoint(router.Get, "/domain", domHandler.getDomains, "Get all domains")
	common.RegisterEndpoint(router.Delete, "/domain/{id}", domHandler.deleteDomain, "Soft deletes a domain (mark it as deleted in db and does not delete it from aws)")
}

type domainHandler struct {
	sesService           *ses.SESService
	domainInfoRepository domain.DomainInfoRepository
	snowflakeNode        *snowflake.Node
}

type createDomainRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type domainRecordResponse struct {
	Purpose    string `json:"purpose"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	RecordType string `json:"record_type"`
	Status     string `json:"status"`
}
type domainResponse struct {
	Name      string                 `json:"name"`
	CreatedAt string                 `json:"created_at"`
	Records   []domainRecordResponse `json:"records"`
	Id        int64                  `json:"id,string"`
	Validated bool                   `json:"validated"`
}

func getCreateDomainResponse(di *domain.DomainInfo) *domainResponse {
	records, err := di.GetDnsRecords()
	if err != nil {
		return nil
	}
	cdr := domainResponse{
		Id:        di.Id,
		Name:      di.DomainName,
		CreatedAt: di.Created.Format(time.RFC3339),
		Validated: di.Validated,
		Records:   make([]domainRecordResponse, 0, len(records)),
	}
	for _, rec := range records {
		sts := "PENDING"
		switch rec.Status {
		case domain.DNS_RECORD_STATUS_FAILED:
			sts = "FAILED"
		case domain.DNS_RECORD_STATUS_VERIFIED:
			sts = "VERIFIED"
		}
		cdr.Records = append(cdr.Records, domainRecordResponse{
			Purpose:    rec.Purpose,
			Name:       rec.Name,
			Value:      rec.Value,
			RecordType: rec.Type,
			Status:     sts,
		})
	}
	return &cdr
}

func (dh *domainHandler) createDomain(w http.ResponseWriter, r *http.Request) {
	var createDomainReq createDomainRequest
	err := json.NewDecoder(r.Body).Decode(&createDomainReq)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	domainInfo := &domain.DomainInfo{
		Id:         dh.snowflakeNode.Generate().Int64(),
		Created:    time.Now().UTC(),
		DomainName: createDomainReq.Name,
		Region:     createDomainReq.Region,
	}
	err = dh.sesService.CreateDomain(r.Context(), domainInfo)
	if err != nil {
		common.ReturnReponse(err.Error(), 500, w)
		return
	}
	err = dh.domainInfoRepository.Save(r.Context(), domainInfo)
	if err != nil {
		dh.sesService.DeleteIdentity(r.Context(), domainInfo)
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	common.ReturnReponse(getCreateDomainResponse(domainInfo), 201, w)
}

func (dh *domainHandler) getDomains(w http.ResponseWriter, r *http.Request) {
	doms, err := dh.domainInfoRepository.GetAll(r.Context())
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	domsResp := make([]domainResponse, 0, len(doms))
	for _, dom := range doms {
		domsResp = append(domsResp, *getCreateDomainResponse(&dom))
	}
	common.ReturnReponse(domsResp, 200, w)
}

func (dh *domainHandler) getDomain(w http.ResponseWriter, r *http.Request) {
	domainId := chi.URLParam(r, "id")
	if domainId == "" {
		common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		return
	}
	domainIdInt, err := strconv.ParseInt(domainId, 10, 64)
	if err != nil {
		common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		return
	}
	dom, err := dh.domainInfoRepository.GetDomainInfoById(r.Context(), domainIdInt)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.ReturnReponse(getCreateDomainResponse(dom), 200, w)
}

func (dh *domainHandler) deleteDomain(w http.ResponseWriter, r *http.Request) {
	domainId := chi.URLParam(r, "id")
	if domainId == "" {
		common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		return
	}
	domainIdInt, err := strconv.ParseInt(domainId, 10, 64)
	if err != nil {
		common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		return
	}
	_, err = dh.domainInfoRepository.GetDomainInfoById(r.Context(), domainIdInt)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	err = dh.domainInfoRepository.DeleteDomain(r.Context(), domainIdInt)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkOperation(w)
}
