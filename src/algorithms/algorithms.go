package algorithms

import (
	"math"
	"math/rand"
	"slices"

	"navigator/containers/queue"
	"navigator/containers/stack"
	"navigator/graph"
	errorspkg "navigator/internal/errors"
)

// GraphAlgorithms stores implementations of graph algorithms.
// It does not depend on the internal graph representation.
type GraphAlgorithms struct{}

// TsmResult is the result for the traveling salesman problem.
type TsmResult struct {
	Vertices []int
	Distance float64
}

const (
	infinityInt = int(math.MaxInt / 4)

	acoPheromoneInfluence       = 1.0
	acoDistanceInfluence        = 1.1
	acoEvaporationRate          = 0.1
	acoPheromoneDepositStrength = 100.0
	acoIterations               = 120
)

// DepthFirstSearch performs a non-recursive DFS starting from startVertex.
// Vertex numbering is 1-based.
func (GraphAlgorithms) DepthFirstSearch(g *graph.Graph, startVertex int) ([]int, error) {
	if g == nil {
		return nil, errorspkg.ErrGraphNotLoaded
	}
	n := g.VerticesCount()
	if startVertex < 1 || startVertex > n {
		return nil, errorspkg.ErrInvalidVertex
	}

	visited := make([]bool, n+1)
	result := make([]int, 0, n)
	s := stack.Stack[int]()
	s.Push(startVertex)

	for s.Len() > 0 {
		v, err := s.Pop()
		if err != nil {
			return nil, err
		}
		if visited[v] {
			continue
		}
		visited[v] = true
		result = append(result, v)

		neighbors, err := g.Neighbors(v)
		if err != nil {
			return nil, err
		}
		// Push in reverse order so traversal is deterministic (ascending).
		for i := len(neighbors) - 1; i >= 0; i-- {
			to := neighbors[i].To
			if !visited[to] {
				s.Push(to)
			}
		}
	}

	return result, nil
}

// BreadthFirstSearch performs a BFS starting from startVertex.
func (GraphAlgorithms) BreadthFirstSearch(g *graph.Graph, startVertex int) ([]int, error) {
	if g == nil {
		return nil, errorspkg.ErrGraphNotLoaded
	}
	n := g.VerticesCount()
	if startVertex < 1 || startVertex > n {
		return nil, errorspkg.ErrInvalidVertex
	}

	visited := make([]bool, n+1)
	result := make([]int, 0, n)
	q := queue.Queue[int]()
	visited[startVertex] = true
	q.Push(startVertex)

	for q.Len() > 0 {
		v, err := q.Pop()
		if err != nil {
			return nil, err
		}
		result = append(result, v)

		neighbors, err := g.Neighbors(v)
		if err != nil {
			return nil, err
		}
		for _, e := range neighbors {
			if !visited[e.To] {
				visited[e.To] = true
				q.Push(e.To)
			}
		}
	}

	return result, nil
}

// GetShortestPathBetweenVertices finds the shortest path distance using Dijkstra's algorithm.
func (GraphAlgorithms) GetShortestPathBetweenVertices(g *graph.Graph, vertex1, vertex2 int) (int, error) {
	_, dist, err := dijkstraPath(g, vertex1, vertex2)
	return dist, err
}

// GetShortestPathsBetweenAllVertices finds shortest paths between all pairs using Floyd-Warshall.
func (GraphAlgorithms) GetShortestPathsBetweenAllVertices(g *graph.Graph) ([][]int, error) {
	if g == nil {
		return nil, errorspkg.ErrGraphNotLoaded
	}
	n := g.VerticesCount()
	if n == 0 {
		return nil, errorspkg.ErrGraphNotLoaded
	}

	dist := make([][]int, n)
	for i := range n {
		dist[i] = make([]int, n)
		for j := range n {
			if i == j {
				dist[i][j] = 0
				continue
			}
			w, err := g.Weight(i+1, j+1)
			if err != nil {
				return nil, err
			}
			if w == 0 {
				dist[i][j] = infinityInt
			} else {
				dist[i][j] = w
			}
		}
	}

	for k := range n {
		for i := range n {
			if dist[i][k] == infinityInt {
				continue
			}
			for j := range n {
				if dist[k][j] == infinityInt {
					continue
				}
				alt := dist[i][k] + dist[k][j]
				if alt < dist[i][j] {
					dist[i][j] = alt
				}
			}
		}
	}

	return dist, nil
}

