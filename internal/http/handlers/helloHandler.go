package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sgatu/ezmail/internal/domain/models/user"
)

type helloHandler struct {
	userRepo user.UserRepository
}

func (h *helloHandler) hello(w http.ResponseWriter, r *http.Request) {
	user, err := h.userRepo.GetById(r.Context(), "bca")
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	fmt.Fprintf(w, "Found user %s, id %s", user.Name, user.Id)
}
