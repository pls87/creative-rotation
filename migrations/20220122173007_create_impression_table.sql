-- +goose Up
-- +goose StatementBegin
CREATE TABLE "impression"
(
    "ID"          serial      NOT NULL,
    "slot_id"     integer     NOT NULL,
    "creative_id" integer     NOT NULL,
    "segment_id"  integer     NOT NULL,
    "time"        timestamptz NOT NULL,
    CONSTRAINT "impression_ID" PRIMARY KEY ("ID")
);
ALTER TABLE ONLY "impression"
    ADD CONSTRAINT "impression_slot_id_fkey" FOREIGN KEY (slot_id) REFERENCES "slot" ("ID") ON DELETE NO ACTION;
ALTER TABLE ONLY "impression"
    ADD CONSTRAINT "impression_creative_id_fkey" FOREIGN KEY (creative_id) REFERENCES "creative" ("ID") ON DELETE RESTRICT;
ALTER TABLE ONLY "impression"
    ADD CONSTRAINT "impression_segment_id_fkey" FOREIGN KEY (segment_id) REFERENCES "segment" ("ID") ON DELETE RESTRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "impression";
-- +goose StatementEnd
