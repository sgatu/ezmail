package domain

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/service/ses"
)

func DomainHandler(ctx *internal_http.AppContext, router chi.Router) {
	domHandler := &domainHandler{
		sesService:           ctx.SESService,
		domainInfoRepository: ctx.DomainInfoRepository,
		snowflakeNode:        ctx.SnowflakeNode,
	}
	router.Post("/domain", domHandler.createDomain)
	router.Get("/domain/{id}", domHandler.getDomain)
	router.Get("/domain", domHandler.getUserDomains)
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
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt string                 `json:"created_at"`
	Records   []domainRecordResponse `json:"records"`
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
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	domainInfo := &domain.DomainInfo{
		Id:         dh.snowflakeNode.Generate().String(),
		Created:    time.Now().UTC(),
		DomainName: createDomainReq.Name,
		Region:     createDomainReq.Region,
		UserId:     currUser.Id,
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

func (dh *domainHandler) getUserDomains(w http.ResponseWriter, r *http.Request) {
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	doms, err := dh.domainInfoRepository.GetAllByUserId(r.Context(), currUser.Id)
	if err != nil {
		if err == domain.ErrDomainInfoNotFound {
			common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		} else {
			common.ErrorResponse(common.InternalServerError(err), w)
		}
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
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	dom, err := dh.domainInfoRepository.GetDomainInfoById(r.Context(), domainId)
	if err != nil {
		if err == domain.ErrDomainInfoNotFound {
			common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		} else {
			common.ErrorResponse(common.InternalServerError(err), w)
		}
		return
	}
	if currUser.Id != dom.UserId {
		common.ErrorResponse(common.EntityNotFoundError("domain"), w)
		return
	}
	common.ReturnReponse(getCreateDomainResponse(dom), 200, w)
}
