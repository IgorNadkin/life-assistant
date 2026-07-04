CREATE TABLE IF NOT EXISTS graph_node (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS graph_edge (
    id SERIAL PRIMARY KEY,
    from_node INT NOT NULL,
    to_node INT NOT NULL,
    condition TEXT NOT NULL
);