// Copyright 2012 Google Inc.
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

// BasicBlock's static members
//
var numBasicBlocks = 0;

function getNumBasicBlocks() {
  return numBasicBlocks;
}

//
// class BasicBlock
//
// BasicBlock only maintains a vector of in-edges and
// a vector of out-edges.
//
function BasicBlock(name)
{
  this.name = name;
  this.inEdges  = [];
  this.outEdges = [];

  numBasicBlocks = numBasicBlocks + 1;
}

BasicBlock.prototype.toString = function() {
  return "BB" + this.name;
}

BasicBlock.prototype.getNumPred = function() {
  return this.inEdges.length;
}
BasicBlock.prototype.getNumSucc = function() {
  return this.outEdges.length;
}
BasicBlock.prototype.addInEdge = function(bb) {
  this.inEdges.push(bb);
}
BasicBlock.prototype.addOutEdge = function(bb) {
  this.outEdges.push(bb);
}

BasicBlock.prototype.dump = function() {
  var res = "  BB#" + this.name;
  if (this.inEdges.length > 0) {
    res += "\tin :";
    for (var i = 0; i < this.inEdges.length; ++i) {
      res += " " + this.inEdges[i];
    }
  }
  if (this.outEdges.length > 0) {
    res += "\tout:";
    for (var i = 0; i < this.outEdges.length; ++i) {
      res += " " + this.outEdges[i];
    }
  }

  print(res);
}



//
// class CFG
//
// CFG maintains a list of nodes, plus a start node.
// That's it.
//
function CFG() {
  this.startNode = null;
  this.basicBlockMap = {};
  this.basicBlockLen = 0;
  this.edgeList      = [];
}

CFG.prototype.createNode = function(name) {
  var node;
  if (name in this.basicBlockMap) {
    node = this.basicBlockMap[name];
  } else {
    node = new BasicBlock(name);
    this.basicBlockMap[name] = node;
    this.basicBlockLen++;
  }

  if (this.getNumNodes() == 1) {
    this.startNode = node;
  }
  return node;
}

CFG.prototype.dump = function() {
  for (var k in this.basicBlockMap) {
    this.basicBlockMap[k].dump();
  }
}

CFG.prototype.addEdge = function(edge) {
  this.edgeList.push(edge);
}
CFG.prototype.getNumNodes = function() {
  return this.basicBlockLen;
}
CFG.prototype.getDst = function(edge) {
  return edge.To;
}
CFG.prototype.getSrc = function(edge) {
  return edge.From;
}

//
// class BasicBlockEdge
//
// These data structures are stubbed out to make the code below easier
// to review.
//
// BasicBlockEdge only maintains two pointers to BasicBlocks.
// Note: from is apparently a keyword in python. Changed to uppercase
//
function BasicBlockEdge(cfg, fromName, toName) {
  this.From = cfg.createNode(fromName);
  this.To   = cfg.createNode(toName);

  this.From.addOutEdge(this.To);
  this.To.addInEdge(this.From);

  cfg.addEdge(this);
}

//
// class SimpleLoop
//
// Basic representation of loops, a loop has an entry point,
// one or more exit edges, a set of basic blocks, and potentially
// an outer loop - a "parent" loop.
//
// Furthermore, it can have any set of properties, e.g.,
// it can be an irreducible loop, have control flow, be
// a candidate for transformations, and what not.
//
function SimpleLoop() {
  this.basicBlocks  = [];
  this.children     = [];

  this.isRoot       = false;
  this.isReducible  = true;
  this.counter      = 0;
  this.nestingLevel = 0;
  this.depthLevel   = 0;
}

SimpleLoop.prototype.addNode = function(bb) {
  this.basicBlocks.push(bb);
}

SimpleLoop.prototype.addChildLoop = function(loop) {
  this.children.push(loop);
}

SimpleLoop.prototype.dump = function(indent) {
  var res = "";
  for (i=0; i < indent; ++i) {
    res += "  ";
  }

  print(res, "loop-", this.counter, " nest: ",
        this.nestingLevel, " depth ", this.depthLevel );
  if (this.isReducible == false) {
    print("  irreducible");
  }
}

SimpleLoop.prototype.setParent = function(parent) {
  this.parent = parent;
  parent.addChildLoop(this);
}

