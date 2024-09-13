package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	Owner = "Created, Published, Canceled"
	All   = "Published"
)

type Offer struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorId   uuid.UUID `json:"authorId"`
	Version    int32     `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateOfferParam struct {
	Name            string
	Description     string
	TenderId        string
	AuthorType      string
	AuthorId        string
	Organization_id string
}

func (q *Queries) CreateOffer(ctx context.Context, param CreateOfferParam) (*Offer, error) {
	sqlquery := `INSERT INTO offer (name, description, tender_id, author_type, creator_id,organization_id)
	VALUES ($1,$2,$3,$4,$5) RETURNING id, name, status, author_type, creator_id, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery,
		param.Name,
		param.Description,
		param.TenderId,
		param.AuthorType,
		param.AuthorId,
		param.Organization_id,
	)
	var i Offer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.AuthorType,
		&i.AuthorId,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err

}

type MyListOffersParams struct {
	Creator_id string
	Offset     int32
	Limit      int32
}

func (q *Queries) MyListOffers(ctx context.Context, params *MyListOffersParams) ([]Offer, error) {
	sqlquery := `SELECT id, name, status, author_type, creator_id, version, created_at
	   FROM offer 
	   WHERE creator_id = $1 ORDER BY name OFFSET $2 LIMIT $3`
	rows, err := q.db.QueryContext(ctx, sqlquery, params.Creator_id, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Offer
	for rows.Next() {
		var i Offer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Status,
			&i.AuthorType,
			&i.AuthorId,
			&i.Version,
			&i.CreatedAt,
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

type TenderListOffersParams struct {
	Tender_id       string
	Offset          int32
	Limit           int32
	Organization_id string
}

func (q *Queries) TenderListOffers(ctx context.Context, params *TenderListOffersParams) ([]Offer, error) {
	sqlquery := `SELECT id, name, status, author_type, creator_id, version, created_at
	   FROM offer 
	   WHERE tender_id = $1 AND (
    (organization_id = $2 AND status IN ('Approved','Created', 'Published', 'Canceled'))
    OR (organization_id != $2 AND status = 'Published')
)  ORDER BY name OFFSET $3 LIMIT $4`
	rows, err := q.db.QueryContext(ctx, sqlquery, params.Tender_id, params.Organization_id, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Offer
	for rows.Next() {
		var i Offer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Status,
			&i.AuthorType,
			&i.AuthorId,
			&i.Version,
			&i.CreatedAt,
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

type OfferFull struct {
	ID              uuid.UUID
	Tender_ID       uuid.UUID
	Creator_ID      uuid.UUID
	Organization_ID uuid.UUID
	Name            string
	Description     string
	Status          string
	AuthorType      string
	AuthorId        uuid.UUID
	Version         int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (q *Queries) GetOffer(ctx context.Context, offer_id string) (*OfferFull, error) {
	sqlquery := `SELECT id, tender_id, creator_id, organization_id,
       author_type, status, version, name, 
       description, created_at, updated_at FROM offer WHERE id = $1 LIMIT 1`
	rows, err := q.db.QueryContext(ctx, sqlquery, offer_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offer OfferFull
	if err := rows.Scan(
		&offer.ID,
		&offer.Tender_ID,
		&offer.Creator_ID,
		&offer.Organization_ID,
		&offer.AuthorType,
		&offer.Status,
		&offer.Version,
		&offer.Name,
		&offer.Description,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &offer, nil
}

func (q *Queries) GetOfferStatus(ctx context.Context, offer_id, author_id string) (string, error) {
	sqlquery := `SELECT status FROM offer WHERE id = $1 AND creator_id = $2 LIMIT 1`
	rows, err := q.db.QueryContext(ctx, sqlquery, offer_id, author_id)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var status string
	if err := rows.Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

func (q *Queries) ChangeOfferStatus(ctx context.Context, offer_id, new_status string) (*Offer, error) {
	sqlquery := `UPDATE offer SET
                  status = $1
                  WHERE id = $2 
                  RETURNING id, name, status, author_type, creator_id, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery, new_status, offer_id)
	var i Offer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.AuthorType,
		&i.AuthorId,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err
}

type OfferChangeParam struct {
	Bid_id      string
	Name        string
	Description string
}

func buildUpdateQueryOffer(param OfferChangeParam) string {
	query := "UPDATE offer SET version = version + 1, "

	var setClauses []string
	if param.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = '%s'", param.Name))
	}
	if param.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = '%s'", param.Description))
	}

	if len(setClauses) > 0 {
		query += strings.Join(setClauses, ", ")
	} else {
		return "UPDATE offer SET "
	}

	query += ",updated_at = $2 WHERE id = $1 id, name, status, author_type, creator_id, version, created_at"

	return query
}

func (q *Queries) EditOffer(ctx context.Context, param OfferChangeParam) (*Offer, error) {
	sqlquery := buildUpdateQueryOffer(param)
	now := time.Now()
	row := q.db.QueryRowContext(ctx, sqlquery, param.Bid_id, now)
	var i Offer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.AuthorType,
		&i.AuthorId,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err
}

type CreateOfferHistoryParams struct {
	Offer_id    string
	Creator_id  string
	Name        string
	Description string
	OldVersion  int32
}

func (q *Queries) CreateOfferHistory(ctx context.Context, params CreateOfferHistoryParams) error {
	sqlquery := `INSERT INTO offer_history (offer_id, creator_id, name, description, version)
	VALUES ($1,$2,$3,$4,$5)`
	_, err := q.db.ExecContext(ctx, sqlquery,
		params.Offer_id,
		params.Creator_id,
		params.Name,
		params.Description,
		params.OldVersion,
	)
	return err
}

type EditOfferWithTxParam struct {
	ChangeOfferParam  OfferChangeParam
	OfferHistoryParam CreateOfferHistoryParams
}

func (q *Queries) EditOfferWithTX(ctx context.Context, params EditOfferWithTxParam) (*Offer, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//редактирование предложения
	offer, err := q.EditOffer(ctx, params.ChangeOfferParam)
	if err != nil {
		return nil, err
	}
	//создание истории изменений
	err = q.CreateOfferHistory(ctx, params.OfferHistoryParam)
	if err != nil {
		return nil, err
	}

	return offer, tx.Commit()
}

type OfferHistory struct {
	Tender_id   string
	Creator_id  string
	Name        string
	Description string
	Version     int32
}

func (q *Queries) GetOfferHistory(ctx context.Context, offer_id string, version int32) (*OfferHistory, error) {
	sqlquery := `SELECT offer_id, 
       creator_id, name, description, version FROM tender_history WHERE offer_id = $1 AND version = $2 LIMIT 1`

	row := q.db.QueryRowContext(ctx, sqlquery, offer_id, version)
	var i OfferHistory
	err := row.Scan(
		&i.Tender_id,
		&i.Creator_id,
		&i.Name,
		&i.Description,
		&i.Version,
	)
	return &i, err
}

func (q *Queries) GetOfferByAuthor(ctx context.Context, tender_id, user_id string) (*Offer, error) {
	sqlquery := `SELECT id, name, status, 
       author_type, creator_id, version, 
       created_at FROM offer WHERE tender_id = $1 AND creator_id = $2 LIMIT 1`
	row := q.db.QueryRowContext(ctx, sqlquery, tender_id, user_id)
	var i Offer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.AuthorType,
		&i.AuthorId,
		&i.Version,
		&i.CreatedAt,
	)
	return &i, err

}
