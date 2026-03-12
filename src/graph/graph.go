package graph

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	errorspkg "navigator/internal/errors"
)

// Edge represents a directed edge from the current vertex to To.
// Vertex numbering is 1-based.
type Edge struct {
	To     int
	Weight int
}

// Graph stores a graph using an adjacency matrix.
// A weight of 0 means no edge.
type Graph struct {
	adj [][]int
}

// LoadGraphFromFile loads a graph from a file in adjacency-matrix format.
// File format:
//
//	n
//	a11 a12 ... a1n
//	...
//	an1 an2 ... ann
//
// Vertex numbering is 1-based in algorithms, but file stores raw matrix.
func (g *Graph) LoadGraphFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var tokens []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		tokens = append(tokens, fields...)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(tokens) == 0 {
		return errorspkg.ErrInvalidGraphFormat
	}

	n, err := strconv.Atoi(tokens[0])
	if err != nil || n <= 0 {
		return errorspkg.ErrInvalidGraphFormat
	}
	expected := 1 + n*n
	if len(tokens) != expected {
		return errorspkg.ErrInvalidGraphFormat
	}

	matrix := make([][]int, n)
	idx := 1
	for i := range n {
		matrix[i] = make([]int, n)
		for j := range n {
			v, convErr := strconv.Atoi(tokens[idx])
			idx++
			if convErr != nil || v < 0 {
				return errorspkg.ErrInvalidGraphFormat
			}
			if v != 0 {
				// Natural numbers only.
				if v < 1 {
					return errorspkg.ErrInvalidWeight
				}
			}
			matrix[i][j] = v
		}
	}

	if !hasAnyEdge(matrix) {
		return errorspkg.ErrEmptyGraph
	}
	if !isConnectedUndirected(matrix) {
		return errorspkg.ErrGraphNotConnected
	}

	g.adj = matrix
	return nil
}

// ExportGraphToDot exports the current graph to a Graphviz dot file.
func (g *Graph) ExportGraphToDot(filename string) error {
	if g.adj == nil {
		return errorspkg.ErrGraphNotLoaded
	}

	directed := g.IsDirected()
	var b strings.Builder
	if directed {
		b.WriteString("digraph G {\n")
	} else {
		b.WriteString("graph G {\n")
	}

	n := len(g.adj)
	for i := 1; i <= n; i++ {
		b.WriteString(fmt.Sprintf("  %d;\n", i))
	}

	connector := " -- "
	if directed {
		connector = " -> "
	}

	for i := range n {
		for j := range n {
			w := g.adj[i][j]
			if w == 0 {
				continue
			}
			if !directed {
				if j < i {
					continue
				}
			}
			from := i + 1
			to := j + 1
			b.WriteString(fmt.Sprintf("  %d%s%d [label=\"%d\"];\n", from, connector, to, w))
		}
	}
	b.WriteString("}\n")

	return os.WriteFile(filename, []byte(b.String()), 0o644)
}

// VerticesCount returns the number of vertices in the graph.
func (g *Graph) VerticesCount() int {
	return len(g.adj)
}

// IsDirected reports whether the graph is directed.
// The graph is considered directed if its adjacency matrix is not symmetric.
func (g *Graph) IsDirected() bool {
	if g.adj == nil {
		return false
	}
	n := len(g.adj)
	for i := range n {
		for j := i + 1; j < n; j++ {
			if g.adj[i][j] != g.adj[j][i] {
				return true
			}
		}
	}
	return false
}

// IsUndirected reports whether the graph is undirected.
func (g *Graph) IsUndirected() bool {
	return !g.IsDirected()
}

// Weight returns the weight of the edge from -> to (0 if absent).
func (g *Graph) Weight(from, to int) (int, error) {
	if g.adj == nil {
		return 0, errorspkg.ErrGraphNotLoaded
	}
	if from < 1 || to < 1 || from > len(g.adj) || to > len(g.adj) {
		return 0, errorspkg.ErrInvalidVertex
	}
	return g.adj[from-1][to-1], nil
}

// Neighbors returns outgoing neighbors from a vertex.
func (g *Graph) Neighbors(vertex int) ([]Edge, error) {
	if g.adj == nil {
		return nil, errorspkg.ErrGraphNotLoaded
	}
	if vertex < 1 || vertex > len(g.adj) {
		return nil, errorspkg.ErrInvalidVertex
	}
	row := g.adj[vertex-1]
	neighbors := make([]Edge, 0, len(row))
	for j, w := range row {
		if w == 0 {
			continue
		}
		neighbors = append(neighbors, Edge{To: j + 1, Weight: w})
	}
	return neighbors, nil
}

// AdjacencyMatrixCopy returns a deep copy of the adjacency matrix.
// func (g *Graph) AdjacencyMatrixCopy() ([][]int, error) {
// 	if g.adj == nil {
// 		return nil, errorspkg.ErrGraphNotLoaded
// 	}
// 	copyM := make([][]int, len(g.adj))
// 	for i := range g.adj {
// 		copyM[i] = append([]int{}, g.adj[i]...)
// 	}
// 	return copyM, nil
// }

func hasAnyEdge(m [][]int) bool {
	for i := range m {
		for j := range m[i] {
			if m[i][j] != 0 {
				return true
			}
		}
	}
	return false
}

// isConnectedUndirected checks connectivity treating the graph as undirected:
// an undirected edge exists if at least one of (i->j) or (j->i) has a weight.
func isConnectedUndirected(m [][]int) bool {
	n := len(m)
	visited := make([]bool, n)
	queue := make([]int, 0, n)
	visited[0] = true
	queue = append(queue, 0)

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for u := range n {
			if visited[u] {
				continue
			}
			if m[v][u] != 0 || m[u][v] != 0 {
				visited[u] = true
				queue = append(queue, u)
			}
		}
	}

	for _, ok := range visited {
		if !ok {
			return false
		}
	}
	return true
}
