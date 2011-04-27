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
#include <map>
#include <set>
#include <vector>
#include <algorithm>

#include "mao-loops.h"

int buildDiamond(MaoCFG *cfg, int start) {
  int bb0 = start;

  new BasicBlockEdge(cfg, bb0, bb0 + 1);
  new BasicBlockEdge(cfg, bb0, bb0 + 2);
  new BasicBlockEdge(cfg, bb0 + 1, bb0 + 3);
  new BasicBlockEdge(cfg, bb0 + 2, bb0 + 3);

  return bb0 + 3;
}

void buildConnect(MaoCFG *cfg, int start, int end) {
  new BasicBlockEdge(cfg, start, end);
}

int buildStraight(MaoCFG *cfg, int start, int n) {
  for (int i = 0; i < n; i++) {
    buildConnect(cfg, start + i, start + i + 1);
  }

  return start + n;
}

int buildBaseLoop(MaoCFG *cfg, int from) {
  int header = buildStraight(cfg, from, 1);
  int diamond1 = buildDiamond(cfg, header);
  int d11 = buildStraight(cfg, diamond1, 1);
  int diamond2 = buildDiamond(cfg, d11);
  int footer = buildStraight(cfg, diamond2, 1);
  buildConnect(cfg, diamond2, d11);
  buildConnect(cfg, diamond1, header);

  buildConnect(cfg, footer, from);
  footer = buildStraight(cfg, footer, 1);
  return footer;
}


int main(int argc, char *argv[]) {
  fprintf(stderr, "Welcome to LoopTesterApp, C++ edition\n");
  fprintf(stderr, "Constructing App...\n");
  MaoCFG cfg;
  LoopStructureGraph lsg;

  fprintf(stderr, "Constructing Simple CFG...\n");
  cfg.CreateNode(0);  // top
  buildBaseLoop(&cfg, 0);
  cfg.CreateNode(1);  // bottom
  new BasicBlockEdge(&cfg, 0,  2);

  fprintf(stderr, "15000 dummy loops\n");
  for (int dummyloops = 0; dummyloops < 15000; ++dummyloops) {
    LoopStructureGraph * lsglocal = new LoopStructureGraph();
    FindHavlakLoops(&cfg, lsglocal);
    delete(lsglocal);
  }

  fprintf(stderr, "Constructing CFG...\n");
  int n = 2;

  for (int parlooptrees = 0; parlooptrees < 10; parlooptrees++) {
    cfg.CreateNode(n + 1);
    buildConnect(&cfg, 2, n + 1);
    n = n + 1;

    for (int i = 0; i < 100; i++) {
      int top = n;
      n = buildStraight(&cfg, n, 1);
      for (int j = 0; j < 25; j++) {
        n = buildBaseLoop(&cfg, n);
      }
      int bottom = buildStraight(&cfg, n, 1);
      buildConnect(&cfg, n, top);
      n = bottom;
    }
    buildConnect(&cfg, n, 1);
  }

  fprintf(stderr, "Performing Loop Recognition\n1 Iteration\n");
  int num_loops = FindHavlakLoops(&cfg, &lsg);

  fprintf(stderr, "Another 50 iterations...\n");
  int sum = 0;
  for (int i = 0; i < 50; i++) {
    LoopStructureGraph lsg;
    fprintf(stderr, ".");
    sum += FindHavlakLoops(&cfg, &lsg);
  }
  fprintf(stderr,
          "\nFound %d loops (including artificial root node)"
          "(%d)\n", num_loops, sum);
  lsg.Dump();
}
