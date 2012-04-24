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
int numBasicBlocks = 0;
int getNumBasicBlocks() => numBasicBlocks;

//
// class BasicBlock
//
// BasicBlock only maintains a vector of in-edges and
// a vector of out-edges.
//
class BasicBlock implements Hashable
{
  int name;
  List<BasicBlock> inEdges;
  List<BasicBlock> outEdges;
  int hashCode() => name;

  BasicBlock(this.name) {
    inEdges = new List<BasicBlock>();
    outEdges = new List<BasicBlock>();
    numBasicBlocks = numBasicBlocks + 1;
  }

  String toString() => "BB" + this.name;
  int    getNumPred() => inEdges.length;
  int    getNumSucc() => outEdges.length;
  void   addInEdge(BasicBlock bb) => inEdges.add(bb);
  void   addOutEdge(BasicBlock bb) => outEdges.add(bb);

  int dump() {
    var res = "  BB#" + this.name;
    if (inEdges.length > 0) {
      res += "\tin :";
      for (var i = 0; i < inEdges.length; ++i) {
        res += " " + inEdges[i].toString();
      }
    }

    if (outEdges.length > 0) {
      res += "\tout:";
      for (var i = 0; i < outEdges.length; ++i) {
        res += " " + outEdges[i].toString();
      }
    }

    print(res);
  }
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
class BasicBlockEdge {
  BasicBlock From, To;

  BasicBlockEdge(CFG cfg, int fromName, int toName) {
    From = cfg.createNode(fromName);
    To   = cfg.createNode(toName);

    From.addOutEdge(To);
    To.addInEdge(From);

    cfg.addEdge(this);
  }
}


//
// class CFG
//
// CFG maintains a list of nodes, plus a start node.
// That's it.
//
class CFG {
  Map<int, BasicBlock> basicBlockMap;
  List<BasicBlockEdge>  edgeList;
  BasicBlock            startNode;

  CFG() {
    basicBlockMap = new Map<int, BasicBlock>();
    edgeList      = new List<BasicBlockEdge>();
    startNode     = null;
  }

  BasicBlock createNode(int name) {
    BasicBlock node;
    if (basicBlockMap.containsKey(name)) {
      node = basicBlockMap[name];
    } else {
      node = new BasicBlock(name);
      basicBlockMap[name] = node;
    }

    if (getNumNodes() == 1) {
      startNode = node;
    }
    return node;
  }

  void dump() {
    for (var k in basicBlockMap.getValues()) {
      k.dump();
    }
  }
  void addEdge(BasicBlockEdge edge) => edgeList.add(edge);
  int getNumNodes() => basicBlockMap.length;
  BasicBlock getDst(edge) => edge.To;
  BasicBlock getSrc(edge) => edge.From;
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
class SimpleLoop implements Hashable {
  List<BasicBlock> basicBlocks;
  List<SimpleLoop> children;
  SimpleLoop      parent;
  BasicBlock      header;

  bool isRoot;
  bool isReducible;
  int  counter;
  int  nestingLevel;
  int  depthLevel;

  int hashCode() => counter;

  SimpleLoop() {
    children     = new List<SimpleLoop>();
    basicBlocks  = new List<BasicBlock>();

    isRoot       = false;
    isReducible  = true;
    counter      = 0;
    nestingLevel = 0;
    depthLevel   = 0;
    parent       = null;
  }

  void addNode(BasicBlock bb) => basicBlocks.add(bb);
  void addChildLoop(SimpleLoop loop) => children.add(loop);

  void dump(int indent) {
    String res = "";
    for (int i=0; i < indent; ++i) {
      res += "  ";
    }

    print("$res loop-$counter, nest: $nestingLevel, depth $depthLevel");
    if (isReducible == false) {
      print("  irreducible");
    }
  }

  void setParent(SimpleLoop p) {
    this.parent = p;
    p.addChildLoop(this);
  }

  void setHeader(BasicBlock bb) {
    basicBlocks.add(bb);
    header = bb;
  }

