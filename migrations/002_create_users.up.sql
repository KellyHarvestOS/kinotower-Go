CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    fio varchar(150) NOT NULL,
    birthday date NULL,
    gender_id int NOT NULL REFERENCES genders(id),
    email varchar(50) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
