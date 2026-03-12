package algorithms

import (
	"path/filepath"
	"reflect"
	"testing"

	"navigator/graph"
)

func loadUndirected4(t *testing.T) *graph.Graph {
	t.Helper()
	var g graph.Graph
	path := filepath.Join("..", "testdata", "graph_undirected_4.txt")
	if err := g.LoadGraphFromFile(path); err != nil {
		t.Fatalf("LoadGraphFromFile() error = %v", err)
	}
	return &g
}

func TestDFS(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	got, err := ga.DepthFirstSearch(g, 1)
	if err != nil {
		t.Fatalf("DepthFirstSearch() error = %v", err)
	}
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DFS = %v, want %v", got, want)
	}
}

func TestBFS(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	got, err := ga.BreadthFirstSearch(g, 1)
	if err != nil {
		t.Fatalf("BreadthFirstSearch() error = %v", err)
	}
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BFS = %v, want %v", got, want)
	}
}

func TestDijkstra(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	dist, err := ga.GetShortestPathBetweenVertices(g, 1, 4)
	if err != nil {
		t.Fatalf("GetShortestPathBetweenVertices() error = %v", err)
	}
	if dist != 6 {
		t.Fatalf("dist = %d, want 6", dist)
	}
}

func TestFloydWarshall(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	m, err := ga.GetShortestPathsBetweenAllVertices(g)
	if err != nil {
		t.Fatalf("GetShortestPathsBetweenAllVertices() error = %v", err)
	}
	want := [][]int{
		{0, 1, 3, 6},
		{1, 0, 2, 5},
		{3, 2, 0, 3},
		{6, 5, 3, 0},
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("matrix = %v, want %v", m, want)
	}
}

func TestPrim(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	mst, err := ga.GetLeastSpanningTree(g)
	if err != nil {
		t.Fatalf("GetLeastSpanningTree() error = %v", err)
	}
	want := [][]int{
		{0, 1, 0, 0},
		{1, 0, 2, 0},
		{0, 2, 0, 3},
		{0, 0, 3, 0},
	}
	if !reflect.DeepEqual(mst, want) {
		t.Fatalf("mst = %v, want %v", mst, want)
	}
}

func TestTSP(t *testing.T) {
	g := loadUndirected4(t)
	ga := GraphAlgorithms{}
	res, err := ga.SolveTravelingSalesmanProblem(g)
	if err != nil {
		t.Fatalf("SolveTravelingSalesmanProblem() error = %v", err)
	}
	if res.Distance != 12 {
		t.Fatalf("distance = %v, want 12", res.Distance)
	}
	wantStartEnd := 1
	if len(res.Vertices) < 2 || res.Vertices[0] != wantStartEnd || res.Vertices[len(res.Vertices)-1] != wantStartEnd {
		t.Fatalf("route must start and end at 1, got %v", res.Vertices)
	}
	// Must contain all vertices at least once.
	seen := make(map[int]bool)
	for _, v := range res.Vertices {
		seen[v] = true
	}
	for v := 1; v <= g.VerticesCount(); v++ {
		if !seen[v] {
			t.Fatalf("vertex %d missing from route %v", v, res.Vertices)
		}
	}
}
