package search

import (
	"container/heap"
	"math"
	"sort"

	gr "github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/set"
	"github.com/gonum/graph/xifo"
)

// Returns an ordered list consisting of the nodes between start and goal. The path will be the shortest path assuming the function heuristicCost is admissible.
// The second return value is the cost, and the third is the number of nodes expanded while searching (useful info for tuning heuristics). Negative Costs will cause
// bad things to happen, as well as negative heuristic estimates.
//
// A heuristic is admissible if, for any node in the graph, the heuristic estimate of the cost between the node and the goal is less than or equal to the true cost.
//
// Performance may be improved by providing a consistent heuristic (though one is not needed to find the optimal path), a heuristic is consistent if its value for a given node is less than (or equal to) the
// actual cost of reaching its neighbors + the heuristic estimate for the neighbor itself. You can force consistency by making your HeuristicCost function
// return max(NonConsistentHeuristicCost(neighbor,goal), NonConsistentHeuristicCost(self,goal) - Cost(self,neighbor)). If there are multiple neighbors, take the max of all of them.
//
// Cost and HeuristicCost take precedence for evaluating cost/heuristic distance. If one is not present (i.e. nil) the function will check the graph's interface for the respective interface:
// Coster for Cost and HeuristicCoster for HeuristicCost. If the correct one is present, it will use the graph's function for evaluation.
//
// Finally, if neither the argument nor the interface is present, the function will assume discrete.UniformCost for Cost and discrete.NullHeuristic for HeuristicCost
//
// To run Uniform Cost Search, run A* with the NullHeuristic
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost function that returns a uniform positive value)
func AStar(start, goal gr.Node, graph gr.Graph, Cost, HeuristicCost func(gr.Node, gr.Node) float64) (path []gr.Node, cost float64, nodesExpanded int) {
	successors, _, _, _, _, _, Cost, HeuristicCost := setupFuncs(graph, Cost, HeuristicCost)

	closedSet := make(map[int]internalNode)
	openSet := &aStarPriorityQueue{nodes: make([]internalNode, 0), indexList: make(map[int]int)}
	heap.Init(openSet)
	node := internalNode{start, 0, HeuristicCost(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]gr.Node)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		nodesExpanded += 1

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.ID()] = curr

		for _, neighbor := range successors(curr.Node) {
			if _, ok := closedSet[neighbor.ID()]; ok {
				continue
			}

			g := curr.gscore + Cost(curr.Node, neighbor)

			if existing, exists := openSet.Find(neighbor.ID()); !exists {
				predecessor[neighbor.ID()] = curr
				node = internalNode{neighbor, g, g + HeuristicCost(neighbor, goal)}
				heap.Push(openSet, node)
			} else if g < existing.gscore {
				predecessor[neighbor.ID()] = curr
				openSet.Fix(neighbor.ID(), g, g+HeuristicCost(neighbor, goal))
			}
		}
	}

	return nil, 0.0, nodesExpanded
}

// BreadthFirstSearch finds a path with a minimal number of edges from from start to goal.
//
// BreadthFirstSearch returns the path found and the number of nodes visited in the search.
// The returned path is nil if no path exists.
func BreadthFirstSearch(start, goal gr.Node, graph gr.Graph) ([]gr.Node, int) {
	path, _, visited := AStar(start, goal, graph, UniformCost, NullHeuristic)
	return path, visited
}

