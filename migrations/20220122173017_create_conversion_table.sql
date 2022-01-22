-- +goose Up
-- +goose StatementBegin
CREATE TABLE "conversion"
(
    "ID"          serial      NOT NULL,
    "slot_id"     integer     NOT NULL,
    "creative_id" integer     NOT NULL,
    "time"        timestamptz NOT NULL,
    CONSTRAINT "conversion_ID" PRIMARY KEY ("ID")
);
ALTER TABLE ONLY "conversion"
    ADD CONSTRAINT "conversion_slot_id_fkey" FOREIGN KEY (slot_id) REFERENCES "slot" ("ID") ON DELETE NO ACTION;
ALTER TABLE ONLY "conversion"
    ADD CONSTRAINT "conversion_creative_id_fkey" FOREIGN KEY (creative_id) REFERENCES "creative" ("ID") ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "conversion";
-- +goose StatementEnd
