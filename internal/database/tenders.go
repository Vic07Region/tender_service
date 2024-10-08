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
		sqlquery = fmt.Sprintf("SELECT id, organization_id, creator_id, status, version, service_type, name, description, created_at, updated_at FROM tender WHERE status = 'Published' AND service_type in ('%s') ORDER BY name OFFSET $1 LIMIT $2", service_type)
	} else {
		sqlquery = "SELECT id, organization_id, creator_id, status, version, service_type, name, description, created_at, updated_at FROM tender WHERE status = 'Published' ORDER BY name OFFSET $1 LIMIT $2"
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
	   WHERE creator_id = $1 ORDER BY name OFFSET $2 LIMIT $3`
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

func (q *Queries) ChangeTenderStatus(ctx context.Context, tender_id, new_status string) (*Tender, error) {
	sqlquery := `UPDATE tender SET
                  status = $1
                  WHERE id = $2 
                  RETURNING id, name, description, status, service_type, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery, new_status, tender_id)
	var i Tender
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Status,
		&i.ServiceType,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err
}

type TenderChangeParam struct {
	Tender_id    string
	Name         string
	Description  string
	Service_type string
}

func buildUpdateQuery(param TenderChangeParam) string {
	query := "UPDATE tender SET version = version + 1, "

	var setClauses []string
	if param.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = '%s'", param.Name))
	}
	if param.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = '%s'", param.Description))
	}
	if param.Service_type != "" {
		setClauses = append(setClauses, fmt.Sprintf("service_type = '%s'", param.Service_type))
	}

	if len(setClauses) > 0 {
		query += strings.Join(setClauses, ", ")
	} else {
		return "UPDATE tender SET "
	}

	query += ",updated_at = $2 WHERE id = $1 RETURNING id, name, description, status, service_type, version, created_at"

	return query
}

func (q *Queries) EditTender(ctx context.Context, param TenderChangeParam) (*Tender, error) {
	sqlquery := buildUpdateQuery(param)
	now := time.Now()
	row := q.db.QueryRowContext(ctx, sqlquery, param.Tender_id, now)
	var i Tender
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Status,
		&i.ServiceType,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err
}

type CreateTenderHistoryParams struct {
	Tender_id   string
	Creator_id  string
	ServiceType string
	Name        string
	Description string
	OldVersion  int32
}

func (q *Queries) CreateTenderHistory(ctx context.Context, params CreateTenderHistoryParams) error {
	sqlquery := `INSERT INTO tender_history (tender_id, creator_id, service_type, name, description, version)
	VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := q.db.ExecContext(ctx, sqlquery,
		params.Tender_id,
		params.Creator_id,
		params.ServiceType,
		params.Name,
		params.Description,
		params.OldVersion,
	)
	return err
}

type EditTenderWithTxParam struct {
	ChangeTenderParam  TenderChangeParam
	TenderHistoryParam CreateTenderHistoryParams
}

func (q *Queries) EditTenderWithTX(ctx context.Context, params EditTenderWithTxParam) (*Tender, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//редактирование тендера
	tender, err := q.EditTender(ctx, params.ChangeTenderParam)
	if err != nil {
		return nil, err
	}
	//создание истории изменений
	err = q.CreateTenderHistory(ctx, params.TenderHistoryParam)
	if err != nil {
		return nil, err
	}

	return tender, tx.Commit()
}

type TenderHistory struct {
	Tender_id   string
	Creator_id  string
	ServiceType string
	Name        string
	Description string
	Version     int32
}

func (q *Queries) GetTenderHistory(ctx context.Context, tender_id string, version int32) (*TenderHistory, error) {
	sqlquery := `SELECT tender_id, 
       creator_id, service_type,
       name, description, version FROM tender_history WHERE tender_id = $1 AND version = $2 LIMIT 1`

	row := q.db.QueryRowContext(ctx, sqlquery, tender_id, version)
	var i TenderHistory
	err := row.Scan(
		&i.Tender_id,
		&i.Creator_id,
		&i.ServiceType,
		&i.Name,
		&i.Description,
		&i.Version,
	)
	return &i, err
}
