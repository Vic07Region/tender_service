-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tender_history" (
    "id" SERIAL PRIMARY KEY,
    "tender_id" UUID NULL,
    "creator_id" UUID NULL,
    "version" INTEGER NOT NULL DEFAULT 1,
    "service_type" VARCHAR(50) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT NULL DEFAULT ' ',
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "offer_history" (
    "id" SERIAL PRIMARY KEY,
    "offer_id" UUID NULL,
    "creator_id" UUID NULL,
    "version" INTEGER NOT NULL DEFAULT 1,
    "service_type" VARCHAR(50) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT NULL DEFAULT ' ',
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE "tender_history" ADD CONSTRAINT "tenderhistory_tender_id_foreign" FOREIGN KEY ("tender_id") REFERENCES "tender" ("id");
ALTER TABLE "tender_history" ADD CONSTRAINT "tenderhistory_employee_id_foreign" FOREIGN KEY("creator_id") REFERENCES "employee"("id");
ALTER TABLE "offer_history" ADD CONSTRAINT "offerhistory_offer_id_foreign" FOREIGN KEY ("offer_id") REFERENCES "offer" ("id");
ALTER TABLE "offer_history" ADD CONSTRAINT "tenderhistory_employee_id_foreign" FOREIGN KEY("creator_id") REFERENCES "employee"("id");
-- +goose StatementEnd

-- +goose Down
DROP TABLE tender_history;
DROP TABLE offer_history;

-- +goose StatementBegin


-- +goose StatementEnd
