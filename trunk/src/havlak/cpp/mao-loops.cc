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

#include <stdio.h>
#include <list>
#include <set>
#include <vector>
#include <algorithm>

#include "mao-loops.h"

//======================================================
// Main Algorithm
//======================================================

//
// Union/Find algorithm after Tarjan, R.E., 1983, Data Structures
// and Network Algorithms.
//
class UnionFindNode {
 public:
  UnionFindNode() : parent_(NULL), bb_(NULL), loop_(NULL), dfs_number_(0) {
  }

  // Initialize this node.
  //
  void Init(BasicBlock *bb, int dfs_number) {
    parent_     = this;
    bb_         = bb;
    dfs_number_ = dfs_number;
  }

  // Union/Find Algorithm - The find routine.
  //
  // Implemented with Path Compression (inner loops are only
  // visited and collapsed once, however, deep nests would still
  // result in significant traversals).
  //
  UnionFindNode *FindSet() {
    typedef std::list<UnionFindNode *> NodeListType;
    NodeListType nodeList;

    UnionFindNode *node = this;
    while (node != node->parent()) {
      if (node->parent() != node->parent()->parent())
        nodeList.push_back(node);
      node = node->parent();
    }

    // Path Compression, all nodes' parents point to the 1st level parent.
    NodeListType::iterator iter = nodeList.begin();
    NodeListType::iterator end  = nodeList.end();
    for (; iter != end; ++iter)
      (*iter)->set_parent(node->parent());

    return node;
  }

  // Union/Find Algorithm - The union routine.
  //
  // We rely on path compression.
  //
  void Union(UnionFindNode *B) {
    set_parent(B);
  }


  // Getters/Setters
  //
  UnionFindNode *parent() const { return parent_; }
  BasicBlock    *bb() const { return bb_; }
  SimpleLoop    *loop() const { return loop_; }
  int            dfs_number() const { return dfs_number_; }

  void           set_parent(UnionFindNode *parent) { parent_ = parent; }
  void           set_loop(SimpleLoop *loop) { loop_ = loop; }

 private:
  UnionFindNode *parent_;
  BasicBlock    *bb_;
  SimpleLoop    *loop_;
  int            dfs_number_;
};

//------------------------------------------------------------------
// Loop Recognition
//
// based on:
//   Paul Havlak, Nesting of Reducible and Irreducible Loops,
//      Rice University.
//
//   We avoid doing tree balancing and instead use path compression
//   to avoid traversing parent pointers over and over.
//
//   Most of the variable names and identifiers are taken literally
//   from this paper (and the original Tarjan paper mentioned above).
//-------------------------------------------------------------------
class HavlakLoopFinder {
 public:
  HavlakLoopFinder(MaoCFG *cfg, LoopStructureGraph *lsg) :
    CFG_(cfg), lsg_(lsg) {
  }

  enum BasicBlockClass {
    BB_TOP,          // uninitialized
    BB_NONHEADER,    // a regular BB
    BB_REDUCIBLE,    // reducible loop
    BB_SELF,         // single BB loop
    BB_IRREDUCIBLE,  // irreducible loop
    BB_DEAD,         // a dead BB
    BB_LAST          // Sentinel
  };

  //
  // Constants
  //
  // Marker for uninitialized nodes.
  static const int kUnvisited = -1;
  // Safeguard against pathologic algorithm behavior.
  static const int kMaxNonBackPreds = (32*1024);

  //
  // Local types used for Havlak algorithm, all carefully
  // selected to guarantee minimal complexity.
  //
  typedef std::vector<UnionFindNode>          NodeVector;
  typedef std::map<BasicBlock*, int>          BasicBlockMap;
  typedef std::list<int>                      IntList;
  typedef std::set<int>                       IntSet;
  typedef std::list<UnionFindNode*>           NodeList;
  typedef std::vector<IntList>                IntListVector;
  typedef std::vector<IntSet>                 IntSetVector;
  typedef std::vector<int>                    IntVector;
  typedef std::vector<char>                   CharVector;

