// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"sort"

	"github.com/jingcheng-WU/gonum/graph"
	"github.com/jingcheng-WU/gonum/graph/internal/ordered"
	"github.com/jingcheng-WU/gonum/graph/iterator"
	"github.com/jingcheng-WU/gonum/mat"
)

var (
	um *UndirectedMatrix

	_ graph.Graph        = um
	_ graph.Undirected   = um
	_ edgeSetter         = um
	_ weightedEdgeSetter = um
)

// UndirectedMatrix represents an undirected graph using an adjacency
// matrix such that all IDs are in a contiguous block from 0 to n-1.
// Edges are stored implicitly as an edge weight, so edges stored in
// the graph are not recoverable.
type UndirectedMatrix struct {
	mat   *mat.SymDense
	nodes []graph.Node

	self   float64
	absent float64
}

// NewUndirectedMatrix creates an undirected dense graph with n nodes.
// All edges are initialized with the weight given by init. The self parameter
// specifies the cost of self connection, and absent specifies the weight
// returned for absent edges.
func NewUndirectedMatrix(n int, init, self, absent float64) *UndirectedMatrix {
	matrix := make([]float64, n*n)
	if init != 0 {
		for i := range matrix {
			matrix[i] = init
		}
	}
	for i := 0; i < len(matrix); i += n + 1 {
		matrix[i] = self
	}
	return &UndirectedMatrix{
		mat:    mat.NewSymDense(n, matrix),
		self:   self,
		absent: absent,
	}
}

// NewUndirectedMatrixFrom creates an undirected dense graph with the given nodes.
// The IDs of the nodes must be contiguous from 0 to len(nodes)-1, but may
// be in any order. If IDs are not contiguous NewUndirectedMatrixFrom will panic.
// All edges are initialized with the weight given by init. The self parameter
// specifies the cost of self connection, and absent specifies the weight
// returned for absent edges.
func NewUndirectedMatrixFrom(nodes []graph.Node, init, self, absent float64) *UndirectedMatrix {
	sort.Sort(ordered.ByID(nodes))
	for i, n := range nodes {
		if int64(i) != n.ID() {
			panic("simple: non-contiguous node IDs")
		}
	}
	g := NewUndirectedMatrix(len(nodes), init, self, absent)
	g.nodes = nodes
	return g
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *UndirectedMatrix) Edge(uid, vid int64) graph.Edge {
	return g.WeightedEdgeBetween(uid, vid)
}

// EdgeBetween returns the edge between nodes x and y.
func (g *UndirectedMatrix) EdgeBetween(uid, vid int64) graph.Edge {
	return g.WeightedEdgeBetween(uid, vid)
}

