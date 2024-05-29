package register

import (
	"encoding/json"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

func RegisterHandler(context *internal_http.AppContext, router *chi.Mux) {
	rh := registerHandler{
		userRepo:      context.UserRepository,
		snowflakeNode: context.SnowflakeNode,
	}
	router.Post("/register", rh.register)
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerHandler struct {
	userRepo      user.UserRepository
	snowflakeNode *snowflake.Node
}

func (rh *registerHandler) register(w http.ResponseWriter, r *http.Request) {
	var registerReq registerRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		common.ReturnError(common.InvalidRequestBodyError(), w)
		return
	}
	usr, err := user.CreateNewUser(rh.snowflakeNode, &user.BcryptPasswordHasher{}, registerReq.Email, registerReq.Password, registerReq.Name)
	if err != nil {
		common.ReturnError(err, w)
		return
	}
	err = rh.userRepo.Save(r.Context(), usr)
	if err != nil {
		common.ReturnError(err, w)
		return
	}
	common.EntityCreated(usr.Id, "user", w)
}
