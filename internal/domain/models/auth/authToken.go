package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/uptrace/bun"
)

type AuthToken struct {
	bun.BaseModel `bun:"table:auth,alias:au"`
	Created       time.Time `bun:",notnull"`
	Expire        bun.NullTime
	Id            string `bun:",pk"`
	Token         string `bun:",notnull"`
	UserId        string `bun:",notnull"`
}

var ErrNoAuthTokenFound error = fmt.Errorf("no auth token found")

func CreateAuthToken(snowflakeNode *snowflake.Node, userId string, expire *time.Time) (*AuthToken, error) {
	id := snowflakeNode.Generate().String()
	tokenExpire := bun.NullTime{}
	if expire != nil {
		tokenExpire.Time = *expire
	}
	token, err := generateToken(TOKEN_TYPE_AUTH)
	if err != nil {
		return nil, err
	}
	return &AuthToken{
		Id:      id,
		UserId:  userId,
		Token:   token,
		Expire:  tokenExpire,
		Created: time.Now(),
	}, nil
}

func (auth *AuthToken) DisableToken() {
	auth.Expire = bun.NullTime{Time: time.Now()}
}

type AuthTokenRepository interface {
	GetAuthTokenById(ctx context.Context, id string) (*AuthToken, error)
	GetAuthTokenByToken(ctx context.Context, token string) (*AuthToken, error)
	Save(ctx context.Context, authToken *AuthToken) error
}
