package config

import (
	"os"
	"strconv"
)

type MbopConfig struct {
	MailerModule           string
	JwtModule              string
	JwkURL                 string
	UsersModule            string
	CognitoAppClientID     string
	CognitoAppClientSecret string
	CognitoScope           string
	OauthTokenURL          string
	AmsURL                 string
	TokenTTL               string
	TokenKID               string
	PrivateKey             string
	PublicKey              string
	DisableCatchall        bool
	IsInternalLabel        string

	StoreBackend     string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
}

var conf *MbopConfig

func Get() *MbopConfig {
	if conf != nil {
		return conf
	}

	disableCatchAll, _ := strconv.ParseBool(fetchWithDefault("DISABLE_CATCHALL", "false"))

	c := &MbopConfig{
		UsersModule:     fetchWithDefault("USERS_MODULE", ""),
		JwtModule:       fetchWithDefault("JWT_MODULE", ""),
		JwkURL:          fetchWithDefault("JWK_URL", ""),
		MailerModule:    fetchWithDefault("MAILER_MODULE", "print"),
		DisableCatchall: disableCatchAll,

		DatabaseHost:     fetchWithDefault("DATABASE_HOST", "localhost"),
		DatabasePort:     fetchWithDefault("DATABASE_PORT", "5432"),
		DatabaseUser:     fetchWithDefault("DATABASE_USER", "postgres"),
		DatabasePassword: fetchWithDefault("DATABASE_PASSWORD", ""),
		DatabaseName:     fetchWithDefault("DATABASE_NAME", "mbop"),
		StoreBackend:     fetchWithDefault("STORE_BACKEND", "memory"),

		CognitoAppClientID:     fetchWithDefault("COGNITO_APP_CLIENT_ID", ""),
		CognitoAppClientSecret: fetchWithDefault("COGNITO_APP_CLIENT_SECRET", ""),
		CognitoScope:           fetchWithDefault("COGNITO_SCOPE", ""),
		OauthTokenURL:          fetchWithDefault("OAUTH_TOKEN_URL", ""),
		AmsURL:                 fetchWithDefault("AMS_URL", ""),
		TokenTTL:               fetchWithDefault("TOKEN_TTL_DURATION", "5m"),
		TokenKID:               fetchWithDefault("TOKEN_KID", ""),
		PrivateKey:             fetchWithDefault("TOKEN_PRIVATE_KEY", ""),
		PublicKey:              fetchWithDefault("TOKEN_PUBLIC_KEY", ""),
		IsInternalLabel:        fetchWithDefault("IS_INTERNAL_LABEL", ""),
	}

	conf = c
	return conf
}

func fetchWithDefault(name, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return defaultValue
}

// TO BE USED FROM TESTING ONLY.
func Reset() {
	conf = nil
}