// Dijkstra's Algorithm is essentially a goalless Uniform Cost Search. That is, its results are roughly equivalent to
// running A* with the Null Heuristic from a single node to every other node in the graph -- though it's a fair bit faster
// because running A* in that way will recompute things it's already computed every call. Note that you won't necessarily get the same path
// you would get for A*, but the cost is guaranteed to be the same (that is, if multiple shortest paths exist, you may get a different shortest path).
//
// Like A*, Dijkstra's Algorithm likely won't run correctly with negative edge weights -- use Bellman-Ford for that instead
//
// Dijkstra's algorithm usually only returns a cost map, however, since the data is available this version will also reconstruct the path to every node
func Dijkstra(source gr.Node, graph gr.Graph, Cost func(gr.Node, gr.Node) float64) (paths map[int][]gr.Node, costs map[int]float64) {
	successors, _, _, _, _, _, Cost, _ := setupFuncs(graph, Cost, nil)

	nodes := graph.NodeList()
	openSet := &aStarPriorityQueue{nodes: make([]internalNode, 0), indexList: make(map[int]int)}
	closedSet := set.NewSet()                 // This is to make use of that same
	costs = make(map[int]float64, len(nodes)) // May overallocate, will change if it becomes a problem
	predecessor := make(map[int]gr.Node, len(nodes))
	nodeIDMap := make(map[int]gr.Node, len(nodes))
	heap.Init(openSet)

	// I don't think we actually need the init step since I use a map check rather than inf to check if we're done
	/*for _, node := range nodes {
		if node == source {
			heap.Push(openSet, internalNode{node, 0, 0})
			costs[node] = 0
		} else {
			heap.Push(openSet, internalNode{node, math.MaxFloat64, math.MaxFloat64})
			predecessor[node] = -1
		}
	}*/

	costs[source.ID()] = 0
	heap.Push(openSet, internalNode{source, 0, 0})

	for openSet.Len() != 0 {
		node := heap.Pop(openSet).(internalNode)
		/* if _, ok := costs[node.int]; !ok {
			 break
		 } */

		if closedSet.Contains(node.ID()) { // As in A*, prevents us from having to slowly search and reorder the queue
			continue
		}

		nodeIDMap[node.ID()] = node

		closedSet.Add(node.ID())

		for _, neighbor := range successors(node) {
			tmpCost := costs[node.ID()] + Cost(node, neighbor)
			if cost, ok := costs[neighbor.ID()]; !ok || tmpCost < cost {
				costs[neighbor.ID()] = tmpCost
				predecessor[neighbor.ID()] = node
				heap.Push(openSet, internalNode{neighbor, tmpCost, tmpCost})
			}
		}
	}

	paths = make(map[int][]gr.Node, len(costs))
	for node, _ := range costs { // Only reconstruct the path if one exists
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
	}
	return paths, costs
}

// The Bellman-Ford Algorithm is the same as Dijkstra's Algorithm with a key difference. They both take a single source and find the shortest path to every other
// (reachable) node in the graph. Bellman-Ford, however, will detect negative edge loops and abort if one is present. A negative edge loop occurs when there is a cycle in the graph
// such that it can take an edge with a negative cost over and over. A -(-2)> B -(2)> C isn't a loop because A->B can only be taken once, but A<-(-2)->B-(2)>C is one because
// A and B have a bi-directional edge, and algorithms like Dijkstra's will infinitely flail between them getting progressively lower costs.
//
// That said, if you do not have a negative edge weight, use Dijkstra's Algorithm instead, because it's faster.
//
// Like Dijkstra's, along with the costs this implementation will also construct all the paths for you. In addition, it has a third return value which will be true if the algorithm was aborted
// due to the presence of a negative edge weight cycle.
func BellmanFord(source gr.Node, graph gr.Graph, Cost func(gr.Node, gr.Node) float64) (paths map[int][]gr.Node, costs map[int]float64, aborted bool) {
	_, _, _, _, _, _, Cost, _ = setupFuncs(graph, Cost, nil)

	predecessor := make(map[int]gr.Node)
	costs = make(map[int]float64)
	nodeIDMap := make(map[int]gr.Node)
	nodeIDMap[source.ID()] = source
	costs[source.ID()] = 0
	nodes := graph.NodeList()
	edges := graph.EdgeList()

	for i := 1; i < len(nodes)-1; i++ {
		for _, edge := range edges {
			weight := Cost(edge.Head(), edge.Tail())
			nodeIDMap[edge.Head().ID()] = edge.Head()
			nodeIDMap[edge.Tail().ID()] = edge.Tail()
			if dist := costs[edge.Head().ID()] + weight; dist < costs[edge.Tail().ID()] {
				costs[edge.Tail().ID()] = dist
				predecessor[edge.Tail().ID()] = edge.Head()
			}
		}
	}

	for _, edge := range edges {
		weight := Cost(edge.Head(), edge.Tail())
		if costs[edge.Head().ID()]+weight < costs[edge.Tail().ID()] {
			return nil, nil, true // Abandoned because a cycle is detected
		}
	}

	paths = make(map[int][]gr.Node, len(costs))
	for node, _ := range costs {
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
	}
	return paths, costs, false
}

