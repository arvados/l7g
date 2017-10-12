# asmukk

Implementation of Ukkonen's approximate string matching algorithm using fixed mismatch and gap cost.

The straight forward dynamic programming algorithm for approximate string alignment (ASM) takes
`O(N*M)` time and `O(N*M)` space, where `N` and `M` are the lengths of the strings to compare.
Though there's good reason to [think that this is optimal](http://arxiv.org/abs/1412.0348), often
the strings we want to compare are very similar and the straight forward dynamic programming
solution does a lot of extra work, considering paths that are far from optimal.

This was noticed by [Ukkonen](http://www.sciencedirect.com/science/article/pii/S0019995885800462)
and led to an "output sensitive" algorithm that runs in `O( max{N,M) * d )`, where `d` is the edit
distance between the two strings.

Note that [Hirschberg's algorithm](https://en.wikipedia.org/wiki/Hirschberg's_algorithm) only
takes `O( max{N,M} )` space but still requires the full `O(N*M)` time.

Though Ukkonen's approximate string matching algorithm is well known, I found it difficult to find
an accessible implementation of it.

# Quick Start

```
$ git clone https://github.com/abeconnelly/asmukk
$ cd asmukk
$ make
$ echo -e "cute\ncat\n" | ./asm_ukk
5
cute
cat-
```

The mismatch cost is fixed as is the gap cost.  By default, the mismatch cost is set to 3 and
the gap cost is set to 2.

# Examples

## C

```c
// gcc asm_ukk.o main.cpp -o main
#include <stdio.h>
#include <stdio.h>
#include "asm_ukk.h"

int main(int argc, char **argv) {
  char *x = "cute", *y = "cat";
  char *X=NULL, *Y=NULL;
  int score;

  int mismatch = 3, gap = 2;
  char gap_char = '-';

  score = asm_ukk_align2(&X, &Y, x, y, mismatch, gap, gap_char);

  printf("got score: %d\n", score);
  printf("orig x: %5s, aligned x: %s\n", x, X);
  printf("orig y: %5s, aligned y: %s\n", y, Y);

  free(X);
  free(Y);
}
```

# References

* ["Algorithms for Approximate String Matching" by E. Ukkonen](http://www.sciencedirect.com/science/article/pii/S0019995885800462)
* ["Edit Distance Cannot Be Computed in Strongly Subquadratic Time (unless SETH is false)" by A. Backurs and P. Indyk](http://arxiv.org/abs/1412.0348)
* [Hirschberg's algorithm](https://en.wikipedia.org/wiki/Hirschberg's_algorithm)

# License

Copyright Curoverse, Inc, AGPLv3.  See the LICENSE file for more details.


