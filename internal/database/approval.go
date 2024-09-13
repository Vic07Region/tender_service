package database

import (
	"context"
)

type NewDecisionParams struct {
	Offer_id string
	User_id  string
	Decision string
}

func (q *Queries) NewDecision(ctx context.Context, params NewDecisionParams) error {
	sqlquery := `INSERT INTO approval (offer_id, user_id, decision) VALUES ($1,$2,$3)`
	_, err := q.db.ExecContext(ctx, sqlquery, params.Offer_id, params.User_id, params.Decision)
	return err
}

func (q *Queries) CountDecision(ctx context.Context, offer_id string) (int32, error) {
	sqlquery := `SELECT COUNT(id) FROM approval WHERE offer_id = $1 AND decision ='Approved'`
	row := q.db.QueryRowContext(ctx, sqlquery, offer_id)
	var count int32
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
