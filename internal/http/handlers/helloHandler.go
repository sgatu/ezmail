package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type helloHandler struct{}

func (h *helloHandler) hello(w http.ResponseWriter, r *http.Request) {
	sqldb, err := sql.Open("mysql", "")
	if err != nil {
		w.Write([]byte("Error sql. " + err.Error()))
		return
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	defer db.Close()
	repo := repositories.NewMysqlUserRepository(db)
	user, err := repo.GetById(r.Context(), "bca")
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
	fmt.Printf("%+v\n", user)
	w.Write([]byte(fmt.Sprintf("Found user %s, id %s", user.Name, user.Id)))
}