SimpleLoop.prototype.setHeader = function(bb) {
  this.basicBlocks.push(bb);
  this.header = bb;
}

SimpleLoop.prototype.setNestingLevel = function(level) {
  this.nestingLevel = level;
  if (level == 0) {
    this.isRoot = true;
  }
}


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
function LSG() {
  this.loopCounter = 0;
  this.loops = [];
  this.root  = new SimpleLoop();
  this.root.setNestingLevel(0);
  this.root.counter = this.loopCounter;
  this.loopCounter += 1;
  this.loops.push(this.root);
}

LSG.prototype.createNewLoop = function() {
  var loop = new SimpleLoop();
  loop.counter = this.loopCounter;
  this.loopCounter += 1;
  return loop;
}

LSG.prototype.addLoop = function(loop) {
  this.loops.push(loop);
}

LSG.prototype.dumpRec = function(loop, indent) {
  loop.dump(indent)

  for (var liter = 0; liter < loop.children.length; ++liter) {
    this.dumpRec(loop.children[liter], indent + 1);
  }
}

LSG.prototype.dump = function() {
  this.dumpRec(this.root,0);
}

LSG.prototype.max = function(a,b) {
  if (a>b)
    return a;
  return b;
}

LSG.prototype.calculateNestingLevelRec = function(loop, depth) {
  loop.depthLevel = depth;
  for (var liter = 0; liter < loop.children.length; ++liter) {
    this.calculateNestingLevelRec(loop.children[liter], depth+1);

    loop.setNestingLevel(this.max(loop.nestingLevel,
                                  1+loop.children[liter].nestingLevel));
  }
}

LSG.prototype.calculateNestingLevel = function() {
  for (var liter = 0; liter < this.loops.length; ++liter) {
    if (!this.loops[liter].isRoot) {
      if (!this.loops[liter].hasOwnProperty('parent')) {
        this.loops[liter].setParent(this.root);
      }
    }
  }

  this.calculateNestingLevelRec(this.root, 0);
}

LSG.prototype.getNumLoops = function() {
  return this.loops.length;
}


//======================================================
// Main Algorithm
//======================================================

//
// class UnionFindNode
//
// The algorithm uses the Union/Find algorithm to collapse
// complete loops into a single node. These nodes and the
// corresponding functionality are implemented with this class
//
function UnionFindNode() {
  this.dfsNumber = 0;
}

// Initialize this node.
//
UnionFindNode.prototype.initNode = function(bb, dfsNumber) {
  this.parent     = this;
  this.bb         = bb;
  this.dfsNumber  = dfsNumber;
}

// Union/Find Algorithm - The find routine.
//
// Implemented with Path Compression (inner loops are only
// visited and collapsed once, however, deep nests would still
// result in significant traversals).
//
UnionFindNode.prototype.findSet = function() {
  var nodeList = [];

  node = this;
  while (node != node.parent) {
    if (node.parent != node.parent.parent)
      nodeList.push(node);

    node = node.parent;
  }

  // Path Compression, all nodes' parents point to the 1st level parent.
  for (iter=0; iter < nodeList.length; ++iter) {
    nodeList[iter].parent = node.parent;
  }

  return node;
}

// Union/Find Algorithm - The union routine.
//
// Trivial. Assigning parent pointer is enough,
// we rely on path compression.
//
UnionFindNode.prototype.union = function(unionFindNode) {
  this.parent = unionFindNode;
}



