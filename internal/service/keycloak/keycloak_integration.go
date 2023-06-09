package keycloak

type KeyCloak interface {
	InitKeycloakConnection() error
	GetAccessToken() (string, error)
}

func NewKeyCloakClient() KeyCloak {
	client := &Client{}

	return client
}
