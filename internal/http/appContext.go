package http

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
	"github.com/sgatu/ezmail/internal/service/ses"
	"github.com/uptrace/bun"
)

type AppContext struct {
	DomainInfoRepository domain.DomainInfoRepository
	SESService           *ses.SESService
	SnowflakeNode        *snowflake.Node
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
		fmt.Printf("No snowflake node id defined (missing env NODE_ID), generated random as %d\n", nodeId)
	}
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	domainInfoRepository := repositories.NewMysqlDomainInfoRepository(db)
	return &AppContext{
		DomainInfoRepository: domainInfoRepository,
		SESService:           ses.NewSESService(domainInfoRepository, awsConfig, snowflakeNode),
		SnowflakeNode:        snowflakeNode,
	}
}
