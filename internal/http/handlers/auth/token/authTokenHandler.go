package token

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

type authTokenHandler struct {
	authTokenRepository auth.AuthTokenRepository
	snowflakeNode       *snowflake.Node
}

func RegisterAuthToken(appContext *internal_http.AppContext, router chi.Router) {
	handler := authTokenHandler{
		authTokenRepository: appContext.AuthTokenRepository,
		snowflakeNode:       appContext.SnowflakeNode,
	}
	router.Get("/tokens", handler.getAllTokens)
	router.Post("/tokens", handler.createNewToken)
	router.Delete("/tokens/{id}", handler.disableToken)
}

type createTokenRequest struct {
	Expire     *string                   `json:"expire"`
	AccessType *auth.AuthTokenAccessType `json:"access_type"`
	Name       string                    `json:"name"`
}

type tokenResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Token      string `json:"token"`
	Expire     string `json:"expire"`
	Created    string `json:"created"`
	AccessType string `json:"access_type"`
}

func tokenToTokenResponse(tok *auth.AuthToken) tokenResponse {
	expire := ""
	if !tok.Expire.IsZero() {
		expire = tok.Expire.Format(time.RFC3339)
	}
	accessType := "FULL_ACCESS"
	if tok.AccessType == auth.AUTH_TYPE_SEND_ONLY {
		accessType = "SEND_ONLY"
	}
	return tokenResponse{
		Id:         tok.Id,
		Name:       tok.Name,
		Expire:     expire,
		AccessType: accessType,
		Token:      tok.Token,
		Created:    tok.Created.Format(time.RFC3339),
	}
}

func (ath *authTokenHandler) createNewToken(w http.ResponseWriter, r *http.Request) {
	var createReq createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	var expire *time.Time = nil
	if createReq.Expire != nil {
		parsedExp, err := time.Parse(time.RFC3339, *createReq.Expire)
		if err != nil {
			common.ErrorResponse(common.InvalidRequestBodyError(), w)
			return
		}
		expire = &parsedExp
	}
	var accessType auth.AuthTokenAccessType = auth.AUTH_TYPE_FULL_ACCESS
	if createReq.AccessType != nil {
		accessType = *createReq.AccessType
	}
	authToken, err := auth.CreateAuthToken(ath.snowflakeNode, currUser.Id, expire, createReq.Name, accessType)
	if err != nil {
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	err = ath.authTokenRepository.Save(r.Context(), authToken)
	if err != nil {
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	common.ReturnReponse(tokenToTokenResponse(authToken), http.StatusCreated, w)
}

func (ath *authTokenHandler) disableToken(w http.ResponseWriter, r *http.Request) {
	tokenId := chi.URLParam(r, "id")
	if tokenId == "" {
		common.ErrorResponse(common.InvalidRequest(), w)
		return
	}
	tok, err := ath.authTokenRepository.GetAuthTokenById(r.Context(), tokenId)
	if err != nil {
		if err == auth.ErrNoAuthTokenFound {
			common.ErrorResponse(common.EntityNotFoundError("token"), w)
		} else {
			common.ErrorResponse(err, w)
		}
		return
	}
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	if tok.UserId != currUser.Id {
		common.ErrorResponse(common.EntityNotFoundError("token"), w)
		return
	}
	tok.DisableToken()
	err = ath.authTokenRepository.Save(r.Context(), tok)
	if err != nil {
		common.ErrorResponse(err, w)
		return
	}
	common.OkOperation(w)
}

func (ath *authTokenHandler) getAllTokens(w http.ResponseWriter, r *http.Request) {
	currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
	toks, err := ath.authTokenRepository.GetAuthTokensByUserId(r.Context(), currUser.Id)
	if err != nil {
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	toksResponse := make([]tokenResponse, 0, len(toks))
	for _, t := range toks {
		toksResponse = append(toksResponse, tokenToTokenResponse(&t))
	}
	common.OkResponse(toksResponse, w)
}
