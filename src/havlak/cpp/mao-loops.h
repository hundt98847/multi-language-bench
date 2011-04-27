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

#ifndef MAO_LOOPS_H_
#define MAO_LOOPS_H_

#include <set>
#include <map>
#include <list>
#include <vector>
#include <algorithm>

// Forward Decls
class BasicBlock;
class MaoCFG;

//--- MOCKING CODE begin -------------------
//
// These data structures are stubbed out to make the code below easier
// to review.
//
// BasicBlockEdge only maintains two pointers to BasicBlocks.
//
class BasicBlockEdge {
 public:
  inline BasicBlockEdge(MaoCFG *cfg, int from, int to);

  BasicBlock *GetSrc() { return from_; }
  BasicBlock *GetDst() { return to_; }

 private:
  BasicBlock *from_, *to_;
};

// BasicBlock only maintains a vector of in-edges and
// a vector of out-edges.
//
class BasicBlock {
 public:
  typedef std::vector<BasicBlock *> EdgeVector;

  explicit BasicBlock(int name) : name_(name) {
  }

  EdgeVector *in_edges() { return &in_edges_; }
  EdgeVector *out_edges() { return &out_edges_; }

  int GetNumPred() { return in_edges_.size(); }
  int GetNumSucc() { return out_edges_.size(); }

  void AddOutEdge(BasicBlock *to) { out_edges_.push_back(to); }
  void AddInEdge(BasicBlock *from) { in_edges_.push_back(from); }

 private:
  EdgeVector in_edges_, out_edges_;
  int name_;
};

// MaoCFG maintains a list of nodes.
//
class MaoCFG {
 public:
  typedef std::map<int, BasicBlock *> NodeMap;
  typedef std::list<BasicBlockEdge *> EdgeList;

  MaoCFG() : start_node_(NULL) {
  }

  ~MaoCFG() {
    for (NodeMap::iterator it = basic_block_map_.begin();
         it != basic_block_map_.end(); ++it)
      delete (*it).second;

    for (EdgeList::iterator edge_it = edge_list_.begin();
         edge_it != edge_list_.end(); ++edge_it)
      delete (*edge_it);
  }

  BasicBlock *CreateNode(int name) {
    BasicBlock *node;

    NodeMap::iterator it = basic_block_map_.find(name);

    if (it == basic_block_map_.end()) {
      node = new BasicBlock(name);
      basic_block_map_[name] = node;
    } else {
      node = (*it).second;
    }

    if (GetNumNodes() == 1)
      start_node_ = node;

    return node;
  }

  void AddEdge(BasicBlockEdge *edge) {
    edge_list_.push_back(edge);
  }

  int GetNumNodes() {
    return basic_block_map_.size();
  }

  BasicBlock *GetStartBasicBlock() {
    return start_node_;
  }

  BasicBlock *GetDst(BasicBlockEdge *edge) {
    return edge->GetDst();
  }

  BasicBlock *GetSrc(BasicBlockEdge *edge) {
    return edge->GetSrc();
  }

  NodeMap *GetBasicBlocks() {
    return &basic_block_map_;
  }

 private:
  NodeMap       basic_block_map_;
  BasicBlock   *start_node_;
  EdgeList      edge_list_;
};

//
//--- MOCKING CODE end  -------------------


//
// SimpleLoop
//
// Basic representation of loops, a loop has an entry point,
// one or more exit edges, a set of basic blocks, and potentially
// an outer loop - a "parent" loop.
//
// Furthermore, it can have any set of properties, e.g.,
// it can be an irreducible loop, have control flow, be
// a candidate for transformations, and what not.
//
//
class SimpleLoop {
 public:
  typedef std::set<BasicBlock *> BasicBlockSet;
  typedef std::set<SimpleLoop *> LoopSet;


  SimpleLoop() : parent_(NULL), is_root_(false), nesting_level_(0),
                 depth_level_(0) {
  }

  void AddNode(BasicBlock *basic_block) {
    basic_blocks_.insert(basic_block);
  }

