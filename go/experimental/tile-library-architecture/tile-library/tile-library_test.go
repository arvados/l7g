// tile-library_test.go is the test file for the tilelibrary package.
// Places where test files and directories would normally be are intentionally left blank.
// Fill them in with wherever data should go.
package tilelibrary

import (
	"crypto/md5"
	"log"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"../structures"
)

// This is a test of the TileExists function.
func TestExists(t *testing.T) {
	testFile := "" // Put your test here.
	l, err := New(testFile, [][md5.Size]byte{[md5.Size]byte{}})
	if err != nil {
		log.Fatal(err)
	}
	hash1 := md5.Sum([]byte("a"))
	hash2 := md5.Sum([]byte("b"))
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 1, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	l.AddTile(0, 0, &tile1, "a")
	_, ok := l.TileExists(0, 0, &tile1)
	if !ok {
		t.Errorf("Test 1 is false, expected true")
	}
	_, ok = l.TileExists(0, 0, &tile2)
	if !ok {
		t.Errorf("Test 2 is false, expected true") // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash.
	}
	_, ok = l.TileExists(0, 0, &tile3)
	if ok {
		t.Errorf("Test 3 is true, expected false")
	}
}

// This is a test of the FindFrequency function.
func TestFrequency(t *testing.T) {
	testFile := "" // Put your test file here.
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	hash1 := md5.Sum([]byte("a"))
	hash2 := md5.Sum([]byte("b"))
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 1, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 0, &tile3, "b")
	l.AddTile(0, 1, &tile1, "a")
	test1 := l.FindFrequency(0, 0, &tile1)
	if test1 != 3 {
		t.Errorf("Test 1 was %v, expected 3", test1) // This is because tile1 and tile2 have the same hash--equality of tiles is only based on hash, so tile1 counts the two instances of tile1 and the instance of tile2 in path 0, step 0.
	}
	test2 := l.FindFrequency(0, 0, &tile2)
	if test2 != 3 {
		t.Errorf("Test 2 was %v, expected 3", test2)
	}
	test3 := l.FindFrequency(0, 0, &tile3)
	if test3 != 1 {
		t.Errorf("Test 3 was %v, expected 1", test3)
	}
	test4 := l.FindFrequency(0, 1, &tile1)
	if test4 != 1 {
		t.Errorf("Test 4 was %v, expected 1", test4)
	}
	test5 := l.FindFrequency(0, 1, &tile3)
	if test5 != 0 {
		t.Errorf("Test 5 was %v, expected 0", test5)
	}
	test6 := l.FindFrequency(0, 2, &tile2)
	if test6 != 0 {
		t.Errorf("Test 6 was %v, expacted 0", test6)
	}
}

// This is a test of the SortLibrary function.
func TestLibrarySorting(t *testing.T) {
	testFile := "" // Put your test file here.
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	hash1 := md5.Sum([]byte("a"))
	hash2 := md5.Sum([]byte("b"))
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 1, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 0, &tile3, "b")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 1, &tile1, "a")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 2, &tile3, "b")
	l.AddTile(0, 2, &tile1, "a")
	l.SortLibrary()
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

// This is a test on AddTile to make sure it writes tile information to disk properly.
func TestAddTile(t *testing.T) {
	testFile := "" // Put your test file here.
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	hash1 := md5.Sum([]byte("a"))
	hash2 := md5.Sum([]byte("b"))
	tile1 := structures.TileVariant{hash1, 1, "", 0, true, &l}
	tile2 := structures.TileVariant{hash1, 1, "testing", 1, true, &l}
	tile3 := structures.TileVariant{hash2, 2, "", 2, false, &l}
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile1, "a")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 0, &tile3, "b")
	l.AddTile(0, 0, &tile2, "a")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 1, &tile1, "a")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 1, &tile3, "b")
	l.AddTile(0, 2, &tile3, "b")
	l.AddTile(0, 2, &tile1, "a")
	l.SortLibrary()
	l.AssignID()
	testDirectory := "" // Put your test directory here.
	err = l.WriteLibraryToSGLFv2(testDirectory)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l1, err := New(testDirectory, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibrarySGLFv2()
	if !l1.Equals(l) {
		t.Fatalf("libraries are not equal")
	}
}

// Note: some of the following tests and benchmarks use random libraries and random tile generation, so they will have a (slightly) variable amount of time taken.

// tileAndBases is a test struct to group a variant and its bases together.
type tileAndBases struct {
	tile  *structures.TileVariant
	bases string
}

