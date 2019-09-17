// Program mergelibraries merges a set of given directories of SGLFv2 files together into one library, and then writes SGLF or SGLFv2 files for the new library to disk.
package main

import (
	"flag"
	"log"
	"os"

	"../tile-library"
)

func main() {
	log.SetFlags(log.Llongfile)
	directoryToWrite := flag.String("dir", "", "The directory to write files to. Leave blank to write to the current directory.")
	version := flag.Int("version", 0, "The SGLF version. Use 0 for normal SGLF files and 1 for SGLFv2 files. Undefined for any other value currently.")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatalf("error in parsing")
	}
	arguments := flag.Args()
	l, err := tilelibrary.New(*directoryToWrite, nil) // Empty
	if err != nil {
		log.Fatal(err)
	}
	l.AssignID()
	for _, directory := range arguments {
		info, err := os.Stat(directory)
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			log.Fatalf("argument provided is not a directory")
		}
		newLibrary, err := tilelibrary.New(directory, nil)
		if err != nil {
			log.Fatal(err)
		}
		err = newLibrary.AddLibrarySGLFv2()
		if err != nil {
			log.Fatal(err)
		}
		newLibrary.AssignID()
		l.MergeLibrariesWithoutCreation(newLibrary)
	}
	if *version == 0 {
		err := l.WriteLibraryToSGLF(*directoryToWrite)
		if err != nil {
			log.Fatal(err)
		}
	} else if *version == 1 {
		err := l.WriteLibraryToSGLFv2(*directoryToWrite)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("version is invalid")
	}
}
