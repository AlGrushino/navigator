package errors

import "errors"

var (
	ErrGraphNotLoaded     = errors.New("graph is not loaded")
	ErrInvalidVertex      = errors.New("invalid vertex")
	ErrInvalidGraphFormat = errors.New("invalid graph file format")
	ErrGraphNotConnected  = errors.New("graph must be connected")
	ErrEmptyGraph         = errors.New("graph must be non-zero")
	ErrInvalidWeight      = errors.New("edge weights must be natural numbers (or 0 for no edge)")
	ErrNoPath             = errors.New("no path between vertices")
	ErrEmptyContainer     = errors.New("container is empty")
	ErrNotUndirected      = errors.New("operation requires an undirected graph")
	ErrTSPUnsolvable      = errors.New("traveling salesman problem cannot be solved for this graph")
)
