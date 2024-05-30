package auth

import (
	"context"
	"crypto/rand"
	"encoding/base32"
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
	token, err := generateNewToken()
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

func generateNewToken() (string, error) {
	token := "ezmail_"
	randomBytes := make([]byte, 20)
	_, err := rand.Reader.Read(randomBytes)
	if err != nil {
		return "", err
	}
	resultEncoded := make([]byte, base32.StdEncoding.EncodedLen(len(randomBytes)))
	base32.StdEncoding.Encode(resultEncoded, randomBytes)
	token = token + string(resultEncoded)
	return token, nil
}

func (auth *AuthToken) DisableToken() {
	auth.Expire = bun.NullTime{Time: time.Now()}
}

type AuthTokenRepository interface {
	GetAuthTokenById(ctx context.Context, id string) (*AuthToken, error)
	GetAuthTokenByToken(ctx context.Context, token string) (*AuthToken, error)
	Save(ctx context.Context, authToken *AuthToken) error
}