// generateRandomTiles generates a set of tiles for each path and step for random libraries to use.
func generateRandomTiles(minTiles, maxTiles, minSteps, maxSteps, minBases, maxBases int) [][][]tileAndBases {
	rand.Seed(time.Now().UnixNano())
	characters := []byte{'a', 'c', 'g', 't'}
	var tileLists [][][]tileAndBases
	tileLists = make([][][]tileAndBases, structures.Paths, structures.Paths)
	for path := 0; path < structures.Paths; path++ {
		numberOfSteps := minSteps + rand.Intn(maxSteps-minSteps)
		tileLists[path] = make([][]tileAndBases, numberOfSteps, numberOfSteps)
		for step := 0; step < numberOfSteps; step++ {
			numberOfTiles := minTiles + rand.Intn(maxTiles-minTiles)
			for i := 0; i < numberOfTiles; i++ {
				tileLength := minBases + rand.Intn(maxBases-minBases)
				var b strings.Builder
				b.Reset()
				for i := 0; i < tileLength; i++ {
					b.WriteByte(characters[rand.Intn(4)])
				}
				bases := b.String()
				hashArray := md5.Sum([]byte(bases))
				newVariants := &structures.TileVariant{Hash: hashArray, Length: 1, Annotation: "", LookupReference: -1}
				tileLists[path][step] = append(tileLists[path][step], tileAndBases{newVariants, bases})
			}
		}
	}
	return tileLists
}

// This benchmarks sorting libraries on a randomly generated library.
func BenchmarkSort(b *testing.B) {
	testFile := ""                                            // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l1.SortLibrary()
		b.StopTimer()
		l1 = generateRandomData(testFile, nil, tiles)
		b.StartTimer()
	}
}

// Generates a random library based on a random set of data previously generated.
func generateRandomData(text string, components [][md5.Size]byte, tileLists [][][]tileAndBases) *Library {
	l, err := New(text, components)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	for path := 0; path < structures.Paths; path++ {
		for step, stepList := range tileLists[path] {
			numberOfTiles := len(stepList)
			numberOfPhases := 20 + rand.Intn(20) // Anywhere between 10 and 20 genomes worth of data here.
			for k := 0; k < numberOfPhases; k++ {
				randomTile := rand.Intn(numberOfTiles)
				randomTileAndBases := tileLists[path][step][randomTile]
				if _, ok := l.TileExists(path, step, randomTileAndBases.tile); ok {
					l.addTileUnsafe(path, step, randomTileAndBases.tile)
				}
				// Above is a silent version of generating random data, which does not write to files.
				// If you would also like to test the writing of random data to files, uncomment the line below.
				// AddTile(path, step, randomTileAndBases.tile, randomTileAndBases.bases) // The fields marked "test" don't affect anything, so they can be named anything.
			}
		}
	}
	return l
}

// Benchmark to see how fast random tiles can be generated.
func BenchmarkGenerateTiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateRandomTiles(10, 15, 500, 1000, 248, 300)
	}
}

// Benchmark to see how fast random libraries can be generated.
func BenchmarkGenerateLibraries(b *testing.B) {
	testFile := "" // Put your test file here.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateRandomData(testFile, nil, tiles)
	}
}

// Test for the libraryCopy function, which copies libraries.
func TestCopy(t *testing.T) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l4, err := New(testFile2, nil)
	if err != nil {
		log.Fatal(err)
	}
	libraryCopy(l4, l1)
	if !l1.Equals(l4) {
		t.Errorf("libraries are not equal, but copy should make them equal.")
	}
}

// Benchmark to see how fast a library can be copied.
func BenchmarkCopy(b *testing.B) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l4, err := New(testFile2, nil)
		if err != nil {
			log.Fatal(err)
		}
		libraryCopy(l4, l1)
	}
}

