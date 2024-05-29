package login

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	ihttp "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

func LoginHandler(appContext *ihttp.AppContext, router *chi.Mux) {
	h := loginHandler{userRepo: appContext.UserRepository}
	router.Post("/login", h.login)
}

type loginHandler struct {
	userRepo user.UserRepository
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *loginHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		common.ReturnError(common.InvalidRequestBodyError(), w)
		return
	}
	usr, err := h.userRepo.FindByEmailAndPassword(r.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		if err == user.ErrUserNotFoundError {
			common.ReturnError(common.EntityNotFoundError("user"), w)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "Found user %s, id %s", usr.Name, usr.Id)
}
