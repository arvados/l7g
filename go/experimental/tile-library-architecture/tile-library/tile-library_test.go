package main

import (
	"crypto/md5"
	"testing"
	"sort"
	"../structures"
)

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
/*
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
*/

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