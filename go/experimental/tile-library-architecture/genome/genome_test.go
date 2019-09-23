// genome_test is a file for running various tests related to the genome package.
// Places where files and directories should be are left blank as variables.
// Fill them in with wherever files and directories should be made.
package genome

import (
	"log"
	"testing"

	"../tile-library"
)

// This tests writing and reading genomes to and from files.
func TestWriteToFile(t *testing.T) {
	testFile := ""         // Put your test tile library file here.
	testGenomeFile1 := ""  // Put your test genome file here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l, err := tilelibrary.New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.AddLibraryFastJ(genomeDirectory2)
	l.SortLibrary()
	l.AssignID()
	g := New(l)
	g.Add(genomeDirectory1)
	g.WriteToFile(testGenomeFile1)
	newGenome, err := ReadGenomeFromFile(testGenomeFile1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if newGenome.libraryID != g.libraryID {
		t.Fatalf("IDs are not equal")
	}
	for i, path := range newGenome.Paths {
		if len(path) != len(g.Paths[i]) {
			t.Fatalf("path lengths not equal")
		}
		for j, phase := range path {
			if len(phase) != len(g.Paths[i][j]) {
				t.Fatalf("phase lengths not equal")
			}
			for k, step := range phase {
				if step != g.Paths[i][j][k] && step != -2 {
					t.Fatalf("steps are not equal")
				}
			}
		}
	}
}

// This tests reading and writing genomes to and from numpy files.
func TestGenomeNumpy(t *testing.T) {
	testFile := ""         // Put your test file here.
	testNumpy1 := ""       // Put your test numpy file here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	path := 0              // Put your test path number here.
	log.SetFlags(log.Llongfile)
	l, err := tilelibrary.New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.AddLibraryFastJ(genomeDirectory2)
	l.SortLibrary()
	l.AssignID()
	g := New(l)
	g.Add(genomeDirectory1)
	testPath := g.Paths[path]
	g.WriteNumpy(testNumpy1, path)
	err = g.ReadGenomePathNumpy(testNumpy1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(g.Paths[path][0]) != len(testPath[0]) || len(g.Paths[path][1]) != len(testPath[1]) {
		t.Fatalf("number of steps is not equal")
	}
	for i := range g.Paths[path][0] {
		if (g.Paths[path][0][i] != testPath[0][i] && testPath[0][i] != -2) || (g.Paths[path][1][i] != testPath[1][i] && testPath[1][i] != -2) {
			t.Fatalf("a step is not equal")
		}
	}
}

// This tests the liftover of genomes from library to another.
func TestGenomeLiftover(t *testing.T) {
	testFile := ""         // Put your test file here.
	testFile2 := ""        // Put your test file here.
	testFile3 := ""        // Put your test file here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l, err := tilelibrary.New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.SortLibrary()
	l.AssignID()
	l1, err := tilelibrary.New(testFile2, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory2)

	l.SortLibrary()
	l1.AssignID()
	l2, err := l1.MergeLibraries(l, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	(*l2).AssignID()
	g := New(l)
	g.Add(genomeDirectory1)
	g1 := New(l2)
	g1.Add(genomeDirectory1)
	g.Liftover(l2)
	for i, path := range g.Paths {
		for j, phase := range path {
			for step, value := range phase {
				if value != g1.Paths[i][j][step] {
					t.Fatalf("index values not the same")
				}
			}
		}
	}
}
