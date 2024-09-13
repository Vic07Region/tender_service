-- +goose Up
-- +goose StatementBegin
ALTER TABLE offer
ADD COLUMN organization_id UUID NULL;

ALTER TABLE offer ADD CONSTRAINT offer_organization_id_foreign FOREIGN KEY (organization_id) REFERENCES organization (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd