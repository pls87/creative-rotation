-- +goose Up
-- +goose StatementBegin
INSERT INTO "creative" ("description")
VALUES ('Lego Technic'),
       ('Lego Friends'),
       ('Lego Duplo'),
       ('Lego Architecture'),
       ('KIA Stinger'),
       ('KIA Sportage'),
       ('Mazda CX-9'),
       ('Chevrolet Tahoe'),
       ('Chanel #5'),
       ('Chanel Chance'),
       ('Dior Homme');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "creative";
-- +goose StatementEnd
