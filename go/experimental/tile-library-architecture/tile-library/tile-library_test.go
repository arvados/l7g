package tilelibrary

import (
	"crypto/md5"
	//"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"testing"
	"../structures"
)
/*
// This is a test of the TileExists function.
func TestExists(t *testing.T) {
	l := InitializeLibrary()
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	tile1 := structures.TileVariant{hash1, 1, ""}
	tile2 := structures.TileVariant{hash1, 3, "testing"}
	tile3 := structures.TileVariant{hash2, 2, ""}
	AddTile(0,0,tile1, "a", &l)
	test1 := TileExists(0,0,tile1, &l)
	if test1==-1 {
		t.Errorf("Test 1 is false, expected true")
	}
	test2 := TileExists(0,0,tile2, &l)
	if test2==-1 {
		t.Errorf("Test 2 is false, expected true") // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash.
	}
	test3 := TileExists(0,0,tile3,&l)
	if test3!=-1 {
		t.Errorf("Test 3 is true, expected false")
	}
}

// This is a test of the FindFrequency function.
func TestFrequency(t *testing.T) {
	l := InitializeLibrary()
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	tile1 := structures.TileVariant{hash1, 1, ""}
	tile2 := structures.TileVariant{hash1, 3, "testing"}
	tile3 := structures.TileVariant{hash2, 2, ""}
	AddTile(0,0,tile1, "a", &l)
	AddTile(0,0,tile1, "g", &l)
	AddTile(0,0,tile2, "c", &l)
	AddTile(0,0,tile3, "a", &l)
	AddTile(0,1,tile1, "t", &l)
	test1 := FindFrequency(0,0,tile1,&l)
	if test1!=3 {
		t.Errorf("Test 1 was %v, expected 3", test1) // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash, so tile1 counts the two instances of tile1 and the instance of tile2 in path 0, step 0.
	}
	test2 := FindFrequency(0,0,tile2, &l)
	if test2!=3 {
		t.Errorf("Test 2 was %v, expected 3", test2)
	}
	test3 := FindFrequency(0,0,tile3, &l)
	if test3!=1 {
		t.Errorf("Test 3 was %v, expected 1", test3)
	}
	test4 := FindFrequency(0,1,tile1, &l)
	if test4!=1 {
		t.Errorf("Test 4 was %v, expected 1",test4)
	}
	test5 := FindFrequency(0,1,tile3, &l)
	if test5!=0 {
		t.Errorf("Test 5 was %v, expected 0",test5)
	}
	test6 := FindFrequency(0,2, tile2, &l)
	if test6 !=0 {
		t.Errorf("Test 6 was %v, expacted 0", test6)
	}
}

func TestLibraryMerging(t *testing.T) {
	l := InitializeLibrary()
	l2 := InitializeLibrary
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	hash2 := md5.Sum([]byte{3})
	tile1 := structures.TileVariant{hash1, 1, ""}
	tile2 := structures.TileVariant{hash2, 3, "testing"}
	tile3 := structures.TileVariant{hash3, 2, ""}
	AddTile(0,0,tile1, "a", &l)
	AddTile(0,0,tile1, "g", &l2)
	AddTile(0,0,tile2, "c", &l2)
	AddTile(0,0,tile3, "a", &l)
	AddTile(0,1,tile1, "t", &l2)
	mergeLibraries()
}


// This is a test of the sortLibrary function.
func TestLibrarySorting(t *testing.T) {
	l := InitializeLibrary()
	hash1 := md5.Sum([]byte{1})
	hash2 := md5.Sum([]byte{2})
	hash3 := md5.Sum([]byte{3})
	tile1 := structures.TileVariant{hash1, 1, ""}
	tile2 := structures.TileVariant{hash2, 3, "testing"}
	tile3 := structures.TileVariant{hash3, 2, ""}
	AddTile(0,0,tile1, "a", &l)
	AddTile(0,0,tile1, "g", &l)
	AddTile(0,0,tile2, "c", &l)
	AddTile(0,0,tile2, "c", &l)
	AddTile(0,0,tile3, "a", &l)
	AddTile(0,0,tile2, "c", &l)
	AddTile(0,1,tile3, "a", &l)
	AddTile(0,1,tile1, "t", &l)
	AddTile(0,1,tile3, "a", &l)
	AddTile(0,1,tile3, "a", &l)
	AddTile(0,2,tile3, "a", &l)
	AddTile(0,2,tile1, "t", &l)
	sortLibrary(&l)
	test1 := sort.IsSorted(sort.Reverse(sort.IntSlice(l[0][0].Counts)))
	if !test1 {
		t.Errorf("Path 0 step 0 is not sorted correctly.")
	}
	test2 := sort.IsSorted(sort.Reverse(sort.IntSlice(l[0][1].Counts)))
	if !test2 {
		t.Errorf("Path 0 step 1 is not sorted correctly.")
	}
	test3 := sort.IsSorted(sort.Reverse(sort.IntSlice(l[0][2].Counts)))
	if !test3 {
		t.Errorf("Path 0 step 2 is not sorted correctly.")
	}
}

// Benchmarks to be put in later
*/

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
			for k :=0; k<numberOfPhases; k++ {
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
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", [][md5.Size]byte{}, tiles)
	l1.AssignID()
	SortLibrary(&l1)
	l2 := generateRandomData("b", [][md5.Size]byte{}, tiles)
	l2.AssignID()
	SortLibrary(&l2)
	l3:= MergeLibraries(&l1,&l2)
	for path := range l1.Paths {
		l1.Paths[path].Lock.RLock()
		for step, stepList := range l1.Paths[path].Variants {
			for _, variant := range (*stepList).List {
				if TileExists(path, step, variant, l3) == -1 {
					t.Fatalf("a tile in library 1 is not in library 3")
				}
			}
		}
		l1.Paths[path].Lock.RUnlock()
	}
	for path := range l2.Paths {
		l2.Paths[path].Lock.RLock()
		for step, stepList := range l2.Paths[path].Variants {
			for _, variant := range (*stepList).List {
				if TileExists(path, step, variant, l3) == -1 {
					t.Fatalf("a tile in library 2 is not in library 3")
				}
			}
		}
		l2.Paths[path].Lock.RUnlock()
	}
	for path := range l3.Paths {
		l3.Paths[path].Lock.Lock()
		for _, stepList := range l3.Paths[path].Variants {
			if !sort.IsSorted(sort.Reverse(sort.IntSlice((*stepList).Counts))) {
				t.Fatalf("a step is not sorted in descending order")
			}
		}
		l3.Paths[path].Lock.Unlock()
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

// Test for the createMapping function, which makes liftover mappings from one library to another.
// TODO: write benchmarks for functions in addition to tests.
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
					t.Fatalf("variants from liftover in library 1 are not equal.")
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
// Also tests IDs.
func TestParseSGLFv2(t *testing.T) {
	l:=InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	SortLibrary(&l)
	l.AssignID()
	WriteLibraryToSGLFv2(&l, "/data-sdc/jc/tile-library", "/data-sdc/jc/tile-library", "test.txt")
	l1:=InitializeLibrary("/data-sdc/jc/tile-library/test.txt", [][md5.Size]byte{})
	AddLibrarySGLFv2("/data-sdc/jc/tile-library", &l1)
	if !l1.Equals(l) || !l.Equals(l1) || l1.ID != l.ID {
		t.Errorf("expected both libraries to be equal, but they aren't")
	}
}