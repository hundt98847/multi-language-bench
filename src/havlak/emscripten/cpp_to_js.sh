#
# emscripten version of Havlak Benchmark
#
# This directory does not contain additional sources. Instead
# it only contains this script, which performs the necessary steps
# to compile the available C++ sources to LLVM bitcodes, and from there
# to JavaScript.
#
# To make this work:
#    - set the path to the tools
#        emscripten
#        python
#        LLVM
#        Closure compiler (optional)
#        V8 shell
#
#    - set the path to the C++ sources
#
#    - you need llvm-3.0 or later

#===============================================================
# Environment
#===============================================================

#
# emscripten home directory
#
EMSCR=/home/rhundt/projects/emscripten

#
# python binary
#
PYTHON=python

#
# LLVM bin directory
#
LLVM=/home/rhundt/projects/llvm-3.0/bin

#
# Closure compiler
#
CLOSURE=/home/rhundt/projects/closure-compiler/compiler.jar

#
# V8 shell
#
D8=/home/rhundt/projects/v8/d8

#
# C++ source directoryu
#
SRC=../cpp


#===============================================================
# Compile
#===============================================================

#
# generate LLVM bitcode, use emcc as driver
#
$PYTHON $EMSCR/emcc $SRC/LoopTesterApp.cc -c -o ./LoopTesterApp.bc
$PYTHON $EMSCR/emcc $SRC/mao-loops.cc -c -o ./mao-loops.bc

#
# invoke LLVM's optimizer
#
$LLVM/opt -O3 ./mao-loops.bc -o=./mao-loops.opt
$LLVM/opt -O3 ./LoopTesterApp.bc -o=./LoopTesterApp.opt

#
# link bitcode files into one
#
$LLVM/llvm-link  LoopTesterApp.opt mao-loops.opt >./loops.bc

#
# generate JS
#
$PYTHON $EMSCR/emscripten.py --dlmalloc loops.bc >./loops.js

#
# invoke closure compiler (optional)
#
java -jar $CLOSURE --js ./loops.js --compilation_level ADVANCED_OPTIMIZATIONS >./loops.cl.opt.js

#
# time the closure and non-closure versions
#
time $D8 loops.cl.opt.js 
time $D8 loops.js