  //
  // IsAncestor
  //
  // As described in the paper, determine whether a node 'w' is a
  // "true" ancestor for node 'v'.
  //
  // Dominance can be tested quickly using a pre-order trick
  // for depth-first spanning trees. This is why DFS is the first
  // thing we run below.
  //
  bool IsAncestor(int w, int v, IntVector *last) {
    return ((w <= v) && (v <= (*last)[w]));
  }

  // Iterators
  //
  typedef BasicBlock::EdgeVector::iterator BasicBlockIter;

  //
  // DFS - Depth-First-Search
  //
  // DESCRIPTION:
  // Simple depth first traversal along out edges with node numbering.
  //
  int DFS(BasicBlock      *current_node,
          NodeVector      *nodes,
          BasicBlockMap   *number,
          IntVector       *last,
          const int       current) {
    (*nodes)[current].Init(current_node, current);
    (*number)[current_node] = current;

    int lastid = current;
    for (BasicBlockIter outedges = current_node->out_edges()->begin();
         outedges != current_node->out_edges()->end(); ++outedges) {
      BasicBlock *target = *outedges;

      if ((*number)[target] == kUnvisited)
        lastid = DFS(target, nodes, number, last, lastid + 1);
    }
    (*last)[(*number)[current_node]] = lastid;
    return lastid;
  }

