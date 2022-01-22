-- +goose Up
-- +goose StatementBegin
CREATE TABLE "stats"
(
    "slot_id"     integer NOT NULL,
    "creative_id" integer NOT NULL,
    "impressions" integer NOT NULL,
    "conversions" integer NOT NULL,
    CONSTRAINT "stats_ID" PRIMARY KEY ("slot_id", "creative_id")
);
ALTER TABLE ONLY "stats"
    ADD CONSTRAINT "stats_slot_id_fkey" FOREIGN KEY (slot_id) REFERENCES "slot" ("ID") ON DELETE NO ACTION;
ALTER TABLE ONLY "stats"
    ADD CONSTRAINT "stats_creative_id_fkey" FOREIGN KEY (creative_id) REFERENCES "creative" ("ID") ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stats";
-- +goose StatementEnd
