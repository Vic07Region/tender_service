-- +goose Up
-- +goose StatementBegin
CREATE TABLE offer(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    tender_id UUID NULL,
    creator_id UUID NULL,
    author_type VARCHAR(15) NOT NULL,
    status VARCHAR(255) CHECK
        (
            status IN('Created', 'Published', 'Canceled')
        ) NOT NULL DEFAULT 'Created',
    version INTEGER NOT NULL DEFAULT '1',
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
-- Комментарий к таблице offer
ALTER TABLE offer ADD PRIMARY KEY (id);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  offer;
-- +goose StatementEnd