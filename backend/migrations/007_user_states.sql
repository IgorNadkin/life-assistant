CREATE TABLE IF NOT EXISTS user_states (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scenario_id BIGINT NOT NULL REFERENCES scenario(id) ON DELETE CASCADE,
    current_node_id BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'in_progress',
    completed_steps BIGINT[] DEFAULT ARRAY[]::BIGINT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_states_user_id ON user_states(user_id);
CREATE INDEX IF NOT EXISTS idx_user_states_scenario_id ON user_states(scenario_id);
CREATE INDEX IF NOT EXISTS idx_user_states_status ON user_states(status);
