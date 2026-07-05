CREATE TABLE IF NOT EXISTS user_actions (
    id BIGSERIAL PRIMARY KEY,
    user_state_id BIGINT NOT NULL REFERENCES user_states(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL,
    action VARCHAR(255) NOT NULL,
    organization VARCHAR(255) NOT NULL,
    deadline TIMESTAMP,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_actions_user_state_id ON user_actions(user_state_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_completed ON user_actions(completed);
