CREATE TABLE IF NOT EXISTS scenario (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    start_node_id INT NOT NULL
);