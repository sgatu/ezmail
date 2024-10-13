package main

import (
	"database/sql"
	"os"
	"os/signal"
	"sync"
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
	wg := &sync.WaitGroup{}
	e := worker.NewExecutor(
		runningContext,
		wg,
		processors.InitNewEmailProcessor(),
		processors.InitRescheduledEmailProcessor(),
	)
	var s *worker.Scheduler = nil
	if runningContext.ScheduledEventsRepo != nil {
		s := worker.NewScheduler(runningContext.ScheduledEventsRepo, runningContext.EventBus, wg)
		s.Run()
	}
	e.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		e.Stop()
		if s != nil {
			s.Stop()
		}
	}()
	wg.Wait()
}
