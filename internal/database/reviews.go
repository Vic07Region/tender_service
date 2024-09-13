package database

import "context"

type NewReviewParams struct {
	Offer_id string
	Content  string
	User_id  string
}

func (q *Queries) NewReview(ctx context.Context, params NewReviewParams) error {
	sqlquery := `INSERT INTO review (creator_id, offer_id, content) VALUES ($1,$2,$3)`
	_, err := q.db.ExecContext(ctx, sqlquery, params.User_id, params.Offer_id, params.Content)
	return err
}

type ListReviewParam struct {
	Offer_id string
	Limit    int32
	Offset   int32
}

type Review struct {
	ID          string
	Description string
	CreatedAt   string
}

func (q *Queries) ListReviw(ctx context.Context, params ListReviewParam) ([]Review, error) {
	sqlquery := `SELECT id, content, created_at FROM review WHERE offer_id = $1 LIMIT $2 OFFSET $3`
	rows, err := q.db.QueryContext(ctx, sqlquery, params.Offer_id, params.Offset, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Review
	for rows.Next() {
		var i Review
		if err := rows.Scan(
			&i.ID,
			&i.Description,
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