// Johnson's Algorithm generates the lowest cost path between every pair of nodes in the graph.
//
// It makes use of Bellman-Ford and a dummy graph. It creates a dummy node containing edges with a cost of zero to every other node. Then it runs Bellman-Ford with this
// dummy node as the source. It then modifies the all the nodes' edge weights (which gets rid of all negative weights).
//
// Finally, it removes the dummy node and runs Dijkstra's starting at every node.
//
// This algorithm is fairly slow. Its purpose is to remove negative edge weights to allow Dijkstra's to function properly. It's probably not worth it to run this algorithm if you have
// all non-negative edge weights. Also note that this implementation copies your whole graph into a GonumGraph (so it can add/remove the dummy node and edges and reweight the graph).
//
// Its return values are, in order: a map from the source node, to the destination node, to the path between them; a map from the source node, to the destination node, to the cost of the path between them;
// and a bool that is true if Bellman-Ford detected a negative edge weight cycle -- thus causing it (and this algorithm) to abort (if aborted is true, both maps will be nil).
func Johnson(graph gr.Graph, Cost func(gr.Node, gr.Node) float64) (nodePaths map[int]map[int][]gr.Node, nodeCosts map[int]map[int]float64, aborted bool) {
	successors, _, _, _, _, _, Cost, _ := setupFuncs(graph, Cost, nil)
	/* Copy graph into a mutable one since it has to be altered for this algorithm */
	dummyGraph := concrete.NewGonumGraph(true)
	for _, node := range graph.NodeList() {
		neighbors := successors(node)
		if !dummyGraph.NodeExists(node) {
			dummyGraph.AddNode(node, neighbors)
			for _, neighbor := range neighbors {
				dummyGraph.SetEdgeCost(concrete.GonumEdge{node, neighbor}, Cost(node, neighbor))
			}
		} else {
			for _, neighbor := range neighbors {
				dummyGraph.AddEdge(concrete.GonumEdge{node, neighbor})
				dummyGraph.SetEdgeCost(concrete.GonumEdge{node, neighbor}, Cost(node, neighbor))
			}
		}
	}

	/* Step 1: Dummy node with 0 cost edge weights to every other node*/
	dummyNode := dummyGraph.NewNode(graph.NodeList())
	for _, node := range graph.NodeList() {
		dummyGraph.SetEdgeCost(concrete.GonumEdge{dummyNode, node}, 0)
	}

	/* Step 2: Run Bellman-Ford starting at the dummy node, abort if it detects a cycle */
	_, costs, aborted := BellmanFord(dummyNode, dummyGraph, nil)
	if aborted {
		return nil, nil, true
	}

	/* Step 3: reweight the graph and remove the dummy node */
	for _, edge := range graph.EdgeList() {
		dummyGraph.SetEdgeCost(edge, Cost(edge.Head(), edge.Tail())+costs[edge.Head().ID()]-costs[edge.Tail().ID()])
	}

	dummyGraph.RemoveNode(dummyNode)

	/* Step 4: Run Dijkstra's starting at every node */
	nodePaths = make(map[int]map[int][]gr.Node, len(graph.NodeList()))
	nodeCosts = make(map[int]map[int]float64)

	for _, node := range graph.NodeList() {
		nodePaths[node.ID()], nodeCosts[node.ID()] = Dijkstra(node, dummyGraph, nil)
	}

	return nodePaths, nodeCosts, false
}

