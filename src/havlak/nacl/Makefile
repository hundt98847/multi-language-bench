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
# Similar to the GWT version, I copied the C++ source code into
# some simple NaCl hello_world benchmark. The dimensions of the benchmark
# are smaller than the original C++ benchmark to allow comparison with
# the JavaScript versions.

# To run this simple benchmark, copy the sources under the 
# examples directory in the NaCl SDK, e.g.:
#
# /nacl_sdk/pepper_16/examples/havlak/...
#
# and follow the instructions on how to run the examples, e.g.:
#   - go to the example directory
#
#   - start the web server with:
#        python httpd.py 5103
#
#   - get to your app in the browser via:
#        http://localhost:5103/havlak/havlak.html
#
# More details are here (excellent documentation):
#   https://developers.google.com/native-client/devguide/tutorial
#

run: 
	./scons

clean:
	rm -rf opt* dbg* havlak_x86* *.nmf *~
