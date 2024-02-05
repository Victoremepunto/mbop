package store

import (
	"time"
)

type inMemoryStore struct {
	db               []Registration
	allowedAddresses []Address
}

func (m *inMemoryStore) All(orgID string, _, _ int) ([]Registration, int, error) {
	out := make([]Registration, 0)
	for i := range m.db {
		if m.db[i].OrgID == orgID {
			out = append(out, m.db[i])
		}
	}
	return out, len(out), nil
}

func (m *inMemoryStore) Find(orgID string, uid string) (*Registration, error) {
	for _, r := range m.db {
		if r.OrgID == orgID && r.UID == uid {
			return &r, nil
		}
	}

	return nil, ErrRegistrationNotFound
}

func (m *inMemoryStore) FindByUID(uid string) (*Registration, error) {
	for _, r := range m.db {
		if r.UID == uid {
			return &r, nil
		}
	}

	return nil, ErrRegistrationNotFound
}

func (m *inMemoryStore) Create(r *Registration) (string, error) {
	for i := range m.db {
		if m.db[i].UID == r.UID {
			return "", ErrRegistrationAlreadyExists{Detail: "uid already exists"}
		}
		if m.db[i].DisplayName == r.DisplayName {
			return "", ErrRegistrationAlreadyExists{Detail: "display_name already exists"}
		}
	}

	r.CreatedAt = time.Now()
	m.db = append(m.db, *r)
	return "", nil
}

func (m *inMemoryStore) Update(r *Registration, update *RegistrationUpdate) error {
	r, err := m.Find(r.OrgID, r.UID)
	if err != nil {
		return err
	}

	r.Extra = *update.Extra

	return nil
}

func (m *inMemoryStore) Delete(orgID string, uid string) error {
	for i := range m.db {
		if m.db[i].OrgID == orgID || m.db[i].UID == uid {
			m.db = append(m.db[:i], m.db[i+1:]...)
			return nil
		}
	}

	return ErrRegistrationNotFound
}

func (m *inMemoryStore) AllowedAddresses(orgID string) ([]Address, error) {
	return m.allowedAddresses, nil
}
func (m *inMemoryStore) AllowedIP(ip *Address) (bool, error) {
	for _, addr := range m.allowedAddresses {
		if ip.IP == addr.IP {
			return true, nil
		}
	}
	return false, nil
}
func (m *inMemoryStore) AllowAddress(ip *Address) error {
	m.allowedAddresses = append(m.allowedAddresses, *ip)
	return nil
}
func (m *inMemoryStore) DenyAddress(ip *Address) error {
	for i := range m.allowedAddresses {
		if m.allowedAddresses[i].OrgID == ip.OrgID && m.allowedAddresses[i].IP == ip.IP {
			m.allowedAddresses = append(m.allowedAddresses[:i], m.allowedAddresses[i+1:]...)
			return nil
		}
	}

	return ErrAddressNotAllowListed
}