// Expands the first node it sees trying to find the destination. Depth First Search is *not* guaranteed to find the shortest path,
// however, if a path exists DFS is guaranteed to find it (provided you don't find a way to implement a Graph with an infinite depth)
func DepthFirstSearch(start, goal gr.Node, graph gr.Graph) []gr.Node {
	successors, _, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)
	closedSet := set.NewSet()
	openSet := xifo.GonumStack([]interface{}{start})
	predecessor := make(map[int]gr.Node)

	for !openSet.IsEmpty() {
		c := openSet.Pop()

		curr := c.(gr.Node)

		if closedSet.Contains(curr.ID()) {
			continue
		}

		if curr == goal {
			return rebuildPath(predecessor, goal)
		}

		closedSet.Add(curr.ID())

		for _, neighbor := range successors(curr) {
			if closedSet.Contains(neighbor.ID()) {
				continue
			}

			predecessor[neighbor.ID()] = curr
			openSet.Push(neighbor)
		}
	}

	return nil
}

// An admissible, consistent heuristic that won't speed up computation time at all.
func NullHeuristic(a, b gr.Node) float64 {
	return 0.0
}

// Assumes all edges in the graph have the same weight (including edges that don't exist!)
func UniformCost(a, b gr.Node) float64 {
	return 1.0
}

/* Slow functions to replace the guarantee of a graph being directed or undirected */

// A slow function that iterates over the entire edge list trying to find all successors and predecessors of a node
func NeighborsFunc(graph gr.Graph) func(node gr.Node) []gr.Node {
	return func(node gr.Node) []gr.Node {
		neighbors := []gr.Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() {
				neighbors = append(neighbors, edge.Tail())
			} else if edge.Tail().ID() == node.ID() {
				neighbors = append(neighbors, edge.Head())
			}
		}

		return neighbors
	}
}

// A slow function that iterates over the entire edge list trying to find all successors of a node
func SuccessorsFunc(graph gr.Graph) func(node gr.Node) []gr.Node {
	return func(node gr.Node) []gr.Node {
		neighbors := []gr.Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() {
				neighbors = append(neighbors, edge.Tail())
			}
		}

		return neighbors
	}
}

// A slow function that iterates over the entire edge list trying to find all predecessors of a node
func PredecessorsFunc(graph gr.Graph) func(node gr.Node) []gr.Node {
	return func(node gr.Node) []gr.Node {
		neighbors := []gr.Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Tail().ID() == node.ID() {
				neighbors = append(neighbors, edge.Head())
			}
		}

		return neighbors
	}
}

// A slow function that iterates over the entire edge list trying to find a matching edge such that the first argument is the Head and the second is Tail
func IsSuccessorFunc(graph gr.Graph) func(gr.Node, gr.Node) bool {
	return func(node, succ gr.Node) bool {
		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() && edge.Tail().ID() == succ.ID() {
				return true
			}
		}

		return false
	}
}

// A slow function that iterates over the entire edge list trying to find a matching edge such that the first argument is the Tail and the second is Head
func IsPredecessorFunc(graph gr.Graph) func(gr.Node, gr.Node) bool {
	return func(node, pred gr.Node) bool {
		for _, edge := range graph.EdgeList() {
			if edge.Tail().ID() == node.ID() && edge.Head().ID() == pred.ID() {
				return true
			}
		}

		return false
	}
}

// A slow function that iterates over the entire edge list trying to find a matching edge such that the first argument is one end and the second argument is the other
func IsNeighborFunc(graph gr.Graph) func(gr.Node, gr.Node) bool {
	return func(node, succ gr.Node) bool {
		for _, edge := range graph.EdgeList() {
			if (edge.Tail().ID() == node.ID() || edge.Head().ID() == node.ID()) && (edge.Tail().ID() == succ.ID() || edge.Head().ID() == succ.ID()) {
				return true
			}
		}

		return false
	}
}

/* Simple operations */

// Copies a graph into the destination; maintaining all node IDs.
func CopyGraph(dst gr.MutableGraph, src gr.Graph) {
	dst.EmptyGraph()
	dst.SetDirected(false)

	var Cost func(gr.Node, gr.Node) float64
	if cgraph, ok := src.(gr.Coster); ok {
		Cost = cgraph.Cost
	}

	for _, edge := range src.EdgeList() {
		if !dst.NodeExists(edge.Head()) {
			dst.AddNode(edge.Head(), []gr.Node{edge.Tail()})
		} else {
			dst.AddEdge(edge)
		}

		if Cost != nil {
			dst.SetEdgeCost(edge, Cost(edge.Head(), edge.Tail()))
		}
	}
}

