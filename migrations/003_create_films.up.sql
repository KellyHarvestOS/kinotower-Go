CREATE TABLE IF NOT EXISTS films (
    id SERIAL PRIMARY KEY,
    name varchar(150) NOT NULL,
    country_id int NOT NULL REFERENCES countries(id),
    duration int NOT NULL CHECK (duration >= 0),
    year_of_issue int NOT NULL CHECK (year_of_issue >= 1888),
    age int NOT NULL CHECK (age >= 0),
    link_img varchar(255) NULL,
    link_kinopoisk varchar(255) NULL,
    link_video varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS categories_films (
    id SERIAL PRIMARY KEY,
    category_id int NOT NULL REFERENCES categories(id),
    film_id int NOT NULL REFERENCES films(id),
    UNIQUE(category_id, film_id)
);

CREATE INDEX IF NOT EXISTS idx_films_name ON films(name);
CREATE INDEX IF NOT EXISTS idx_films_country ON films(country_id);
