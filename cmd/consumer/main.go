package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

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

	runningContext, cleanup, err := processors.SetupRunningContext(db)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	e, wg := worker.NewExecutor(
		runningContext,
		processors.InitNewEmailProcessor(),
		processors.InitRescheduledEmailProcessor(),
	)
	e.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		e.Stop()
	}()
	wg.Wait()
}
