// Createlibrary is a command line function that parses directories of FastJ files into a tile library and writes files to a specified directory.
// It can write SGLF or SGLFv2 files, and also allows a choice of where the intermediate data can go.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"../tile-library"
)

// Main function for creating and exporting a tile library.
func main() {
	directoryToWrite := flag.String("dir", "", "The directory to write files to. Leave blank to write to the current directory.")
	textFile := flag.String("text", "test.txt", "The file to write the intermediate text file to.")
	version := flag.Int("version", 0, "The SGLF version. Use 0 for normal SGLF files and 1 for SGLFv2 files. Undefined for any other value currently.")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatalf("error in parsing")
	}
	arguments := flag.Args()
	l, err := tilelibrary.CompileDirectoriesToLibrary(arguments, *textFile, true)
	if err != nil {
		log.Fatal(err)
	}
	if *version == 1 {
		err = l.WriteLibraryToSGLFv2(*directoryToWrite)
	} else if *version == 0 {
		err = l.WriteLibraryToSGLF(*directoryToWrite)
	} else {
		log.Fatalf("invalid version number")
	}
	if err != nil {
		log.Fatal(err)
	}
	absolutePath, err := filepath.Abs(*directoryToWrite)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(absolutePath)
}
