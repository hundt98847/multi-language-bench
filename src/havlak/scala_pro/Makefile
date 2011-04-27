# Copyright 2011 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Scala allows building via two compilers:
#
#    scalac, which is "loaded" on each invocation
#
#    fsc, which is "loaded" on first invocation, and 
#         kept in memory after that.
#
# Compilations with fsc are therefore on 2nd invocation. To compare,
# start make with either:
#
#    make COMPILER=fsc
#    make COMPILER=scalac
#
COMPILER=fsc

#
# Point to scala-library (change path according to your setting):
#
LIB=./scala-library.jar

#
# Point to Java binary (change path according to your setting):
#
JAVA=./java

LoopTesterApp.jar: LoopTesterApp.class
	jar -cvf LoopTesterApp.jar `find . -name \*.class`

LoopTesterApp.class: LoopTesterApp.scala
	$(COMPILER) LoopTesterApp.scala

run: 
	$(JAVA)  \
	  -XX:+UseCompressedOops -server -Xss15500k -cp .:$(LIB) LoopTesterApp

clean:
	rm -f *.class
	rm -f LoopTesterApp.jar
	rm -f *~
