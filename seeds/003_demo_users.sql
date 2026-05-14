INSERT INTO users (id, fio, birthday, gender_id, email, password, created_at) VALUES
    (1, 'Ivanov Ivan', '2006-04-15', 1, 'ivanov@ivan.kz', '$2a$10$mFoA8XtFgSFjoN25HuFEXuV5jbZ8hIhB0HdV0UAN2rbppmtWw5xoK', now())
ON CONFLICT (id) DO NOTHING;

INSERT INTO reviews (film_id, user_id, message, is_approved, created_at) VALUES
    (1, 1, 'Очень атмосферный сериал.', TRUE, now()),
    (2, 1, 'Визуально мощный фильм.', TRUE, now())
ON CONFLICT DO NOTHING;

INSERT INTO ratings (film_id, user_id, ball, created_at) VALUES
    (1, 1, 5, now()),
    (2, 1, 5, now()),
    (3, 1, 4, now())
ON CONFLICT (film_id, user_id) DO NOTHING;
