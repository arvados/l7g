package main

import (
	"crypto/md5"
	"errors"
	//"fmt"
	"log"
	"math/rand"
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
// generateRandomTiles generates a set of tiles for each path and step for random libraries to use.
func generateRandomTiles(minTiles, maxTiles, minSteps, maxSteps, minBases, maxBases int) [][][]*structures.TileVariant {
	rand.Seed(time.Now().UnixNano())
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
					character := rand.Intn(4)
					if character == 0 {
						b.WriteByte('a')
					} else if character == 1 {
						b.WriteByte('c')
					} else if character == 2 {
						b.WriteByte('g')
					} else {
						b.WriteByte('t')
					}
				}
				bases := b.String()
				hashArray := md5.Sum([]byte(bases))
				newVariants := &structures.TileVariant{hashArray, 1, "", -1}
				tileLists[path][step] = append(tileLists[path][step], newVariants)
			}
		}
	}
	return tileLists
}

//Generate a random set of data based on a random set of data previously generated.
func generateRandomData(text string, components []string, tileLists [][][]*structures.TileVariant) Library {
	l := InitializeLibrary(text, components)
	rand.Seed(time.Now().UnixNano())
	for path:=0; path<structures.Paths; path++ {
		for step, stepList := range tileLists[path] {
			numberOfTiles := len(stepList)
			numberOfPhases := 20 + rand.Intn(20) // Anywhere between 10 and 20 genomes worth of data here.
			for k :=0; k<numberOfPhases; k++ {
				randomTile := rand.Intn(numberOfTiles)
				AddTile(path, step, -1, tileLists[path][step][randomTile], "test", "test", &l) // The fields marked "test" don't affect anything, so they can named anything.
			}
		}
	}
	return l
}

func BenchmarkLibrary(b *testing.B) {
	l := InitializeLibrary("test", []string{})
	newTile := structures.TileVariant{md5.Sum([]byte{1}), 1, "", -1}
	AddTile(0, 0, -1, &newTile, "test", "gcat", &l)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		l.Paths[0].Variants[0].List[0].LookupReference = 2
	}
}

func BenchmarkPointer(b *testing.B) {
	l := InitializeLibrary("test", []string{})
	newTile := &structures.TileVariant{md5.Sum([]byte{1}), 1, "", -1}
	info := baseInfo{"gcat", md5.Sum([]byte{1}), newTile}
	AddTile(0, 0, -1, newTile, "test", "gcat", &l)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		info.variant.LookupReference = 2
	}
}

// Test for the libraryCopy function, which copies libraries.
func TestCopy(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(3, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", []string{}, tiles)
	l4 := InitializeLibrary("b", []string{})
	libraryCopy(&l4, &l1)
	if !l1.Equals(l4) {
		t.Errorf("libraries are not equal, but copy should make them equal.")
	}
}

func TestMerging(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(3, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", []string{}, tiles)
	l2 := generateRandomData("b", []string{}, tiles)
	l3,_ := mergeLibraries("test",&l1,&l2)
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
}

func TestReference(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(3, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", []string{}, tiles)
	l2 := generateRandomData("b", []string{}, tiles)
	l3,references := mergeLibraries("test",&l1,&l2)
	index1 := -1
	for i, libraryString := range (*l3).Components {
		if libraryString == l1.Text {
			index1 = i
			break
		}
	}
	if index1 == -1 { // Destination was not made from the source--can't guarantee a mapping here.
		log.Fatal(errors.New("source library is not part of the destination library"))
	}

	index2 := -1
	for i, libraryString := range (*l3).Components {
		if libraryString == l2.Text {
			index2 = i
			break
		}
	}
	if index2 == -1 { // Destination was not made from the source--can't guarantee a mapping here.
		log.Fatal(errors.New("source library is not part of the destination library"))
	}
	for path := range references {
		for step := range references[path] {
			for index, variant := range references[path][step] {
				if variant[index1]!= -1 {
					if !(*l1.Paths[path].Variants[step].List[variant[index1]]).Equals(*l3.Paths[path].Variants[step].List[index]) {
						t.Fatalf("tiles in reference mapping don't match to library 1")
					}
				}
				if variant[index2]!= -1 {
					if !(*l2.Paths[path].Variants[step].List[variant[index2]]).Equals(*l3.Paths[path].Variants[step].List[index]) {
						t.Fatalf("tiles in reference mapping don't match to library 2")
					}
				}
			}
		}
	}
}

// Test for the createMapping function, which makes liftover mappings from one library to another.
// TODO: write benchmarks for functions in addition to tests.
func TestLiftover(t *testing.T) {
	rand.Seed(1) // Seed for randomness can be changed.
	tiles := generateRandomTiles(3, 15, 500, 1000, 248, 300) // Notice that this only creates around 600000 steps, rather than the more realistic 10 million.
	l1 := generateRandomData("a", []string{}, tiles)
	l2 := generateRandomData("b", []string{}, tiles)
	l3,_ := mergeLibraries("test",&l1,&l2)
	newMapping := createMapping(&l1, l3)
	for path, stepMapping := range newMapping.Mapping {
		for step, variantMapping := range stepMapping {
			for oldIndex, newIndex := range variantMapping {
				if !(*l1.Paths[path].Variants[step].List[oldIndex]).Equals(*l3.Paths[path].Variants[step].List[newIndex]) {
					t.Fatalf("Variants from liftover are not equal.")
				}
			} 
		}
	}
}