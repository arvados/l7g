package genome

import (
	"fmt"
	"crypto/md5"
	"log"
	"testing"
	"../tile-library"
)

// This tests writing and reading genomes to and from files.
func TestWriteToFile(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	WriteGenomeToFile("/data-sdc/jc/testGenome.genome", &g)
	paths, err := ReadGenomeFromFile("/data-sdc/jc/testGenome.genome")
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

// Note: some genomes have different number of steps within a path (e.g. path 811), so this test is likely to fail.
// This also means that writing multiple genomes to the same numpy will likely fail.
// This test makes sure that genomes have the same path length in all scenarios. Mostly important for writing the path of multiple genomes to a numpy array.
func TestGenomePathConsistency(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	g1 := InitializeGenome(&l)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000037847-ASM", &g1)
	for path := range g.Paths {
		if len(g.Paths[path][0]) + len(g.Paths[path][1]) != len(g1.Paths[path][0]) + len(g1.Paths[path][1]) {
			fmt.Println(g.Paths[path][0], g.Paths[path][1])
			fmt.Println(g1.Paths[path][0], g1.Paths[path][1])
			t.Fatalf("path lengths are not the same at path %v, %v %v", path, len(g.Paths[path][0]) + len(g.Paths[path][1]), len(g1.Paths[path][0]) + len(g1.Paths[path][1]))
		}
	}
}

func TestGenomeNumpy(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	g := InitializeGenome(&l)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	g.WriteNumpy("/data-sdc/jc/testGenome.npy", 24)
	testPath, err := ReadGenomeNumpy("/data-sdc/jc/testGenome.npy")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(g.Paths[24][0]) != len(testPath[0]) || len(g.Paths[24][1]) != len(testPath[1]) {
		t.Fatalf("number of steps is not equal")
	}
	for i := range g.Paths[24][0] {
		if (g.Paths[24][0][i] != testPath[0][i] && testPath[0][i] != -2) || (g.Paths[24][1][i] != testPath[1][i] && testPath[1][i] != -2)  {
			t.Fatalf("a step is not equal")
		}
	}
}

func TestGenomeLiftover(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	tilelibrary.SortLibrary(&l)
	l.AssignID()
	l1:=tilelibrary.InitializeLibrary("/data-sdc/jc/tile-library/testing/test.txt", [][md5.Size]byte{})
	tilelibrary.AddLibraryFastJ("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu02C8E3_masterVarBeta-GS000036653-ASM", "/data-sdc/jc/tile-library/testing/test.txt",&l1)

	tilelibrary.SortLibrary(&l1)
	l1.AssignID()
	l2, err :=tilelibrary.MergeLibraries(&l, &l1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	(*l2).AssignID()
	g := InitializeGenome(&l)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", &g)
	g1 := InitializeGenome(l2)
	CreateGenome("../../../../../keep/by_id/su92l-4zz18-2hxdqjw6cbrnr7s/hu03E3D2_masterVarBeta-GS000038659-ASM", &g1)
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