/* Basic Graph tests */

// Also known as Tarjan's Strongly Connected Components Algorithm. This returns all the strongly connected components in the graph.
//
// A strongly connected component of a graph is a set of vertices where it's possible to reach any vertex in the set from any other (meaning there's a cycle between them)
//
// Generally speaking, a directed graph where the number of strongly connected components is equal to the number of nodes is acyclic, unless you count reflexive edges as a cycle (which requires only a little extra testing)
//
// An undirected graph should end up with as many SCCs as there are "islands" (or subgraphs) of connections, meaning having more than one strongly connected component implies that your graph is not fully connected.
func Tarjan(graph gr.Graph) (sccs [][]gr.Node) {
	index := 0
	vStack := &xifo.GonumStack{}
	stackSet := set.NewSet()
	sccs = make([][]gr.Node, 0)

	nodes := graph.NodeList()
	lowlinks := make(map[int]int, len(nodes))
	indices := make(map[int]int, len(nodes))

	successors, _, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)

	var strongconnect func(gr.Node) []gr.Node

	strongconnect = func(node gr.Node) []gr.Node {
		indices[node.ID()] = index
		lowlinks[node.ID()] = index
		index += 1

		vStack.Push(node)
		stackSet.Add(node.ID())

		for _, succ := range successors(node) {
			if _, ok := indices[succ.ID()]; !ok {
				strongconnect(succ)
				lowlinks[node.ID()] = int(math.Min(float64(lowlinks[node.ID()]), float64(lowlinks[succ.ID()])))
			} else if stackSet.Contains(succ) {
				lowlinks[node.ID()] = int(math.Min(float64(lowlinks[node.ID()]), float64(lowlinks[succ.ID()])))
			}
		}

		if lowlinks[node.ID()] == indices[node.ID()] {
			scc := make([]gr.Node, 0)
			for {
				v := vStack.Pop()
				stackSet.Remove(v.(gr.Node).ID())
				scc = append(scc, v.(gr.Node))
				if v.(gr.Node).ID() == node.ID() {
					return scc
				}
			}
		}

		return nil
	}

	for _, n := range nodes {
		if _, ok := indices[n.ID()]; !ok {
			sccs = append(sccs, strongconnect(n))
		}
	}

	return sccs
}

