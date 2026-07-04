CREATE TABLE IF NOT EXISTS user_state (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    current_node_id INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);