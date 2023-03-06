package store

type Store interface {
	All() ([]Registration, error)
	// Find a registration that both the org ID + UID match
	Find(orgID, uid string) (*Registration, error)
	// lookup a registration by uid only
	FindByUID(uid string) (*Registration, error)
	Create(r *Registration) (string, error)
	Update(r *Registration, update *RegistrationUpdate) error
	Delete(orgID, uid string) error
}