// GetLeastSpanningTree finds the least spanning tree using Prim's algorithm.
// Returns an adjacency matrix for the LST.
func (GraphAlgorithms) GetLeastSpanningTree(g *graph.Graph) ([][]int, error) {
	if g == nil {
		return nil, errorspkg.ErrGraphNotLoaded
	}
	if !g.IsUndirected() {
		return nil, errorspkg.ErrNotUndirected
	}
	n := g.VerticesCount()
	if n == 0 {
		return nil, errorspkg.ErrGraphNotLoaded
	}

	lst := make([][]int, n)
	for i := range lst {
		lst[i] = make([]int, n)
	}

	used := make([]bool, n)
	minEdge := make([]int, n)
	selectedFrom := make([]int, n)
	for i := range minEdge {
		minEdge[i] = infinityInt
		selectedFrom[i] = -1
	}
	minEdge[0] = 0

	for range n {
		v := -1
		for i := range n {
			if used[i] {
				continue
			}
			if v == -1 || minEdge[i] < minEdge[v] {
				v = i
			}
		}
		if v == -1 || minEdge[v] == infinityInt {
			return nil, errorspkg.ErrGraphNotConnected
		}
		used[v] = true
		if selectedFrom[v] != -1 {
			from := selectedFrom[v]
			w, err := g.Weight(from+1, v+1)
			if err != nil {
				return nil, err
			}
			lst[from][v] = w
			lst[v][from] = w
		}

		for to := range n {
			if used[to] || to == v {
				continue
			}
			w, err := g.Weight(v+1, to+1)
			if err != nil {
				return nil, err
			}
			if w == 0 {
				continue
			}
			if w < minEdge[to] {
				minEdge[to] = w
				selectedFrom[to] = v
			}
		}
	}

	return lst, nil
}

// SolveTravelingSalesmanProblem solves TSP using a simple ant colony algorithm.
// The route visits every vertex at least once and returns to the start vertex.
func (GraphAlgorithms) SolveTravelingSalesmanProblem(g *graph.Graph) (TsmResult, error) {
	if g == nil {
		return TsmResult{}, errorspkg.ErrGraphNotLoaded
	}
	n := g.VerticesCount()
	if n == 0 {
		return TsmResult{}, errorspkg.ErrGraphNotLoaded
	}

	pheromones := make([][]float64, n)
	for i := range pheromones {
		pheromones[i] = make([]float64, n)
		for j := range pheromones[i] {
			pheromones[i][j] = 1.0
		}
	}

	rng := rand.New(rand.NewSource(1))
	ants := n
	best := TsmResult{Distance: math.Inf(1)}

	for range acoIterations {
		iterationBest := TsmResult{Distance: math.Inf(1)}

		for range ants {
			start := 1
			path, dist, ok := constructTour(rng, g, pheromones, start)
			if !ok {
				continue
			}
			res := TsmResult{Vertices: path, Distance: dist}
			if res.Distance < iterationBest.Distance {
				iterationBest = res
			}
			if res.Distance < best.Distance {
				best = res
			}
		}

		evaporate(pheromones)
		if !math.IsInf(iterationBest.Distance, 1) {
			deposit(pheromones, iterationBest.Vertices, iterationBest.Distance, g.IsUndirected())
		}
		if !math.IsInf(best.Distance, 1) {
			deposit(pheromones, best.Vertices, best.Distance, g.IsUndirected())
		}
	}

	if math.IsInf(best.Distance, 1) || len(best.Vertices) == 0 {
		return TsmResult{}, errorspkg.ErrTSPUnsolvable
	}
	return best, nil
}

