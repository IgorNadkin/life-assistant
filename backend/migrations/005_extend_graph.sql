-- Add new columns to graph_node
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS scenario_id BIGINT;
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS action_type VARCHAR(255);
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS organization VARCHAR(255);
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS deadline VARCHAR(255);
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS reference_links TEXT;
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS consequences TEXT;
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS "order" INT DEFAULT 0;
ALTER TABLE graph_node ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();

-- Add new columns to graph_edge
ALTER TABLE graph_edge ADD COLUMN IF NOT EXISTS logic VARCHAR(255);

-- Add category to scenario
ALTER TABLE scenario ADD COLUMN IF NOT EXISTS category VARCHAR(255);
