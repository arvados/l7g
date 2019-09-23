// Program liftovergenome takes a genome, a source library, and a destination library, along with a destination filepath and a boolean.
// It performs a liftover operation on the genome from the source library to the destination library, and writes the result in the specified path, and writes in either a text file or multiple numpy arrays.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"../genome"
	"../structures"
	"../tile-library"
)

func main() {
	genomeFile := flag.String("genome", "", "The text file to read the genome from.")
	newGenomePath := flag.String("path", "newgenome.txt", "The filepath and name to write the new genome to. Writes to the current directory in a file called newgenome.txt by default.")
	source := flag.String("source", "", "The directory of SGLFv2 files for the source library. Uses the current directory by default.")
	destination := flag.String("destination", "", "The directory of SGLFv2 files for the destination library. Uses the current directory by default.")
	numpy := flag.Bool("npy", false, "Writes the output as numpy arrays instead of a text file. -path must be a directory if -npy is true. Default is false (writes a text file)")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatalf("error in parsing")
	}
	g, err := genome.ReadGenomeFromFile(*genomeFile)
	if err != nil {
		log.Fatal(err)
	}
	sourceLibrary, err := tilelibrary.New(*source, nil)
	if err != nil {
		log.Fatal(err)
	}
	destinationLibrary, err := tilelibrary.New(*destination, nil)
	if err != nil {
		log.Fatal(err)
	}
	sourceLibrary.AddLibrarySGLFv2()
	destinationLibrary.AddLibrarySGLFv2()
	err = g.AssignLibrary(sourceLibrary)
	if err != nil {
		log.Fatal(err)
	}
	g.Liftover(destinationLibrary)
	if *numpy {
		info, err := os.Stat(*newGenomePath)
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			log.Fatal(errors.New("specified path is not a directory, cannot write numpy arrays"))
		}
		for path := 0; path < structures.Paths; path++ {
			g.WriteNumpy(*newGenomePath, path)
		}
	} else {
		g.WriteToFile(*newGenomePath)
	}
	absolutePath, err := filepath.Abs(*newGenomePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(absolutePath)
}
