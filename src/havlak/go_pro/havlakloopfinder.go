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

//======================================================
// Main Algorithm
//======================================================

//
// The Havlak loop finding algorithm.
//
package havlakloopfinder

import "./basicblock"
import "./lsg"

// Basic Blocks and Loops are being classified as regular, irreducible,
// and so on. This enum contains a symbolic name for all these classifications
//
const (
	_             = iota // Go has an interesting iota concept
	bbTop                // uninitialized
	bbNonHeader          // a regular BB
	bbReducible          // reducible loop
	bbSelf               // single BB loop
	bbIrreducible        // irreducible loop
	bbDead               // a dead BB
	bbLast               // sentinel
)

// UnionFindNode is used in the Union/Find algorithm to collapse
// complete loops into a single node. These nodes and the
// corresponding functionality are implemented with this class
//
type UnionFindNode struct {
	parent    *UnionFindNode
	bb        *cfg.BasicBlock
	loop      *lsg.SimpleLoop
	dfsNumber int
}

// Init explicitly initializes UnionFind nodes.
//
func (u *UnionFindNode) Init(bb *cfg.BasicBlock, dfsNumber int) {
	u.parent = u
	u.bb = bb
	u.dfsNumber = dfsNumber
	u.loop = nil
}

// FindSet implements the Find part of the Union/Find Algorithm
//
// Implemented with Path Compression (inner loops are only
// visited and collapsed once, however, deep nests would still
// result in significant traversals).
//
func (u *UnionFindNode) FindSet() *UnionFindNode {
	var nodeList []*UnionFindNode
	node := u

	for ; node != node.parent; node = node.parent {
		if node.parent != node.parent.parent {
			nodeList = append(nodeList, node)
		}

	}

	// Path Compression, all nodes' parents point to the 1st level parent.
	for _, ll := range nodeList {
		ll.parent = node.parent
	}

	return node
}

// Union relies on path compression.
//
func (u *UnionFindNode) Union(B *UnionFindNode) {
	u.parent = B
}


// Constants
//
// Marker for uninitialized nodes.
const unvisited = -1

// Safeguard against pathological algorithm behavior.
const maxNonBackPreds = 32 * 1024

// IsAncestor
//
// As described in the paper, determine whether a node 'w' is a
// "true" ancestor for node 'v'.
//
// Dominance can be tested quickly using a pre-order trick
// for depth-first spanning trees. This is why DFS is the first
// thing we run below.
//
// Go comment: Parameters can be written as w,v int, inlike in C, where
//   each parameter needs its own type.
//
func isAncestor(w, v int, last []int) bool {
	return ((w <= v) && (v <= last[w]))
}


// listContainsNode
//
// Check whether a list contains a specific element. 
//
func listContainsNode(l []*UnionFindNode, u *UnionFindNode) bool {
	for _, ll := range l {
		if ll == u {
			return true
		}
	}
	return false
}

// DFS - Depth-First-Search and node numbering.
//
func DFS(currentNode *cfg.BasicBlock, nodes []*UnionFindNode, number map[*cfg.BasicBlock]int, last []int, current int) int {
	nodes[current].Init(currentNode, current)
	number[currentNode] = current

	lastid := current
	for _, target := range currentNode.OutEdges {
		if number[target] == unvisited {
			lastid = DFS(target, nodes, number, last, lastid+1)
		}
	}
	last[number[currentNode]] = lastid
	return lastid
}

