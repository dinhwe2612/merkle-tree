CREATE TABLE IF NOT EXISTS merkle_trees (
    id SERIAL PRIMARY KEY,
    issuer_did VARCHAR(255) NOT NULL,
    node_count INT NOT NULL,
    need_sync BOOLEAN NOT NULL DEFAULT true,
    node_count_sync INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS merkle_nodes (
    id SERIAL PRIMARY KEY,
    tree_id INT NOT NULL,
    node_id INT NOT NULL,
    data BYTEA NOT NULL,
    FOREIGN KEY (tree_id) REFERENCES merkle_trees(id)
);