// Test for merging two random libraries.
func TestMerging(t *testing.T) {
	testFile := ""  // Put your test file here.
	testFile2 := "" // Put your test file here.
	testFile3 := "" // Put your test file here.
	testFile4 := "" // Put your test file here.
	log.SetFlags(log.Llongfile)
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.AssignID()
	l1.SortLibrary()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.AssignID()
	l2.SortLibrary()
	l3, err := l2.MergeLibraries(l1, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l3.AssignID()
	l4, err := New(testFile4, nil) // Created as a copy of l3 to check if the count of each tile is correct.
	if err != nil {
		log.Fatal(err)
	}
	libraryCopy(l4, l3)
	for path := range l1.Paths {
		l1.Paths[path].Lock.RLock()
		for step, stepList := range l1.Paths[path].Variants {
			if stepList != nil {
				for i, variant := range (*stepList).List {
					if index1, ok := l3.TileExists(path, step, variant); !ok {
						t.Fatalf("a tile in library 1 is not in library 3")
					} else {
						l4.Paths[path].Lock.Lock()
						l4.Paths[path].Variants[step].Counts[index1] -= stepList.Counts[i]
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
					if index2, ok := l3.TileExists(path, step, variant); !ok {
						t.Fatalf("a tile in library 2 is not in library 3")
					} else {
						l4.Paths[path].Lock.Lock()
						l4.Paths[path].Variants[step].Counts[index2] -= stepList.Counts[i]
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

// Benchmark to see how fast the library merging operation is.
func BenchmarkMerging(b *testing.B) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	testFile3 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.SortLibrary()
	l1.AssignID()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.SortLibrary()
	l2.AssignID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l2.MergeLibraries(l1, testFile3)
	}
}

// Test for the CreateMapping function, which makes liftover mappings from one library to another.
func TestLiftover(t *testing.T) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	testFile3 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.SortLibrary()
	l1.AssignID()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.SortLibrary()
	l2.AssignID()
	l3, err := l2.MergeLibraries(l1, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newMapping, err := CreateMapping(l1, l3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l1.Paths[path].Variants[step].List[oldIndex]).Equals(*l3.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 1 are not equal.")
				}
			}
		}
	}

	newMapping2, err := CreateMapping(l2, l3)
	if err != nil {
		t.Fatalf(err.Error())
	}
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

// Benchmark to see how quickly LiftoverMappings can be made.
func BenchmarkLiftover(b *testing.B) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	testFile3 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.SortLibrary()
	l1.AssignID()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.SortLibrary()
	l2.AssignID()
	l3, err := l2.MergeLibraries(l1, testFile3)
	if err != nil {
		b.Fatalf(err.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateMapping(l1, l3)
	}
}

// Benchmark to see how fast a library ID can be given to a random library.
func BenchmarkID(b *testing.B) {
	testFile := "" // Put your test file here.
	rand.Seed(1)
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300)
	l1 := generateRandomData(testFile, nil, tiles)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l1.AssignID()
	}
}

// Tests the ParseSGLFv2 and AddLibrarySGLFv2 functions by writing a library to SGLFv2 files, converting it back, and testing if both libraries are equal.
// Also tests IDs and references created from AddLibrarySGLFv2.
func TestParseSGLFv2(t *testing.T) {
	testFile := ""         // Put your test file here.
	testFile2 := ""        // Put your test file here.
	testFile3 := ""        // Put your test file here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.AddLibraryFastJ(genomeDirectory2)
	l.SortLibrary()
	l.AssignID()
	l.WriteLibraryToSGLFv2(testFile2) // Writes it to disk and gets it back to make sure that the libraries will be the same
	l1, err := New(testFile2, nil)
	if err != nil {
		log.Fatal(err)
	}
	l1.AddLibrarySGLFv2()
	if !l1.Equals(l) || !l.Equals(l1) || l1.ID != l.ID {
		t.Errorf("expected both libraries to be equal, but they aren't")
	}
	l1.WriteLibraryToSGLFv2(testFile3) // Writes it to disk again to test that writing from a directory of sglfv2 files works.
	l2, err := New(testFile3, nil)
	if err != nil {
		log.Fatal(err)
	}
	l2.AddLibrarySGLFv2()
	if !l1.Equals(l2) || !l2.Equals(l1) || l1.ID != l2.ID {
		t.Errorf("expected second and third libraries to be equal, but they aren't")
	}
}

// Test for merging libraries, based on real data.
func TestParseSGLFv2WithMerge(t *testing.T) {
	testFile := ""         // Put your test file here.
	testFile2 := ""        // Put your test file here.
	testFile3 := ""        // Put your test file here.
	testFile4 := ""        // Put your test file here.
	testDirectory1 := ""   // Put your test directory here.
	testDirectory2 := ""   // Put your test directory here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.SortLibrary()
	l.AssignID()
	l1, err := New(testFile2, nil)
	if err != nil {
		log.Fatal(err)
	}
	l1.AddLibraryFastJ(genomeDirectory2)

	l1.SortLibrary()
	l1.AssignID()
	l2, err := l1.MergeLibraries(l, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	(*l2).AssignID()
	newMapping, err := CreateMapping(l, l2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l.Paths[path].Variants[step].List[oldIndex]).Equals(*l2.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 1 are not equal.")
				}
			}
		}
	}

	newMapping2, err := CreateMapping(l1, l2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for path, stepMapping := range newMapping2.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l1.Paths[path].Variants[step].List[oldIndex]).Equals(*l2.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 2 are not equal.")
				}
			}
		}
	}
	l2.WriteLibraryToSGLFv2(testDirectory1)
	l3, err := New(testDirectory1, nil)
	if err != nil {
		log.Fatal(err)
	}
	l3.AddLibrarySGLFv2() // Verification that the library is written correctly.
	if !l3.Equals(l2) || !l2.Equals(l3) || l3.ID != l2.ID {
		t.Errorf("expected second and third libraries to be equal, but they aren't")
	}
	l4, err := New(testFile4, nil) // Comparing a library made from merging to a library made without merging
	if err != nil {
		log.Fatal(err)
	}
	l4.AddLibraryFastJ(genomeDirectory1)
	l4.AddLibraryFastJ(genomeDirectory2)
	l4.SortLibrary()
	l4.AssignID()
	l4.WriteLibraryToSGLFv2(testDirectory2)
	if !l4.Equals(l2) || !l2.Equals(l4) || l4.ID != l2.ID {
		t.Errorf("expected merged and unmerged libraries to be equal, but they aren't") // check why test is failing here.
	}
}