func constructTour(rng *rand.Rand, g *graph.Graph, pheromones [][]float64, start int) ([]int, float64, bool) {
	n := g.VerticesCount()
	visited := make([]bool, n+1)
	visited[start] = true
	path := make([]int, 0, n+1)
	path = append(path, start)
	current := start
	distance := 0.0

	for len(path) < n {
		next, w, ok := chooseNext(rng, g, pheromones, current, visited)
		if !ok {
			return nil, 0, false
		}
		visited[next] = true
		path = append(path, next)
		distance += float64(w)
		current = next
	}

	// Return to start by the shortest path to satisfy "at least once" even when there's no direct edge.
	backPath, backDist, err := dijkstraPath(g, current, start)
	if err != nil {
		return nil, 0, false
	}
	if len(backPath) < 2 {
		return nil, 0, false
	}
	// backPath includes current and start.
	for i := 1; i < len(backPath); i++ {
		path = append(path, backPath[i])
	}
	distance += float64(backDist)

	return path, distance, true
}

func chooseNext(rng *rand.Rand, g *graph.Graph, pheromones [][]float64, current int, visited []bool) (int, int, bool) {
	neighbors, err := g.Neighbors(current)
	if err != nil {
		return 0, 0, false
	}

	type candidate struct {
		to             int
		weight         int
		selectionScore float64
	}
	cands := make([]candidate, 0, len(neighbors))

	totalScore := 0.0
	for _, e := range neighbors {
		if visited[e.To] {
			continue
		}
		edgeDesirability := 1.0 / float64(e.Weight)
		selectionScore := math.Pow(pheromones[current-1][e.To-1], acoPheromoneInfluence) * math.Pow(edgeDesirability, acoDistanceInfluence)
		cands = append(cands, candidate{to: e.To, weight: e.Weight, selectionScore: selectionScore})
		totalScore += selectionScore
	}
	if len(cands) == 0 {
		return 0, 0, false
	}
	if totalScore == 0 {
		// Fallback: choose the smallest weight.
		best := cands[0]
		for _, cand := range cands[1:] {
			if cand.weight < best.weight {
				best = cand
			}
		}
		return best.to, best.weight, true
	}

	targetScore := rng.Float64() * totalScore
	cumulativeScore := 0.0
	for _, cand := range cands {
		cumulativeScore += cand.selectionScore
		if targetScore <= cumulativeScore {
			return cand.to, cand.weight, true
		}
	}
	last := cands[len(cands)-1]
	return last.to, last.weight, true
}

func evaporate(pheromones [][]float64) {
	keep := 1.0 - acoEvaporationRate
	for i := range pheromones {
		for j := range pheromones[i] {
			pheromones[i][j] *= keep
			if pheromones[i][j] < 1e-12 {
				pheromones[i][j] = 1e-12
			}
		}
	}
}

func deposit(pheromones [][]float64, route []int, distance float64, undirected bool) {
	if distance <= 0 || len(route) < 2 {
		return
	}
	delta := acoPheromoneDepositStrength / distance
	for i := 0; i < len(route)-1; i++ {
		from := route[i] - 1
		to := route[i+1] - 1
		pheromones[from][to] += delta
		if undirected {
			pheromones[to][from] += delta
		}
	}
}

// dijkstraPath returns a shortest path (as vertices) and its distance.
func dijkstraPath(g *graph.Graph, start, target int) ([]int, int, error) {
	if g == nil {
		return nil, 0, errorspkg.ErrGraphNotLoaded
	}
	n := g.VerticesCount()
	if start < 1 || target < 1 || start > n || target > n {
		return nil, 0, errorspkg.ErrInvalidVertex
	}

	dist := make([]int, n+1)
	prev := make([]int, n+1)
	used := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = infinityInt
		prev[i] = -1
	}
	dist[start] = 0

	for {
		v := -1
		for i := 1; i <= n; i++ {
			if used[i] {
				continue
			}
			if v == -1 || dist[i] < dist[v] {
				v = i
			}
		}
		if v == -1 || dist[v] == infinityInt {
			break
		}
		if v == target {
			break
		}
		used[v] = true

		neighbors, err := g.Neighbors(v)
		if err != nil {
			return nil, 0, err
		}
		for _, e := range neighbors {
			if dist[v] > infinityInt-e.Weight {
				continue
			}
			if dist[v]+e.Weight < dist[e.To] {
				dist[e.To] = dist[v] + e.Weight
				prev[e.To] = v
			}
		}
	}

	if dist[target] == infinityInt {
		return nil, 0, errorspkg.ErrNoPath
	}

	path := make([]int, 0)
	for v := target; v != -1; v = prev[v] {
		path = append(path, v)
	}
	slices.Reverse(path)
	return path, dist[target], nil
}
