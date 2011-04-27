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
// Scaffold Code
//======================================================

// A simple class simulating the concept of Basic Blocks
//
package cfg

import "container/list"
import "fmt"

type BasicBlock struct {
	name     int
	inEdges  list.List
	outEdges list.List
}

func NewBasicBlock(name int) *BasicBlock {
	return &BasicBlock{name: name}
}

func (bb *BasicBlock) Dump() {
	fmt.Printf("BB#%03d: ", bb.Name())
	if bb.NumPred() > 0 {
		fmt.Printf("in : ")
		for iter := bb.InEdges().Front(); iter != nil; iter = iter.Next() {
			fmt.Printf("BB#%03d ", iter.Value.(*BasicBlock).Name())
		}
	}
	if bb.NumSucc() > 0 {
		fmt.Print("out: ")
		for iter := bb.OutEdges().Front(); iter != nil; iter = iter.Next() {
			fmt.Printf("BB#%03d ", iter.Value.(*BasicBlock).Name())
		}
	}
	fmt.Printf("\n")
}

func (bb *BasicBlock) Name() int {
	return bb.name
}

func (bb *BasicBlock) InEdges() *list.List {
	return &bb.inEdges
}

func (bb *BasicBlock) OutEdges() *list.List {
	return &bb.outEdges
}

func (bb *BasicBlock) NumPred() int {
	return bb.inEdges.Len()
}

func (bb *BasicBlock) NumSucc() int {
	return bb.outEdges.Len()
}

func (bb *BasicBlock) AddInEdge(from *BasicBlock) {
	bb.inEdges.PushBack(from)
}

func (bb *BasicBlock) AddOutEdge(to *BasicBlock) {
	bb.outEdges.PushBack(to)
}

//-----------------------------------------------------------

type CFG struct {
	bb        map[int]*BasicBlock
	startNode *BasicBlock
}

func NewCFG() *CFG {
	return &CFG{bb: make(map[int]*BasicBlock)}
}

func (cfg *CFG) BasicBlocks() map[int]*BasicBlock {
	return cfg.bb
}

func (cfg *CFG) NumNodes() int {
	return len(cfg.bb)
}

func (cfg *CFG) CreateNode(node int) *BasicBlock {
	if bblock := cfg.bb[node]; bblock != nil {
		return bblock
	}
	bblock := NewBasicBlock(node)
	cfg.bb[node] = bblock

	if cfg.NumNodes() == 1 {
		cfg.startNode = bblock
	}

	return bblock
}

func (cfg *CFG) Dump() {
	for _, n := range cfg.bb {
		n.Dump()
	}
}

func (cfg *CFG) StartBasicBlock() *BasicBlock {
	return cfg.startNode
}

func (cfg *CFG) Dst(edge *BasicBlockEdge) *BasicBlock {
	return edge.Dst()
}

func (cfg *CFG) Src(edge *BasicBlockEdge) *BasicBlock {
	return edge.Src()
}

//-----------------------------------------------------------

type BasicBlockEdge struct {
	to   *BasicBlock
	from *BasicBlock
}

func (edge *BasicBlockEdge) Dst() *BasicBlock {
	return edge.to
}

func (edge *BasicBlockEdge) Src() *BasicBlock {
	return edge.from
}

func NewBasicBlockEdge(cfg *CFG, from int, to int) *BasicBlockEdge {
	self := new(BasicBlockEdge)
	self.to = cfg.CreateNode(to)
	self.from = cfg.CreateNode(from)

	self.from.AddOutEdge(self.to)
	self.to.AddInEdge(self.from)

	return self
}

//-----------------------------------------------------------
