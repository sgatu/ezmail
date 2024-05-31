package http

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
	"github.com/uptrace/bun"
)

type AppContext struct {
	UserRepository      user.UserRepository
	AuthTokenRepository auth.AuthTokenRepository
	AuthManager         authManager
	SnowflakeNode       *snowflake.Node
}

type authManager struct {
	authTokenRepository auth.AuthTokenRepository
	userRepository      user.UserRepository
}

type currentUserKey struct{}

var CurrentUserKey currentUserKey = currentUserKey{}

func (am *authManager) ValidateToken(ctx context.Context, token string) *user.User {
	tok, err := am.authTokenRepository.GetAuthTokenByToken(ctx, token)
	if err != nil {
		return nil
	}
	isValid := tok.Expire.IsZero() || tok.Expire.After(time.Now())
	if !isValid {
		return nil
	}
	user, err := am.userRepository.GetById(ctx, tok.UserId)
	if err != nil {
		return nil
	}
	return user
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
	authRepository := repositories.NewMysqlAuthTokenRepository(db)
	userRepository := repositories.NewMysqlUserRepository(db)
	return &AppContext{
		UserRepository:      userRepository,
		AuthTokenRepository: authRepository,
		AuthManager:         authManager{authTokenRepository: authRepository, userRepository: userRepository},
		SnowflakeNode:       snowflakeNode,
	}
}