  void AddChildLoop(SimpleLoop *loop) {
    children_.insert(loop);
  }

  void Dump() {
    // Simplified for readability purposes.
    fprintf(stderr, "loop-%d, nest: %d, depth: %d\n",
            counter_, nesting_level_, depth_level_);
  }

  LoopSet *GetChildren() {
    return &children_;
  }

  // Getters/Setters
  SimpleLoop  *parent() { return parent_; }
  int          nesting_level() const { return nesting_level_; }
  int          depth_level() const { return depth_level_; }
  int          counter() const { return counter_; }
  bool         is_root() const { return is_root_; }

  void set_parent(SimpleLoop *parent) {
    parent_ = parent;
    parent->AddChildLoop(this);
  }

  void set_is_root() { is_root_ = true; }
  void set_counter(int value) { counter_ = value; }
  void set_nesting_level(int level) {
    nesting_level_ = level;
    if (level == 0)
      set_is_root();
  }
  void set_depth_level(int level) { depth_level_ = level; }

 private:
  BasicBlockSet          basic_blocks_;
  std::set<SimpleLoop *> children_;
  SimpleLoop            *parent_;

  bool         is_root_: 1;
  int          counter_;
  int          nesting_level_;
  int          depth_level_;
};



//
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
class LoopStructureGraph {
 public:
  typedef std::list<SimpleLoop *> LoopList;

  LoopStructureGraph() : root_(new SimpleLoop()),
                         loop_counter_(0) {
    root_->set_nesting_level(0);  // make it the root node
    root_->set_counter(loop_counter_++);
    AddLoop(root_);
  }

  ~LoopStructureGraph() {
    KillAll();
  }

  SimpleLoop *CreateNewLoop() {
    SimpleLoop *loop = new SimpleLoop();
    loop->set_counter(loop_counter_++);
    return loop;
  }

  void KillAll() {
    for (LoopList::iterator it = loops_.begin(); it != loops_.end(); ++it)
      delete (*it);
  }

  void AddLoop(SimpleLoop *loop) {
    loops_.push_back(loop);
  }

  void Dump() {
    DumpRec(root_, 0);
  }

  void DumpRec(SimpleLoop *loop, int indent) {
    // Simplified for readability purposes.
    loop->Dump();

    for (SimpleLoop::LoopSet::iterator liter = loop->GetChildren()->begin();
         liter != loop->GetChildren()->end(); ++liter)
      DumpRec(*liter,  indent+1);
  }

  void CalculateNestingLevel() {
    // link up all 1st level loops to artificial root node.
    for (LoopList::iterator liter = loops_.begin();
         liter != loops_.end(); ++liter) {
      SimpleLoop *loop = *liter;
      if (loop->is_root()) continue;
      if (!loop->parent()) loop->set_parent(root_);
    }

    // recursively traverse the tree and assign levels.
    CalculateNestingLevelRec(root_, 0);
  }


  void CalculateNestingLevelRec(SimpleLoop *loop, int depth) {
    loop->set_depth_level(depth);
    for (SimpleLoop::LoopSet::iterator liter = loop->GetChildren()->begin();
         liter != loop->GetChildren()->end(); ++liter) {
      CalculateNestingLevelRec(*liter, depth+1);

      loop->set_nesting_level(std::max(loop->nesting_level(),
                                       1+(*liter)->nesting_level()));
    }
  }

  int GetNumLoops() const { return loops_.size(); }

  SimpleLoop *root() const { return root_; }

 private:
  SimpleLoop   *root_;
  LoopList      loops_;
  int           loop_counter_;
};

inline
BasicBlockEdge::BasicBlockEdge(MaoCFG     *cfg,
                               int         from_name,
                               int         to_name) {
  from_ = cfg->CreateNode(from_name);
  to_ = cfg->CreateNode(to_name);

  from_->AddOutEdge(to_);
  to_->AddInEdge(from_);

  cfg->AddEdge(this);
}


// External entry point.
int FindHavlakLoops(MaoCFG *CFG, LoopStructureGraph *LSG);

#endif  // MAO_LOOPS_H_
