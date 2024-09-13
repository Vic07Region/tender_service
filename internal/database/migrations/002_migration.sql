-- +goose Up
-- +goose StatementBegin
CREATE TABLE approval(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    proposal_id UUID NULL,
    user_id UUID NULL,
    decision VARCHAR(255) CHECK
        (decision IN('reject', 'approve')) NULL,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
-- Комментарий к таблице approval
ALTER TABLE approval ADD PRIMARY KEY (id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  approval;

-- +goose StatementEnd