function HavlakLoopFinder(cfgParm, lsgParm) {
  this.cfg = cfgParm;
  this.lsg = lsgParm;

  //
  // enum BasicBlockClass
  //
  // Basic Blocks and Loops are being classified as regular, irreducible,
  // and so on. This enum contains a symbolic name for all these
  // classifications. Python doesn't have enums, so we just create values.
  //
  var BB_TOP          = 0; // uninitialized
  var BB_NONHEADER    = 1; // a regular BB
  var BB_REDUCIBLE    = 2; // reducible loop
  var BB_SELF         = 3; // single BB loop
  var BB_IRREDUCIBLE  = 4; // irreducible loop
  var BB_DEAD         = 5; // a dead BB
  var BB_LAST         = 6; // Sentinel

  //
  // Constants
  //
  // Marker for uninitialized nodes.
  var UNVISITED = -1;

  // Safeguard against pathologic algorithm behavior.
  var MAXNONBACKPREDS = (32 * 1024);

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
  function isAncestor(w, v, last) {
    return (w <= v) && (v <= last[w]);
  }


  var depth = 0;
  var maxdepth = 0;
  var total = 0;
  //
  // DFS - Depth-First-Search
  //
  // DESCRIPTION:
  // Simple depth first traversal along out edges with node numbering.
  //
  function DFS(currentNode, nodes, number, last, current) {
    ++total;
    nodes[current].initNode(currentNode, current);
    number[currentNode] = current;

    var lastid = current;
    for (target = 0; target < currentNode.outEdges.length; target++) {
      if (number[currentNode.outEdges[target]] == UNVISITED) {
        depth += 1;
        if (depth > maxdepth)
          maxdepth = depth;
        lastid = DFS(currentNode.outEdges[target], nodes, number,
                     last, lastid + 1);
        depth -= 1;
      }
    }

    last[number[currentNode]] = lastid;
    return lastid;
  }

  //
  // findLoops
  //
  // Find loops and build loop forest using Havlak's algorithm, which
  // is derived from Tarjan. Variable names and step numbering has
  // been chosen to be identical to the nomenclature in Havlak's
  // paper (which, in turn, is similar to the one used by Tarjan).
  //
  this.findLoops = function() {
    var size = this.cfg.getNumNodes();

    var nonBackPreds    = [];
    var backPreds       = [];
    var number          = {};
    var header          = [];
    var types           = [];
    var last            = [];
    var nodes           = [];

    for (var i=0; i < size; ++i) {
      nonBackPreds.push([]);
      backPreds.push([]);
      header.push(0);
      types.push(0);
      last.push(0);
      nodes.push(new UnionFindNode());
    }

    // Step a:
    //   - initialize all nodes as unvisited.
    //   - depth-first traversal and numbering.
    //   - unreached BB's are marked as dead.
    //
    for (var k in this.cfg.basicBlockMap) {
      number[this.cfg.basicBlockMap[k]] = UNVISITED;
    }

    DFS(cfg.startNode, nodes, number, last, 0);

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
    for (var w = 0; w < size; ++w) {
      header[w] = 0;
      types[w]  = BB_NONHEADER;

      var nodeW = nodes[w].bb;
      if (nodeW == undefined) {
        types[w] = BB_DEAD;
      } else {
        if (nodeW.getNumPred() > 0) {
          for (var nv = 0; nv < nodeW.inEdges.length; ++nv) {
            var nodeV = nodeW.inEdges[nv];
            var v = number[nodeV];
            if (v != UNVISITED) {
              if (isAncestor(w, v, last)) {
                backPreds[w].push(v);
              } else {
                nonBackPreds[w].push(v);
              }
            }
          }
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
    for (var w = size-1; w >=0; --w) {
      // this is 'P' in Havlak's paper
      var nodePool = [];

      var nodeW = nodes[w].bb;
      if (nodeW == undefined) {
        continue;
      }

      // Step d:
      for (var vi = 0; vi < backPreds[w].length; ++vi) {
        var v = backPreds[w][vi];
        if (v != w) {
          nodePool.push(nodes[v].findSet());
        } else {
          types[w] = BB_SELF;
        }
      }

      // Copy nodePool to workList.
      //
      var workList = [];
      for (var n = 0; n < nodePool.length; ++n) {
        workList.push(nodePool[n]);
      }

      if (nodePool.length != 0) {
          types[w] = BB_REDUCIBLE;
      }
      // work the list...
      //
      while (workList.length) {
        var x = workList.shift();
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
        var nonBackSize = nonBackPreds[x.dfsNumber].length;
        if (nonBackSize > MAXNONBACKPREDS) {
          return 0;
        }

        for (var iter=0; iter < nonBackPreds[x.dfsNumber].length; ++iter) {
          var y = nodes[nonBackPreds[x.dfsNumber][iter]];
          var ydash = y.findSet();

          if (!isAncestor(w, ydash.dfsNumber, last)) {
            types[w] = BB_IRREDUCIBLE;
            nonBackPreds[w].push(ydash.dfsNumber);
          } else {
            if (ydash.dfsNumber != w) {
              if (nodePool.indexOf(ydash) == -1) {
                workList.push(ydash);
                nodePool.push(ydash);
              }
            }
          }
        }
      }

      // Collapse/Unionize nodes in a SCC to a single node
      // For every SCC found, create a loop descriptor and link it in.
      //
      if ((nodePool.length > 0) || (types[w] == BB_SELF)) {
        var loop = this.lsg.createNewLoop();

        loop.setHeader(nodeW);
        if (types[w] == BB_IRREDUCIBLE) {
          loop.isReducible = true;
        } else {
          loop.isReducible = false;
        }

        // At this point, one can set attributes to the loop, such as:
        //
        // the bottom node:
        //    iter  = backPreds(w).begin();
        //    loop bottom is: nodes(iter).node;
        //
        // the number of backedges:
        //    backPreds(w).size()
        //
        // whether this loop is reducible:
        //    types(w) != BB_IRREDUCIBLE
        //
        nodes[w].loop = loop;

        for (var np = 0; np < nodePool.length; ++np) {
          var node = nodePool[np];

          // Add nodes to loop descriptor.
          header[node.dfsNumber] = w;
          node.union(nodes[w]);

          // Nested loops are not added, but linked together.
          if (node.loop != undefined) {
            node.loop.setParent(loop);
          } else {
            loop.addNode(node.bb);
          }
        }
        this.lsg.addLoop(loop);
      } // nodePool.length
    } // Step c
    return this.lsg.getNumLoops();
  } // findLoops
} // HavlakLoopFinder


//======================================================
// Testing Code
//======================================================

function buildDiamond(cfg, start) {
  var bb0 = start;
  new BasicBlockEdge(cfg, bb0, bb0 + 1);
  new BasicBlockEdge(cfg, bb0, bb0 + 2);
  new BasicBlockEdge(cfg, bb0 + 1, bb0 + 3);
  new BasicBlockEdge(cfg, bb0 + 2, bb0 + 3);
  return bb0 + 3;
}


function buildConnect(cfg, start, end) {
  new BasicBlockEdge(cfg, start, end);
}

function buildStraight(cfg, start, n) {
  for (var i=0; i < n; i++) {
    buildConnect(cfg, start + i, start + i + 1);
  }
  return start + n;
}

function buildBaseLoop(cfg, From) {
  var header   = buildStraight(cfg, From, 1);
  var diamond1 = buildDiamond(cfg, header);
  var d11      = buildStraight(cfg, diamond1, 1);
  var diamond2 = buildDiamond(cfg, d11);
  var footer   = buildStraight(cfg, diamond2, 1);
  buildConnect(cfg, diamond2, d11);
  buildConnect(cfg, diamond1, header);

  buildConnect(cfg, footer, From);
  footer = buildStraight(cfg, footer, 1);
  return footer;
}

var cfg = new CFG();
var lsg = new LSG();

cfg.createNode(0);  // top
//buildBaseLoop(cfg, 0);
cfg.createNode(1);  //s bottom
buildConnect(cfg, 0, 2);

// execute loop recognition 15000 times to force compilation
print("15000 dummy loops");

for (var dummyloop = 0; dummyloop < 15000; ++dummyloop) {
  var lsglocal = new LSG();
  var finder = new HavlakLoopFinder(cfg, lsglocal);
  var x = finder.findLoops();
}

var n = 2;

for (var parlooptrees=0; parlooptrees < 10; parlooptrees++) {
  cfg.createNode(n + 1);
  buildConnect(cfg, n, n + 1);
  n = n + 1;
  for (var i=0; i < 2; ++i) {
    var topNode = n;
    n = buildStraight(cfg, n, 1);
    for (var j=0; j < 25; j++) {
      n = buildBaseLoop(cfg, n);
    }

    var bottom = buildStraight(cfg, n, 1);
    buildConnect(cfg, n, topNode);
    n = bottom;
  }
}
print("Performing Loop Recognition\n1 Iteration");

finder = new HavlakLoopFinder(cfg, lsg);
var num_loops = finder.findLoops();
lsg.calculateNestingLevel();
print("Another 100 iterations...")

for (var i = 0; i < 100; i++) {
  //  print(i);
  var lsglocal = new LSG();
  var finder = new HavlakLoopFinder(cfg, lsglocal);
  num_loops = finder.findLoops();
}

print("Found ", num_loops, " loops (including artificial root node)" );
