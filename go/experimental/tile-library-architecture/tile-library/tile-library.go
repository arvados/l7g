/*Package tilelibrary is a package for implementing tile libraries in Go.
It is assumed that the tile information provided beforehand is imputed, or that having nocalls is okay in SGLF files.
The library does check for completeness of tiles, but doesn't modify them to be complete or imputed before writing them to files.
Various functions to merge, liftover, import, export, and modify libraries are provided.
Should be used in conjunction with the structures package.
Note: Do not add tiles to libraries made from SGLFv2 files (since it's not clear how the tile information would be included)
Note: adding tiles to a library at any point will require sorting that library before writing it to a file.*/
package tilelibrary

import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"../structures"
)

// The following are various possible errors that can occur.

// ErrInconsistentHash is an error for when the hash of a set of bases does not match the TileVariant hash.
var ErrInconsistentHash = errors.New("bases and hash do not match")

// ErrCannotAddTile is an error when trying to add tiles to libraries built from SGLFv2 files, since it's not clear how to add tiles properly to them.
var ErrCannotAddTile = errors.New("library was built from sglfv2 files--cannot add new tile")

// ErrInvalidReferenceLibrary is an error when the ReferenceLibrary field of a TileVariant is not a pointer to a Library.
var ErrInvalidReferenceLibrary = errors.New("reference library field is not a library pointer")

// ErrTileContradiction is an error that occurs when a tile that is known to be in the library was not found (usually this is used after adding that tile to the library)
var ErrTileContradiction = errors.New("a tile that was added is not found in the library")

// ErrBadSource is an error for when the origin of a liftover mapping is not a subset of the destination library.
var ErrBadSource = errors.New("source library is not part of the destination library")

// ErrBadLiftover is an error for when a file is not a liftover mapping.
var ErrBadLiftover = errors.New("not a valid mapping file")

// ErrIncorrectSourceText is an error for when the source file/directory doesn't match the function it's used in.
var ErrIncorrectSourceText = errors.New("library doesn't have the right intermediate file(s) as its Text field")

