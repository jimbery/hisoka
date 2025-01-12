CREATE TABLE IF NOT EXISTS anime (
    id SERIAL PRIMARY KEY,
    mal_id INT,
    name VARCHAR(255),
    dub_vote INT,
    sub_vote INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