// FindLoops
//
// Find loops and build loop forest using Havlak's algorithm, which
// is derived from Tarjan. Variable names and step numbering has
// been chosen to be identical to the nomenclature in Havlak's
// paper (which, in turn, is similar to the one used by Tarjan).
//
func FindLoops(cfgraph *cfg.CFG, lsgraph *lsg.LSG) {
	if cfgraph.Start == nil {
		return
	}

	size := cfgraph.NumNodes()

	nonBackPreds := make([]map[int]bool, size)
	backPreds := make([][]int, size)

	number := make(map[*cfg.BasicBlock]int)
	header := make([]int, size, size)
	types := make([]int, size, size)
	last := make([]int, size, size)
	nodes := make([]*UnionFindNode, size, size)

	for i := 0; i < size; i++ {
		nodes[i] = new(UnionFindNode)
	}

	// Step a:
	//   - initialize all nodes as unvisited.
	//   - depth-first traversal and numbering.
	//   - unreached BB's are marked as dead.
	//
	for i, bb := range cfgraph.Blocks {
		number[bb] = unvisited
		nonBackPreds[i] = make(map[int]bool)
	}

	DFS(cfgraph.Start, nodes, number, last, 0)

	// Step b:
	//   - iterate over all nodes.
	//
	//   A backedge comes from a descendant in the DFS tree, and non-backedges
	//   from non-descendants (following Tarjan).
	//
	//   - check incoming edges 'v' and add them to either
	//     - the list of backedges (backPreds) or
	//     - the list of non-backedges (nonBackPreds)
	//
	for w := 0; w < size; w++ {
		header[w] = 0
		types[w] = bbNonHeader

		nodeW := nodes[w].bb
		if nodeW == nil {
			types[w] = bbDead
			continue // dead BB
		}

		if nodeW.NumPred() > 0 {
			for _, nodeV := range nodeW.InEdges {
				v := number[nodeV]
				if v == unvisited {
					continue // dead node
				}

				if isAncestor(w, v, last) {
					backPreds[w] = append(backPreds[w], v)
				} else {
					nonBackPreds[w][v] = true
				}
			}
		}
	}

	// Start node is root of all other loops.
	header[0] = 0

	// Step c:
	//
	// The outer loop, unchanged from Tarjan. It does nothing except
	// for those nodes which are the destinations of backedges.
	// For a header node w, we chase backward from the sources of the
	// backedges adding nodes to the set P, representing the body of
	// the loop headed by w.
	//
	// By running through the nodes in reverse of the DFST preorder,
	// we ensure that inner loop headers will be processed before the
	// headers for surrounding loops.
	//
	for w := size - 1; w >= 0; w-- {
		// this is 'P' in Havlak's paper
		var nodePool []*UnionFindNode

		nodeW := nodes[w].bb
		if nodeW == nil {
			continue // dead BB
		}

		// Step d:
		for _, v := range backPreds[w] {
			if v != w {
				nodePool = append(nodePool, nodes[v].FindSet())
			} else {
				types[w] = bbSelf
			}
		}

		// Copy nodePool to workList.
		//
		workList := append([]*UnionFindNode(nil), nodePool...)

		if len(nodePool) != 0 {
			types[w] = bbReducible
		}

		// work the list...
		//
		for len(workList) > 0 {
			x := workList[0]
			workList = workList[1:]

			// Step e:
			//
			// Step e represents the main difference from Tarjan's method.
			// Chasing upwards from the sources of a node w's backedges. If
			// there is a node y' that is not a descendant of w, w is marked
			// the header of an irreducible loop, there is another entry
			// into this loop that avoids w.
			//

			// The algorithm has degenerated. Break and
			// return in this case.
			//
			nonBackSize := len(nonBackPreds[x.dfsNumber])
			if nonBackSize > maxNonBackPreds {
				return
			}

			for iter := range nonBackPreds[x.dfsNumber] {
				y := nodes[iter]
				ydash := y.FindSet()

				if !isAncestor(w, ydash.dfsNumber, last) {
					types[w] = bbIrreducible
					nonBackPreds[w][ydash.dfsNumber] = true
				} else {
					if ydash.dfsNumber != w {
						if !listContainsNode(nodePool, ydash) {
							workList = append(workList, ydash)
							nodePool = append(nodePool, ydash)
						}
					}
				}
			}
		}

		// Collapse/Unionize nodes in a SCC to a single node
		// For every SCC found, create a loop descriptor and link it in.
		//
		if (len(nodePool) > 0) || (types[w] == bbSelf) {
			loop := lsgraph.NewLoop()

			loop.SetHeader(nodeW)
			if types[w] != bbIrreducible {
				loop.IsReducible = true
			}

			// At this point, one can set attributes to the loop, such as:
			//
			// the bottom node:
			//    iter  = backPreds[w].begin();
			//    loop bottom is: nodes[iter].node);
			//
			// the number of backedges:
			//    backPreds[w].size()
			//
			// whether this loop is reducible:
			//    type[w] != BasicBlockClass.bbIrreducible
			//
			nodes[w].loop = loop

			for _, node := range nodePool {
				// Add nodes to loop descriptor.
				header[node.dfsNumber] = w
				node.Union(nodes[w])

				// Nested loops are not added, but linked together.
				if node.loop != nil {
					node.loop.Parent = loop
				} else {
					loop.AddNode(node.bb)
				}
			}

			lsgraph.AddLoop(loop)
		} // nodePool.size
	} // Step c


}

// External entry point.
func FindHavlakLoops(cfgraph *cfg.CFG, lsgraph *lsg.LSG) int {
	FindLoops(cfgraph, lsgraph)
	return lsgraph.NumLoops()
}
