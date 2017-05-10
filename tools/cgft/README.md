CGFT
===

Compact Genome Format Tool (CGFT), a swiss army knife tool to manipulate and inspect
Compact Genome Format (CGF) files.


Quick Start
---

This will compile `cgft`, create a test `cgf` file and
print out the contents of it's header in human readable format:

```bash
git clone https://github.com/curoverse/l7g
cd l7g/tools/cgft
make
./cgft -C test.cgf3
./cgft -H test.cgf3
```

Note that `cgft` requires the Succinct Data Structure Library, [sdsl-lite](https://github.com/simongog/sdsl-lite).
See their installation instructions on how to [compile and install sdsl-lite](https://github.com/simongog/sdsl-lite#installation).
The `Makefile` assumes `sdsl-lite` is installed in your home directory under `include` and `lib` (the default installation location).


Usage
---

```
usage: cgft [-H] [-b tilepath] [-e tilepath] [-i ifn] [-o ofn] [-h] [-v] [-V] [ifn]

  [-H|--header]               show header
  [-C|--create-container]     create empty container
  [-I|--info]                 print basic information about CGF file
  [-b|--band tilepath]        output band for tilepath
  [-e|--encode tilepath]      input tilepath band and add it to file, overwriting if it already exists
  [-i|--input ifn]            input file (CGF)
  [-o|--output ofn]           output file (CGF)
  [-A|--show-all]             show all tilepaths
  [-h|--help]                 show help (this screen)
  [-v|--version]              show version
  [-V|--verbose]              set verbose level
  [-t|--tilemap tfn]          use tilemap file (instead of default)
  [-Z|--ez-print]             print "ez" structure information
  [-T|--version-opt vopt]     CGF version option.  Must be one of "default" or "noc-inv"
  [-L|--library-version lopt] CGF library version option.  Overwrites default value if specified
  [-U|--update-header]        Update header only
```


Copyright and License
---

Copyright Curoverse Inc., licensed under AGPLv3.

Where applicable, data is licensed under CC0.
