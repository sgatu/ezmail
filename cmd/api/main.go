package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sgatu/ezmail/cmd/api/server"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	server := server.NewServer()
	sqldb, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	defer db.Close()
	appContext := internal_http.SetupAppContext(db)
	handlers.SetupMiddlewares(server, appContext)
	handlers.SetupRoutes(server, appContext)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), server)
}
