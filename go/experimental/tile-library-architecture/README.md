# Golang Tile Library Architecture

---

In this directory is a series of packages useful for creating and managing tile libraries and genomes in Go.

The following Go packages are included:
- `structures`, a package with the basic `TileVariant` structure important to tile libraries, along with an equality method, and the number of paths in the genome, for convenience.
- `tilelibrary`, a package containing ways to create, add to, compare, export, and import tile libraries. Requires the `structures` package.
- `genome`, a package containing ways to create, export, and import genomes *relative to a tile library*. Requires the `tilelibrary` and the `structures` packages.

In addition, a couple programs are provided to create and merge libraries directly from the command line.

## Installation

---

Once this project is on GitHub, use the `go get` command like with any other Go package on GitHub. Make sure to get every package, since `tilelibrary` and `genome` rely on `structures`, and `genome` relies on `tilelibrary`.

## CreateLibrary Program Usage

---

The `createlibrary` program has the functionality to combine a series of directories of genome files containing FastJs and write the resulting library to disk.

Within the `createlibrary` folder, the command to run the program is `./createlibrary`. Outside of this directory, the command will be the filepath to the `createlibrary` program in the `createlibrary` folder. All non-flag arguments given to the command should be after all flags, and should be directories of FastJ files to add to the library.

The following flags apply:
- `-version`: specifies the version of the SGLF to be outputted. A value 0 outputs regular SGLF files, and 1 outputs SGLFv2 files. Behavior is not defined when the `-version` flag is given any other number. Default behavior is 0 (regular SGLF files).
- `-textfile`: specifies the location and name of the text file where the command is allowed to write intermediate data. Default behavior is to write the intermediate data to a text file called `test.txt` in the current directory. Files are created if they do not exist.
- `-dir`: specifies the directory to write the output files to. If the directory does not exist, it will be created. Default behavior is to write the output files to the current directory.

A successful run will print out the directory that was written to as an absolute path, if other programs need to access the created files.

## MergeLibraries Program Usage

---

The `mergelibraries` program has the functionality to merge libraries written to disk the in the form of SGLFv2 files, and output either SGLF or SGLFv2 files from the merging.

Within the `mergelibraries` folder, the command to run the program is `./mergelibraries`. Outside of this directory, the command will be the filepath to the `mergelibraries` program in the `mergelibraries` folder. All non-flag arguments given to the command should be after all flags, and should be directories of SGLFv2 files to add to the library.

The following flags apply:
- `-version`: specifies the version of the SGLF to be outputted. A value 0 outputs regular SGLF files, and 1 outputs SGLFv2 files. Behavior is not defined when the `-version` flag is given any other number. Default behavior is 0 (regular SGLF files).
- `-dir`: specifies the directory to write the output files to. If the directory does not exist, it will be created. Default behavior is to write the output files to the current directory.

A successful run will create files silently. Any errors that are encountered will exit the program and be printed out.

## LiftoverGenome Program Usage

---

The `liftovergenome` program has the functionality to liftover a genome from a source library to a destination library and write out a text file or numpy files to disk.

Within the `liftovergenome` folder, the command to run the program is `./liftovergenome`. Outside of this directory, the command will be the filepath to the `liftovergenome` program in the `liftovergenome` folder.

The following flags apply:
- `-genome`: specifies the location of the **text file** of the genome.
- `-source`: specifies the directory of the SGLFv2 files for the source library. Uses the current directory by default.
- `-destination`: specifies the directory of the SGLFv2 files for the destination library. Uses the current directory by default.
- `-path`: specifies the location where the file should be outputted. If `npy` is true, this must be a directory. Writes a textfile called newgenome.txt in the current directory by default.
- `-npy`: specifies what type of files to output. True outputs a numpy array for each path--false will output one text file for the genome. Default is false.

A successful run will print out the path that was written to as an absolute path, if other programs need to access the created files.

## GenomesToNumpy Program Usage

---

The `genomestonumpy` program will take a series of genomes, a source library, a directory, and a path number, and map those genomes to the library, and write the numpy array for that path for all genomes (if the path number is given), or writes out all numpy arrays for all paths, if no specified path number is given.

Within the `genomestonumpy` folder, the command to run the program is `./genomestonumpy`. Outside of this directory, the command will be the filepath to the `genomestonumpy` program in the `genomestonumpy` folder. All non-flag arguments given to the command should be after all flags, and should be directories of FastJ files for each genome that should be included.

