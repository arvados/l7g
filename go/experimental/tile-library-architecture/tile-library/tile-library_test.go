package tilelibrary

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	//"os"
	"sort"
	"strings"
	"time"
	"testing"
	"../structures"
)

// This is a test of the TileExists function.
func TestExists(t *testing.T) {
	l := InitializeLibrary("test", [][md5.Size]byte{[md5.Size]byte{}})
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 3, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	AddTile(0,0,&tile1, &l)
	test1 := TileExists(0,0,&tile1, &l)
	if test1==-1 {
		t.Errorf("Test 1 is false, expected true")
	}
	test2 := TileExists(0,0,&tile2, &l)
	if test2==-1 {
		t.Errorf("Test 2 is false, expected true") // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash.
	}
	test3 := TileExists(0,0,&tile3,&l)
	if test3!=-1 {
		t.Errorf("Test 3 is true, expected false")
	}
}

// This is a test of the FindFrequency function.
func TestFrequency(t *testing.T) {
	l := InitializeLibrary("test", [][md5.Size]byte{[md5.Size]byte{}})
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 3, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	AddTile(0,0,&tile1, &l)
	AddTile(0,0,&tile1, &l)
	AddTile(0,0,&tile2, &l)
	AddTile(0,0,&tile3, &l)
	AddTile(0,1,&tile1, &l)
	test1 := FindFrequency(0,0,&tile1,&l)
	if test1!=3 {
		t.Errorf("Test 1 was %v, expected 3", test1) // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash, so tile1 counts the two instances of tile1 and the instance of tile2 in path 0, step 0.
	}
	test2 := FindFrequency(0,0,&tile2, &l)
	if test2!=3 {
		t.Errorf("Test 2 was %v, expected 3", test2)
	}
	test3 := FindFrequency(0,0,&tile3, &l)
	if test3!=1 {
		t.Errorf("Test 3 was %v, expected 1", test3)
	}
	test4 := FindFrequency(0,1,&tile1, &l)
	if test4!=1 {
		t.Errorf("Test 4 was %v, expected 1",test4)
	}
	test5 := FindFrequency(0,1,&tile3, &l)
	if test5!=0 {
		t.Errorf("Test 5 was %v, expected 0",test5)
	}
	test6 := FindFrequency(0,2, &tile2, &l)
	if test6 !=0 {
		t.Errorf("Test 6 was %v, expacted 0", test6)
	}
}


// This is a test of the sortLibrary function.
func TestLibrarySorting(t *testing.T) {
	l := InitializeLibrary("test", [][md5.Size]byte{[md5.Size]byte{}})
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 3, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	AddTile(0,0,&tile1, &l)
	AddTile(0,0,&tile1, &l)
	AddTile(0,0,&tile2, &l)
	AddTile(0,0,&tile2, &l)
	AddTile(0,0,&tile3, &l)
	AddTile(0,0,&tile2, &l)
	AddTile(0,1,&tile3, &l)
	AddTile(0,1,&tile1, &l)
	AddTile(0,1,&tile3, &l)
	AddTile(0,1,&tile3, &l)
	AddTile(0,2,&tile3, &l)
	AddTile(0,2,&tile1, &l)
	SortLibrary(&l)
	test1 := sort.IsSorted(sort.Reverse(sort.IntSlice(l.Paths[0].Variants[0].Counts)))
	if !test1 {
		t.Errorf("Path 0 step 0 is not sorted correctly.")
	}
	test2 := sort.IsSorted(sort.Reverse(sort.IntSlice(l.Paths[0].Variants[1].Counts)))
	if !test2 {
		t.Errorf("Path 0 step 1 is not sorted correctly.")
	}
	test3 := sort.IsSorted(sort.Reverse(sort.IntSlice(l.Paths[0].Variants[2].Counts)))
	if !test3 {
		t.Errorf("Path 0 step 2 is not sorted correctly.")
	}
}


// Note: benchmarks involving random libraries and tile generation will have a (slightly) variable amount of time taken.

