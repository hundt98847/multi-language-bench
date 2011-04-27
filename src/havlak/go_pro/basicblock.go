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

import "fmt"

type BasicBlock struct {
	Name     int
	InEdges  []*BasicBlock
	OutEdges []*BasicBlock
}

func NewBasicBlock(name int) *BasicBlock {
	return &BasicBlock{Name: name}
}

func (bb *BasicBlock) Dump() {
	fmt.Printf("BB#%03d: ", bb.Name)
	if len(bb.InEdges) > 0 {
		fmt.Printf("in : ")
		for _, iter := range bb.InEdges {
			fmt.Printf("BB#%03d ", iter.Name)
		}
	}
	if len(bb.OutEdges) > 0 {
		fmt.Print("out: ")
		for _, iter := range bb.OutEdges {
			fmt.Printf("BB#%03d ", iter.Name)
		}
	}
	fmt.Printf("\n")
}

func (bb *BasicBlock) NumPred() int {
	return len(bb.InEdges)
}

func (bb *BasicBlock) NumSucc() int {
	return len(bb.OutEdges)
}

func (bb *BasicBlock) AddInEdge(from *BasicBlock) {
	bb.InEdges = append(bb.InEdges, from)
}

func (bb *BasicBlock) AddOutEdge(to *BasicBlock) {
	bb.OutEdges = append(bb.OutEdges, to)
}

//-----------------------------------------------------------

type CFG struct {
	Blocks map[int]*BasicBlock
	Start  *BasicBlock
}

func NewCFG() *CFG {
	return &CFG{Blocks: make(map[int]*BasicBlock)}
}

func (cfg *CFG) NumNodes() int {
	return len(cfg.Blocks)
}

func (cfg *CFG) CreateNode(node int) *BasicBlock {
	if bblock := cfg.Blocks[node]; bblock != nil {
		return bblock
	}
	bblock := NewBasicBlock(node)
	cfg.Blocks[node] = bblock

	if len(cfg.Blocks) == 1 {
		cfg.Start = bblock
	}

	return bblock
}

func (cfg *CFG) Dump() {
	for _, n := range cfg.Blocks {
		n.Dump()
	}
}

//-----------------------------------------------------------

type BasicBlockEdge struct {
	Dst *BasicBlock
	Src *BasicBlock
}

func NewBasicBlockEdge(cfg *CFG, from int, to int) *BasicBlockEdge {
	self := new(BasicBlockEdge)
	self.Dst = cfg.CreateNode(to)
	self.Src = cfg.CreateNode(from)

	self.Src.AddOutEdge(self.Dst)
	self.Dst.AddInEdge(self.Src)

	return self
}

//-----------------------------------------------------------
