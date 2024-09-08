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
