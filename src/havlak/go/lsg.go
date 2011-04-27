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

// Support data structures for Havlak's algorithm:
//    Simple Loop representation (SimpleLoop)
//    A loop structure graph (LSG)
//
package lsg

import "container/list"
import "fmt"
import "./basicblock"

//======================================================
// Scaffold Code
//======================================================

// Basic representation of loops, a loop has an entry point,
// one or more exit edges, a set of basic blocks, and potentially
// an outer loop - a "parent" loop.
//
// Furthermore, it can have any set of properties, e.g.,
// it can be an irreducible loop, have control flow, be
// a candidate for transformations, and what not.
//
type SimpleLoop struct {
	// No set, use map to bool
	basicBlocks map[*cfg.BasicBlock]bool
	children    map[*SimpleLoop]bool
	parent      *SimpleLoop
	header      *cfg.BasicBlock

	isRoot       bool
	isReducible  bool
	counter      int
	nestingLevel int
	depthLevel   int
}

func (loop *SimpleLoop) AddNode(bb *cfg.BasicBlock) {
	loop.basicBlocks[bb] = true
}

func (loop *SimpleLoop) AddChildLoop(child *SimpleLoop) {
	loop.children[child] = true
}

func (loop *SimpleLoop) Dump(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("  ")
	}

	// No ? operator ?
	fmt.Printf("loop-%d nest: %d depth %d ",
		loop.counter, loop.nestingLevel, loop.depthLevel)
	if !loop.isReducible {
		fmt.Printf("(Irreducible) ")
	}

	// must have > 0
	if len(loop.children) > 0 {
		fmt.Printf("Children: ")
		for ll, _ := range loop.Children() {
			fmt.Printf("loop-%d", ll.Counter())
		}
	}
	if len(loop.basicBlocks) > 0 {
		fmt.Printf("(")
		for bb, _ := range loop.basicBlocks {
			fmt.Printf("BB#%03d ", bb.Name())
			if loop.header == bb {
				fmt.Printf("*")
			}
		}
		fmt.Printf("\b)")
	}
	fmt.Printf("\n")
}

func (loop *SimpleLoop) Children() map[*SimpleLoop]bool {
	return loop.children
}

func (loop *SimpleLoop) Parent() *SimpleLoop {
	return loop.parent
}

func (loop *SimpleLoop) NestingLevel() int {
	return loop.nestingLevel
}

func (loop *SimpleLoop) DepthLevel() int {
	return loop.depthLevel
}

func (loop *SimpleLoop) Counter() int {
	return loop.counter
}

func (loop *SimpleLoop) IsRoot() bool {
	return loop.isRoot
}

func (loop *SimpleLoop) SetParent(parent *SimpleLoop) {
	loop.parent = parent
	loop.parent.AddChildLoop(loop)
}

func (loop *SimpleLoop) SetHeader(bb *cfg.BasicBlock) {
	loop.AddNode(bb)
	loop.header = bb
}

func (loop *SimpleLoop) SetIsRoot() {
	loop.isRoot = true
}

func (loop *SimpleLoop) SetNestingLevel(level int) {
	loop.nestingLevel = level
}

func (loop *SimpleLoop) SetDepthLevel(level int) {
	loop.depthLevel = level
}

func (loop *SimpleLoop) SetIsReducible(isReducible bool) {
	loop.isReducible = isReducible
}

func (loop *SimpleLoop) SetCounter(value int) {
	loop.counter = value
}


//------------------------------------
// Helper (No templates or such)
//
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// LoopStructureGraph
//
// Maintain loop structure for a given CFG.
//
// Two values are maintained for this loop graph, depth, and nesting level.
// For example:
//
// loop        nesting level    depth
//----------------------------------------
// loop-0      2                0
//   loop-1    1                1
//   loop-3    1                1
//     loop-2  0                2
//
var loopCounter = 0

type LSG struct {
	root  *SimpleLoop
	loops list.List
}

func NewLSG() *LSG {
	lsg := new(LSG)
	lsg.root = lsg.NewLoop()
	lsg.root.SetNestingLevel(0)

	return lsg
}

func (lsg *LSG) NewLoop() *SimpleLoop {
	loop := new(SimpleLoop)
	loop.basicBlocks = make(map[*cfg.BasicBlock]bool)
	loop.children = make(map[*SimpleLoop]bool)
	loop.parent = nil
	loop.header = nil

	loop.SetCounter(loopCounter)
	loopCounter++
	return loop
}

func (lsg *LSG) AddLoop(loop *SimpleLoop) {
	lsg.loops.PushBack(loop)
}

func (lsg *LSG) Dump() {
	lsg.dump(lsg.root, 0)
}

func (lsg *LSG) dump(loop *SimpleLoop, indent int) {
	loop.Dump(indent)

	for ll, _ := range loop.children {
		lsg.dump(ll, indent+1)
	}
}

func (lsg *LSG) CalculateNestingLevel() {
	for ll := lsg.loops.Front(); ll != nil; ll = ll.Next() {
		sl := ll.Value.(*SimpleLoop)
		if sl.IsRoot() {
			continue
		}
		if sl.Parent() == nil {
			sl.SetParent(lsg.root)
		}
	}
	lsg.calculateNestingLevel(lsg.root, 0)
}

func (lsg *LSG) calculateNestingLevel(loop *SimpleLoop, depth int) {
	loop.SetDepthLevel(depth)
	for ll, _ := range loop.children {
		lsg.calculateNestingLevel(ll, depth+1)

		ll.SetNestingLevel(max(loop.NestingLevel(),
			ll.NestingLevel()+1))
	}
}

func (lsg *LSG) NumLoops() int {
	return lsg.loops.Len()
}

func (lsg *LSG) Root() *SimpleLoop {
	return lsg.root
}
