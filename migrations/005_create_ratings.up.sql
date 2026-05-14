CREATE TABLE IF NOT EXISTS ratings (
    id SERIAL PRIMARY KEY,
    film_id int NOT NULL REFERENCES films(id),
    user_id int NOT NULL REFERENCES users(id),
    ball int NOT NULL CHECK (ball >= 1 AND ball <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(film_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_film ON ratings(film_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user ON ratings(user_id);
