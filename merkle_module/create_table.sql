CREATE TABLE IF NOT EXISTS merkle_trees (
    id SERIAL PRIMARY KEY,
    tree_id INT,
    issuer_did VARCHAR(255) NOT NULL,
    node_count INT NOT NULL DEFAULT 0,
    CONSTRAINT unique_tree_per_issuer UNIQUE (tree_id, issuer_did)
);

CREATE TABLE IF NOT EXISTS merkle_nodes (
    id SERIAL PRIMARY KEY,
    tree_id INT NOT NULL DEFAULT -1,
    node_id INT NOT NULL DEFAULT 0,
    data BYTEA NOT NULL,
    issuer_did VARCHAR(255) NOT null,
    CONSTRAINT unique_data_per_issuer UNIQUE (data, issuer_did)
);