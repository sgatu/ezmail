package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	fmt.Println("### [ SETUP ROUTES ] ###")
	handlers.SetupRoutes(server, appContext)
	fmt.Printf("Server listening on :%s\n", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), server)
}
