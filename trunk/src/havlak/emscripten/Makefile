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
# This is sort of a dummy Makefile, so that all versions have one and
# also a single method to run the benchmark.
#
#
# Also - this version only includes the sources from the C++ version.
# However, the dimensions of the C++ version are too large. To
# actually run the generated JS, the C++ benchmark needs to be reduced
# in size.
#

loops.cl.opt.js:  ./cpp_to_js.sh
	./cpp_to_js.sh

run: 
	d8 loops.cl.opt.js

clean:
	rm -f *~
	rm -f *.js
	rm -f *.opt
	rm -f *.bc