// generateRandomTiles generates a set of tiles for each path and step for random libraries to use.
func generateRandomTiles(minTiles, maxTiles, minSteps, maxSteps, minBases, maxBases int) [][][]*structures.TileVariant {
	rand.Seed(time.Now().UnixNano())
	characters := []byte{'a', 'c', 'g', 't'}
	var tileLists [][][]*structures.TileVariant
	tileLists = make([][][]*structures.TileVariant, structures.Paths, structures.Paths)
	for path:=0; path<structures.Paths; path++ {
		numberOfSteps := minSteps + rand.Intn(maxSteps-minSteps)
		tileLists[path] = make([][]*structures.TileVariant, numberOfSteps, numberOfSteps)
		for step:=0; step<numberOfSteps; step++ {
			numberOfTiles := minTiles + rand.Intn(maxTiles-minTiles)
			for i:=0; i<numberOfTiles; i++ {
				tileLength := minBases + rand.Intn(maxBases-minBases)
				var b strings.Builder
				b.Reset()
				for i:=0; i<tileLength; i++ {
					b.WriteByte(characters[rand.Intn(4)])
				}
				bases := b.String()
				hashArray := md5.Sum([]byte(bases))
				newVariants := &structures.TileVariant{Hash: hashArray, Length: 1, Annotation: "", LookupReference: -1}
				tileLists[path][step] = append(tileLists[path][step], newVariants)
			}
		}
	}
	return tileLists
}

func BenchmarkSort(b *testing.B) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		SortLibrary(&l1)
		b.StopTimer()
		l1 = generateRandomData("a", [][md5.Size]byte{}, tiles)
		b.StartTimer()
	}
}

//Generate a random set of data based on a random set of data previously generated.
func generateRandomData(text string, components [][md5.Size]byte, tileLists [][][]*structures.TileVariant) Library {
	l := InitializeLibrary(text, components)
	rand.Seed(time.Now().UnixNano())
	for path:=0; path<structures.Paths; path++ {
		for step, stepList := range tileLists[path] {
			numberOfTiles := len(stepList)
			numberOfPhases := 20 + rand.Intn(20) // Anywhere between 10 and 20 genomes worth of data here.
			for k := 0; k<numberOfPhases; k++ {
				randomTile := rand.Intn(numberOfTiles)
				AddTile(path, step, tileLists[path][step][randomTile], &l) // The fields marked "test" don't affect anything, so they can be named anything.
			}
		}
	}
	return l
}

func BenchmarkGenerateTiles(b *testing.B) {
	for i:=0; i<b.N; i++ {
		generateRandomTiles(10, 15, 500, 1000, 248, 300)
	}
}

func BenchmarkGenerateLibraries(b *testing.B) {
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		generateRandomData("test", [][md5.Size]byte{}, tiles)
	}
}

// Test for the libraryCopy function, which copies libraries.
func TestCopy(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	l4 := InitializeLibrary("b", [][md5.Size]byte{})
	libraryCopy(&l4, &l1)
	if !l1.Equals(l4) {
		t.Errorf("libraries are not equal, but copy should make them equal.")
	}
}

func BenchmarkCopy(b *testing.B) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		l4 := InitializeLibrary("b", [][md5.Size]byte{})
		libraryCopy(&l4, &l1)
	}
}

