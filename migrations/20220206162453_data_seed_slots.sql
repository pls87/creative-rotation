-- +goose Up
-- +goose StatementBegin
INSERT INTO "slot" ("description")
VALUES ('drom.ru - homepage'),
       ('ozon.ru - homepage'),
       ('ozon.ru - toys'),
       ('letu.ru - homepage'),
       ('letu.ru - man perfume');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "slot";
-- +goose StatementEnd
