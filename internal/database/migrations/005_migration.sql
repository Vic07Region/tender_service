-- +goose Up
-- +goose StatementBegin
ALTER TABLE offer DROP CONSTRAINT offer_creator_id_foreign;

ALTER TABLE offer ADD CONSTRAINT offer_creator_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd