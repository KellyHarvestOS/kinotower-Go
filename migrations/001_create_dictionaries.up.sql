CREATE TABLE IF NOT EXISTS genders (
    id SERIAL PRIMARY KEY,
    name varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    name varchar(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name varchar(150) NOT NULL,
    parent_id int NULL REFERENCES categories(id),
    deleted_at TIMESTAMPTZ NULL
);
