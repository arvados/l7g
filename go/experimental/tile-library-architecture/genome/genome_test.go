// genome_test is a file for running various tests related to the genome package.
package genome

import (
	"log"
	"testing"

	"../tile-library"
)

// This tests writing and reading genomes to and from files.
func TestWriteToFile(t *testing.T) {
	testFile := ""         // Put your test file here.
	testGenomeFile1 := ""  // Put your test genome file here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l := tilelibrary.InitializeLibrary(testFile, nil)
	tilelibrary.AddLibraryFastJ(genomeDirectory1, &l)
	tilelibrary.AddLibraryFastJ(genomeDirectory2, &l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	CreateGenome(genomeDirectory1, &g)
	WriteGenomeToFile(testGenomeFile1, &g)
	paths, err := ReadGenomeFromFile(testGenomeFile1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for i, path := range paths {
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
	path := 0              // Put your test path number here
	log.SetFlags(log.Llongfile)
	l := tilelibrary.InitializeLibrary(testFile, nil)
	tilelibrary.AddLibraryFastJ(genomeDirectory1, &l)
	tilelibrary.AddLibraryFastJ(genomeDirectory2, &l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	CreateGenome(genomeDirectory1, &g)
	g.WriteNumpy(testNumpy1, path)
	testPath, err := ReadGenomeNumpy(testNumpy1)
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
	l := tilelibrary.InitializeLibrary(testFile, nil)
	tilelibrary.AddLibraryFastJ(genomeDirectory1, &l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	l1 := tilelibrary.InitializeLibrary(testFile2, nil)
	tilelibrary.AddLibraryFastJ(genomeDirectory2, &l1)

	tilelibrary.SortLibrary(&l1)
	l1.AssignID()
	l2, err := tilelibrary.MergeLibraries(&l, &l1, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	(*l2).AssignID()
	g := InitializeGenome(&l)
	CreateGenome(genomeDirectory1, &g)
	g1 := InitializeGenome(l2)
	CreateGenome(genomeDirectory1, &g1)
	LiftoverGenome(&g, l2)
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
