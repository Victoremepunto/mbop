package store

type inMemoryStore struct {
	db []Registration
}

func (m *inMemoryStore) All(orgID string) ([]Registration, error) {
	out := make([]Registration, 0)
	for i := range m.db {
		if m.db[i].OrgID == orgID {
			out = append(out, m.db[i])
		}
	}
	return out, nil
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
	x, _ := m.Find(r.OrgID, r.UID)
	if x != nil {
		return "", ErrRegistrationAlreadyExists
	}

	for i := range m.db {
		if m.db[i].UID == r.UID {
			return "", ErrRegistrationAlreadyExists
		}
		if m.db[i].DisplayName == r.DisplayName {
			return "", ErrRegistrationAlreadyExists
		}
	}

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
