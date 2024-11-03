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
	"github.com/sgatu/ezmail/internal/thirdparty"
	"github.com/uptrace/bun"
)

type RescheduleConfig struct {
	Retries     int8
	RetryTimeMs int64
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
	mainBusRedis := os.Getenv("REDIS")
	if mainBusRedis == "" {
		mainBusRedis = "localhost:6379"
	}
	redisCli := thirdparty.RedisClient{Client: redis.NewClient(&redis.Options{
		Addr:                  mainBusRedis,
		Password:              "",
		DB:                    0,
		ContextTimeoutEnabled: true,
	})}
	maxLenEventsStr := os.Getenv("REDIS_EVENTS_MAX_LEN")
	maxLenEvents, err := strconv.ParseInt(maxLenEventsStr, 10, 64)
	if err != nil {
		slog.Warn("Could not load redis max events, defaulting to 2500")
		maxLenEvents = 2500
	}
	eventsTopic := os.Getenv("EVENTS_TOPIC")
	if eventsTopic == "" {
		eventsTopic = "topic:email.events"
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
	scheduleKey := os.Getenv("SCHEDULING_KEY")
	var scheduledEvRepo *redis.RedisScheduledEventRepository = nil
	if scheduleKey != "" {
		scheduledEvRepo = redis.NewRedisScheduledEventRepository(redisCli, scheduleKey)
	}
	rescheduleRetries := os.Getenv("RESCHEDULE_RETRIES")
	rescheduleTimeMs := os.Getenv("RESCHEDULE_TIME_MS")
	var rescheduleConfig *RescheduleConfig = nil
	if rescheduleRetries != "" && rescheduleTimeMs != "" {
		rsRetries, errRetries := strconv.ParseInt(rescheduleRetries, 10, 8)
		rsTime, errTime := strconv.ParseInt(rescheduleTimeMs, 10, 64)
		if errRetries == nil && errTime == nil {
			rescheduleConfig = &RescheduleConfig{
				Retries:     int8(rsRetries),
				RetryTimeMs: rsTime,
			}
		} else {
			slog.Warn("Invalid reschedule configuration")
		}

	}
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	sesClient := sesv2.NewFromConfig(awsConfig)
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
			EmailerService: infra_services.NewSesEmailer(
				&thirdparty.AWSSesV2Client{Client: sesClient},
				snowflakeNode,
			),
			Rc:        rescheduleConfig,
			BusReader: eventbus.NewRedisEventsReader(redisCli, eventsTopic, "mainReader", false),
		}, func() {
			redisCli.Close()
		}, nil
}