// Test for merging libraries without creating a new library, based on real data.
func TestParseSGLFv2WithMergeWithoutCreation(t *testing.T) {
	testFile := ""         // Put your test file here.
	testFile2 := ""        // Put your test file here.
	testDirectory1 := ""   // Put your test directory here.
	genomeDirectory1 := "" // Put your genome directory here.
	genomeDirectory2 := "" // Put your genome directory here.
	log.SetFlags(log.Llongfile)
	l, err := New(testFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	l.AddLibraryFastJ(genomeDirectory1)
	l.SortLibrary()
	l.AssignID()
	l1, err := New(testFile2, nil)
	if err != nil {
		log.Fatal(err)
	}
	l1.AddLibraryFastJ(genomeDirectory2)

	l1.SortLibrary()
	l1.AssignID()
	l2, err := l1.MergeLibrariesWithoutCreation(l)
	if err != nil {
		t.Fatalf(err.Error())
	}
	(*l2).AssignID()
	newMapping, err := CreateMapping(l, l2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l.Paths[path].Variants[step].List[oldIndex]).Equals(*l2.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("variants from liftover in library 1 are not equal.")
				}
			}
		}
	}

	l2.WriteLibraryToSGLFv2(testDirectory1)
	l3, err := New(testDirectory1, nil)
	if err != nil {
		log.Fatal(err)
	}
	l3.AddLibrarySGLFv2() // Verification that the library is written correctly.
	if !l3.Equals(l2) || !l2.Equals(l3) || l3.ID != l2.ID {
		t.Errorf("expected second and third libraries to be equal, but they aren't")
	}
}

// Test to determine if sorting is not dependent on the order of libraries given.
func TestIdenticalSort(t *testing.T) {
	testFile := ""                                            // Put your test file here.
	testFile2 := ""                                           // Put your test file here.
	testFile3 := ""                                           // Put your test file here.
	testFile4 := ""                                           // Put your test file here.
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.AssignID()
	l1.SortLibrary()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.AssignID()
	l2.SortLibrary()
	l3, err := l2.MergeLibraries(l1, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l4, err := l1.MergeLibraries(l2, testFile4)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !l3.Equals(l4) {
		t.Errorf("sort is different for different ordering of merge")
	}
}

// Test to create liftover mappings on random libraries.
func TestMapping(t *testing.T) {
	testFile := ""    // Put your test file here.
	testFile2 := ""   // Put your test file here.
	testFile3 := ""   // Put your test file here.
	mappingFile := "" // Put the name of the mapping here.
	log.SetFlags(log.Llongfile)
	rand.Seed(1)                                              // Seed for randomness can be changed.
	tiles := generateRandomTiles(10, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData(testFile, nil, tiles)
	l1.AssignID()
	l1.SortLibrary()
	l2 := generateRandomData(testFile2, nil, tiles)
	l2.AssignID()
	l2.SortLibrary()
	l3, err := l2.MergeLibraries(l1, testFile3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l3.AssignID()
	newMapping, err := CreateMapping(l1, l3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	WriteMapping(mappingFile, newMapping)
	mapping, source, destination, err := ReadMapping(mappingFile)
	if err != nil {
		t.Fatalf(err.Error())
	}
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
