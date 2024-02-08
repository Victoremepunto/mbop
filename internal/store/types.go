package store

import "time"

/*
Registration represents an instance of a satellite that is registered via:
- OrgID; comes from keycloak
- Uid; the CN from the satellite certificate
- DisplayName; the name to display in the UI for the satellite

ID is a generated UUID
Extra is just a jsonb column if we want to store some extra metadata someday
*/
type Registration struct {
	ID          string
	OrgID       string
	Username    string
	UID         string
	DisplayName string
	Extra       map[string]interface{}
	CreatedAt   time.Time
}

type RegistrationUpdate struct {
	Extra *map[string]interface{}
}

type AllowlistBlock struct {
	IPBlock   string
	OrgID     string
	CreatedAt time.Time
}
