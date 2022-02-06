-- +goose Up
-- +goose StatementBegin
INSERT INTO "creative" ("ID", "description")
VALUES (1, 'Lego Technic'),
       (2, 'Lego Friends'),
       (3, 'Lego Duplo'),
       (4, 'Lego Architecture'),
       (5, 'KIA Stinger'),
       (6, 'KIA Sportage'),
       (7, 'Mazda CX-9'),
       (8, 'Chevrolet Tahoe'),
       (9, 'Chanel #5'),
       (10, 'Chanel Chance'),
       (11, 'Dior Homme');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "creative";
-- +goose StatementEnd
