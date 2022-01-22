-- +goose Up
-- +goose StatementBegin
CREATE TABLE "segment"
(
    "ID"          SERIAL NOT NULL,
    "description" TEXT,
    CONSTRAINT "segment_ID" PRIMARY KEY ("ID")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "segment";
-- +goose StatementEnd