  void setNestingLevel(int level) {
    nestingLevel = level;
    if (level == 0) {
      isRoot = true;
    }
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
class LSG {
  int loopCounter;
  List<SimpleLoop> loops;
  SimpleLoop       root;

  LSG() {
    loopCounter = 0;
    loops = new List<SimpleLoop>();
    root  = new SimpleLoop();
    root.setNestingLevel(0);
    root.counter = loopCounter;
    loopCounter += 1;
    loops.add(this.root);
  }

  SimpleLoop createNewLoop() {
    SimpleLoop loop = new SimpleLoop();
    loop.counter = loopCounter;
    loopCounter += 1;
    return loop;
  }

  void addLoop(SimpleLoop loop) => loops.add(loop);

  void dumpRec(SimpleLoop loop, int indent) {
    loop.dump(indent);

    for (int liter = 0; liter < loop.children.length; ++liter) {
      dumpRec(loop.children[liter], indent + 1);
    }
  }

  void dump() => dumpRec(root,0);

  int max(int a, int b) {
    if (a>b)
      return a;
    return b;
  }

  calculateNestingLevelRec(SimpleLoop loop, int depth) {
    loop.depthLevel = depth;
    for (int liter = 0; liter < loop.children.length; ++liter) {
      calculateNestingLevelRec(loop.children[liter], depth+1);

      loop.setNestingLevel(max(loop.nestingLevel,
                               1+loop.children[liter].nestingLevel));
    }
  }

  void calculateNestingLevel() {
    for (int liter = 0; liter < this.loops.length; ++liter) {
      if (!loops[liter].isRoot) {
        if (loops[liter].parent == null) {
          loops[liter].setParent(this.root);
        }
      }
    }

    this.calculateNestingLevelRec(this.root, 0);
  }

  int getNumLoops() => loops.length;
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
class UnionFindNode {
  int           dfsNumber;
  UnionFindNode parent;
  BasicBlock    bb;
  SimpleLoop    loop;

  UnionFindNode() {
    dfsNumber = 0;
    bb        = null;
    parent    = null;
    loop      = null;
  }

  // Initialize this node.
  //
  void initNode(BasicBlock bb_arg, int dfsNumber_arg) {
    this.parent     = this;
    this.bb         = bb_arg;
    this.dfsNumber  = dfsNumber_arg;
  }

  // Union/Find Algorithm - The find routine.
  //
  // Implemented with Path Compression (inner loops are only
  // visited and collapsed once, however, deep nests would still
  // result in significant traversals).
  //
  UnionFindNode findSet() {
    List<UnionFindNode> nodeList = new List<UnionFindNode>();

    UnionFindNode node = this;
    while (node != node.parent) {
      if (node.parent != node.parent.parent)
        nodeList.add(node);

      node = node.parent;
    }

    // Path Compression, all nodes' parents point to the 1st level parent.
    for (int iter=0; iter < nodeList.length; ++iter) {
      nodeList[iter].parent = node.parent;
    }

    return node;
  }

  // Union/Find Algorithm - The union routine.
  //
  // Trivial. Assigning parent pointer is enough,
  // we rely on path compression.
  //
  void union(UnionFindNode unionFindNode) {
     parent = unionFindNode;
  }
  SimpleLoop setLoop(SimpleLoop l) => loop = l;
}



class HavlakLoopFinder {
  CFG cfg;
  LSG lsg;
  int BB_TOP, BB_NONHEADER, BB_REDUCIBLE, BB_SELF, BB_IRREDUCIBLE, BB_DEAD, BB_LAST,
    UNVISITED, MAXNONBACKPREDS;

  HavlakLoopFinder(CFG cfgParm, LSG lsgParm) {
    cfg = cfgParm;
    lsg = lsgParm;

    //
    // enum BasicBlockClass
    //
    // Basic Blocks and Loops are being classified as regular, irreducible,
    // and so on. This enum contains a symbolic name for all these
    // classifications. Python doesn't have enums, so we just create values.
    //
    BB_TOP          = 0; // uninitialized
    BB_NONHEADER    = 1; // a regular BB
    BB_REDUCIBLE    = 2; // reducible loop
    BB_SELF         = 3; // single BB loop
    BB_IRREDUCIBLE  = 4; // irreducible loop
    BB_DEAD         = 5; // a dead BB
    BB_LAST         = 6; // Sentinel

    //
    // Constants
    //
    // Marker for uninitialized nodes.
    UNVISITED = -1;

    // Safeguard against pathologic algorithm behavior.
    MAXNONBACKPREDS = (32 * 1024);
  }

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
  bool isAncestor(int w, int v, List<int> last) {
    return (w <= v) && (v <= last[w]);
  }

  //
  // DFS - Depth-First-Search
  //
  // DESCRIPTION:
  // Simple depth first traversal along out edges with node numbering.
  //
  int DFS(BasicBlock currentNode,
          List<UnionFindNode> nodes,
          Map<BasicBlock,int> number,
          List<int> last, int current) {
    nodes[current].initNode(currentNode, current);
    number[currentNode] = current;

    int lastid = current;
    for (int target = 0; target < currentNode.outEdges.length; target++) {
      if (number[currentNode.outEdges[target]] == UNVISITED)
        lastid = DFS(currentNode.outEdges[target], nodes, number,
                     last, lastid + 1);
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
  int findLoops() {
    if (cfg.startNode == null) {
      return 0;
    }

    int size = cfg.getNumNodes();

    List<List<int>> nonBackPreds = new List<List<int>>(size);
    List<List<int>> backPreds = new List<List<int>>(size);
    Map<BasicBlock, int> number = new Map<BasicBlock, int>();
    List<int> header            = new List<int>(size);
    List<int> types             = new List<int>(size);
    List<int> last              = new List<int>(size);
    List<UnionFindNode> nodes   = new List<UnionFindNode>(size);

    for (int i=0; i < size; ++i) {
      nonBackPreds[i] = new List<int>();
      backPreds[i] = new List<int>();
      header[i]    = 0;
      types[i]     = 0;
      last[i]      = 0;
      nodes[i]     = new UnionFindNode();
    }

    // Step a:
    //   - initialize all nodes as unvisited.
    //   - depth-first traversal and numbering.
    //   - unreached BB's are marked as dead.
    //
    for (BasicBlock bb in cfg.basicBlockMap.getValues()) {
      number[bb] = UNVISITED;
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
    for (int w = 0; w < size; ++w) {
      header[w] = 0;
      types[w]  = BB_NONHEADER;

      BasicBlock nodeW = nodes[w].bb;
      if (nodeW == null) {
        types[w] = BB_DEAD;
      } else {
        if (nodeW.getNumPred() > 0) {
          for (int nv = 0; nv < nodeW.inEdges.length; ++nv) {
            BasicBlock nodeV = nodeW.inEdges[nv];
            int v = number[nodeV];
            if (v != UNVISITED) {
              if (isAncestor(w, v, last)) {
                backPreds[w].add(v);
              } else {
                nonBackPreds[w].add(v);
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
    for (int w = size-1; w >=0; --w) {
      // this is 'P' in Havlak's paper
      List<UnionFindNode> nodePool = new List<UnionFindNode>();

      BasicBlock nodeW = nodes[w].bb;
      if (nodeW == null) {
        continue;
      }

      // Step d:
      for (int vi = 0; vi < backPreds[w].length; ++vi) {
        var v = backPreds[w][vi];
        if (v != w) {
          nodePool.add(nodes[v].findSet());
        } else {
          types[w] = BB_SELF;
        }
      }

      // Copy nodePool to workList.
      //
      List<UnionFindNode> workList = new List<UnionFindNode>();
      for (int n = 0; n < nodePool.length; ++n) {
        workList.add(nodePool[n]);
      }

      if (nodePool.length != 0) {
        types[w] = BB_REDUCIBLE;
      }
      // work the list...
      //
      while (workList.length > 0) {
        UnionFindNode x = workList[0];
        workList = workList.getRange(0, workList.length-1);

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
        int nonBackSize = nonBackPreds[x.dfsNumber].length;
        if (nonBackSize > MAXNONBACKPREDS) {
          return 0;
        }

        for (int iter=0; iter < nonBackPreds[x.dfsNumber].length; ++iter) {
          UnionFindNode y = nodes[nonBackPreds[x.dfsNumber][iter]];
          UnionFindNode ydash = y.findSet();

          if (!isAncestor(w, ydash.dfsNumber, last)) {
            types[w] = BB_IRREDUCIBLE;
            nonBackPreds[w].add(ydash.dfsNumber);
          } else {
            if (ydash.dfsNumber != w) {
              if (nodePool.indexOf(ydash) == -1) {
                workList.add(ydash);
                nodePool.add(ydash);
              }
            }
          }
        }
      }

      // Collapse/Unionize nodes in a SCC to a single node
      // For every SCC found, create a loop descriptor and link it in.
      //
      if ((nodePool.length > 0) || (types[w] == BB_SELF)) {
        SimpleLoop loop = lsg.createNewLoop();

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

        for (int np = 0; np < nodePool.length; ++np) {
          UnionFindNode node = nodePool[np];

          // Add nodes to loop descriptor.
          header[node.dfsNumber] = w;
          node.union(nodes[w]);

          // Nested loops are not added, but linked together.
          if (node.loop != null) {
            node.loop.setParent(loop);
          } else {
            loop.addNode(node.bb);
          }
        }
        lsg.addLoop(loop);
      } // nodePool.length
    } // Step c

    return lsg.getNumLoops();
  } // findLoops
} // HavlakLoopFinder


//======================================================
// Testing Code
//======================================================

int buildDiamond(CFG cfg, int start) {
  var bb0 = start;
  new BasicBlockEdge(cfg, bb0, bb0 + 1);
  new BasicBlockEdge(cfg, bb0, bb0 + 2);
  new BasicBlockEdge(cfg, bb0 + 1, bb0 + 3);
  new BasicBlockEdge(cfg, bb0 + 2, bb0 + 3);
  return bb0 + 3;
}


void buildConnect(CFG cfg, int start, int end) {
  new BasicBlockEdge(cfg, start, end);
}

int buildStraight(CFG cfg, int start, int n) {
  for (int i=0; i < n; i++) {
    buildConnect(cfg, start + i, start + i + 1);
  }
  return start + n;
}

int buildBaseLoop(CFG cfg, int from) {
  int header   = buildStraight(cfg, from, 1);
  int diamond1 = buildDiamond(cfg, header);
  int d11      = buildStraight(cfg, diamond1, 1);
  int diamond2 = buildDiamond(cfg, d11);
  int footer   = buildStraight(cfg, diamond2, 1);
  buildConnect(cfg, diamond2, d11);
  buildConnect(cfg, diamond1, header);

  buildConnect(cfg, footer, from);
  footer = buildStraight(cfg, footer, 1);
  return footer;
}

main() {
  print("Welcome to LoopTesterApp, Dart edition");

  CFG cfg = new CFG();
  LSG lsg = new LSG();

  print("Constructing Simple CFG...");

  cfg.createNode(0);  // top
  buildBaseLoop(cfg, 0);
  cfg.createNode(1);  //s bottom
  buildConnect(cfg, 0, 2);

  // execute loop recognition 15000 times to force compilation
  print("15000 dummy loops");

  for (int dummyloop = 0; dummyloop < 1; ++dummyloop) {
    var lsglocal = new LSG();
    var finder = new HavlakLoopFinder(cfg, lsglocal);
    var x = finder.findLoops();
    print("found $x");
  }

  print("Constructing CFG...");
  var n = 2;

  for (int parlooptrees=0; parlooptrees < 10; parlooptrees++) {
    cfg.createNode(n + 1);
    buildConnect(cfg, n, n + 1);
    n = n + 1;
    for (int i=0; i < 2; ++i) {
      var top = n;
      n = buildStraight(cfg, n, 1);
      for (int j=0; j < 25; j++) {
        n = buildBaseLoop(cfg, n);
      }

      var bottom = buildStraight(cfg, n, 1);
      buildConnect(cfg, n, top);
      n = bottom;
    }
  }

  print("Performing Loop Recognition\n1 Iteration");

  HavlakLoopFinder finder = new HavlakLoopFinder(cfg, lsg);
  int num_loops = finder.findLoops();
  lsg.calculateNestingLevel();

  print("Found: $num_loops.\nAnother 100 iterations...");

  for (var i = 0; i < 100; i++) {
    LSG lsglocal = new LSG();
    finder = new HavlakLoopFinder(cfg, lsglocal);
    num_loops = finder.findLoops();
  }

  print("Found $num_loops loops (including artificial root node)" );
  lsg.calculateNestingLevel();
}
