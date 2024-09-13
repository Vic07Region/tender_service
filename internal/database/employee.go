package database

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) FetchUserID(ctx context.Context, username string) (string, error) {
	sqlquery := "SELECT ID FROM employee WHERE username = $1 LIMIT 1"
	row := q.db.QueryRowContext(ctx, sqlquery, username)

	var name string
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

type User struct {
	ID        uuid.UUID
	Username  string
	Firstname string
	Lastname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FetchUser(ctx context.Context, user_id string) (*User, error) {
	sqlquery := "SELECT * FROM employee WHERE id = $1 LIMIT 1"
	row := q.db.QueryRowContext(ctx, sqlquery, user_id)

	var user User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (q *Queries) IsResponsible(ctx context.Context, org_id, user_id string) (bool, error) {
	sqlquery := "SELECT user_id FROM organization_responsible WHERE organization_id = $1 LIMIT 1"
	row := q.db.QueryRowContext(ctx, sqlquery, org_id)

	var userid string
	if err := row.Scan(&userid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return userid == user_id, nil
}

func (q *Queries) GetUserOrganization(ctx context.Context, user_id string) (string, error) {
	sqlquery := "SELECT organization_id FROM organization_responsible WHERE user_id = $1 LIMIT 1"
	row := q.db.QueryRowContext(ctx, sqlquery, user_id)

	var org_id string
	if err := row.Scan(&org_id); err != nil {
		return "", err
	}
	return org_id, nil
}

func (q *Queries) ResponsibleUserCount(ctx context.Context, org_id string) (int32, error) {
	sqlquery := "SELECT COUNT(id) FROM organization_responsible WHERE organization_id = $1"
	row := q.db.QueryRowContext(ctx, sqlquery, org_id)
	var count int32
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
