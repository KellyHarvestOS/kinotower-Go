INSERT INTO genders (id, name) VALUES (1, 'Мужской'), (2, 'Женский') ON CONFLICT (id) DO NOTHING;

INSERT INTO countries (id, name) VALUES
    (1, 'США'),
    (2, 'Великобритания'),
    (3, 'Южная Корея'),
    (4, 'Япония')
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories (id, name, parent_id) VALUES
    (1, 'Фантастика', NULL),
    (2, 'Драма', NULL),
    (3, 'Триллер', NULL),
    (4, 'Киберпанк', 1),
    (5, 'Детектив', NULL)
ON CONFLICT (id) DO NOTHING;
