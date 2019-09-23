/*Package structures is a basic package to hold basic structures, methods, and functions for tile libraries and genomes.
Most important is the TileVariant structure, which holds all the necessary data to represent a tile variant.
Equality on TileVariants is defined by hash.
In addition, the default number of paths in a genome is provided here for convenience.*/
package structures

import (
	"crypto/md5"
)

// TileVariant is a struct for representing a tile variant using the hash, length, and any annotation(s) on the variant.
// In this case, a tile variant is a set of a bases that could exist in between two tags on the genome, spanning any number of steps.
// We use the hash of the bases as a standin for the actual bases, to save space. We also keep the span of the tile, since it tells us how many steps are in the path of the genome.
// The Annotation field holds a string for any notes about this tile.
// The LookupReference field tells at what location in the reference text file that this tile can be found. This is majorly for the use of the tilelibrary package.
// Complete is a field that determines whether the tile variant contains nocalls or not. Used for convenience in both the genome and tilelibrary packages.
// ReferenceLibrary is a field that holds a *Library, used in package tilelibrary. It is an interface{} to avoid circular imports.
// Putting any other type into ReferenceLibrary may result in errors.
// Note: no explicit TileVariant constructor is provided. However, all fields are exported, so other packages can construct and modify TileVariants freely.
type TileVariant struct {
	Hash             VariantHash // The hash of the tile variant's bases.
	Length           int         // The length (span) of the tile
	Annotation       string      // Any notes about this tile (by default, no comments)
	LookupReference  int64       // The lookup reference for the bases of this variant in the text file for the corresponding library.
	Complete         bool        // Tells if this variant is complete or not (complete meaning no nocalls). Mostly for genomes, to quickly tell which tiles are complete and which are not.
	ReferenceLibrary interface{} // A way of referencing the library this variant is from. (Will be a *Library).
}

// VariantHash is a type hash for the bases of a tile variant.
// Currently, the hash algorithm is MD5.
type VariantHash [md5.Size]byte

// Equals checks for equality of variants based on hash.
// This works based on the assumption that no two tile variants in the same path and step have the same hash.
func (t TileVariant) Equals(t2 TileVariant) bool {
	return (t.Hash == t2.Hash)
}

// Paths is the number of paths a genome has, for convenience.
const Paths int = 863 // Constant because we know that the number of paths is always the same.
