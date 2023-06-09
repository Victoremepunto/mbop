package keycloakUserService

import (
	"fmt"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
)

type KeyCloakUserService interface {
	InitKeycloakUserServiceConnection() error
	GetUsers(token string, users models.UserBody, q models.UserV1Query) (models.Users, error)
	GetAccountV3Users(orgID string, token string, q models.UserV3Query) (models.Users, error)
}

// re-declaring keycloak constant here to avoid circular module importing
const keyCloakModule = "keycloak"

func NewKeyCloakUserServiceClient() (KeyCloakUserService, error) {
	var client KeyCloakUserService

	switch config.Get().UsersModule {
	case keyCloakModule:
		client = &UserServiceClient{}
	default:
		return nil, fmt.Errorf("unsupported users module %q", config.Get().UsersModule)
	}

	return client, nil
}
