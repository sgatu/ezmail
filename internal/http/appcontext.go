package http

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
	"github.com/sgatu/ezmail/internal/infrastructure/eventbus"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories/mysql"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories/redis"
	infra_services "github.com/sgatu/ezmail/internal/infrastructure/services"
	"github.com/uptrace/bun"
)

type AppContext struct {
	DomainInfoRepository domain.DomainInfoRepository
	EmailStoreService    services.EmailStoreService
	IdentityManager      services.IdentityManager
	SnowflakeNode        *snowflake.Node
	TemplateRepository   email.TemplateRepository
	EventsBus            events.EventBus
}

func (ac *AppContext) ValidateToken(ctx context.Context, token string) error {
	if token != os.Getenv("AUTH_TOKEN") {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func SetupAppContext(db *bun.DB) (*AppContext, func()) {
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
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	domainInfoRepository := mysql.NewMysqlDomainInfoRepository(db)
	emailRepository := mysql.NewMysqlEmailRepository(db)
	templateRepository := mysql.NewMysqlTemplateRepository(db)

	mainBusRedis := os.Getenv("REDIS")
	if mainBusRedis == "" {
		mainBusRedis = "localhost:6379"
	}
	redisCli := redis.NewClient(&redis.Options{
		Addr:     mainBusRedis,
		Password: "",
		DB:       0,
	})
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
	redisEventBus := eventbus.NewRedisEventBus(redisCli, maxLenEvents, eventsTopic)
	scheduleKey := os.Getenv("SCHEDULING_KEY")
	var scheduledEvRepo *redis.RedisScheduledEventRepository = nil
	if scheduleKey != "" {
		scheduledEvRepo = redis.NewRedisScheduledEventRepository(redisCli, scheduleKey)
	}
	emailService := services.NewDefaultEmailStoreService(
		emailRepository,
		templateRepository,
		domainInfoRepository,
		redisEventBus,
		snowflakeNode,
		scheduledEvRepo,
	)
	return &AppContext{
			DomainInfoRepository: domainInfoRepository,
			IdentityManager:      infra_services.NewSESIdentityManager(domainInfoRepository, awsConfig, snowflakeNode),
			SnowflakeNode:        snowflakeNode,
			EmailStoreService:    emailService,
			TemplateRepository:   templateRepository,
			EventsBus:            redisEventBus,
		}, func() {
			redisCli.Close()
		}
}
