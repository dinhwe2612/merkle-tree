CREATE TABLE IF NOT EXISTS merkle_trees (
    id SERIAL PRIMARY KEY,
    issuer_did VARCHAR(255) NOT NULL,
    node_count INT NOT NULL
);

CREATE TABLE IF NOT EXISTS merkle_nodes (
    id SERIAL PRIMARY KEY,
    tree_id INT NOT NULL,
    node_id INT NOT NULL,
    value VARCHAR(255) NOT NULL,
    FOREIGN KEY (tree_id) REFERENCES merkle_trees(id),
    CONSTRAINT unique_value_per_tree UNIQUE (tree_id, value)
);
