// Program genomestonumpy takes a directory to write to, a directory for a source library, a path number, and any number of directories for genomes.
// It creates Genome structures for each genome relative to the source library, and then writes the numpy array for the path given (or all paths if no path is given) to the provided directory.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"../genome"
	"../structures"
	"../tile-library"
)

func main() {
	directoryToWrite := flag.String("dir", "", "The directory to write numpy arrays to. Writes to current directory by default.")
	source := flag.String("source", "", "The directory of SGLFv2 files for the source library. Uses the current directory by default.")
	path := flag.Int("path", -1, "The path to make a numpy array for. By default, numpy arrays are created for all paths.")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatalf("error in parsing")
	}
	genomeDirectories := flag.Args()
	sourceLibrary, err := tilelibrary.New(*source, nil)
	if err != nil {
		log.Fatal(err)
	}
	sourceLibrary.AddLibrarySGLFv2()
	var genomeList []*genome.Genome
	for _, genomeDirectory := range genomeDirectories {
		newGenome := genome.New(sourceLibrary)
		newGenome.Add(genomeDirectory)
		genomeList = append(genomeList, newGenome)
	}
	if *path < 0 {
		for genomePath := 0; genomePath < structures.Paths; genomePath++ {
			genome.WriteGenomesPathToNumpy(genomeList, *directoryToWrite, genomePath)
		}
	} else {
		genome.WriteGenomesPathToNumpy(genomeList, *directoryToWrite, *path)
	}

	absolutePath, err := filepath.Abs(*directoryToWrite)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(absolutePath)
}