The following flags apply:
- `-dir`: specifies the directory to output files to. Uses the current directory by default.
- `-source`: specifies the directory of the SGLFv2 files for the source library. Uses the current directory by default.
- `-path`: specifies the path number of the numpy array that should be outputted. Outputs all numpy arrays instead by default.

A successful run will print out the path that was written to as an absolute path, if other programs need to access the created files.

## Go package descriptions

---

For package `structures`, the major structure here is the TileVariant, which contains information regarding the variant and the library it belongs to. Fields of the TileVariant can be called and modified directly if needed. An equality method based on the hash is provided.

Package `tilelibrary` contains methods for working with, modifying, and exporting and importing libraries. This allows for writing to SGLF and SGLFv2 files, along with creating libraries from SGLFv2 files. Tiles can be added and annotated, along with finding frequencies and existence of specific tiles. Libraries can be merged, and liftover mappings from one library to another can be created if genomes attached to a specific library need to reference another library. Tile libraries are given IDs based on the MD5 hash algorithm, to make them easy to refer to. Libraries are safe for concurrent use.

Package `genome` contains methods for creating, exporting, and importing genomes *relative to a specific tile library*. This includes writing to and reading from numpy arrays and text files, and also allows for lifting over genomes from one tile library to another. Genomes refer to a specific tile by using its tile variant number in the reference tile library.

## SGLFv2 Specification

This package creates a new type of file to keep track of libraries, being the SGLFv2 file. The format of the SGLFv2 file is as follows:

First, each name must be a 4 digit hexadecimal number between 0 and the number of paths in the genome representing the path, followed by the suffix .sglfv2. Every variant in that path for this library must be in that file.

The first line of each file follows the following format:
```
ID:LibraryID,Components:ComponentID1,ComponentID2,ComponentID3...
```
where the current library ID is the first ID, and the IDs of any components are separated by commas.

The following lines are of the following format, where each line contains exactly one tile variant's information:
```
PATH.01.STEP.VARIANT+COUNT+LENGTH,HASH,BASES
```
where PATH and STEP are 4 digit hexadecimal representations of the tile's path and step, VARIANT is the variant number (ordered from most common to least common, sorted by increasing hash in ties) as a 3 digit hexadecimal number, LENGTH is the tile span in hexadecimal, COUNT is the frequency of the tile, as an 8 digit hexadecimal number, HASH is the hash of the bases of the tile, and BASES is the string of bases of the tile (may include nocalls).

The appearance of lines is done first increasing by step, and then increasing in variant number within each step.

The current hash algorithm for determining IDs and for hashing tile variants is MD5.

## Documentation

---

Once this project is on GitHub, documentation of the packages should appear on Godoc (godoc.org). For now, documentation for Go packages can be found within the .go files, and documentation for programs can be found here or by using the help flag (-h).

## Tests

---

Tests are provided for the `genome` and the `tilelibrary` packages, under the files `genome_test.go` and `tile-library_test.go`. In both files, there are variable fields with empty strings for where various directories or file names would go. When running these tests, make sure to replace these empty strings with whatever directories and file names to use. Both relative and absolute file paths work for these tests.

## Notes

---

- Make sure to initialize libraries and genomes for use--using the provided New functions will automatically set up the necessary structures.
- Check to make sure libraries are initialized with the correct reference paths--importing a library from SGLFv2 files is only allowed if the reference path is the directory of those SGLFv2 files.
- Leftover files from libraries that need to be removed should be removed by the caller. The `tilelibrary` package does not delete intermediate files.
- TileVariants are compared by **hash only**. Even if two TileVariants might have different fields elsewhere (for example, different annotations), equality is determined only by the hash of both variants.
- Adding tiles directly to a library created from SGLFv2 files is not valid, since adding to an SGLFv2 file directly would cause all lookup reference numbers for tiles to be shifted over. One workaround is by merging this library with an empty library, as adding tiles to merged libraries is allowed.
- In genomes, the number -1 represents a skipped step location because of a spanning tile. The number -2, which appears in text file and numpy array representations of genomes, represents an incomplete tile (that is, it contains a nocall).

## Future goals/features

---

- Parallelization of merging. This will speed up the merge process considerably, especially when merging more than two libraries together.
- Parallelization by path in library creation.