// readAllGZ is a function to open gzipped files and return the corresponding slice of bytes of the data.
// Mostly important for gzipped FastJs, but any gzipped file can be opened too.
func readAllGZ(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	data, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// KnownVariants is a struct to hold the known variants in a specific step.
// KnownVariants will also keep track of the count of each tile--the variant at List[i] has been seen Counts[i] times.
type KnownVariants struct {
	List   [](*structures.TileVariant) // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index after sorting)
	Counts []int                       // Counts of each variant so far
}

// concurrentPath is a type to represent a path, while also being safe for concurrent use.
// It contains a Variants field, where Variants[i] is a pointer to the set of known variants at step i of this path.
type concurrentPath struct {
	Lock     sync.RWMutex     // The read/write lock used for concurrency within a path.
	Variants []*KnownVariants // The list of steps, where each entry contains the known variants at that step.
}

/*
Library is a struct to represent a library of tile variants from a set of genomes.

A Library is separated into paths, in the Paths field, represented as a slice of concurrentPaths. This makes Libraries safe for concurrent use in terms of modification of tiles.

Libraries have IDs for easy discussion and reference. Currently, the IDs are calculated by the MD5 hash algorithm.
It hashes all of the location and tile information (except annotations) in order, by path, then step, then in the order of variants by increasing variant number.
Upon reaching a new path or step, the uint32 form, separated into 4 bytes, of the path/step is added to the hash.
Then, for each variant, infomation is added to the list of bytes to be hashed in the following order:
	-Variant number, in uint32 form, separated into bytes
	-Total count of this variant
	-Tile length, in terms of steps.
	-The hash of the tile variant, as bytes.

Hashing everything including location infomation ensures that two libraries with the same tiles and counts but with some tiles in different locations would not be considered the same libraries.

The Components of a library determine specifically which libraries are allowed to liftover to that library (since without being part of the components, it's impossible to know easily if you can make a liftover mapping)

In terms of usage, create a new library using New, which will set up the Paths of the Library for you, and will set the reference text file and any component libraries.

Notes: files and directories of tile libraries will not be automatically deleted. If files or directories must be deleted, you must add this functionality (e.g. by using os.Remove).
In addition, the ID is not automatically updated when using AddTile, since AssignID is not quick.
The caller must use AssignID on the library to update its ID after adding all tiles.
*/
type Library struct {
	Paths      []concurrentPath // The paths of the library.
	ID         [md5.Size]byte   // The ID of a library.
	text       string           // The path of the text file containing the bases. As a special case, if Text is a directory, it refers to the sglf/sglfv2 files there.
	isDir      bool             // Field to tell if the "reference text" is a directory, for convenience.
	Components [][md5.Size]byte // the IDs of the libraries that this library is composed of (empty if this library was not made from others)
	// This includes libraries that directly made this library and indirectly made this library (e.g. through making libraries that directly made this library)
}

// Helper function used to make byte slices from unsigned int32s.
// Used as a helper to make a hash for each library.
func uint32ToByteSlice(integer uint32) []byte {
	slice := make([]byte, 4)
	binary.LittleEndian.PutUint32(slice, integer)
	return slice
}

// AssignID assigns a library its ID.
// Current method takes the path, step, variant number, variant count, variant hash, and variant length into account when making the ID.
func (l *Library) AssignID() {
	var byteList []byte
	byteList = make([]byte, 0, 9000000000) // Allocation of 9 gigabytes, a little over the size of what would be expected of a library with 10 tiles per step.
	// While this is costly, in most cases this saves time by not needing to reallocate much.
	for path := range l.Paths {
		byteList = append(byteList, uint32ToByteSlice(uint32(path))...)
		l.Paths[path].Lock.RLock()
		for step, variants := range l.Paths[path].Variants {
			byteList = append(byteList, uint32ToByteSlice(uint32(step))...)
			if variants != nil {
				for variantNumber, variant := range (*variants).List {
					byteList = append(byteList, uint32ToByteSlice(uint32(variantNumber))...)
					byteList = append(byteList, uint32ToByteSlice(uint32((*variants).Counts[variantNumber]))...)
					byteList = append(byteList, uint32ToByteSlice(uint32((*variant).Length))...)
					byteList = append(byteList, (*variant).Hash[:]...)
				}
			}
		}
		l.Paths[path].Lock.RUnlock()
	}
	(*l).ID = md5.Sum(byteList)
}

// Equals checks for equality between two libraries. It does not check similarity in text or components, and tiles are checked by hash.
// HashEquals will generally be a faster way of checking equality--this is best used when you need to be completely sure about library equality (or inequality)
func (l *Library) Equals(l2 *Library) bool {
	for path := range l.Paths {
		l.Paths[path].Lock.RLock()
		l2.Paths[path].Lock.RLock()
		if len(l.Paths[path].Variants) != len(l2.Paths[path].Variants) {
			return false
		}
		for step, stepList := range l.Paths[path].Variants {
			if stepList != nil && l2.Paths[path].Variants[step] != nil {
				if len((*stepList).List) != len(l2.Paths[path].Variants[step].List) {
					return false
				}
				for i, variant := range (*stepList).List {
					if !variant.Equals((*l2.Paths[path].Variants[step].List[i])) || l.Paths[path].Variants[step].Counts[i] != l2.Paths[path].Variants[step].Counts[i] {
						return false
					}
				}
			} else if stepList != nil || l2.Paths[path].Variants[step] != nil {
				return false
			}
		}
		l2.Paths[path].Lock.RUnlock()
		l.Paths[path].Lock.RUnlock()
	}
	return true
}

// HashEquals is a simpler way of checking library equality, since two libraries with the same ID are almost certainly equal.
// It's faster to use than Equals, given that the IDs have been calculated already.
func (l *Library) HashEquals(l2 *Library) bool {
	return l.ID == l2.ID
}

// md5Compare is a comparison function between md5sums. True means that hash1 is smaller.
// Used in SortLibrary to create a consistent order of tiles in a step.
func md5Compare(hash1, hash2 [md5.Size]byte) bool {
	for i := 0; i < md5.Size; i++ {
		if hash1[i] != hash2[i] {
			return hash1[i] < hash2[i]
		}
	}
	return false // the two md5sums are equal, so hash1 is not smaller.
}

// SortLibrary is a function to sort the library once all initial genomes are done being added.
// This function should only be used once after initial setup of the library, after all tiles have been added, since it sorts everything.
// The sort function compares tile counts and hashes, so the order in which tiles are added doesn't matter.
func (l *Library) SortLibrary() {
	type sortStruct struct { // Temporary struct that groups together the variant and the count for sorting purposes.
		variant *structures.TileVariant
		count   int
	}
	for i := range (*l).Paths {
		(*l).Paths[i].Lock.Lock()
		for _, steplist := range (*l).Paths[i].Variants {
			if steplist != nil {
				var sortStructList []sortStruct
				sortStructList = make([]sortStruct, len((*steplist).List))
				for i := 0; i < len((*steplist).List); i++ {
					sortStructList[i] = sortStruct{(*steplist).List[i], (*steplist).Counts[i]}
				}
				sort.Slice(sortStructList, func(i, j int) bool {
					if sortStructList[i].count != sortStructList[j].count {
						return sortStructList[i].count > sortStructList[j].count
					}
					return md5Compare(sortStructList[i].variant.Hash, sortStructList[j].variant.Hash)
				})
				for j := 0; j < len((*steplist).List); j++ {
					(*steplist).List[j], (*steplist).Counts[j] = sortStructList[j].variant, sortStructList[j].count
				}
			}
		}
		(*l).Paths[i].Lock.Unlock()
	}
}

// TileExists is a function to check if a specific tile exists at a specific path and step in a library.
// Returns the index of the variant and the boolean true, if found--otherwise, returns 0 and false, meaning not found.
// It creates more room for new steps and variants, if necessary.
func (l *Library) TileExists(path, step int, toCheck *structures.TileVariant) (int, bool) {
	(*l).Paths[path].Lock.Lock()
	defer (*l).Paths[path].Lock.Unlock()
	for len((*l).Paths[path].Variants) <= step+toCheck.Length-1 { // Makes enough room so that there are step+1 elements in Paths[path].Variants
		(*l).Paths[path].Variants = append((*l).Paths[path].Variants, nil)
	}
	if len((*l).Paths[path].Variants) > step && (*l).Paths[path].Variants[step] != nil { // Safety to make sure that the KnownVariants struct has been created
		for i, value := range (*l).Paths[path].Variants[step].List {
			if toCheck.Equals(*value) {
				return i, true
			}
		}
		return 0, false
	}
	newKnownVariants := &KnownVariants{make([](*structures.TileVariant), 0, 1), make([]int, 0, 1)}
	(*l).Paths[path].Variants[step] = newKnownVariants
	return 0, false
}

// AddTile is a function to add a tile (without sorting).
// Safe to use without checking existence of the tile beforehand (since the function will do that for you).
// Will return any error encountered.
// Note: AddTile will write any new tiles to disk in an intermediate file.
func (l *Library) AddTile(genomePath, step int, new *structures.TileVariant, bases string) error {
	index, ok := l.TileExists(genomePath, step, new)
	if !ok { // If the tile doesn't exist, add it and write it to a file.
		if md5.Sum([]byte(bases)) != new.Hash { // Check to make sure the bases and the tile variant hash do not conflict.
			return ErrInconsistentHash
		}
		new.ReferenceLibrary = l
		(*l).Paths[genomePath].Lock.Lock()
		defer (*l).Paths[genomePath].Lock.Unlock()
		(*l).Paths[genomePath].Variants[step].List = append((*l).Paths[genomePath].Variants[step].List, new)
		(*l).Paths[genomePath].Variants[step].Counts = append((*l).Paths[genomePath].Variants[step].Counts, 1)
		// Added new tile--write the hash and bases to a file.
		info, err := os.Stat(l.text)
		if err != nil {
			file, err := os.Create(l.text)
			if err != nil {
				return err
			}
			file.Close()
			info, err = os.Stat(l.text) // updates the file information, since it didn't exist before.
			if err != nil {
				return err
			}
		}
		if l.isDir {
			return ErrCannotAddTile
		}
		new.LookupReference = info.Size()
		file, err := os.OpenFile(l.text, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		bufferedWriter := bufio.NewWriter(file)
		bufferedWriter.WriteString(hex.EncodeToString(new.Hash[:]))
		bufferedWriter.WriteString(",")
		bufferedWriter.WriteString(bases)
		bufferedWriter.WriteString("\n")
		bufferedWriter.Flush()
		return nil
	}
	(*l).Paths[genomePath].Lock.RLock()
	defer (*l).Paths[genomePath].Lock.RUnlock()
	(*l).Paths[genomePath].Variants[step].Counts[index]++ // Adds 1 to the count of the tile (since it's already in the library)
	// Nothing to write to disk here--can just return.
	return nil
}

// AddTileUnsafe is a function to add a tile without sorting.
// Unsafe because it doesn't check if the tile is already in the library, unlike AddTile.
// Be careful to check if the tile already exists before using this function to avoid repeats in the library.
// This will NOT write anything to disk, as most functions using this will probably write to disk somewhere else or don't need to write to disk at all.
// If you need to write to disk, you can use AddTile or manually add to disk.
func (l *Library) addTileUnsafe(genomePath, step int, new *structures.TileVariant) {
	(*l).Paths[genomePath].Lock.Lock()
	defer (*l).Paths[genomePath].Lock.Unlock()
	(*l).Paths[genomePath].Variants[step].List = append((*l).Paths[genomePath].Variants[step].List, new)
	(*l).Paths[genomePath].Variants[step].Counts = append((*l).Paths[genomePath].Variants[step].Counts, 1)
}

// FindFrequency is a function to find the frequency of a specific tile at a specific path and step.
// A tile that is not found at a specific location has a frequency of 0.
func (l *Library) FindFrequency(path, step int, toFind *structures.TileVariant) int {
	if index, ok := l.TileExists(path, step, toFind); ok {
		return (*l).Paths[path].Variants[step].Counts[index]
	}
	return 0
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified through the returned boolean.
func (l *Library) Annotate(path, step int, hash structures.VariantHash, annotation string) bool {
	(*l).Paths[path].Lock.Lock()
	defer (*l).Paths[path].Lock.Unlock()
	for _, tile := range (*l).Paths[path].Variants[step].List {
		if hash == (*tile).Hash {
			tile.Annotation = annotation
			return true
		}
	}
	return false
}

// baseInfo is a temporary struct to pass around information about a tile's bases.
type baseInfo struct {
	bases   string
	hash    structures.VariantHash
	variant *(structures.TileVariant)
}

// A Builder for readFastJ to use repeatedly without needing to create more Builders.
var tileBuilder strings.Builder

// readFastJ reads a FastJ file and adds its tiles to the provided library.
// Allows for gzipped FastJ files and regular FastJ files.
// Will return any error encountered.
func (l *Library) readFastJ(fastJFilepath string) error {
	var wg sync.WaitGroup
	var baseChannel chan baseInfo
	baseChannel = make(chan baseInfo, 16) // Put information about bases of tiles in here while they need to be processed.
	writeChannel := make(chan bool)
	wg.Add(2)
	go writeTileBases(l.text, baseChannel, writeChannel, &wg)
	file := path.Base(fastJFilepath)      // The name of the file.
	splitpath := strings.Split(file, ".") // This is used to make sure the file is in the right format.
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		return errors.New("invalid hex file name") // Makes sure the file title is four digits of hexadecimal
	}
	pathNumber := 256*int(pathHex[0]) + int(pathHex[1]) // Conversion from hex into decimal--this is the path
	var data []byte
	var err error
	if strings.HasSuffix(file, ".gz") {
		data, err = readAllGZ(fastJFilepath)
	} else {
		data, err = ioutil.ReadFile(fastJFilepath)
	}
	if err != nil {
		return err
	}
	text := string(data)
	tiles := strings.Split(text, "\n\n") // The divider between two tiles is two newlines.

	for _, line := range tiles {
		if strings.HasPrefix(line, ">") { // Makes sure that a "line" starts with the correct character ">"
			stepInHex := line[20:24] // These are the indices where the step is located.
			stepBytes, err := hex.DecodeString(stepInHex)
			if err != nil {
				return err
			}
			step := 256*int(stepBytes[0]) + int(stepBytes[1])
			hashString := line[40:72] // These are the indices where the hash is located.
			hash, err := hex.DecodeString(hashString)
			if err != nil {
				return err
			}
			var hashArray structures.VariantHash
			copy(hashArray[:], hash)
			var lengthString string
			commaCounter := 0
			for i, character := range line {
				if character == ',' {
					commaCounter++
				}
				if commaCounter == 6 { // This is dependent on the location of the length field. In FastJs, the tile length is before the 6th comma.
					k := 1                 // This accounts for the possibility of the length of a tile spanning at least 16 tiles.
					for line[i-k] != ':' { // Goes back a few characters until it knows the string of the tile length.
						k++
					}
					lengthString = line[(i - k + 1):i]
					break
				}
			}

			length, err := strconv.Atoi(lengthString) // Length is provided here in base 10
			if err != nil {
				return err
			}
			baseData := strings.Split(line, "\n")[1:] // Data after the first part of the line are the bases of the tile variant.
			tileBuilder.Reset()
			if pathNumber == 811 { // Grows the buffer of tileBuilder to 2^22 bytes if the path is 811, which seems to contain very large (roughly 2.7 million bases per tile) tiles.
				// Not sure of the exact reasons why this is needed, but SGLFv2 construction will fail on paths 811 and 812 without these lines.
				tileBuilder.Grow(4194304)
			}
			for i := range baseData {
				tileBuilder.WriteString(baseData[i])
			}
			bases := tileBuilder.String()
			newTile := &structures.TileVariant{Hash: hashArray, Length: length, Annotation: "", LookupReference: -1, Complete: isComplete(bases), ReferenceLibrary: l}
			if tileIndex, ok := l.TileExists(pathNumber, step, newTile); !ok {
				l.addTileUnsafe(pathNumber, step, newTile)
				baseChannel <- baseInfo{bases, hashArray, newTile}
			} else {
				(*l).Paths[pathNumber].Variants[step].Counts[tileIndex]++ // Increments the count of the tile variant if it is found.
			}
		}
	}
	wg.Done()
	var nilHash [md5.Size]byte
	baseChannel <- baseInfo{"", nilHash, nil} // Sends a "nil" baseInfo to signal that there are no more tiles.
	<-writeChannel                            // Make sure the writer is done.
	close(baseChannel)
	splitpath, data, tiles = nil, nil, nil // Clears most things in memory that were used here, to free up memory.
	wg.Wait()
	return nil
}

// writeTileBases writes bases and hashes of tiles to the given text file.
// To be used in conjunction with readFastJ.
// Will return any error encountered.
func writeTileBases(libraryTextFile string, channel chan baseInfo, writeChannel chan bool, group *sync.WaitGroup) error {
	err := os.MkdirAll(path.Dir(libraryTextFile), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(libraryTextFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriter(file)
	for bases := range channel {
		if bases.variant != nil { // Checks if there are any more tiles
			currentPosition, err := file.Seek(0, 2)
			if err != nil {
				return err
			}

			(bases.variant).LookupReference = currentPosition + int64(bufferedWriter.Buffered())
			hashString := hex.EncodeToString(bases.hash[:])
			bufferedWriter.WriteString(hashString)
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(bases.bases)
			bufferedWriter.WriteString("\n")
		} else {
			writeChannel <- true // No more tiles, so we are done writing.
		}
	}
	bufferedWriter.Flush()

	file.Close()
	group.Done()
	return nil
}

// writePathToSGLF writes an SGLF for an entire path given a library to a specific directory.
// This assumes that the library has been sorted beforehand.
// Will return any error encountered.
func (l *Library) writePathToSGLF(genomePath int, directoryToWriteTo string) error {
	pathFileMap := make(map[string]*os.File, 0)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglf")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		return err
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo, b.String()), os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriterSize(sglfFile, 4194304) // 4MB buffer

	var fileReader *bufio.Reader
	(*l).Paths[genomePath].Lock.RLock()
	defer (*l).Paths[genomePath].Lock.RUnlock()
	for step := range (*l).Paths[genomePath].Variants {
		if (*l).Paths[genomePath].Variants[step] != nil {
			for i := range (*(*l).Paths[genomePath].Variants[step]).List {
				referenceLibrary, ok := (*(*l).Paths[genomePath].Variants[step]).List[i].ReferenceLibrary.(*Library)
				if !ok {
					return ErrInvalidReferenceLibrary
				}

				file, fileOk := pathFileMap[referenceLibrary.text]
				if !fileOk {
					var fileToOpen string
					if l.isDir {
						fileToOpen = path.Join(referenceLibrary.text, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
					} else {
						fileToOpen = path.Join(referenceLibrary.text)
					}
					textFile, err := os.OpenFile(fileToOpen, os.O_RDONLY, 0644)
					if err != nil {
						return err
					}
					defer textFile.Close()
					pathFileMap[referenceLibrary.text] = textFile
					file = textFile
					if fileReader == nil {
						fileReader = bufio.NewReader(textFile)
					}
				}
				_, err = file.Seek((*(*l).Paths[genomePath].Variants[step]).List[i].LookupReference, 0)
				if err != nil {
					return err
				}
				fileReader.Reset(file)
				tileString, err := fileReader.ReadString('\n') // This includes the newline at the end.
				if err != nil && err != io.EOF {
					return err
				}
				stepHex := fmt.Sprintf("%04x", step)

				bufferedWriter.WriteString(fmt.Sprintf("%04x", genomePath)) // path
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%02v", 0)) // Version
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(stepHex) // Step
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%03x", i)) // Tile variant number
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(strconv.FormatInt(int64((*(*l).Paths[genomePath].Variants[step]).List[i].Length), 16)) // Tile length
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(tileString) // Hash and bases of tile.
				// Newline is at the end of tileString, so no newline needs to be put here.
				// This also means that every SGLF file ends with a newline.
			}
		}
	}
	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}
	err = sglfFile.Close()
	if err != nil {
		return err
	}

	for _, file := range pathFileMap {
		file.Close()
	}
	pathFileMap = nil
	return nil
}

// WriteLibraryToSGLF writes the contents of a library to SGLF files to a specified directory.
// Will return any error encountered.
func (l *Library) WriteLibraryToSGLF(directoryToWriteTo string) error {
	var emptyID [16]byte
	if l.ID == emptyID { // Ensures that the library will have an ID before writing everything out.
		// In the rare case the library's ID is the emptyID, it's still good to double-check.
		l.AssignID()
	}
	for path := 0; path < structures.Paths; path++ {
		err := l.writePathToSGLF(path, directoryToWriteTo)
		if err != nil {
			return err
		}
	}
	return nil
}

// isComplete determines if a set of bases is complete (has no nocalls).
func isComplete(bases string) bool {
	return !strings.ContainsRune(bases, 'n')
}

// writePathToSGLFv2 writes an SGLFv2 for an entire path given a library to a specific directory.
// This assumes that the library has been sorted beforehand.
// Will return any error encountered.
func (l *Library) writePathToSGLFv2(genomePath int, directoryToWriteTo string) error {
	pathFileMap := make(map[string]*os.File, 0)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglfv2")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		return err
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo, b.String()), os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriterSize(sglfFile, 4194304) // 4MB buffer

	var fileReader *bufio.Reader
	bufferedWriter.WriteString("ID:")
	bufferedWriter.WriteString(hex.EncodeToString(l.ID[:]))
	bufferedWriter.WriteString(";Components:")
	if len(l.Components) > 0 {
		bufferedWriter.WriteString(hex.EncodeToString(l.Components[0][:]))
		for _, id := range l.Components[1:] {
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(hex.EncodeToString(id[:]))
		}
	}
	bufferedWriter.WriteString("\n")

	(*l).Paths[genomePath].Lock.RLock()
	defer (*l).Paths[genomePath].Lock.RUnlock()
	for step := range (*l).Paths[genomePath].Variants {
		if (*l).Paths[genomePath].Variants[step] != nil {
			for i := range (*(*l).Paths[genomePath].Variants[step]).List {
				referenceLibrary, ok := (*(*l).Paths[genomePath].Variants[step]).List[i].ReferenceLibrary.(*Library)
				if !ok {
					return ErrInvalidReferenceLibrary
				}

				file, fileOk := pathFileMap[referenceLibrary.text]
				if !fileOk {
					var fileToOpen string
					if l.isDir {
						fileToOpen = path.Join(referenceLibrary.text, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
					} else {
						fileToOpen = path.Join(referenceLibrary.text)
					}
					textFile, err := os.OpenFile(fileToOpen, os.O_RDONLY, 0644)
					if err != nil {
						return err
					}
					defer textFile.Close()
					pathFileMap[referenceLibrary.text] = textFile
					file = textFile
					if fileReader == nil {
						fileReader = bufio.NewReader(textFile)
					}
				}
				_, err = file.Seek((*(*l).Paths[genomePath].Variants[step]).List[i].LookupReference, 0)
				if err != nil {
					return err
				}
				fileReader.Reset(file)
				tileString, err := fileReader.ReadString('\n') // This includes the newline at the end.
				if err != nil && err != io.EOF {
					return err
				}
				stepHex := fmt.Sprintf("%04x", step)

				bufferedWriter.WriteString(fmt.Sprintf("%04x", genomePath)) // path
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%02v", 1)) // Version
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(stepHex) // Step
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%03x", i)) // Tile variant number
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(fmt.Sprintf("%08x", (*(*l).Paths[genomePath].Variants[step]).Counts[i])) // The count of this tile.
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(strconv.FormatInt(int64((*(*l).Paths[genomePath].Variants[step]).List[i].Length), 16)) // Tile length
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(tileString) // Hash and bases of tile.
				// Newline is at the end of tileString, so no newline needs to be put here.
				// This also means that every SGLFv2 file ends with a newline.
			}
		}
	}
	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}
	err = sglfFile.Close()
	if err != nil {
		return err
	}

	for _, file := range pathFileMap {
		file.Close()
	}
	pathFileMap = nil
	return nil
}

// WriteLibraryToSGLFv2 writes the contents of a library to SGLFv2 files to a specified directory.
// Will return any error encountered.
func (l *Library) WriteLibraryToSGLFv2(directoryToWriteTo string) error {
	var emptyID [16]byte
	if l.ID == emptyID { // Ensures that the library will have an ID before writing everything out.
		// In the rare case the library's ID is the emptyID, it's still good to double-check.
		l.AssignID()
	}
	for path := 0; path < structures.Paths; path++ {
		err := l.writePathToSGLFv2(path, directoryToWriteTo)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddLibraryFastJ adds a directory of gzipped FastJ files to a specific library.
// Will return any error encountered.
func (l *Library) AddLibraryFastJ(directory string) error {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") || strings.HasSuffix(file.Name(), ".fj") { // Checks the directory for the right file types.
			err = l.readFastJ(path.Join(directory, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddPathFromDirectories parses the same path for all genomes, represented by a list of directories, and puts the information in a Library.
// Will return any error encountered.
func (l *Library) AddPathFromDirectories(directories []string, genomePath int, gzipped bool) error {
	var filename string
	if gzipped {
		filename = fmt.Sprintf("%04x.fj.gz", genomePath)
	} else {
		filename = fmt.Sprintf("%04x.fj", genomePath)
	}
	for _, directory := range directories {
		err := l.readFastJ(path.Join(directory, filename))
		if err != nil {
			return err
		}
	}
	return nil
}

// AddDirectories adds information from a list of directories for genomes into a library, but parses by path.
// Will return any error encountered.
func (l *Library) AddDirectories(directories []string, gzipped bool) error {
	for path := 0; path < structures.Paths; path++ {
		err := l.AddPathFromDirectories(directories, path, gzipped)
		if err != nil {
			return err
		}
	}
	return nil
}

// CompileDirectoriesToLibrary creates a new Library based on the directories given, sorts it, and gives it its ID (so this library is ready for use).
// Returns the library pointer and an error, if any (nil if no error was encounted)
func CompileDirectoriesToLibrary(directories []string, libraryTextFile string, gzipped bool) (*Library, error) {
	l, err := New(libraryTextFile, nil)
	if err != nil {
		return nil, err
	}
	err = l.AddDirectories(directories, gzipped)
	if err != nil {
		return nil, err
	}
	l.SortLibrary()
	(l).AssignID()
	return l, nil
}

// SequentialCompileDirectoriesToLibrary creates a new Library based on the directories given, sorts it, and gives it its ID (so this library is ready for use).
// This adds each directory sequentially (so each genome is done one at a time, rather than doing one path of all genomes all at once)
// Returns the library pointer and an error, if any (nil if no error was encounted)
func SequentialCompileDirectoriesToLibrary(directories []string, libraryTextFile string) (*Library, error) {
	l, err := New(libraryTextFile, nil)
	if err != nil {
		return nil, err
	}
	for _, directory := range directories {
		err := l.AddLibraryFastJ(directory)
		if err != nil {
			return nil, err
		}
	}
	l.SortLibrary()
	(l).AssignID()
	return l, nil
}

// New sets up the basic structure for a library and returns a pointer to a new library.
// For consistency, it's best to use an absolute path for the text file. Relative paths will still work, but they are not recommended.
func New(textFile string, componentLibraries [][md5.Size]byte) (*Library, error) {
	var newLibraryPaths []concurrentPath
	newLibraryPaths = make([]concurrentPath, structures.Paths, structures.Paths)
	for i := range newLibraryPaths {
		var newLock sync.RWMutex
		newLibraryPaths[i] = concurrentPath{newLock, make([]*KnownVariants, 0, 1)} // Lock is copied, but hasn't been used yet, so this is fine.
	}
	info, err := os.Stat(textFile)
	if err != nil {
		file, err := os.OpenFile(textFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		info, err = os.Stat(textFile)
		if err != nil {
			return nil, err
		}
	}
	return &Library{Paths: newLibraryPaths, text: textFile, isDir: info.IsDir(), Components: componentLibraries}, nil
}

// Function to copy the contents of the source library into the new library.
// Note that the normal copy function can't be used on the entire library, since locks should not be copied after use.
func libraryCopy(destination, source *Library) {
	destination.ID = source.ID
	for i := range source.Paths {
		source.Paths[i].Lock.RLock()                                                          // Locked since we're reading from path.Variants when we copy.
		destination.Paths[i].Variants = make([]*KnownVariants, len(source.Paths[i].Variants)) // This is to make sure all elements are copied over.
		for step := range destination.Paths[i].Variants {
			if source.Paths[i].Variants[step] != nil {
				length := len(source.Paths[i].Variants[step].List)
				destination.Paths[i].Variants[step] = &KnownVariants{List: make([]*structures.TileVariant, length), Counts: make([]int, length)}
				copy(destination.Paths[i].Variants[step].List, source.Paths[i].Variants[step].List)
				copy(destination.Paths[i].Variants[step].Counts, source.Paths[i].Variants[step].Counts)
				for j, variant := range source.Paths[i].Variants[step].List {
					variantCopy := *variant // copies variants without reusing, so variants in source still have a pointer to source while variants in destination point to destination.
					//variantCopy.ReferenceLibrary = destination
					destination.Paths[i].Variants[step].List[j] = &variantCopy
				}
				// copying List and Counts manually instead of copying Variants[step] makes sure the pointers for List and Counts aren't reused.
			}
		}
		source.Paths[i].Lock.RUnlock()
	}
}

// MergeLibraries is a function to merge the library given with the base library.
// This version creates a new library.
// Returns the library pointer and an error, if any (nil if no error was encounted)
func (l *Library) MergeLibraries(libraryToMerge *Library, textFile string) (*Library, error) {
	var emptyID [16]byte
	if l.ID == emptyID { // If the ID for either library has not been assigned, assign the IDs now.
		l.AssignID()
	}
	if libraryToMerge.ID == emptyID {
		libraryToMerge.AssignID()
	}
	allComponents := append([][md5.Size]byte{libraryToMerge.ID, l.ID}, libraryToMerge.Components...)
	allComponents = append(allComponents, l.Components...)
	newLibrary, err := New(textFile, allComponents)
	if err != nil {
		return nil, err
	}
	libraryCopy(newLibrary, l)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			err := newLibrary.mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j])
			if err != nil {
				return nil, err
			}
		}
	}
	newLibrary.SortLibrary()
	newLibrary.AssignID()
	return newLibrary, nil
}

// MergeKnownVariants puts the contents of a KnownVariants at a specific path and step into another library.
// Will return any encountered, if any.
func (l *Library) mergeKnownVariants(genomePath, step int, variantsToMerge *KnownVariants) error {
	if variantsToMerge != nil {
		for i, variant := range variantsToMerge.List {
			if index, ok := l.TileExists(genomePath, step, variant); !ok {
				l.addTileUnsafe(genomePath, step, variant)
				newIndex, ok := l.TileExists(genomePath, step, variant)
				if ok {
					(*l).Paths[genomePath].Variants[step].Counts[newIndex] += variantsToMerge.Counts[i] - 1
				} else {
					return ErrTileContradiction
				}
			} else {
				(*l).Paths[genomePath].Variants[step].Counts[index] += variantsToMerge.Counts[i]
			}
		}
	}
	return nil
}

// MergeLibrariesWithoutCreation merges libraries without creating a new one, using the "mainLibrary" instead.
// Returns the library pointer and an error, if any (nil if no error was encounted)
func (l *Library) MergeLibrariesWithoutCreation(libraryToMerge *Library) (*Library, error) {
	var emptyID [16]byte
	if l.ID == emptyID { // If the ID for either library has not been assigned, assign the IDs now.
		l.AssignID()
	}
	if libraryToMerge.ID == emptyID {
		libraryToMerge.AssignID()
	}
	l.Components = append(l.Components, libraryToMerge.ID, l.ID) // This is okay since this involves the mainLibrary's old ID.
	l.Components = append(l.Components, libraryToMerge.Components...)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			err := l.mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j])
			if err != nil {
				return nil, err
			}
		}
	}
	l.SortLibrary()
	l.AssignID()
	return l, nil
}

