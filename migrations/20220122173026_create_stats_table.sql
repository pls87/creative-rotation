-- +goose Up
-- +goose StatementBegin
CREATE TABLE "stats"
(
    "slot_id"     integer NOT NULL,
    "creative_id" integer NOT NULL,
    "segment_id"  integer NOT NULL,
    "impressions" integer DEFAULT 0,
    "conversions" integer DEFAULT 0,
    CONSTRAINT "stats_ID" PRIMARY KEY ("slot_id", "creative_id", "segment_id")
);
ALTER TABLE ONLY "stats"
    ADD CONSTRAINT "stats_slot_id_fkey" FOREIGN KEY (slot_id) REFERENCES "slot" ("ID") ON DELETE RESTRICT;
ALTER TABLE ONLY "stats"
    ADD CONSTRAINT "stats_creative_id_fkey" FOREIGN KEY (creative_id) REFERENCES "creative" ("ID") ON DELETE RESTRICT;
ALTER TABLE ONLY "stats"
    ADD CONSTRAINT "stats_segment_id_fkey" FOREIGN KEY (segment_id) REFERENCES "segment" ("ID") ON DELETE RESTRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stats";
-- +goose StatementEnd
