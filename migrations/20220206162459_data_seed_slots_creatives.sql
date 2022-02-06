-- +goose Up
-- +goose StatementBegin
-- show all cars creatives on drom.ru homepage
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (1, 5),
       (1, 6),
       (1, 7),
       (1, 8);

-- show all except cars creatives on ozon.ru homepage
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (2, 1),
       (2, 2),
       (2, 3),
       (2, 4),
       (2, 9),
       (2, 10),
       (2, 11);

-- show lego creatives on ozon.ru toys
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (3, 1),
       (3, 2),
       (3, 3),
       (3, 4);

-- show perfume creatives on letu.ru homepage
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (4, 9),
       (4, 10),
       (4, 11);

-- show perfume and cars creatives on letu.ru man parfime
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (6, 9),
       (6, 10),
       (6, 11),
       (6, 5),
       (6, 6),
       (6, 7),
       (6, 8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE "slot_creative";
-- +goose StatementEnd
