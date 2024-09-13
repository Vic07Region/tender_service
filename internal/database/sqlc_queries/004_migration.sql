-- +goose Up
-- +goose StatementBegin
ALTER TABLE "offer" DROP COLUMN "organization_id";
ALTER TABLE "offer" ADD COLUMN "author_type" VARCHAR(15) NOT NULL;


-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd





CREATE TABLE review(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    proposal_id UUID NULL,
    creator_id UUID NULL,
    content TEXT NULL,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
-- Комментарий к таблице review
ALTER TABLE review ADD PRIMARY KEY (id);

CREATE TABLE tender(
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organization_id UUID NULL,
    creator_id UUID NULL,
    status VARCHAR(255) CHECK
        (
            status IN('Created', 'Published', 'Closed')
        ) NOT NULL DEFAULT 'Created',
    version INTEGER NOT NULL DEFAULT '1',
    service_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL DEFAULT ' ',
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);
-- Комментарий к таблице tender
ALTER TABLE tender ADD PRIMARY KEY (id);

CREATE TABLE tender_history (
    id SERIAL PRIMARY KEY,
    tender_id UUID NULL,
    creator_id UUID NULL,
    version INTEGER NOT NULL DEFAULT 1,
    service_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL DEFAULT ' ',
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE offer_history (
    id SERIAL PRIMARY KEY,
    offer_id UUID NULL,
    creator_id UUID NULL,
    version INTEGER NOT NULL DEFAULT 1,
    service_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL DEFAULT ' ',
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NULL DEFAULT CURRENT_TIMESTAMP
);



ALTER TABLE approval ADD CONSTRAINT approval_proposal_id_foreign FOREIGN KEY (proposal_id) REFERENCES offer (id);
ALTER TABLE tender ADD CONSTRAINT tender_organization_id_foreign FOREIGN KEY (organization_id) REFERENCES organization (id);
ALTER TABLE review ADD CONSTRAINT review_proposal_id_foreign FOREIGN KEY (proposal_id) REFERENCES offer (id);
ALTER TABLE offer ADD CONSTRAINT offer_organization_id_foreign FOREIGN KEY (organization_id) REFERENCES organization (id);
ALTER TABLE offer ADD CONSTRAINT offer_tender_id_foreign FOREIGN KEY (tender_id) REFERENCES tender (id);
ALTER TABLE review ADD CONSTRAINT review_creator_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);
ALTER TABLE tender ADD CONSTRAINT tender_creator_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);
ALTER TABLE approval ADD CONSTRAINT approval_user_id_foreign FOREIGN KEY (user_id) REFERENCES employee (id);
ALTER TABLE offer ADD CONSTRAINT offer_creator_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);
ALTER TABLE tender_history ADD CONSTRAINT tenderhistory_tender_id_foreign FOREIGN KEY (tender_id) REFERENCES tender (id);
ALTER TABLE tender_history ADD CONSTRAINT tenderhistory_employee_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);
ALTER TABLE offer_history ADD CONSTRAINT offerhistory_offer_id_foreign FOREIGN KEY (offer_id) REFERENCES offer (id);
ALTER TABLE offer_history ADD CONSTRAINT tenderhistory_employee_id_foreign FOREIGN KEY (creator_id) REFERENCES employee (id);
