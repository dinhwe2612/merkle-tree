CREATE TABLE IF NOT EXISTS merkle_tree (
    id SERIAL PRIMARY KEY,
    issuer_did VARCHAR(255) NOT NULL,
    tree_count INT NOT NULL
);

CREATE TABLE IF NOT EXISTS merkle_node (
    id SERIAL PRIMARY KEY,
    node_id INT NOT NULL,
    value VARCHAR(255) NOT NULL,
    tree_id INT NOT NULL,
    FOREIGN KEY (tree_id) REFERENCES merkle_tree(id)
);