func TestMerging(t *testing.T) {
	log.SetFlags(log.Llongfile)
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	l1.AssignID()
	SortLibrary(&l1)
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	l2.AssignID()
	SortLibrary(&l2)
	l3:= MergeLibraries(&l1,&l2)
	l3.AssignID()
	l4 := InitializeLibrary("c", [][md5.Size]byte{}) // Created as a copy of l3 to check if the count of each tile is correct.
	libraryCopy(&l4, l3)
	for path := range l1.Paths {
		l1.Paths[path].Lock.RLock()
		for step, stepList := range l1.Paths[path].Variants {
			if stepList != nil {
				for i, variant := range (*stepList).List {
					if index1 := TileExists(path, step, variant, l3); index1 == -1 {
						t.Fatalf("a tile in library 1 is not in library 3")
					} else {
						l4.Paths[path].Lock.Lock()
						l4.Paths[path].Variants[step].Counts[index1]-=stepList.Counts[i]
						l4.Paths[path].Lock.Unlock()
					}
				}
			}
		}
		l1.Paths[path].Lock.RUnlock()
	}
	for path := range l2.Paths {
		l2.Paths[path].Lock.RLock()
		for step, stepList := range l2.Paths[path].Variants {
			if stepList != nil {
				for i, variant := range (*stepList).List {
					if index2 := TileExists(path, step, variant, l3); index2 == -1 {
						t.Fatalf("a tile in library 2 is not in library 3")
					} else {
						l4.Paths[path].Lock.Lock()
						l4.Paths[path].Variants[step].Counts[index2]-=stepList.Counts[i]
						l4.Paths[path].Lock.Unlock()
					}
				}
			}
		}
		l2.Paths[path].Lock.RUnlock()
	}
	for path := range l3.Paths {
		l3.Paths[path].Lock.Lock()
		for _, stepList := range l3.Paths[path].Variants {
			if stepList != nil {
				if !sort.IsSorted(sort.Reverse(sort.IntSlice((*stepList).Counts))) {
					t.Fatalf("a step is not sorted in descending order")
				}
			}
		}
		l3.Paths[path].Lock.Unlock()
	}
	for path := range l4.Paths {
		l4.Paths[path].Lock.Lock()
		for _, stepList := range l4.Paths[path].Variants {
			if stepList != nil {
				for _, count := range stepList.Counts {
					if count != 0 {
						t.Fatalf("counts are not correct in merged library")
					}
				}
			}
		}
		l4.Paths[path].Lock.Unlock()
	}
	if l1.ID != l3.Components[0] || l2.ID != l3.Components[1] {
		t.Errorf("components are not correct")
	}
}

func BenchmarkMerging(b *testing.B) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	SortLibrary(&l1)
	l1.AssignID()
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	SortLibrary(&l2)
	l2.AssignID()
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		MergeLibraries(&l1,&l2)
	}
}

// Test for the CreateMapping function, which makes liftover mappings from one library to another.
func TestLiftover(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	SortLibrary(&l1)
	l1.AssignID()
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	SortLibrary(&l2)
	l2.AssignID()
	l3 := MergeLibraries(&l1,&l2)
	newMapping := CreateMapping(&l1, l3)
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l1.Paths[path].Variants[step].List[oldIndex]).Equals(*l3.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 1 are not equal.")
				}
			} 
		}
	}

	newMapping2 := CreateMapping(&l2, l3)
	for path, stepMapping := range newMapping2.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l2.Paths[path].Variants[step].List[oldIndex]).Equals(*l3.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 2 are not equal.")
				}
			} 
		}
	}
}

func BenchmarkLiftover(b *testing.B) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	SortLibrary(&l1)
	l1.AssignID()
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	SortLibrary(&l2)
	l2.AssignID()
	l3 := MergeLibraries(&l1,&l2)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		CreateMapping(&l1, l3)
	}
}

func BenchmarkID(b *testing.B) {
	rand.Seed(1)
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300)
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		l1.AssignID()
	}
}

// Tests the ParseSGLFv2 and AddLibrarySGLFv2 functions by writing a library to SGLFv2 files, converting it back, and testing if both libraries are equal.
// Also tests IDs and references created from AddLibrarySGLFv2.
func TestParseSGLFv2(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	SortLibrary(&l)
	l.AssignID()
	WriteLibraryToSGLFv2(&l, "/data-sdc/jc/tile-library") // Writes it to disk and gets it back to make sure that the libraries will be the same
	l1:=InitializeLibrary("/data-sdc/jc/tile-library", [][md5.Size]byte{})
	AddLibrarySGLFv2("/data-sdc/jc/tile-library", &l1)
	if !l1.Equals(l) || !l.Equals(l1) || l1.ID != l.ID {
		t.Errorf("expected both libraries to be equal, but they aren't")
	}
	WriteLibraryToSGLFv2(&l1, "/data-sdc/jc/tile-library/test") // Writes it to disk again to test that writing from a directory of sglfv2 files works.
	l2:=InitializeLibrary("/data-sdc/jc/tile-library/test", [][md5.Size]byte{})
	AddLibrarySGLFv2("/data-sdc/jc/tile-library/test", &l2)
	if !l1.Equals(l2) || !l2.Equals(l1) || l1.ID != l2.ID {
		t.Errorf("expected second and third libraries to be equal, but they aren't")
	}
}

