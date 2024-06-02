package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"os"
	"strings"
)

type TokenType int

const (
	TOKEN_TYPE_INVALID = iota
	TOKEN_TYPE_AUTH
	TOKEN_TYPE_SESSION
)

func GetTokenType(token string) TokenType {
	tokenPrefix := os.Getenv("TOKEN_PREFIX")
	prefixAuthToken := fmt.Sprintf("%st_", tokenPrefix)
	if strings.HasPrefix(token, prefixAuthToken) {
		return TOKEN_TYPE_AUTH
	}
	prefixSession := fmt.Sprintf("%ss_", tokenPrefix)
	if strings.HasPrefix(token, prefixSession) {
		return TOKEN_TYPE_SESSION
	}
	return TOKEN_TYPE_INVALID
}

func generateToken(tokenType TokenType) (string, error) {
	if tokenType == TOKEN_TYPE_INVALID {
		return "", fmt.Errorf("invalid token type")
	}
	randomBytes := make([]byte, 20)
	_, err := rand.Reader.Read(randomBytes)
	if err != nil {
		return "", err
	}
	resultEncoded := make([]byte, base32.StdEncoding.EncodedLen(len(randomBytes)))
	base32.StdEncoding.Encode(resultEncoded, randomBytes)
	var token string
	if tokenType == TOKEN_TYPE_SESSION {
		token = fmt.Sprintf("%ss_%s", os.Getenv("TOKEN_PREFIX"), string(resultEncoded))
	} else {
		token = fmt.Sprintf("%st_%s", os.Getenv("TOKEN_PREFIX"), string(resultEncoded))
	}
	return token, nil
}