// LiftoverMapping is a representation of a liftover from one library to another, essentially becoming a translation of variants from the source to the destination.
// If a = LiftoverMapping.Mapping[b][c][d], then in path b, step c, variant d in the first library maps to variant a in path b and step c in the second.
type LiftoverMapping struct {
	Mapping            [][][]int // The actual mapping between the two libraries
	SourceLibrary      *Library  // The source library to map from.
	DestinationLibrary *Library  // The destination library to map to.
}

// CreateMapping creates a liftover mapping from the source library to the destination library.
// Returns the mapping and an error, if any (nil if no error was encounted)
func CreateMapping(source, destination *Library) (LiftoverMapping, error) {
	index := -1
	for i, libraryID := range destination.Components {
		if libraryID == source.ID {
			index = i
			break
		}
	}
	if index == -1 { // Destination was not made from the source--can't guarantee a mapping here.
		return LiftoverMapping{nil, nil, nil}, ErrBadSource
	}
	var mapping [][][]int
	mapping = make([][][]int, structures.Paths, structures.Paths)
	for path := range (*source).Paths {
		(*source).Paths[path].Lock.RLock()
		mapping[path] = make([][]int, len((*source).Paths[path].Variants)) // Number of steps.
		for step, variants := range (*source).Paths[path].Variants {
			if variants != nil {
				for _, variant := range (*variants).List {
					index, ok := destination.TileExists(path, step, variant)
					if ok {
						mapping[path][step] = append(mapping[path][step], index)
					} else {
						return LiftoverMapping{nil, nil, nil}, ErrTileContradiction
					}
				}
			}
		}
		(*source).Paths[path].Lock.RUnlock()
	}
	return LiftoverMapping{mapping, source, destination}, nil
}