// Test for merging libraries.
func TestParseSGLFv2WithMerge(t *testing.T) {
	log.SetFlags(log.Llongfile)
	l:=InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	SortLibrary(&l)
	l.AssignID()
	l1:=InitializeLibrary("/data-sdc/jc/tile-library/testing/test.txt", [][md5.Size]byte{})
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM", "/data-sdc/jc/tile-library/testing/test.txt",&l1)

	SortLibrary(&l1)
	l1.AssignID()
	l2:=MergeLibraries(&l, &l1)
	(*l2).AssignID()
	newMapping := CreateMapping(&l, l2)
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l.Paths[path].Variants[step].List[oldIndex]).Equals(*l2.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 1 are not equal.")
				}
			} 
		}
	}

	newMapping2 := CreateMapping(&l1, l2)
	for path, stepMapping := range newMapping2.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l1.Paths[path].Variants[step].List[oldIndex]).Equals(*l2.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 2 are not equal.")
				}
			} 
		}
	}
	fmt.Println("verified that merged library is valid")
	WriteLibraryToSGLFv2(l2, "/data-sdc/jc/tile-library/test2")
	l3:=InitializeLibrary("/data-sdc/jc/tile-library/test2", [][md5.Size]byte{})
	AddLibrarySGLFv2("/data-sdc/jc/tile-library/test2", &l3) // Verification that the library is written correctly.
	if !l3.Equals(*l2) || !l2.Equals(l3) || l3.ID != l2.ID {
		t.Errorf("expected second and third libraries to be equal, but they aren't")
	}
	l4 := InitializeLibrary("/data-sdc/jc/tile-library/test3/test.txt", [][md5.Size]byte{}) // Comparing a library made from merging to a library made without merging
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test3/test.txt",&l4)
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM", "/data-sdc/jc/tile-library/test3/test.txt",&l4)
	SortLibrary(&l4)
	l4.AssignID()
	WriteLibraryToSGLFv2(&l4, "/data-sdc/jc/tile-library/test3")
	if !l4.Equals(*l2) || !l2.Equals(l4) || l4.ID != l2.ID {
		t.Errorf("expected merged and unmerged libraries to be equal, but they aren't") // check why test is failing here.
	}
}

// Test to determine if sorting is not dependent on the order of libraries given.
func TestIdenticalSort(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	l1.AssignID()
	SortLibrary(&l1)
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	l2.AssignID()
	SortLibrary(&l2)
	l3 := MergeLibraries(&l1,&l2)
	l4 := MergeLibraries(&l2, &l1)
	if !l3.Equals(*l4) {
		t.Errorf("sort is different for different ordering of merge")
	}
}

func BenchmarkMD5(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	lst := make([]byte, 6000000000, 6000000000)
	rand.Read(lst)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		md5.Sum(lst)
	}
}

func TestMapping(t *testing.T) {
	log.SetFlags(log.Llongfile)
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	l1.AssignID()
	SortLibrary(&l1)
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	l2.AssignID()
	SortLibrary(&l2)
	l3:= MergeLibraries(&l1,&l2)
	l3.AssignID()
	newMapping := CreateMapping(&l1, l3)
	WriteMapping("/data-sdc/jc/test.sglfmapping", newMapping)
	mapping, source, destination := ReadMapping("/data-sdc/jc/test.sglfmapping")
	for path := range mapping {
		if len(mapping[path]) != len(newMapping.Mapping[path]) {
			t.Fatalf("path %v doesn't have same length", path)
		}
		for step := range mapping[path] {
			if mapping[path][step] != nil && newMapping.Mapping[path][step] != nil {
				if len(mapping[path][step]) != len(newMapping.Mapping[path][step]) {
					t.Fatalf("step %v doesn't have same length", step)
				}
				for index := range mapping[path][step] {
					if mapping[path][step][index] != newMapping.Mapping[path][step][index] {
						t.Fatalf("an index/value combination does not match")
					}
				}
			} else if mapping[path][step] != nil || newMapping.Mapping[path][step] != nil {
				t.Fatalf("exactly one path/step combination is non-nil")
			}
		}
	}
	if source != l1.ID {
		t.Errorf("source IDs are not the same")
	}
	if destination != l3.ID {
		t.Errorf("destination IDs are not the same")
	}
}