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

/**
 * The Havlak loop finding algorithm.
 *
 * @author rhundt
 */
package lsg;

import cfg.BasicBlock;

import java.util.*;

/**
 * class SimpleLoop
 *
 * Basic representation of loops, a loop has an entry point,
 * one or more exit edges, a set of basic blocks, and potentially
 * an outer loop - a "parent" loop.
 *
 * Furthermore, it can have any set of properties, e.g.,
 * it can be an irreducible loop, have control flow, be
 * a candidate for transformations, and what not.
 */
public class SimpleLoop {
  public static class ObjectSet<T> implements Iterable<T> {
    Object[] arr;
    int size;
    static int total = 0;
    static int resize = 0;
    static {
      Runtime.getRuntime().addShutdownHook(new Thread() {
          @Override public void run() {
            System.out.println(total + " = total");
            System.out.println(resize + " = resize");
          }
        });
    }
    ObjectSet() {
      total++;
      arr = new Object[2];
      size = 0;
    }
    boolean add(T e) {
      for (int i = 0; i < size; i++) {
        if (arr[i] == e) {
          return false;
        }
      }
      if (size == arr.length) {
        resize++;
        Object[] old = arr;
        arr = new Object[arr.length << 1];
        System.arraycopy(old, 0, arr, 0, old.length);
      }
      arr[size] = e;
      size++;
      return true;
    }

    int size() { return size; }
    @Override public Iterator<T> iterator() {
      return new Iterator<T>() {
        int curr = 0;
        @Override public boolean hasNext() {
          return curr != size;
        }
        @SuppressWarnings("unchecked") @Override public T next() {
          return (T) arr[curr++];
        }
        @Override public void remove() {
          throw new UnsupportedOperationException();
        }
      };
    }
  }
  public SimpleLoop() {
    parent = null;
    isRoot = false;
    isReducible  = true;
    nestingLevel = 0;
    depthLevel   = 0;
    basicBlocks  = new ObjectSet<BasicBlock>();
    children     = new ObjectSet<SimpleLoop>();
  }

  public void addNode(BasicBlock bb) {
    basicBlocks.add(bb);
  }

  public void addChildLoop(SimpleLoop loop) {
    children.add(loop);
  }

  public void dump(int indent) {
    for (int i = 0; i < indent; i++)
      System.out.format("  ");

    System.out.format("loop-%d nest: %d depth %d %s",
                      counter, nestingLevel, depthLevel,
                      isReducible ? "" : "(Irreducible) ");
    if (getChildren().size() != 0) {
      System.out.format("Children: ");
      for (SimpleLoop loop : getChildren()) {
        System.out.format("loop-%d ", loop.getCounter());
      }
    }
    if (basicBlocks.size() != 0) {
      System.out.format("(");
      for (BasicBlock bb : basicBlocks) {
        System.out.format("BB#%d%s", bb.getName(), header == bb ? "* " : " ");
      }
      System.out.format("\b)");
    }
    System.out.format("\n");
  }

  // Getters/Setters
  public ObjectSet<SimpleLoop> getChildren() {
    return children;
  }
  public SimpleLoop   getParent() {
    return parent;
  }
  public int          getNestingLevel(){
    return nestingLevel;
  }
  public int          getDepthLevel() {
    return depthLevel;
  }
  public int          getCounter() {
    return counter;
  }
  public boolean      isRoot() {   // Note: fct and var are same!
    return isRoot;
  }
  public void setParent(SimpleLoop parent) {
    this.parent = parent;
    this.parent.addChildLoop(this);
  }
  public void setHeader(BasicBlock bb) {
    basicBlocks.add(bb);
    header = bb;
  }
  public void setIsRoot() {
    isRoot = true;
  }
  public void setCounter(int value) {
    counter = value;
  }
  public void setNestingLevel(int level) {
    nestingLevel = level;
    if (level == 0) {
      setIsRoot();
    }
  }
  public void setDepthLevel(int level) {
    depthLevel = level;
  }
  public void setIsReducible(boolean isReducible) {
    this.isReducible = isReducible;
  }

  private ObjectSet<BasicBlock>        basicBlocks;
  private ObjectSet<SimpleLoop>        children;
  private SimpleLoop             parent;
  private BasicBlock             header;

  private boolean      isRoot;
  private boolean      isReducible;
  private int          counter;
  private int          nestingLevel;
  private int          depthLevel;
};