// WriteMapping writes a LiftoverMapping to a specified file.
// The format is path/step/source1,destination1;source2,destination2;...
// Current suffix for mappings: .sglfmapping (make sure all filenames end with this suffix)
// Returns any error encountered, or nil if there's no error.
func WriteMapping(filename string, mapping LiftoverMapping) error {
	textFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriter(textFile)
	bufferedWriter.WriteString(hex.EncodeToString(mapping.SourceLibrary.ID[:]))
	bufferedWriter.WriteString(",")
	bufferedWriter.WriteString(hex.EncodeToString(mapping.DestinationLibrary.ID[:]))
	bufferedWriter.WriteString("\n")
	for path := range mapping.Mapping {
		for step := range mapping.Mapping[path] {
			if mapping.Mapping[path][step] != nil {
				bufferedWriter.WriteString(fmt.Sprintf("%04x", path))
				bufferedWriter.WriteString("/")
				bufferedWriter.WriteString(fmt.Sprintf("%04x", step))
				bufferedWriter.WriteString("/")
				bufferedWriter.WriteString("0,")
				bufferedWriter.WriteString(strconv.Itoa(mapping.Mapping[path][step][0]))
				for index, value := range mapping.Mapping[path][step][1:] {
					bufferedWriter.WriteString(";")
					bufferedWriter.WriteString(strconv.Itoa(index))
					bufferedWriter.WriteString(",")
					bufferedWriter.WriteString(strconv.Itoa(value))
				}
				bufferedWriter.WriteString("\n")
			}
		}
	}
	bufferedWriter.Flush()
	textFile.Close()
	return nil
}

