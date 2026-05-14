CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    film_id int NOT NULL REFERENCES films(id),
    user_id int NOT NULL REFERENCES users(id),
    message text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_reviews_film ON reviews(film_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
