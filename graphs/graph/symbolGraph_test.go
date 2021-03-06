package graph

import (
	"testing"
)

func TestNewSymbolGraph(t *testing.T) {
	sg, err := NewSymbolGraph("./graph.txt")
	if err != nil {
		panic(err)
	}
	t.Log(sg.NumV())
	t.Log(sg.NumEdge())
	if b, _ := sg.HasEdge("A", "E"); b {
		t.Fatal()
	}
	if err := sg.AddEdge("A", "E"); err != nil {
		t.Fatal(err)
	}
	t.Log(sg.NumEdge())
	if b, _ := sg.HasEdge("A", "E"); !b {
		t.Fatal()
	}
	adj, _ := sg.Adjacent("A")
	t.Log("A -> ", adj)
}

func TestNewGraph(t *testing.T) {
	g := NewGraph(7)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 5)
	g.AddEdge(2, 1)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	bfs, err := NewBFS(g, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bfs.IsMarked(3) {
		t.Fatal()
	}
	if bfs.IsMarked(6) {
		t.Fatal()
	}
	t.Log("shortest path of 0->3: ", bfs.ShortestPathTo(3))
	cd := NewConnectivity(g)
	if b, _ := cd.IsConnected(3, 6); b {
		t.Fatal()
	}
	if b, _ := cd.IsConnected(3, 1); !b {
		t.Fatal()
	}
	t.Log("number of subGraph:", cd.NumSubGraph())
	id, _ := cd.SubGraphIDOf(1)
	t.Logf("subGraph ID of %v: %v", 1, id)
	id, _ = cd.SubGraphIDOf(6)
	t.Logf("subGraph ID of %v: %v", 6, id)

	if !HasCycle(g) {
		t.Fatal()
	}

	if IsBipartiteGraph(g) {
		t.Fatal()
	}
}
