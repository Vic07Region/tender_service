package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type ListTendersParams struct {
	Service_type []string
	Offset       int32
	Limit        int32
}

type Tender struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatorID      uuid.UUID `json:"creator_id"`
	Status         string    `json:"status"`
	Version        int32     `json:"version"`
	ServiceType    string    `json:"service_type"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (q *Queries) PublishedListTenders(ctx context.Context, params ListTendersParams) ([]Tender, error) {
	var sqlquery string
	if len(params.Service_type) > 0 {
		service_type := strings.Join(params.Service_type, "', '")
		sqlquery = fmt.Sprintf("SELECT id, organization_id, creator_id, status, version, service_type, name, description, created_at, updated_at FROM tender WHERE status = 'Published' AND service_type in ('%s') OFFSET $1 LIMIT $2", service_type)
	} else {
		sqlquery = "SELECT id, organization_id, creator_id, status, version, service_type, name, description, created_at, updated_at FROM tender WHERE status = 'Published' OFFSET $1 LIMIT $2"
	}
	rows, err := q.db.QueryContext(ctx, sqlquery, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tender
	for rows.Next() {
		var i Tender
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.CreatorID,
			&i.Status,
			&i.Version,
			&i.ServiceType,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type CreateTenderParams struct {
	OrganizationID string `json:"organization_id"`
	CreatorID      string `json:"creator_id"`
	Status         string `json:"status"`
	ServiceType    string `json:"service_type"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

type CreateTenderRow struct {
	ID        string    `json:"id"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) CreateTender(ctx context.Context, params CreateTenderParams) (CreateTenderRow, error) {
	sqlquery := `INSERT INTO tender (organization_id, creator_id, status, service_type, name, description)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery,
		params.OrganizationID,
		params.CreatorID,
		params.Status,
		params.ServiceType,
		params.Name,
		params.Description,
	)
	var i CreateTenderRow
	err := row.Scan(&i.ID, &i.Version, &i.CreatedAt)
	return i, err
}

type MyListTendersParams struct {
	User_id string
	Offset  int32
	Limit   int32
}

func (q *Queries) GetTender(ctx context.Context, tender_id string) (Tender, error) {
	sqlquery := `SELECT id, 
       organization_id, creator_id, status, version, service_type, name, 
       description, created_at, updated_at 
	   FROM tender 
	   WHERE id = $1 LIMIT 1`
	row := q.db.QueryRowContext(ctx, sqlquery, tender_id)
	var t Tender
	err := row.Scan(
		&t.ID,
		&t.OrganizationID,
		&t.CreatorID,
		&t.Status,
		&t.Version,
		&t.ServiceType,
		&t.Name,
		&t.Description,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return t, err
}

func (q *Queries) MyListTenders(ctx context.Context, params *MyListTendersParams) ([]Tender, error) {
	sqlquery := `SELECT id, 
       organization_id, creator_id, status, version, service_type, name, 
       description, created_at, updated_at 
	   FROM tender 
	   WHERE creator_id = $1 OFFSET $2 LIMIT $3`
	rows, err := q.db.QueryContext(ctx, sqlquery, params.User_id, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tender
	for rows.Next() {
		var i Tender
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.CreatorID,
			&i.Status,
			&i.Version,
			&i.ServiceType,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) CheckTenderStatus(ctx context.Context, tender_id string) (string, error) {
	sqlquery := `SELECT status FROM tender WHERE id = $1 LIMIT 1`
	row := q.db.QueryRowContext(ctx, sqlquery, tender_id)
	var status string
	err := row.Scan(&status)
	return status, err
}

type TenderSmResponce struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Service_type string    `json:"serviceType"`
	Version      int32     `json:"version"`
	Created_at   time.Time `json:"createdAt"`
}

func (q *Queries) ChangeTenderStatus(ctx context.Context, tender_id, new_status string) (TenderSmResponce, error) {
	sqlquery := `UPDATE tender SET
                  status = $1
                  WHERE id = $2 
                  RETURNING id, name, description, status, service_type, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery, new_status, tender_id)
	var i TenderSmResponce
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Status,
		&i.Service_type,
		&i.Version,
		&i.Created_at,
	)
	return i, err
}