// ReadMapping gets the information from a mapping given its filepath.
// It also returns the hashes for the source and destination libraries, in that order.
func ReadMapping(filepath string) (mapping [][][]int, sourceID, destinationID [md5.Size]byte, err error) {
	info, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, [16]byte{}, [16]byte{}, err
	}
	lines := strings.Split(string(info), "\n")
	libraryIDs := strings.Split(lines[0], ",")
	sourceLibraryID, err := hex.DecodeString(libraryIDs[0])
	if err != nil {
		return nil, [16]byte{}, [16]byte{}, err
	}
	var sourceHash [md5.Size]byte
	copy(sourceHash[:], sourceLibraryID)
	destinationLibraryID, err := hex.DecodeString(libraryIDs[1])
	if err != nil {
		return nil, [16]byte{}, [16]byte{}, err
	}
	var destinationHash [md5.Size]byte
	copy(destinationHash[:], destinationLibraryID)
	newMapping := make([][][]int, structures.Paths, structures.Paths)
	for _, stepMapping := range lines[1:] {
		if stepMapping != "" {
			stepInfo := strings.Split(stepMapping, "/")
			path, err := strconv.ParseInt(stepInfo[0], 16, 0)
			if err != nil {
				return nil, [16]byte{}, [16]byte{}, err
			}
			step, err := strconv.ParseInt(stepInfo[1], 16, 0)
			if err != nil {
				return nil, [16]byte{}, [16]byte{}, err
			}
			for len(newMapping[path]) <= int(step) {
				newMapping[path] = append(newMapping[path], nil)
			}
			newMapping[path][step] = make([]int, 0, 1)
			numberMappingInfo := strings.Split(stepInfo[2], ";")
			for _, pair := range numberMappingInfo {
				if pair != "" {
					indexAndValue := strings.Split(pair, ",")
					value, err := strconv.Atoi(indexAndValue[1])
					if err != nil {
						return nil, [16]byte{}, [16]byte{}, err
					}
					newMapping[path][step] = append(newMapping[path][step], value)
				}
			}
		}
	}
	return newMapping, sourceHash, destinationHash, nil
}

