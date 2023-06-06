package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type KeycloakToken struct {
	AccessToken      string `json:"access_token,omitempty"`
	ExpiresIn        int32  `json:"expires_in,omitempty"`
	RefreshExpiresIn int32  `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int32  `json:"not-before-policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type Token struct {
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

func (t Token) Create(ttl time.Duration, xrhid identity.Identity) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.PrivateKey)
	if err != nil {
		fmt.Println(string(t.PrivateKey))
		return "", fmt.Errorf("Failed to parse key: %w", err)
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["org_id"] = xrhid.OrgID
	claims["username"] = xrhid.User.Username
	claims["is_org_admin"] = xrhid.User.OrgAdmin

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = config.Get().TokenKID
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Failed to sign token: %w", err)
	}

	return tokenStr, nil
}
