-- +goose Up
-- +goose StatementBegin
CREATE TABLE "approval"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "proposal_id" UUID NULL,
    "user_id" UUID NULL,
    "decision" VARCHAR(255) CHECK
        ("decision" IN('reject', 'approve')) NULL,
        "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "approval" ADD PRIMARY KEY("id");
CREATE TABLE "offer"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "tender_id" UUID NULL,
    "creator_id" UUID NULL,
    "organization_id" UUID NULL,
    "status" VARCHAR(255) CHECK
        (
            "status" IN('CREATED', 'PUBLISHED', 'CANCELED')
        ) NOT NULL DEFAULT 'CREATED',
        "version" INTEGER NOT NULL DEFAULT '1',
        "name" VARCHAR(100) NOT NULL,
        "description" TEXT NULL,
        "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "offer" ADD PRIMARY KEY("id");
CREATE TABLE "organization"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT NULL,
    "type" VARCHAR(10) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "organization" ADD PRIMARY KEY("id");
CREATE TABLE "review"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "proposal_id" UUID NULL,
    "creator_id" UUID NULL,
    "content" TEXT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "review" ADD PRIMARY KEY("id");
CREATE TABLE "tender"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "organization_id" UUID NULL,
    "creator_id" UUID NULL,
    "status" VARCHAR(255) CHECK
        (
            "status" IN('CREATED', 'PUBLISHED', 'CLOSED')
        ) NOT NULL DEFAULT 'CREATED',
        "version" INTEGER NOT NULL DEFAULT '1',
    "service_type" VARCHAR(50) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "tender" ADD PRIMARY KEY("id");
CREATE TABLE "employee"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "username" VARCHAR(50) NOT NULL,
    "first_name" VARCHAR(50) NULL,
    "last_name" VARCHAR(50) NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "employee" ADD PRIMARY KEY("id");
ALTER TABLE
    "employee" ADD CONSTRAINT "employee_username_unique" UNIQUE("username");
CREATE TABLE "organization_responsible"(
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "organization_id" UUID NULL,
    "user_id" UUID NULL
);
ALTER TABLE
    "organization_responsible" ADD PRIMARY KEY("id");
ALTER TABLE
    "approval" ADD CONSTRAINT "approval_proposal_id_foreign" FOREIGN KEY("proposal_id") REFERENCES "offer"("id");
ALTER TABLE
    "tender" ADD CONSTRAINT "tender_organization_id_foreign" FOREIGN KEY("organization_id") REFERENCES "organization"("id");
ALTER TABLE
    "review" ADD CONSTRAINT "review_proposal_id_foreign" FOREIGN KEY("proposal_id") REFERENCES "offer"("id");
ALTER TABLE
    "offer" ADD CONSTRAINT "offer_organization_id_foreign" FOREIGN KEY("organization_id") REFERENCES "organization"("id");
ALTER TABLE
    "offer" ADD CONSTRAINT "offer_tender_id_foreign" FOREIGN KEY("tender_id") REFERENCES "tender"("id");
ALTER TABLE
    "organization_responsible" ADD CONSTRAINT "organization_responsible_organization_id_foreign" FOREIGN KEY("organization_id") REFERENCES "organization"("id");
ALTER TABLE
    "review" ADD CONSTRAINT "review_creator_id_foreign" FOREIGN KEY("creator_id") REFERENCES "employee"("id");
ALTER TABLE
    "organization_responsible" ADD CONSTRAINT "organization_responsible_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "employee"("id");
ALTER TABLE
    "tender" ADD CONSTRAINT "tender_creator_id_foreign" FOREIGN KEY("creator_id") REFERENCES "employee"("id");
ALTER TABLE
    "approval" ADD CONSTRAINT "approval_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "employee"("id");
ALTER TABLE
    "offer" ADD CONSTRAINT "offer_creator_id_foreign" FOREIGN KEY("creator_id") REFERENCES "employee"("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  "approval";
DROP TABLE  "offers";
DROP TABLE  "tender";
DROP TABLE  "review";
DROP TABLE  "organization_responsible";
DROP TABLE  "employee";
DROP TABLE  "organization";

-- +goose StatementEnd
