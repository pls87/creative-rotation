-- +goose Up
-- +goose StatementBegin
INSERT INTO "segment" ("ID", "description")
VALUES (1, 'Girls < 10'),
       (2, 'Girls 10+'),
       (3, 'Boys 10+'),
       (4, 'Boys <10'),
       (5, 'Children 10+'),
       (6, 'Children <10'),
       (7, 'Man <30'),
       (8, 'Man 30+'),
       (9, 'Man 60+'),
       (10, 'Woman <30'),
       (11, 'Woman 30+'),
       (12, 'Woman 60+');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "segment";
-- +goose StatementEnd
