package keycloak

import (
	"fmt"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
)

type KeyCloak interface {
	InitKeycloakConnection() error
	GetAccessToken() (string, error)
	GetUsers(token string, users models.UserBody, q models.UserV1Query) (models.Users, error)
	GetAccountV3Users(orgID string, token string, q models.UserV3Query) (models.Users, error)
}

// re-declaring keycloak constant here to avoid circular module importing
const keyCloakModule = "keycloak"

func NewKeyCloakClient() (KeyCloak, error) {
	var client KeyCloak

	switch config.Get().UsersModule {
	case keyCloakModule:
		client = &Client{}
	default:
		return nil, fmt.Errorf("unsupported users module %q", config.Get().UsersModule)
	}

	return client, nil
}
