Standardized benchmarks, implemented in multiple languages, using the language's idiomatic, default containers, iteration methods, type methods and so. An effort to allow an almost fair comparison of core characteristics of the various languages.

An initial comparison of a loop recognition algorithm, implemented in C++, Java, Go, and Scala, has been published at the Scala Days 2011 Workshop at Stanford:  `https://code.google.com/p/multi-language-bench/source/browse/trunk/doc/scaladays2011-paper/main.pdf`

**Publicity**

This paper got picked up by [The Register](http://www.theregister.co.uk/2011/06/03/google_paper_on_cplusplus_java_scala_go/) with a slightly more snappy title, and it was all over the Geek Blogs soon after that, e.g., [Slashdot](http://developers.slashdot.org/story/11/06/15/0242237/C-the-Clear-Winner-In-Googles-Language-Performance-Tests), or [Reddit](http://www.reddit.com/r/programming/duplicates/hqkwk/google_paper_comparing_performance_of_c_java/)

**Responses**

Several people have pointed out flaws in the paper or suggested improvements. All of those are greatly appreciated and should help to better understand the various language issues.

**GO**

Russ Cox published an excellent [Blog Post](http://blog.golang.org/2011/06/profiling-go-programs.html) on how to improve performance. In the end, they even beat the tuned C++ version

**Java**

Jeremy Manson published an interesting comment on his [Blog](http://jeremymanson.blogspot.com/2011/06/scala-java-shootout.html)

**Cawen**

Cawen is a new (and not yet published language) developed by Thomas Samain and Gwenael Chailleu. They ran experiments similar to the C++ experiments described in the paper and published their results [here](http://melvenn.com/cawen/the-google-benchmark/). As the compiler will become available, I will add the sources to this repository.

**Energy Efficient Computing**

Interesting comparison of the various benchmarks in regards
[Energy Efficient Computing](http://www.greencodelab.fr/en/content/are-some-languages-%E2%80%8B%E2%80%8Bgreener-than-others)

# News #
  * We've added a manually constructed ` JavaScript ` version
  * We've added a manually constructed ` Dart ` version. It includes instructions on how to generate `JavaScript` from it.
  * We've added an ` emscripten ` C++ to Javascript converter script.
  * We've added a GWT version (basically, a wrapper around the Java version). This version is difficult to run and requires some manual massaging. This version also allows to generate ` JavaScript `.
  * We've added a Native Client (`NaCl`) version. This version uses a similar technique as the GWT version, but can be run more easily (on Chrome).