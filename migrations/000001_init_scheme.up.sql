CREATE TABLE IF NOT EXISTS graph (
    id serial PRIMARY KEY, -- Primary key for the graph table
    identity varchar, -- Identity of the graph
    name varchar -- Name of the graph
);
CREATE TABLE IF NOT EXISTS node (
    id serial PRIMARY KEY, -- Primary key for the node table
    identity varchar NOT NULL, -- Identity of the node
    name varchar NOT NULL, -- Name of the node
    graph_id integer NOT NULL, -- Foreign key referencing the graph table
    FOREIGN KEY(graph_id) REFERENCES graph(id), -- Relationship to the graph table
    CONSTRAINT node_key UNIQUE (identity, graph_id) -- Unique constraint on identity and graph_id
);
CREATE TABLE IF NOT EXISTS edge (
    id serial PRIMARY KEY, -- Primary key for the edge table
    identity varchar, -- Identity of the edge
    from_id integer NOT NULL, -- Foreign key referencing the node table (start node)
    from_identity varchar NOT NULL, -- Identity of the start node
    to_id integer NOT NULL, -- Foreign key referencing the node table (end node)
    to_identity varchar NOT NULL, -- Identity of the end node
    cost numeric(10, 2) NOT NULL, -- Cost associated with the edge
    graph_id integer NOT NULL, -- Foreign key referencing the graph table
    FOREIGN KEY(from_id) REFERENCES node(id), -- Relationship to the node table (start node)
    FOREIGN KEY(to_id) REFERENCES node(id), -- Relationship to the node table (end node)
    FOREIGN KEY(graph_id) REFERENCES graph(id), -- Relationship to the graph table
    CONSTRAINT edge_key UNIQUE (identity, graph_id), -- Unique constraint on identity and graph_id
    CONSTRAINT edge_key2 UNIQUE (from_id, to_id, graph_id) -- Unique constraint on from_id, to_id, and graph_id
);  