// Edges returns all the edges in the graph.
func (g *UndirectedMatrix) Edges() graph.Edges {
	var edges []graph.Edge
	r, _ := g.mat.Dims()
	for i := 0; i < r; i++ {
		for j := i + 1; j < r; j++ {
			if w := g.mat.At(i, j); !isSame(w, g.absent) {
				edges = append(edges, WeightedEdge{F: g.Node(int64(i)), T: g.Node(int64(j)), W: w})
			}
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedEdges(edges)
}

// From returns all nodes in g that can be reached directly from n.
func (g *UndirectedMatrix) From(id int64) graph.Nodes {
	if !g.has(id) {
		return graph.Empty
	}
	var nodes []graph.Node
	r := g.mat.Symmetric()
	for i := 0; i < r; i++ {
		if int64(i) == id {
			continue
		}
		// id is not greater than maximum int by this point.
		if !isSame(g.mat.At(int(id), i), g.absent) {
			nodes = append(nodes, g.Node(int64(i)))
		}
	}
	if len(nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(nodes)
}

// HasEdgeBetween returns whether an edge exists between nodes x and y.
func (g *UndirectedMatrix) HasEdgeBetween(uid, vid int64) bool {
	if !g.has(uid) {
		return false
	}
	if !g.has(vid) {
		return false
	}
	// uid and vid are not greater than maximum int by this point.
	return uid != vid && !isSame(g.mat.At(int(uid), int(vid)), g.absent)
}

// Matrix returns the mat.Matrix representation of the graph.
func (g *UndirectedMatrix) Matrix() mat.Matrix {
	// Prevent alteration of dimensions of the returned matrix.
	m := *g.mat
	return &m
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *UndirectedMatrix) Node(id int64) graph.Node {
	if !g.has(id) {
		return nil
	}
	if g.nodes == nil {
		return Node(id)
	}
	return g.nodes[id]
}

// Nodes returns all the nodes in the graph.
func (g *UndirectedMatrix) Nodes() graph.Nodes {
	if g.nodes != nil {
		nodes := make([]graph.Node, len(g.nodes))
		copy(nodes, g.nodes)
		return iterator.NewOrderedNodes(nodes)
	}
	r := g.mat.Symmetric()
	// Matrix graphs must have at least one node.
	return iterator.NewImplicitNodes(0, r, newSimpleNode)
}

// RemoveEdge removes the edge with the given end point IDs from the graph, leaving the terminal
// nodes. If the edge does not exist it is a no-op.
func (g *UndirectedMatrix) RemoveEdge(fid, tid int64) {
	if !g.has(fid) {
		return
	}
	if !g.has(tid) {
		return
	}
	// fid and tid are not greater than maximum int by this point.
	g.mat.SetSym(int(fid), int(tid), g.absent)
}

// SetEdge sets e, an edge from one node to another with unit weight. If the ends of the edge are
// not in g or the edge is a self loop, SetEdge panics. SetEdge will store the nodes of
// e in the graph if it was initialized with NewUndirectedMatrixFrom.
func (g *UndirectedMatrix) SetEdge(e graph.Edge) {
	g.setWeightedEdge(e, 1)
}

// SetWeightedEdge sets e, an edge from one node to another. If the ends of the edge are not in g
// or the edge is a self loop, SetWeightedEdge panics. SetWeightedEdge will store the nodes of
// e in the graph if it was initialized with NewUndirectedMatrixFrom.
func (g *UndirectedMatrix) SetWeightedEdge(e graph.WeightedEdge) {
	g.setWeightedEdge(e, e.Weight())
}

func (g *UndirectedMatrix) setWeightedEdge(e graph.Edge, weight float64) {
	from := e.From()
	fid := from.ID()
	to := e.To()
	tid := to.ID()
	if fid == tid {
		panic("simple: set illegal edge")
	}
	if int64(int(fid)) != fid {
		panic("simple: unavailable from node ID for dense graph")
	}
	if int64(int(tid)) != tid {
		panic("simple: unavailable to node ID for dense graph")
	}
	if g.nodes != nil {
		g.nodes[fid] = from
		g.nodes[tid] = to
	}
	// fid and tid are not greater than maximum int by this point.
	g.mat.SetSym(int(fid), int(tid), weight)
}

// Weight returns the weight for the edge between x and y if Edge(x, y) returns a non-nil Edge.
// If x and y are the same node or there is no joining edge between the two nodes the weight
// value returned is either the graph's absent or self value. Weight returns true if an edge
// exists between x and y or if x and y have the same ID, false otherwise.
func (g *UndirectedMatrix) Weight(xid, yid int64) (w float64, ok bool) {
	if xid == yid {
		return g.self, true
	}
	if g.HasEdgeBetween(xid, yid) {
		// xid and yid are not greater than maximum int by this point.
		return g.mat.At(int(xid), int(yid)), true
	}
	return g.absent, false
}

// WeightedEdge returns the weighted edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *UndirectedMatrix) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	return g.WeightedEdgeBetween(uid, vid)
}

// WeightedEdgeBetween returns the weighted edge between nodes x and y.
func (g *UndirectedMatrix) WeightedEdgeBetween(uid, vid int64) graph.WeightedEdge {
	if g.HasEdgeBetween(uid, vid) {
		// uid and vid are not greater than maximum int by this point.
		return WeightedEdge{F: g.Node(uid), T: g.Node(vid), W: g.mat.At(int(uid), int(vid))}
	}
	return nil
}

// WeightedEdges returns all the edges in the graph.
func (g *UndirectedMatrix) WeightedEdges() graph.WeightedEdges {
	var edges []graph.WeightedEdge
	r, _ := g.mat.Dims()
	for i := 0; i < r; i++ {
		for j := i + 1; j < r; j++ {
			if w := g.mat.At(i, j); !isSame(w, g.absent) {
				edges = append(edges, WeightedEdge{F: g.Node(int64(i)), T: g.Node(int64(j)), W: w})
			}
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedWeightedEdges(edges)
}

func (g *UndirectedMatrix) has(id int64) bool {
	r := g.mat.Symmetric()
	return 0 <= id && id < int64(r)
}
