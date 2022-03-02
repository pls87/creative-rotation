-- +goose Up
-- +goose StatementBegin
CREATE TABLE "creative"
(
    "ID"          SERIAL NOT NULL,
    "description" TEXT,
    CONSTRAINT "creative_ID" PRIMARY KEY ("ID")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "creative";
-- +goose StatementEnd
