package genome

import (
	"crypto/md5"
	"log"
	"testing"
	"../tile-library"
)

// This tests writing and reading genomes to and from files.
func TestWriteToFile(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	CreateGenome("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	WriteGenomeToFile("/data-sdc/jc/testGenome.genome", &g)
	paths := ReadGenomeFromFile("/data-sdc/jc/testGenome.genome")
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

// This test makes sure that genomes have the same path length in all scenarios. Mostly important for writing the path of multiple genomes to a numpy array.
func TestGenomePathConsistency(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	g1 := InitializeGenome(&l)
	CreateGenome("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	CreateGenome("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", &g1)
	for path := range g.Paths {
		if len(g.Paths[path][0]) + len(g.Paths[path][1]) != len(g1.Paths[path][0]) + len(g1.Paths[path][1]) {
			t.Fatalf("path lengths are not the same at path %v, %v %v", path, len(g.Paths[path][0]) + len(g.Paths[path][1]), len(g1.Paths[path][0]) + len(g1.Paths[path][1]))
		}
	}
}