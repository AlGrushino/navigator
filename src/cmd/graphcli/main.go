package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"navigator/algorithms"
	"navigator/graph"
	errorspkg "navigator/internal/errors"
	"navigator/internal/format"
)

const (
	menuLoadGraph = iota + 1
	menuBFS
	menuDFS
	menuShortestPath
	menuAllPairsShortestPaths
	menuMST
	menuTSP
	menuExportDot
	menuExit
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var g graph.Graph
	loaded := false
	ga := algorithms.GraphAlgorithms{}

	for {
		printMenu(loaded)
		choice, ok := readInt(scanner)
		if !ok {
			fmt.Println("bye")
			return
		}
		switch choice {
		case menuLoadGraph:
			fmt.Print("Enter file path: ")
			path := readLine(scanner)
			if err := g.LoadGraphFromFile(strings.TrimSpace(path)); err != nil {
				fmt.Printf("error: %v\n", err)
				loaded = false
				continue
			}
			loaded = true
			fmt.Println("graph loaded")
		case menuBFS:
			if !requireLoaded(loaded) {
				continue
			}
			start := askVertex(scanner, "Start vertex")
			res, err := ga.BreadthFirstSearch(&g, start)
			printResultOrError(res, err)
		case menuDFS:
			if !requireLoaded(loaded) {
				continue
			}
			start := askVertex(scanner, "Start vertex")
			res, err := ga.DepthFirstSearch(&g, start)
			printResultOrError(res, err)
		case menuShortestPath:
			if !requireLoaded(loaded) {
				continue
			}
			v1 := askVertex(scanner, "Vertex 1")
			v2 := askVertex(scanner, "Vertex 2")
			dist, err := ga.GetShortestPathBetweenVertices(&g, v1, v2)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Printf("shortest distance: %d\n", dist)
		case menuAllPairsShortestPaths:
			if !requireLoaded(loaded) {
				continue
			}
			m, err := ga.GetShortestPathsBetweenAllVertices(&g)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Println(format.Matrix(m))
		case menuMST:
			if !requireLoaded(loaded) {
				continue
			}
			mst, err := ga.GetLeastSpanningTree(&g)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Println(format.Matrix(mst))
		case menuTSP:
			if !requireLoaded(loaded) {
				continue
			}
			res, err := ga.SolveTravelingSalesmanProblem(&g)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Printf("route: %v\n", res.Vertices)
			fmt.Printf("distance: %.0f\n", res.Distance)
		case menuExportDot:
			if !requireLoaded(loaded) {
				continue
			}
			fmt.Print("Enter output dot path: ")
			out := readLine(scanner)
			if err := g.ExportGraphToDot(strings.TrimSpace(out)); err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			fmt.Println("exported")
		case menuExit:
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown option")
		}
	}
}

func printMenu(loaded bool) {
	fmt.Println()
	fmt.Println("1) Load graph from file")
	fmt.Println("2) Breadth-first search")
	fmt.Println("3) Depth-first search")
	fmt.Println("4) Shortest path between two vertices")
	fmt.Println("5) Shortest paths between all vertices")
	fmt.Println("6) Minimum spanning tree (Prim)")
	fmt.Println("7) Traveling salesman problem (ant colony)")
	fmt.Println("8) Export graph to DOT")
	fmt.Println("9) Exit")
	if !loaded {
		fmt.Println("(graph not loaded)")
	}
	fmt.Print("Select: ")
}

func requireLoaded(loaded bool) bool {
	if loaded {
		return true
	}
	fmt.Printf("error: %v\n", errorspkg.ErrGraphNotLoaded)
	return false
}

func askVertex(scanner *bufio.Scanner, label string) int {
	for {
		fmt.Printf("%s (1-based): ", label)
		v, ok := readInt(scanner)
		if ok {
			return v
		}
		fmt.Println("invalid number")
	}
}

func readInt(scanner *bufio.Scanner) (int, bool) {
	line := readLine(scanner)
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, false
	}
	v, err := strconv.Atoi(line)
	if err != nil {
		return 0, false
	}
	return v, true
}

func readLine(scanner *bufio.Scanner) string {
	if !scanner.Scan() {
		return ""
	}
	return scanner.Text()
}

func printResultOrError(vertices []int, err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println(vertices)
}