// addSGLFv2 is a function to put SGLFv2 data back into a library.
// Allows for gzipped SGLFv2 files and regular SGLFv2 files.
// The following bug is believed to be fixed, but just in case it is not and an error occurs here:
// Note: the creation of an SGLFv2 may be incorrect and may put a series of bases before the hash, when the input tile from the FastJ file
// was extremely long (e.g. a couple million bases). This will result in an error here when decoding the hash. To fix this, you can instruct
// the tileBuilder in readFastJ to Grow before constructing the tiles in that path. So far, this error is known to happen for path 811,
// which will create 032b.sglfv2 and 032c.sglfv2 incorrectly.
// It's also possible that the cause is that the two goroutines for reading and writing end too early, leaving some tiles with
// LookupReferences of -1, which results in the incorrect results. This could happen if readFastJ ended first.
// Returns any error encountered, or nil if there's no error.
func (l *Library) addSGLFv2(filepath string) error {
	file := path.Base(filepath)
	splitpath := strings.Split(file, ".")
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		return errors.New("invalid hex file name") // Makes sure the title of the file is four digits of hexadecimal
	}
	pathNumber := 256*int(pathHex[0]) + int(pathHex[1]) // conversion into an integer--this is the path number
	var data []byte
	var err error
	if strings.HasSuffix(file, ".gz") {
		data, err = readAllGZ(filepath)
	} else {
		data, err = ioutil.ReadFile(filepath)
	}
	if err != nil {
		return err
	}
	text := string(data)
	tiles := strings.Split(text, "\n")
	referenceCounter := len(tiles[0]) + 1 // length of the first line plus the newline
	for lineNumber, line := range tiles {
		if line != "" && lineNumber > 0 {
			fields := strings.Split(line, ",")
			hashString := fields[1]
			lineInfo := strings.Split(fields[0], ".")
			stepString := lineInfo[2]
			tileInfo := strings.Split(lineInfo[3], "+")
			tileCountString := tileInfo[1]
			tileLengthString := tileInfo[2]
			step, err := strconv.ParseInt(stepString, 16, 0)
			if err != nil {
				return err
			}
			count, err := strconv.ParseInt(tileCountString, 16, 0)
			if err != nil {
				return err
			}
			length, err := strconv.ParseInt(tileLengthString, 16, 0)
			if err != nil {
				return err
			}
			hash, err := hex.DecodeString(hashString)
			if err != nil {
				return err
			}
			var hashArray [16]byte
			copy(hashArray[:], hash)
			newVariant := structures.TileVariant{Hash: hashArray, Length: int(length), Annotation: "", LookupReference: 27 + int64(len(tileLengthString)) + int64(referenceCounter), Complete: isComplete(fields[2]), ReferenceLibrary: l}
			index, ok := l.TileExists(pathNumber, int(step), &newVariant)
			if !ok {
				l.addTileUnsafe(pathNumber, int(step), &newVariant)
			}
			index, ok = l.TileExists(pathNumber, int(step), &newVariant)
			if ok {
				(*l).Paths[pathNumber].Variants[int(step)].Counts[index] += int(count - 1)
			} else {
				return ErrTileContradiction
			}
			// Adding count-1 instead of count since addTileUnsafe will already add 1 to the count of the new tile.
			referenceCounter += len(line) + 1 // +1 to account for the newline.
		} else if line != "" && len(l.Components) == 0 { // This refers to the first line, which contains ID and Component information.
			idSlice := strings.Split(line, ";")
			idString := strings.Split(idSlice[0], ":")[1]
			libraryHash, err := hex.DecodeString(idString)
			if err != nil {
				return err
			}
			var hashArray [md5.Size]byte
			copy(hashArray[:], libraryHash)
			l.ID = hashArray
			components := strings.Split(idSlice[1], ":")[1]
			if components != "" {
				componentStrings := strings.Split(components, ",")
				for _, component := range componentStrings {
					componentHash, err := hex.DecodeString(component)
					if err != nil {
						return err
					}
					var componentHashArray [md5.Size]byte
					copy(componentHashArray[:], componentHash)
					l.Components = append(l.Components, componentHashArray)
				}
			}
		}
	}
	return nil
}

// AddLibrarySGLFv2 adds a directory of SGLFv2 files to a library.
// Library should be initialized with this directory as the Text field, so that text files of bases and directories aren't mixed together.
// Returns any error encountered, or nil if there's no error.
func (l *Library) AddLibrarySGLFv2() error {

	if !l.isDir {
		return ErrIncorrectSourceText
	}

	sglfv2Files, err := ioutil.ReadDir(l.text)
	if err != nil {
		return err
	}
	for _, file := range sglfv2Files {
		if strings.HasSuffix(file.Name(), ".sglfv2") || strings.HasSuffix(file.Name(), ".gz") { // Checks if a file is an sglfv2 file.
			err = l.addSGLFv2(path.Join(l.text, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
