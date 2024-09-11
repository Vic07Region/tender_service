package database

import "context"

type Offer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    int32  `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type CreateOfferParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}

func (q *Queries) CreateOffer(ctx context.Context, param CreateOfferParam) (*Offer, error) {
	sqlquery := `INSERT INTO offer (name, description, tender_id, author_type, creator_id)
	VALUES ($1,$2,$3,$4,$5) RETURNING id, name, status, author_type, creator_id, version, created_at`
	row := q.db.QueryRowContext(ctx, sqlquery,
		param.Name,
		param.Description,
		param.TenderId,
		param.AuthorType,
		param.AuthorId,
	)
	var i Offer
	err := row.Scan(&i.ID, &i.Version, &i.CreatedAt)
	return &i, err

}

type MyListOffersParams struct {
	User_id string
	Offset  int32
	Limit   int32
}

func (q *Queries) MyListOffers(ctx context.Context, params *MyListOffersParams) ([]Offer, error) {
	sqlquery := `SELECT id, name, status, author_type, creator_id, version, created_at
	   FROM offer 
	   WHERE creator_id = $1 ORDER BY name OFFSET $2 LIMIT $3`
	rows, err := q.db.QueryContext(ctx, sqlquery, params.User_id, params.Offset, params.Limit)
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
