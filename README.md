# Tucow Demo Project - Guohao Ma

## Setup Instructions

1. Clone the repository:
    ```sh
    git clone https://github.com/GuohaoMa/tucowDemo.git
    ```
2. Navigate to the project directory:
    ```sh
    cd tucowDemo
    ```
3. Replace data and config if needed:
   - Test graph: default as `/data/exampleTest.xml`.
   - Note: The whole service need restart if test data and config is changed. Since there is no requirement for an api for download the xml, the service always save a brand new graph based on `/data/exampleTest.xml` in db and used it as the target graph for path finding function later on.
4. Start the service in docker on default port `8080`.
    ```sh
    make run
    ```

## Project Overview
The project builds up a service to deal with graphs, nodes, edges in XML format. It runs in docker with default port `8080`, which includes a backend go service using GIN framework and a database using PostgreSQL. 

## XML Validation
The validation rules are added in file `validation/validate.go`.

## Handler Explanation

### Request Handler
Handlers are defined in `handlers/findPathHandler.go`. For this demo project, only one end point `POST localhost:8080/graphs/paths` is registered.

**Example Request:**
```json
{
    "queries": [
        {
            "paths": {
                "start": "a",
                "end": "e"
            }
        },
        {
            "cheapest": {
                "start": "a",
                "end": "e"
            }
        },
        {
            "cheapest": {
                "start": "a",
                "end": "h"
            }
        }
    ]
}
```

**Example Response:**

```json
{
    "answers": [
        {
            "paths": {
                "from": "a",
                "to": "e",
                "paths": [
                    [
                        "a",
                        "e"
                    ],
                    [
                        "a",
                        "b",
                        "e"
                    ]
                ]
            }
        },
        {
            "cheapest": {
                "from": "a",
                "to": "e",
                "paths": [
                    "a",
                    "b",
                    "e"
                ]
            }
        },
        {
            "cheapest": {
                "from": "a",
                "to": "h",
                "paths": false
            }
        }
    ]
}
```

### Functions Explanation
**findAllPath:** 
This function finds all possible paths between the source and destination nodes if no cycles exists. It is located at `handlers/findPathHandler.go`. The core algorithms is based on recursively dfs with a stack as temp path. Keep backtracing and store the result if there is a path through edges from start node to end node.

**findCheapestPath:**
Similarly to findAllPath function, the core algorithms is based on recursively dfs with a stack as temp path. It is located at `handlers/findPathHandler.go`. Keep backtracing and compare the cost if there is a path between start node and end node. Otherwise return false.



### Database Schema
```
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
```

**SQL query for finding cycles**
The SQL query is written in SQL and used in function `FindCycles` in file `model/graph`
```
WITH RECURSIVE cte AS (
    -- Anchor member: start from each edge
		SELECT 
			from_identity, 
			to_identity, 
			',' || from_identity || ',' || to_identity || ',' AS nodes,  -- Concatenate with commas
			1 AS lev, 
			CASE WHEN from_identity = to_identity THEN 1 ELSE 0 END AS has_cycle
		FROM edge e where from_identity <> to_identity AND graph_id = $1
		
		UNION ALL
		
		-- Recursive member: find the next edge
		SELECT 
			cte.from_identity, 
			e.to_identity,
			cte.nodes || e.to_identity || ',' AS nodes,  -- Append the target to the path
			lev + 1,
			CASE WHEN cte.nodes LIKE '%' || e.to_identity || '%' THEN 1 ELSE 0 END AS has_cycle
		FROM cte 
		JOIN edge e ON e.from_identity = cte.to_identity
		WHERE cte.has_cycle = 0 AND e.from_identity <> e.to_identity AND graph_id = $1
	)
	SELECT *
	FROM cte
	WHERE has_cycle = 1;
```

### Reason for Using JSON Library
The JSON library `encoding/json` is used for parsing and generating JSON data as a pretty standard practice in GO. It supports encoding/decoding well with json tag in go struct.

### Reason for Using XML Library
Similarly, the XML library `encoding/xml` is used for parsing and generating XML data as a pretty standard practice in GO. It supports encoding/decoding well with xml tag in go struct.


## Thanks for reviewing and playing with this demo project. Have a good one!