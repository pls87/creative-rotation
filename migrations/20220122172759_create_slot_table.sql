-- +goose Up
-- +goose StatementBegin
CREATE TABLE "slot"
(
    "ID"          SERIAL NOT NULL,
    "description" TEXT,
    CONSTRAINT "slot_ID" PRIMARY KEY ("ID")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "slot"
-- +goose StatementEnd
