/*Package structures is a basic package to hold basic structures, methods, and functions for tile libraries and genomes.
Most important is the TileVariant structure, which holds all the necessary data for a tile variant.
Equality on TileVariants is defined by hash.
In addition, the default number of paths in a genome is provided here for convenience.*/
package structures

import (
	"crypto/md5"
)

// TileVariant is a struct for representing a tile variant using the hash, length, and any annotation(s) on the variant.
type TileVariant struct {
	Hash             VariantHash // The hash of the tile variant's bases.
	Length           int         // The length (span) of the tile
	Annotation       string      // Any notes about this tile (by default, no comments)
	LookupReference  int64       // The lookup reference for the bases of this variant in the text file for the corresponding library.
	Complete         bool        // Tells if this variant is complete or not (complete meaning no nocalls). Mostly for genomes, to quickly tell which tiles are complete and which are not.
	ReferenceLibrary interface{} // A way of referencing the library this variant is from. (Will be a *Library).
}

// TileCreator is a small function to create a new tile given information about it.
func TileCreator(hash VariantHash, length int, annotation string, reference int64) TileVariant {
	return TileVariant{Hash: hash, Length: length, Annotation: annotation, LookupReference: reference}
}

// VariantHash is a hash for the bases of a tile variant.
// Currently, the hash algorithm is MD5.
type VariantHash [md5.Size]byte

// Equals checks for equality of variants based on hash.
// This works based on the assumption that no two tile variants in the same path and step have the same hash.
func (t TileVariant) Equals(t2 TileVariant) bool {
	return (t.Hash == t2.Hash)
}

// Paths is the number of paths a genome has.
const Paths int = 863 // Constant because we know that the number of paths is always the same.
