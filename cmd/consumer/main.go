package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/sgatu/ezmail/internal/worker"
	"github.com/sgatu/ezmail/internal/worker/processors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	sqldb, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())
	defer db.Close()

	runningContext, err := processors.SetupRunningContext()
	if err != nil {
		panic(err)
	}
	e := worker.NewExecutor(
		runningContext,
		processors.InitNewEmailProcessor(),
		processors.InitRescheduledEmailProcessor(),
	)
	e.Run()
}
