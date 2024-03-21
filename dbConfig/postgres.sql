CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    full_url TEXT NOT NULL,
    short_url TEXT NOT NULL,
    UNIQUE(full_url),
    UNIQUE(short_url)
)