  //
  // FindLoops
  //
  // Find loops and build loop forest using Havlak's algorithm, which
  // is derived from Tarjan. Variable names and step numbering has
  // been chosen to be identical to the nomenclature in Havlak's
  // paper (which is similar to the one used by Tarjan).
  //
  void FindLoops() {
    if (!CFG_->GetStartBasicBlock()) return;

    int                size = CFG_->GetNumNodes();

    IntSetVector       non_back_preds(size);
    IntListVector      back_preds(size);
    IntVector          header(size);
    CharVector         type(size);
    IntVector          last(size);
    NodeVector         nodes(size);
    BasicBlockMap      number;

    // Step a:
    //   - initialize all nodes as unvisited.
    //   - depth-first traversal and numbering.
    //   - unreached BB's are marked as dead.
    //
    for (MaoCFG::NodeMap::iterator bb_iter =
           CFG_->GetBasicBlocks()->begin();
         bb_iter != CFG_->GetBasicBlocks()->end(); ++bb_iter) {
      number[(*bb_iter).second] = kUnvisited;
    }

    DFS(CFG_->GetStartBasicBlock(), &nodes, &number, &last, 0);

    // Step b:
    //   - iterate over all nodes.
    //
    //   A backedge comes from a descendant in the DFS tree, and non-backedges
    //   from non-descendants (following Tarjan).
    //
    //   - check incoming edges 'v' and add them to either
    //     - the list of backedges (back_preds) or
    //     - the list of non-backedges (non_back_preds)
    //
    for (int w = 0; w < size; w++) {
      header[w] = 0;
      type[w] = BB_NONHEADER;

      BasicBlock *node_w = nodes[w].bb();
      if (!node_w) {
        type[w] = BB_DEAD;
        continue;  // dead BB
      }

      if (node_w->GetNumPred()) {
        for (BasicBlockIter inedges = node_w->in_edges()->begin();
             inedges != node_w->in_edges()->end(); ++inedges) {
          BasicBlock     *node_v = *inedges;

          int v = number[ node_v ];
          if (v == kUnvisited) continue;  // dead node

          if (IsAncestor(w, v, &last))
            back_preds[w].push_back(v);
          else
            non_back_preds[w].insert(v);
        }
      }
    }

    // Start node is root of all other loops.
    header[0] = 0;

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
    for (int w = size-1; w >= 0; w--) {
      NodeList    node_pool;  // this is 'P' in Havlak's paper
      BasicBlock *node_w = nodes[w].bb();
      if (!node_w) continue;  // dead BB

      // Step d:
      IntList::iterator back_pred_iter  = back_preds[w].begin();
      IntList::iterator back_pred_end   = back_preds[w].end();
      for (; back_pred_iter != back_pred_end; back_pred_iter++) {
        int v = *back_pred_iter;
        if (v != w)
          node_pool.push_back(nodes[v].FindSet());
        else
          type[w] = BB_SELF;
      }

      // Copy node_pool to worklist.
      //
      NodeList worklist;
      NodeList::iterator niter  = node_pool.begin();
      NodeList::iterator nend   = node_pool.end();
      for (;  niter != nend; ++niter)
        worklist.push_back(*niter);

      if (!node_pool.empty())
        type[w] = BB_REDUCIBLE;

      // work the list...
      //
      while (!worklist.empty()) {
        UnionFindNode x = *worklist.front();
        worklist.pop_front();

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
        size_t non_back_size = non_back_preds[x.dfs_number()].size();
        if (non_back_size > kMaxNonBackPreds) {
          lsg_->KillAll();
          return;
        }

        IntSet::iterator non_back_pred_iter =
          non_back_preds[x.dfs_number()].begin();
        IntSet::iterator non_back_pred_end  =
          non_back_preds[x.dfs_number()].end();
        for (; non_back_pred_iter != non_back_pred_end; non_back_pred_iter++) {
          UnionFindNode  y     = nodes[*non_back_pred_iter];
          UnionFindNode *ydash = y.FindSet();

          if (!IsAncestor(w, ydash->dfs_number(), &last)) {
            type[w] = BB_IRREDUCIBLE;
            non_back_preds[w].insert(ydash->dfs_number());
          } else {
            if (ydash->dfs_number() != w) {
              NodeList::iterator nfind = find(node_pool.begin(),
                                              node_pool.end(), ydash);
              if (nfind == node_pool.end()) {
                worklist.push_back(ydash);
                node_pool.push_back(ydash);
              }
            }
          }
        }
      }

      // Collapse/Unionize nodes in a SCC to a single node
      // For every SCC found, create a loop descriptor and link it in.
      //
      if (!node_pool.empty() || (type[w] == BB_SELF)) {
        SimpleLoop* loop = lsg_->CreateNewLoop();

        // At this point, one can set attributes to the loop, such as:
        //
        // the bottom node:
        //    IntList::iterator iter  = back_preds[w].begin();
        //    loop bottom is: nodes[*backp_iter].node);
        //
        // the number of backedges:
        //    back_preds[w].size()
        //
        // whether this loop is reducible:
        //    type[w] != BB_IRREDUCIBLE
        //
        // TODO(rhundt): Define those interfaces in the Loop Forest.
        //
        nodes[w].set_loop(loop);

        for (niter = node_pool.begin(); niter != node_pool.end(); niter++) {
          UnionFindNode  *node = (*niter);

          // Add nodes to loop descriptor.
          header[node->dfs_number()] = w;
          node->Union(&nodes[w]);

          // Nested loops are not added, but linked together.
          if (node->loop())
            node->loop()->set_parent(loop);
          else
            loop->AddNode(node->bb());
        }

        lsg_->AddLoop(loop);
      }  // node_pool.size
    }  // Step c
  }  // FindLoops

 private:
  MaoCFG             *CFG_;      // current control flow graph.
  LoopStructureGraph *lsg_;      // loop forest.
};  // HavlakLoopFinder


// Constant instantiations.
//
const int HavlakLoopFinder::kUnvisited;
const int HavlakLoopFinder::kMaxNonBackPreds;

// External entry point.
int FindHavlakLoops(MaoCFG *CFG, LoopStructureGraph *LSG) {
  HavlakLoopFinder finder(CFG, LSG);
  finder.FindLoops();
  return LSG->GetNumLoops();
}
