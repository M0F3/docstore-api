package auth

import (
	"encoding/json"
	"time"

	"github.com/M0F3/docstore-api/internal/models"
	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func Init(secret string) {
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

type TokenClaim struct {
	User models.User `json:"user"`
	Exp int64 `json:"exp"`
}

func GenerateToken(user models.User, expiry time.Duration) (string, error) {
	claims :=TokenClaim{
		User: user,
		Exp:  time.Now().Add(expiry).Unix(),
	}
	b, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	var m map[string]interface{
	}

	if err = json.Unmarshal(b, &m); err != nil {
		return "" , err
	}
	_, tokenString, err := TokenAuth.Encode(m)
	return tokenString, err
}