// Returns true if, starting at path[0] and ending at path[len(path)-1], all nodes between are valid neighbors. That is, for each element path[i], path[i+1] is a valid successor
//
// Special case: a nil or zero length path is considered valid (true), a path of length 1 (only one node) is the trivial case, but only if the node listed in path exists.
func IsPath(path []gr.Node, graph gr.Graph) bool {
	_, _, _, isSuccessor, _, _, _, _ := setupFuncs(graph, nil, nil)
	if path == nil || len(path) == 0 {
		return true
	} else if len(path) == 1 {
		return graph.NodeExists(path[0])
	}

	for i := 0; i < len(path)-1; i++ {
		if !isSuccessor(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

/* Implements minimum-spanning tree algorithms; puts the resulting minimum spanning tree in the dst graph */

// Generates a minimum spanning tree with sets.
//
// As with other algorithms that use Cost, the order of precedence is Argument > Interface > UniformCost
func Prim(dst gr.MutableGraph, graph gr.Graph, Cost func(gr.Node, gr.Node) float64) {
	if Cost == nil {
		if cgraph, ok := graph.(gr.Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	dst.EmptyGraph()
	dst.SetDirected(false)

	nlist := graph.NodeList()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0], nil)
	remainingNodes := set.NewSet()
	for _, node := range nlist[1:] {
		remainingNodes.Add(node.ID())
	}

	edgeList := graph.EdgeList()
	for remainingNodes.Cardinality() != 0 {
		edgeWeights := make(edgeSorter, 0)
		for _, edge := range edgeList {
			if dst.NodeExists(edge.Head()) && remainingNodes.Contains(edge.Tail().ID()) {
				edgeWeights = append(edgeWeights, WeightedEdge{Edge: edge, Weight: Cost(edge.Head(), edge.Tail())})
			}
		}

		sort.Sort(edgeWeights)
		myEdge := edgeWeights[0]

		if !dst.NodeExists(myEdge.Head()) {
			dst.AddNode(myEdge.Head(), []gr.Node{myEdge.Tail()})
		} else {
			dst.AddEdge(myEdge.Edge)
		}
		dst.SetEdgeCost(myEdge.Edge, myEdge.Weight)

		remainingNodes.Remove(myEdge.Edge.Head())
	}

}

// Generates a minimum spanning tree for a graph using discrete.DisjointSet
//
// As with other algorithms with Cost, the precedence goes Argument > Interface > UniformCost
func Kruskal(dst gr.MutableGraph, graph gr.Graph, Cost func(gr.Node, gr.Node) float64) {
	if Cost == nil {
		if cgraph, ok := graph.(gr.Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	dst.EmptyGraph()
	dst.SetDirected(false)

	edgeList := graph.EdgeList()
	edgeWeights := make(edgeSorter, 0, len(edgeList))
	for _, edge := range edgeList {
		edgeWeights = append(edgeWeights, WeightedEdge{Edge: edge, Weight: Cost(edge.Head(), edge.Tail())})
	}

	sort.Sort(edgeWeights)

	ds := set.NewDisjointSet()
	for _, node := range graph.NodeList() {
		ds.MakeSet(node.ID())
	}

	for _, edge := range edgeWeights {
		if s1, s2 := ds.Find(edge.Edge.Head().ID()), ds.Find(edge.Edge.Tail().ID); s1 != s2 {
			ds.Union(s1, s2)
			if !dst.NodeExists(edge.Edge.Head()) {
				dst.AddNode(edge.Edge.Head(), []gr.Node{edge.Edge.Tail()})
			} else {
				dst.AddEdge(edge.Edge)
			}
			dst.SetEdgeCost(edge.Edge, edge.Weight)
		}
	}
}

/* Control flow graph stuff */

// A dominates B if and only if the only path through B travels through A
//
// This returns all possible dominators for all nodes, it does not prune for strict dominators, immediate dominators etc
//
// The int map[int]*set.Set is the node's ID
func Dominators(start gr.Node, graph gr.Graph) map[int]*set.Set {
	allNodes := set.NewSet()
	nlist := graph.NodeList()
	dominators := make(map[int]*set.Set, len(nlist))
	for _, node := range nlist {
		allNodes.Add(node.ID())
	}

	_, predecessors, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)

	for _, node := range nlist {
		dominators[node.ID()] = set.NewSet()
		if node.ID() == start.ID() {
			dominators[node.ID()].Add(start.ID())
		} else {
			dominators[node.ID()].Copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == start.ID() {
				continue
			}
			preds := predecessors(node)
			if len(preds) == 0 {
				continue
			}
			tmp := set.NewSet().Copy(dominators[preds[0].ID()])
			for _, pred := range preds[1:] {
				tmp.Intersection(tmp, dominators[pred.ID()])
			}

			dom := set.NewSet()
			dom.Add(node.ID())

			dom.Union(dom, tmp)
			if !set.Equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

// A Postdominates B if and only if all paths from B travel through A
//
// This returns all possible post-dominators for all nodes, it does not prune for strict postdominators, immediate postdominators etc
func PostDominators(end gr.Node, graph gr.Graph) map[int]*set.Set {
	successors, _, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)
	allNodes := set.NewSet()
	nlist := graph.NodeList()
	dominators := make(map[int]*set.Set, len(nlist))
	for _, node := range nlist {
		allNodes.Add(node.ID())
	}

	for _, node := range nlist {
		dominators[node.ID()] = set.NewSet()
		if node.ID() == end.ID() {
			dominators[node.ID()].Add(end.ID())
		} else {
			dominators[node.ID()].Copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == end.ID() {
				continue
			}
			succs := successors(node)
			if len(succs) == 0 {
				continue
			}
			tmp := set.NewSet().Copy(dominators[succs[0].ID()])
			for _, succ := range succs[1:] {
				tmp.Intersection(tmp, dominators[succ.ID()])
			}

			dom := set.NewSet()
			dom.Add(node.ID())

			dom.Union(dom, tmp)
			if !set.Equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}
