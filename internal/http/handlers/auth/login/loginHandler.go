package login

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	ihttp "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

func LoginHandler(appContext *ihttp.AppContext, router *chi.Mux) {
	h := loginHandler{userRepo: appContext.UserRepository, sessionRepo: appContext.SessionRepository, snowflakeNode: appContext.SnowflakeNode}
	router.Post("/login", h.login)
}

type loginHandler struct {
	userRepo      user.UserRepository
	sessionRepo   auth.SessionRepository
	snowflakeNode *snowflake.Node
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginOkResponse struct {
	SessionId string `json:"session_token"`
	Expire    string `json:"expire"`
	Created   string `json:"created"`
}

func (h *loginHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		common.ErrorResponse(common.InvalidRequestBodyError(), w)
		return
	}
	usr, err := h.userRepo.FindByEmailAndPassword(r.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		if err == user.ErrUserNotFoundError {
			common.ErrorResponse(common.UnauthorizedError(), w)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	session, err := auth.CreateSession(h.snowflakeNode, usr.Id, auth.GetDefaultSessionExpire())
	if err != nil {
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	err = h.sessionRepo.Save(r.Context(), session)
	if err != nil {
		common.ErrorResponse(common.InternalServerError(err), w)
		return
	}
	common.OkResponse(loginOkResponse{
		SessionId: session.SessionId,
		Expire:    session.Expire.Format(time.RFC3339),
		Created:   session.Created.Format(time.RFC3339),
	}, w)
}
