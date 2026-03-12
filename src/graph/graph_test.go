package graph

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	errorspkg "navigator/internal/errors"
)

func TestLoadGraphFromFile_OK(t *testing.T) {
	var g Graph
	path := filepath.Join("..", "testdata", "graph_undirected_4.txt")
	if err := g.LoadGraphFromFile(path); err != nil {
		t.Fatalf("LoadGraphFromFile() error = %v", err)
	}
	if got := g.VerticesCount(); got != 4 {
		t.Fatalf("VerticesCount() = %d, want 4", got)
	}
	w, err := g.Weight(2, 3)
	if err != nil {
		t.Fatalf("Weight() error = %v", err)
	}
	if w != 2 {
		t.Fatalf("Weight(2,3) = %d, want 2", w)
	}
}

func TestLoadGraphFromFile_Disconnected(t *testing.T) {
	var g Graph
	path := filepath.Join("..", "testdata", "graph_disconnected.txt")
	if err := g.LoadGraphFromFile(path); err == nil {
		t.Fatalf("expected error")
	} else if err != errorspkg.ErrGraphNotConnected {
		t.Fatalf("error = %v, want %v", err, errorspkg.ErrGraphNotConnected)
	}
}

func TestExportGraphToDot(t *testing.T) {
	var g Graph
	path := filepath.Join("..", "testdata", "graph_directed_4.txt")
	if err := g.LoadGraphFromFile(path); err != nil {
		t.Fatalf("LoadGraphFromFile() error = %v", err)
	}
	tmp := t.TempDir()
	out := filepath.Join(tmp, "graph.dot")
	if err := g.ExportGraphToDot(out); err != nil {
		t.Fatalf("ExportGraphToDot() error = %v", err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(b)
	if !strings.Contains(text, "digraph G") {
		t.Fatalf("dot header missing")
	}
	if !strings.Contains(text, "1 -> 2") {
		t.Fatalf("expected edge missing")
	}
}

func TestNeighbors_InvalidVertex(t *testing.T) {
	var g Graph
	path := filepath.Join("..", "testdata", "graph_undirected_4.txt")
	if err := g.LoadGraphFromFile(path); err != nil {
		t.Fatalf("LoadGraphFromFile() error = %v", err)
	}
	if _, err := g.Neighbors(0); err != errorspkg.ErrInvalidVertex {
		t.Fatalf("error = %v, want %v", err, errorspkg.ErrInvalidVertex)
	}
}
