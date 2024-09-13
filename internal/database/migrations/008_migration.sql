-- +goose Up
-- +goose StatementBegin

ALTER TABLE approval DROP CONSTRAINT approval_proposal_id_foreign;
ALTER TABLE approval DROP COLUMN proposal_id;

ALTER TABLE approval
ADD COLUMN offer_id UUID NULL;
ALTER TABLE approval ADD CONSTRAINT approval_offer_id_foreign FOREIGN KEY (offer_id) REFERENCES offer (id);

ALTER TABLE review DROP CONSTRAINT review_proposal_id_foreign;
ALTER TABLE review DROP COLUMN proposal_id;

ALTER TABLE review
ADD COLUMN offer_id UUID NULL;
ALTER TABLE review ADD CONSTRAINT review_offer_id_foreign FOREIGN KEY (offer_id) REFERENCES offer (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd