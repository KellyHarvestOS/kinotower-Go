INSERT INTO films (id, name, country_id, duration, year_of_issue, age, link_img, link_kinopoisk, link_video, created_at) VALUES
    (1, 'Black Mirror', 2, 150, 2011, 18, NULL, NULL, 'https://www.imdb.com/title/tt2085059/', now()),
    (2, 'Blade Runner 2049', 1, 164, 2017, 18, NULL, NULL, 'https://www.imdb.com/title/tt1856101/', now()),
    (3, 'Parasite', 3, 132, 2019, 18, NULL, NULL, 'https://www.imdb.com/title/tt6751668/', now())
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories_films (category_id, film_id) VALUES
    (1, 1), (3, 1), (1, 2), (4, 2), (2, 3), (3, 3)
ON CONFLICT (category_id, film_id) DO NOTHING;
