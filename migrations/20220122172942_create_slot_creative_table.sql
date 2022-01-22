-- +goose Up
-- +goose StatementBegin
CREATE TABLE "slot_creative"
(
    "slot_id"     integer NOT NULL,
    "creative_id" integer NOT NULL,
    CONSTRAINT "slot_creative_ID" PRIMARY KEY ("slot_id", "creative_id")
);
ALTER TABLE ONLY "slot_creative"
    ADD CONSTRAINT "slot_id_fkey" FOREIGN KEY (slot_id) REFERENCES "slot" ("ID") ON DELETE CASCADE;
ALTER TABLE ONLY "slot_creative"
    ADD CONSTRAINT "creative_id_fkey" FOREIGN KEY (creative_id) REFERENCES "creative" ("ID") ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "slot_creative";
-- +goose StatementEnd
