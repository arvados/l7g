package structures
/*
	Structures is a basic package to hold basic structures, methods, and functions for tile libraries and genomes.
*/
import (
	"compress/gzip"
	"crypto/md5"
	"io/ioutil"
	"log"
	"os"
)

// TileVariant is a struct for representing a tile variant using the hash, length, and any annotation(s) on the variant.
type TileVariant struct {
	Hash       VariantHash 	  // The hash of the tile variant's bases.
	Length     int            // The length (span) of the tile
	Annotation string         // Any notes about this tile (by default, no comments)
	// File string // The path to the file for which the variant is from.
}

// TileCreator is a small function to create a new tile given information about it.
func TileCreator(hash VariantHash, length int, annotation string) TileVariant {
	return TileVariant{hash, length, annotation}
}

// VariantHash is a hash for a tile variant--currently the hash algorithm is MD5
type VariantHash [md5.Size]byte

// Equals checks for equality of variants based on hash.
// This works based on the assumption that no two tile variants in the same path and step have the same hash.
func (t TileVariant) Equals(t2 TileVariant) bool {
	return (t.Hash == t2.Hash)
}

// Paths is the number of paths a genome has.
const Paths int = 863 // Constant because we know that the number of paths is always the same.

// OpenGZ is a function to open gzipped files and return the corresponding slice of bytes.
// Mostly important for gzipped FastJs, but other gzipped files can be opened too.
func OpenGZ(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	gz, err2 := gzip.NewReader(file)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer gz.Close()
	data, err3 := ioutil.ReadAll(gz)
	if err3 != nil {
		log.Fatal(err3)
	}

	return data
}