package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sgatu/ezmail/cmd/api/server"
	"github.com/sgatu/ezmail/internal/http/handlers"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
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
	userRepo := repositories.NewMysqlUserRepository(db)
	handlers.SetupRoutes(server, userRepo)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), server)
}
