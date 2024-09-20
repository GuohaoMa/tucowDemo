package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGraph_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	graph := &Graph{
		Db:       db,
		Identity: "graph-1",
		Name:     "Test Graph",
		Nodes: []Node{
			{Identity: "node-1", Name: "Node 1"},
			{Identity: "node-2", Name: "Node 2"},
		},
		Edges: []Edge{
			{Identity: "edge-1", FromIdentity: "node-1", ToIdentity: "node-2", Cost: 1.0},
		},
	}

	mock.ExpectQuery("insert into graph").
		WithArgs(graph.Identity, graph.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	for i, node := range graph.Nodes {
		mock.ExpectQuery("insert into node").
			WithArgs(node.Identity, node.Name, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
	}

	mock.ExpectQuery("select id, identity from node where identity = \\$1 and graph_id = \\$2").
		WithArgs("node-1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identity"}).AddRow(1, "node-1"))

	mock.ExpectQuery("select id, identity from node where identity = \\$1 and graph_id = \\$2").
		WithArgs("node-2", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identity"}).AddRow(2, "node-2"))

	mock.ExpectQuery("insert into edge").
		WithArgs("edge-1", 1, "node-1", 2, "node-2", 1.0, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = graph.Create()
	assert.NoError(t, err)
	assert.Equal(t, 1, graph.Id)
	assert.Equal(t, 1, graph.Nodes[0].Id)
	assert.Equal(t, 2, graph.Nodes[1].Id)
	assert.Equal(t, 1, graph.Edges[0].Id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGraph_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	graph := &Graph{
		Db: db,
		Id: 1,
	}

	mock.ExpectQuery("select id, identity, name from graph where id = \\$1").
		WithArgs(graph.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identity", "name"}).AddRow(1, "graph-1", "Test Graph"))

	mock.ExpectQuery("select id, identity, name from node where graph_id = \\$1").
		WithArgs(graph.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identity", "name"}).
			AddRow(1, "node-1", "Node 1").
			AddRow(2, "node-2", "Node 2"))

	mock.ExpectQuery("select id, identity, from_id, from_identity, to_id, to_identity, cost from edge where graph_id = \\$1").
		WithArgs(graph.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identity", "from_id", "from_identity", "to_id", "to_identity", "cost"}).
			AddRow(1, "edge-1", 1, "node-1", 2, "node-2", 1.0))

	err = graph.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, graph.Id)
	assert.Equal(t, "graph-1", graph.Identity)
	assert.Equal(t, "Test Graph", graph.Name)
	assert.Len(t, graph.Nodes, 2)
	assert.Equal(t, 1, graph.Nodes[0].Id)
	assert.Equal(t, "node-1", graph.Nodes[0].Identity)
	assert.Equal(t, "Node 1", graph.Nodes[0].Name)
	assert.Equal(t, 2, graph.Nodes[1].Id)
	assert.Equal(t, "node-2", graph.Nodes[1].Identity)
	assert.Equal(t, "Node 2", graph.Nodes[1].Name)
	assert.Len(t, graph.Edges, 1)
	assert.Equal(t, 1, graph.Edges[0].Id)
	assert.Equal(t, "edge-1", graph.Edges[0].Identity)
	assert.Equal(t, 1, graph.Edges[0].FromId)
	assert.Equal(t, "node-1", graph.Edges[0].FromIdentity)
	assert.Equal(t, 2, graph.Edges[0].ToId)
	assert.Equal(t, "node-2", graph.Edges[0].ToIdentity)
	assert.Equal(t, 1.0, graph.Edges[0].Cost)

	assert.NoError(t, mock.ExpectationsWereMet())
}
