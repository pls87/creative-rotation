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
INSERT INTO "slot" ("description")
VALUES ('drom.ru - homepage'),
       ('ozon.ru - homepage'),
       ('ozon.ru - toys'),
       ('letu.ru - homepage'),
       ('letu.ru - man perfume');
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
VALUES (5, 9),
       (5, 10),
       (5, 11),
       (5, 5),
       (5, 6),
       (5, 7),
       (5, 8);