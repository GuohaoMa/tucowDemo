package model

import (
	"database/sql"
	"encoding/xml"

	"github.com/lib/pq"
)

type Graph struct {
	Db       *sql.DB
	XMLName  xml.Name `xml:"graph"`
	Id       int
	Identity string `xml:"id"`
	Name     string `xml:"name"`
	Nodes    []Node `xml:"nodes>node"`
	Edges    []Edge `xml:"edges>node"`
}

func (g *Graph) Create() error {
	err := g.Db.QueryRow("insert into graph (identity, name) values ($1, $2) returning id", g.Identity, g.Name).Scan(&g.Id)
	if err != nil {
		return err
	}
	if len(g.Nodes) > 0 {
		for i, n := range g.Nodes {
			err := g.Db.QueryRow("insert into node (identity, name, graph_id) values ($1, $2, $3) returning id", n.Identity, n.Name, g.Id).Scan(&g.Nodes[i].Id)
			if err != nil {
				return err
			}
		}
	}
	if len(g.Edges) > 0 {
		for i, e := range g.Edges {
			var fromNodeId, toNodeId int
			var fromNodeIdentity, toNodeIdentity string
			err := g.Db.QueryRow("select id, identity from node where identity = $1 and graph_id = $2", e.FromIdentity, g.Id).Scan(&fromNodeId, &fromNodeIdentity)
			if err != nil {
				return err
			}
			err = g.Db.QueryRow("select id, identity from node where identity = $1 and graph_id = $2", e.ToIdentity, g.Id).Scan(&toNodeId, &toNodeIdentity)
			if err != nil {
				return err
			}
			err = g.Db.QueryRow("insert into edge (identity, from_id, from_identity, to_id, to_identity, cost, graph_id) values ($1, $2, $3, $4, $5, $6, $7) returning id", e.Identity, fromNodeId, fromNodeIdentity, toNodeId, toNodeIdentity, e.Cost, g.Id).Scan(&g.Edges[i].Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Graph) Get() error {
	err := g.Db.QueryRow("select id, identity, name from graph where id = $1", g.Id).Scan(&g.Id, &g.Identity, &g.Name)
	if err != nil {
		return err
	}
	rows, err := g.Db.Query("select id, identity, name from node where graph_id = $1", g.Id)
	if err != nil {
		return err
	}
	for rows.Next() {
		n := Node{}
		err = rows.Scan(&n.Id, &n.Identity, &n.Name)
		if err != nil {
			return err
		}
		g.Nodes = append(g.Nodes, n)
	}
	rows.Close()

	r, err := g.Db.Query("select id, identity, from_id, from_identity, to_id, to_identity, cost from edge where graph_id = $1", g.Id)
	if err != nil {
		return err
	}
	for r.Next() {
		e := Edge{}
		err = r.Scan(&e.Id, &e.Identity, &e.FromId, &e.FromIdentity, &e.ToId, &e.ToIdentity, &e.Cost)
		if err != nil {
			return err
		}
		g.Edges = append(g.Edges, e)
	}
	r.Close()
	return nil
}

func (g *Graph) FindCycles() ([][]string, error) {
	result := [][]string{}
	rows, err := g.Db.Query(`WITH RECURSIVE cte AS (
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
	WHERE has_cycle = 1;`,
		g.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := []string{}
		err = rows.Scan(pq.Array(&r))
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	rows.Close()

	return result, nil
}
