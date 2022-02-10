INSERT INTO "creative" ("description")
VALUES ('Lego Technic'),    --1
       ('Lego Friends'),    --2
       ('KIA Soul'),        --3
       ('Chevrolet Tahoe'), --4
       ('Chanel Chance'),   --5
       ('Dior Homme'); --6
INSERT INTO "segment" ("description")
VALUES ('Girl'), --1
       ('Boy'),  --2
       ('Man'),  --3
       ('Woman'); --4
INSERT INTO "slot" ("description")
VALUES ('drom.ru'),--1
       ('ozon.ru'),--2
       ('toys.ru'),--3
       ('letu.ru');
--4
-- show cars creatives on drom.ru homepage
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (1, 3),
       (1, 4);
-- show all except cars creatives on ozon.ru homepage
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (2, 1),
       (2, 2),
       (2, 5),
       (2, 6);
-- show lego creatives on toys.ru
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (3, 1),
       (3, 2);
-- show perfume and car creatives on letu.ru
INSERT INTO "slot_creative" (slot_id, creative_id)
VALUES (4, 3),
       (4, 4),
       (4, 5),
       (4, 6);