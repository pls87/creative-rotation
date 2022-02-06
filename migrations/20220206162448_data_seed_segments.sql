-- +goose Up
-- +goose StatementBegin
INSERT INTO "segment" ("description")
VALUES ('Girls < 10'),
       ('Girls 10+'),
       ('Boys 10+'),
       ('Boys <10'),
       ('Children 10+'),
       ('Children <10'),
       ('Man <30'),
       ('Man 30+'),
       ('Man 60+'),
       ('Woman <30'),
       ('Woman 30+'),
       ('Woman 60+');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "segment";
-- +goose StatementEnd
