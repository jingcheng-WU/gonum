// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen_test

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

// IDRange is an IDer that provides a set of IDs in [First, Last].
type IDRange struct{ First, Last int64 }

func (r IDRange) Len() int       { return int(r.Last - r.First + 1) }
func (r IDRange) ID(i int) int64 { return r.First + int64(i) }

func ExampleStar_undirectedRange() {
	dst := simple.NewUndirectedGraph()
	err := gen.Star(dst, 0, IDRange{First: 1, Last: 6})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "star", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict graph star {
	// 	// Node definitions.
	// 	0;
	// 	1;
	// 	2;
	// 	3;
	// 	4;
	// 	5;
	// 	6;
	//
	// 	// Edge definitions.
	// 	0 -- 1;
	// 	0 -- 2;
	// 	0 -- 3;
	// 	0 -- 4;
	// 	0 -- 5;
	// 	0 -- 6;
	// }
}

func ExampleWheel_directedRange() {
	dst := simple.NewDirectedGraph()
	err := gen.Wheel(dst, 0, IDRange{First: 1, Last: 6})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "wheel", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph wheel {
	// 	// Node definitions.
	// 	0;
	// 	1;
	// 	2;
	// 	3;
	// 	4;
	// 	5;
	// 	6;
	//
	// 	// Edge definitions.
	// 	0 -> 1;
	// 	0 -> 2;
	// 	0 -> 3;
	// 	0 -> 4;
	// 	0 -> 5;
	// 	0 -> 6;
	// 	1 -> 2;
	// 	2 -> 3;
	// 	3 -> 4;
	// 	4 -> 5;
	// 	5 -> 6;
	// 	6 -> 1;
	// }
}

// IDSet is an IDer providing an explicit set of IDs.
type IDSet []int64

func (s IDSet) Len() int       { return len(s) }
func (s IDSet) ID(i int) int64 { return s[i] }

func ExamplePath_directedSet() {
	dst := simple.NewDirectedGraph()
	err := gen.Path(dst, IDSet{2, 4, 5, 9})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "path", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph path {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	4 -> 5;
	// 	5 -> 9;
	// }
}

func ExampleComplete_directedSet() {
	dst := simple.NewDirectedGraph()
	err := gen.Complete(dst, IDSet{2, 4, 5, 9})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	2 -> 5;
	// 	2 -> 9;
	// 	4 -> 5;
	// 	4 -> 9;
	// 	5 -> 9;
	// }
}

// Bidirected allows bidirectional directed graph construction.
type Bidirected struct {
	*simple.DirectedGraph
}

func (g Bidirected) SetEdge(e graph.Edge) {
	g.DirectedGraph.SetEdge(e)
	g.DirectedGraph.SetEdge(e.ReversedEdge())
}

func ExampleComplete_biDirectedSet() {
	dst := simple.NewDirectedGraph()
	err := gen.Complete(Bidirected{dst}, IDSet{2, 4, 5, 9})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict digraph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -> 4;
	// 	2 -> 5;
	// 	2 -> 9;
	// 	4 -> 2;
	// 	4 -> 5;
	// 	4 -> 9;
	// 	5 -> 2;
	// 	5 -> 4;
	// 	5 -> 9;
	// 	9 -> 2;
	// 	9 -> 4;
	// 	9 -> 5;
	// }
}

func ExampleComplete_undirectedSet() {
	dst := simple.NewUndirectedGraph()
	err := gen.Complete(dst, IDSet{2, 4, 5, 9})
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(dst, "complete", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// Output:
	// strict graph complete {
	// 	// Node definitions.
	// 	2;
	// 	4;
	// 	5;
	// 	9;
	//
	// 	// Edge definitions.
	// 	2 -- 4;
	// 	2 -- 5;
	// 	2 -- 9;
	// 	4 -- 5;
	// 	4 -- 9;
	// 	5 -- 9;
	// }
}
