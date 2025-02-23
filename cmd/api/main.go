package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sgatu/ezmail/cmd/api/server"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func main() {
	setupLog()
	server := server.NewServer()
	sqldb, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	defer db.Close()
	appContext, cleanup := internal_http.SetupAppContext(db)
	defer cleanup()
	handlers.SetupMiddlewares(server, appContext)
	slog.Debug("Setting up routes", "Source", "Api-Main")
	handlers.SetupRoutes(server, appContext)
	slog.Info(fmt.Sprintf("Server listening on :%s", os.Getenv("PORT")), "Source", "Api-Main")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), server)
}

func setupLog() {
	level := slog.Level(1)
	err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	if err != nil {
		return
	}
	handler := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(handler)
}
