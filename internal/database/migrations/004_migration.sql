-- +goose Up
-- +goose StatementBegin
ALTER TABLE "offer" DROP COLUMN "organization_id";
ALTER TABLE "offer" ADD COLUMN "author_type" VARCHAR(15) NOT NULL;


-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
