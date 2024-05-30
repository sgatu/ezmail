package http

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
	"github.com/uptrace/bun"
)

type AppContext struct {
	UserRepository      user.UserRepository
	AuthTokenRepository auth.AuthTokenRepository
	SnowflakeNode       *snowflake.Node
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
	return &AppContext{
		UserRepository:      repositories.NewMysqlUserRepository(db),
		AuthTokenRepository: repositories.NewMysqlAuthTokenRepository(db),
		SnowflakeNode:       snowflakeNode,
	}
}
