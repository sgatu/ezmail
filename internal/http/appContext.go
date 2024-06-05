package http

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bwmarrin/snowflake"
	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	"github.com/sgatu/ezmail/internal/infrastructure/repositories"
	"github.com/sgatu/ezmail/internal/service/ses"
	"github.com/uptrace/bun"
)

type AppContext struct {
	UserRepository       user.UserRepository
	AuthTokenRepository  auth.AuthTokenRepository
	SessionRepository    auth.SessionRepository
	DomainInfoRepository domain.DomainInfoRepository
	SESService           *ses.SESService
	AuthManager          authManager
	SnowflakeNode        *snowflake.Node
}

type authManager struct {
	authTokenRepository auth.AuthTokenRepository
	sessionRepository   auth.SessionRepository
	userRepository      user.UserRepository
}

type currentUserKey struct{}

var CurrentUserKey currentUserKey = currentUserKey{}

func (am *authManager) ValidateToken(ctx context.Context, token string) *user.User {
	userId := ""
	tokenType := auth.GetTokenType(token)
	switch tokenType {
	case auth.TOKEN_TYPE_AUTH:
		tok, err := am.authTokenRepository.GetAuthTokenByToken(ctx, token)
		// we didn't find a token or expire > 0 but before now
		if err != nil || (!tok.Expire.IsZero() && tok.Expire.Before(time.Now())) {
			return nil
		}
		userId = tok.UserId
	case auth.TOKEN_TYPE_SESSION:
		session, err := am.sessionRepository.GetSessionBySessionId(ctx, token)
		if err != nil || session.Expire.Before(time.Now()) {
			return nil
		}
		userId = session.UserId
	default:
		return nil
	}
	user, err := am.userRepository.GetById(ctx, userId)
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
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	authRepository := repositories.NewMysqlAuthTokenRepository(db)
	userRepository := repositories.NewMysqlUserRepository(db)
	sessionRepository := repositories.NewMysqlSessionRepository(db)
	domainInfoRepository := repositories.NewMysqlDomainInfoRepository(db)
	return &AppContext{
		UserRepository:       userRepository,
		AuthTokenRepository:  authRepository,
		SessionRepository:    sessionRepository,
		DomainInfoRepository: domainInfoRepository,
		AuthManager:          authManager{authTokenRepository: authRepository, userRepository: userRepository, sessionRepository: sessionRepository},
		SESService:           ses.NewSESService(domainInfoRepository, awsConfig, snowflakeNode),
		SnowflakeNode:        snowflakeNode,
	}
}
