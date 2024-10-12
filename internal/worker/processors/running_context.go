package processors

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
	"github.com/sgatu/ezmail/internal/infrastructure/eventbus"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories/mysql"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories/redis"
	infra_services "github.com/sgatu/ezmail/internal/infrastructure/services"
	"github.com/uptrace/bun"
)

type RescheduleConfig struct {
	Retries  int8
	RetrySec int32
}
type RunningContext struct {
	EventBus            events.EventBus
	EmailStoreService   services.EmailStoreService
	EmailerService      services.Emailer
	ScheduledEventsRepo events.ScheduledEventRepository
	BusReader           events.BusReader
	Rc                  *RescheduleConfig
}

func SetupRunningContext(db *bun.DB) (*RunningContext, func(), error) {
	mainBusRedis := os.Getenv("COMMON_BUS_REDIS")
	if mainBusRedis == "" {
		mainBusRedis = "localhost:6379"
	}
	redisCli := redis.NewClient(&redis.Options{
		Addr:                  mainBusRedis,
		Password:              "",
		DB:                    0,
		ContextTimeoutEnabled: true,
	})
	maxLenEventsStr := os.Getenv("REDIS_EVENTS_MAX_LEN")
	maxLenEvents, err := strconv.ParseInt(maxLenEventsStr, 10, 64)
	if err != nil {
		slog.Warn("Could not load redis max events, defaulting to 2500")
		maxLenEvents = 2500
	}
	eventsTopic := os.Getenv("EVENTS_TOPIC")
	if eventsTopic == "" {
		eventsTopic = "queue:email_events"
	}
	eventBus := eventbus.NewRedisEventBus(redisCli, maxLenEvents, eventsTopic)
	nodeIdStr := os.Getenv("NODE_ID")
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		nodeId = rand.Int63()
		slog.Info(fmt.Sprintf("No snowflake node id defined (missing env NODE_ID), generated random as %d", nodeId))
	}
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
	emailStoreRepo := mysql.NewMysqlEmailRepository(db)
	templateRepo := mysql.NewMysqlTemplateRepository(db)
	domainRepo := mysql.NewMysqlDomainInfoRepository(db)
	scheduledEvRepo := redis.NewRedisScheduledEventRepository(redisCli)
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return &RunningContext{
			EventBus:            eventBus,
			ScheduledEventsRepo: scheduledEvRepo,
			EmailStoreService: services.NewDefaultEmailStoreService(
				emailStoreRepo,
				templateRepo,
				domainRepo,
				eventBus,
				snowflakeNode,
				scheduledEvRepo,
			),
			EmailerService: infra_services.NewSesEmailer(sesv2.NewFromConfig(awsConfig), snowflakeNode),
			Rc:             nil,
			BusReader:      eventbus.NewRedisEventsReader(redisCli, eventsTopic, "mainReader", false),
		}, func() {
			redisCli.Close()
		}, nil
}
