// Copyright 2011 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Test Program for the Havlak loop finder.
//
// This program constructs a fairly large control flow
// graph and performs loop recognition. This is the Go
// version.
//
package main

import "fmt"
import "./basicblock"
import "./lsg"
import "./havlakloopfinder"

//======================================================
// Testing Code
//======================================================

func buildDiamond(cfgraph *cfg.CFG, start int) int {
	bb0 := start
	cfg.NewBasicBlockEdge(cfgraph, bb0, bb0+1)
	cfg.NewBasicBlockEdge(cfgraph, bb0, bb0+2)
	cfg.NewBasicBlockEdge(cfgraph, bb0+1, bb0+3)
	cfg.NewBasicBlockEdge(cfgraph, bb0+2, bb0+3)

	return bb0 + 3
}

func buildConnect(cfgraph *cfg.CFG, start int, end int) {
	cfg.NewBasicBlockEdge(cfgraph, start, end)
}

func buildStraight(cfgraph *cfg.CFG, start int, n int) int {
	for i := 0; i < n; i++ {
		buildConnect(cfgraph, start+i, start+i+1)
	}
	return start + n
}

func buildBaseLoop(cfgraph *cfg.CFG, from int) int {
	header := buildStraight(cfgraph, from, 1)
	diamond1 := buildDiamond(cfgraph, header)
	d11 := buildStraight(cfgraph, diamond1, 1)
	diamond2 := buildDiamond(cfgraph, d11)
	footer := buildStraight(cfgraph, diamond2, 1)
	buildConnect(cfgraph, diamond2, d11)
	buildConnect(cfgraph, diamond1, header)

	buildConnect(cfgraph, footer, from)
	footer = buildStraight(cfgraph, footer, 1)
	return footer
}

func main() {
	fmt.Printf("Welcome to LoopTesterApp, Go edition\n")

	lsgraph := lsg.NewLSG()
	cfgraph := cfg.NewCFG()

	fmt.Printf("Constructing Simple CFG...\n")

	cfgraph.CreateNode(0) // top
	buildBaseLoop(cfgraph, 0)
	cfgraph.CreateNode(1) // bottom
	cfg.NewBasicBlockEdge(cfgraph, 0, 2)

	fmt.Printf("15000 dummy loops\n")
	for dummyloop := 0; dummyloop < 15000; dummyloop++ {
		havlakloopfinder.FindHavlakLoops(cfgraph, lsg.NewLSG())
	}

	fmt.Printf("Constructing CFG...\n")
	n := 2

	for parlooptrees := 0; parlooptrees < 10; parlooptrees++ {
		cfgraph.CreateNode(n + 1)
		buildConnect(cfgraph, 2, n+1)
		n = n + 1

		for i := 0; i < 100; i++ {
			top := n
			n = buildStraight(cfgraph, n, 1)
			for j := 0; j < 25; j++ {
				n = buildBaseLoop(cfgraph, n)
			}
			bottom := buildStraight(cfgraph, n, 1)
			buildConnect(cfgraph, n, top)
			n = bottom
		}
		buildConnect(cfgraph, n, 1)
	}

	fmt.Printf("Performing Loop Recognition\n1 Iteration\n")
	havlakloopfinder.FindHavlakLoops(cfgraph, lsgraph)

	fmt.Printf("Another 50 iterations...\n")
	for i := 0; i < 50; i++ {
		fmt.Printf(".")
		havlakloopfinder.FindHavlakLoops(cfgraph, lsg.NewLSG())
	}

	fmt.Printf("\n")
	fmt.Printf("# of loops: %d (including 1 artificial root node)\n", lsgraph.NumLoops())
	lsgraph.CalculateNestingLevel()
}
