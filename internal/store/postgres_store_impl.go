package store

import (
	"database/sql"
	"encoding/json"
	"time"

	// the pgx driver for the database
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	l "github.com/redhatinsights/mbop/internal/logger"
)

type postgresStore struct {
	db *sql.DB
}

func (p *postgresStore) All(orgID string, limit, offset int) ([]Registration, int, error) {
	rows, err := p.db.Query(`select
	id, org_id, uid, display_name, extra, created_at
	from registrations
	where org_id = $1
	order by created_at desc
	limit $2
	offset $3`,
		orgID,
		limit,
		offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]Registration, 0)
	for rows.Next() {
		r, err := scanRegistration(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *r)
	}

	var count int
	row := p.db.QueryRow(`select count(id) from registrations where org_id = $1`, orgID)
	if err := row.Scan(&count); err != nil {
		return nil, 0, err
	}

	return out, count, nil
}

func (p *postgresStore) Find(orgID, uid string) (*Registration, error) {
	rows := p.db.QueryRow(
		`select id, org_id, uid, display_name, extra, created_at from registrations where org_id = $1 and uid = $2 limit 1`,
		orgID,
		uid,
	)
	return scanRegistration(rows)
}

func (p *postgresStore) FindByUID(uid string) (*Registration, error) {
	rows := p.db.QueryRow(`select id, org_id, uid, display_name, extra, created_at from registrations where uid = $1 limit 1`, uid)
	return scanRegistration(rows)
}

func (p *postgresStore) Create(r *Registration) (string, error) {
	res := p.db.QueryRow(
		`insert into registrations (org_id, uid, display_name, extra) values ($1, $2, $3, $4) returning id`,
		r.OrgID,
		r.UID,
		r.DisplayName,
		r.Extra,
	)

	var id string
	err := res.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		// constraint violation == 23505
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", ErrRegistrationAlreadyExists{Detail: pgErr.Detail}
			}
		} else {
			return "", err
		}
	}

	l.Log.Info("Created registration", "id", id, "org_id", r.OrgID, "uid", r.UID, "display_name", r.DisplayName)
	return id, nil
}

func (p *postgresStore) Update(r *Registration, update *RegistrationUpdate) error {
	//TODO: maybe more fields someday, not sure.
	_, err := p.db.Exec(
		`update registrations set extra = $1 where org_id = $2 and uid = $3`,
		update.Extra,
		r.OrgID,
		r.UID,
	)

	return err
}

func (p *postgresStore) Delete(orgID, uid string) error {
	res, err := p.db.Exec(
		`delete from registrations where org_id = $1 and uid = $2`,
		orgID,
		uid,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return ErrRegistrationNotFound
	}

	l.Log.Info("Deleted registration", "orgID", orgID, "uid", uid)
	return nil
}

// implement our own teeny scanner interface so we can use both sql.Row and/or sql.Rows
type scanner interface {
	Scan(dest ...any) error
}

func scanRegistration(row scanner) (*Registration, error) {
	var (
		id          string
		orgID       string
		uid         string
		displayName string
		extra       []byte
		createdAt   time.Time
	)
	err := row.Scan(&id, &orgID, &uid, &displayName, &extra, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRegistrationNotFound
		}
		return nil, err
	}

	var e map[string]any
	if extra != nil {
		err := json.Unmarshal(extra, &e)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal extra json")
		}
	}

	return &Registration{
		ID:          id,
		OrgID:       orgID,
		UID:         uid,
		DisplayName: displayName,
		Extra:       e,
		CreatedAt:   createdAt,
	}, nil
}
