-- +goose Up
-- +goose StatementBegin
ALTER TABLE "offer"
DROP CONSTRAINT "offer_status_check";

ALTER TABLE "offer"
ADD CONSTRAINT "offer_status_check" CHECK (
    "status" IN ('Approved','Created', 'Published', 'Canceled')
);

ALTER TABLE "approval"
DROP CONSTRAINT "approval_decision_check";

ALTER TABLE "approval"
ADD CONSTRAINT "approval_decision_check" CHECK (
    "decision" IN ('Approved','Rejected')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd