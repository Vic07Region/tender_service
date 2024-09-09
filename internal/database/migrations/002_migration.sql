-- +goose Up
-- +goose StatementBegin
-- Изменяем проверку "status"
SET CONSTRAINTS ALL DEFERRED;

ALTER TABLE "tender"
DROP CONSTRAINT "tender_status_check";

UPDATE "tender" SET "status" = 'Created' WHERE "status" = 'CREATED';
UPDATE "tender" SET "status" = 'Published' WHERE "status" = 'PUBLISHED';
UPDATE "tender" SET "status" = 'Closed' WHERE "status" = 'CLOSED';

ALTER TABLE "tender"
ADD CONSTRAINT "tender_status_check" CHECK (
    "status" IN ('Created', 'Published', 'Closed')
);

-- Изменяем значение по умолчанию для "status"
ALTER TABLE "tender"
ALTER COLUMN "status" SET DEFAULT 'Created';


-- Изменяем проверку "status"
ALTER TABLE "offer"
DROP CONSTRAINT "offer_status_check";

UPDATE "offer" SET "status" = 'Created' WHERE "status" = 'CREATED';
UPDATE "offer" SET "status" = 'Published' WHERE "status" = 'PUBLISHED';
UPDATE "offer" SET "status" = 'Canceled' WHERE "status" = 'CANCELED';

ALTER TABLE "offer"
ADD CONSTRAINT "offer_status_check" CHECK (
    "status" IN ('Created', 'Published', 'Canceled')
);

-- Изменяем значение по умолчанию для "status"
ALTER TABLE "offer"
ALTER COLUMN "status" SET DEFAULT 'Created';



SET CONSTRAINTS ALL IMMEDIATE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
