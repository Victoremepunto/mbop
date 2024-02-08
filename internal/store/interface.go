package store

type Store interface {
	RegistrationStore
	AllowlistStore
}

type RegistrationStore interface {
	All(orgID string, limit, offset int) ([]Registration, int, error)
	// Find a registration that both the org ID + UID match
	Find(orgID, uid string) (*Registration, error)
	// lookup a registration by uid only
	FindByUID(uid string) (*Registration, error)
	Create(r *Registration) (string, error)
	Update(r *Registration, update *RegistrationUpdate) error
	Delete(orgID, uid string) error
}

type AllowlistStore interface {
	AllowedAddresses(orgID string) ([]AllowlistBlock, error)
	AllowedIP(ip, orgID string) (bool, error)
	AllowAddress(ip *AllowlistBlock) error
	DenyAddress(ip *AllowlistBlock) error
}
