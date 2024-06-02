package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:session,alias:sess"`
	Created       time.Time `bun:",notnull"`
	Expire        time.Time `bun:",notnull"`
	Id            string    `bun:",pk"`
	SessionId     string    `bun:",notnull"`
	UserId        string    `bun:",notnull"`
}

var ErrNoSessionFound error = fmt.Errorf("no session found")

func GetDefaultSessionExpire() time.Time {
	return time.Now().Add(time.Duration(time.Hour * 6))
}

func CreateSession(snowflakeNode *snowflake.Node, userId string, expire time.Time) (*Session, error) {
	id := snowflakeNode.Generate().String()
	sessId, err := generateToken(TOKEN_TYPE_SESSION)
	if err != nil {
		return nil, err
	}
	return &Session{
		Id:        id,
		UserId:    userId,
		SessionId: sessId,
		Expire:    expire,
		Created:   time.Now(),
	}, nil
}

func (sess *Session) ExpireSession() {
	sess.Expire = time.Now()
}

type SessionRepository interface {
	GetSessionById(ctx context.Context, id string) (*Session, error)
	GetSessionBySessionId(ctx context.Context, sessionId string) (*Session, error)
	Save(ctx context.Context, session *Session) error
}
