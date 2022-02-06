-- +goose Up
-- +goose StatementBegin
INSERT INTO "slot" ("ID", "description")
VALUES (1, 'drom.ru - homepage'),
       (2, 'ozon.ru - homepage'),
       (3, 'ozon.ru - toys'),
       (4, 'letu.ru - homepage'),
       (5, 'letu.ru - woman perfume'),
       (6, 'letu.ru - man perfume');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "slot";
-- +goose StatementEnd
