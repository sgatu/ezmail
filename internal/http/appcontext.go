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
	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/email"
	"github.com/sgatu/ezmail/internal/domain/models/events"
	"github.com/sgatu/ezmail/internal/domain/services"
	"github.com/sgatu/ezmail/internal/infrastructure/eventbus"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories/mysql"
	infra_services "github.com/sgatu/ezmail/internal/infrastructure/services"
	"github.com/uptrace/bun"
)

type AppContext struct {
	DomainInfoRepository domain.DomainInfoRepository
	EmailStoreService    services.EmailStoreService
	IdentityManager      services.IdentityManager
	SnowflakeNode        *snowflake.Node
	TemplateRepository   email.TemplateRepository
	CommonEventsBus      events.EventBus
}

func (ac *AppContext) ValidateToken(ctx context.Context, token string) error {
	if token != os.Getenv("AUTH_TOKEN") {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func SetupAppContext(db *bun.DB) *AppContext {
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

	mainBusRedis := os.Getenv("COMMON_BUS_REDIS")
	if mainBusRedis == "" {
		mainBusRedis = "localhost:6379"
	}
	redisCli := redis.NewClient(&redis.Options{
		Addr:     mainBusRedis,
		Password: "",
		DB:       0,
	})
	commonEventBus := eventbus.NewCommonEventsEventBus(redisCli)

	emailService := services.NewDefaultEmailStoreService(emailRepository, templateRepository, domainInfoRepository, commonEventBus, snowflakeNode)
	return &AppContext{
		DomainInfoRepository: domainInfoRepository,
		IdentityManager:      infra_services.NewSESIdentityManager(domainInfoRepository, awsConfig, snowflakeNode),
		SnowflakeNode:        snowflakeNode,
		EmailStoreService:    emailService,
		TemplateRepository:   templateRepository,
		CommonEventsBus:      commonEventBus,
	}
}
