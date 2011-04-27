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

6.out: basicblock.6 lsg.6 havlaklookfinder.6 looptesterapp.6
	6l looptesterapp.6

basicblock.6: basicblock.go
	6g basicblock.go

lsg.6: lsg.go
	6g lsg.go

havlaklookfinder.6: havlakloopfinder.go
	6g havlakloopfinder.go

looptesterapp.6: looptesterapp.go
	6g looptesterapp.go


run: 
	./6.out

clean:
	rm -f *6 ./6.out
	rm -f *~
