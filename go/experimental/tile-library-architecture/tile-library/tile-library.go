/*Package tilelibrary is a package for implementing tile libraries in Go.
It is assumed that the tile information provided beforehand is imputed--the library does not check for completeness of tiles before writing them to files.
Various functions to merge, liftover, import, export, and modify libraries are provided.
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

// ErrRemoveIntermediateDirectory is an error returned when an intermediate directory of SGLFv2 files is attempted to be removed.
var ErrRemoveIntermediateDirectory = errors.New("intermediate files are sglfv2 files; nothing removed")

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

// openGZ is a function to open gzipped files and return the corresponding slice of bytes of the data.
// Mostly important for gzipped FastJs, but any gzipped file can be opened too.
func openGZ(filepath string) ([]byte, error) {
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

// openFile is a method to get the data of a file and return the corresponding slice of bytes.
// Available mostly as a way to be flexible with files, since combined with openGZ a gzipped file or a non-gzipped file can be read and used.
func openFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// KnownVariants is a struct to hold the known variants in a specific step.
type KnownVariants struct {
	List   [](*structures.TileVariant) // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index after sorting)
	Counts []int                       // Counts of each variant so far
}

// Library is a type to represent a library of tile variants.
type Library struct {
	Paths []concurrentPath // The paths of the library.
	ID    [md5.Size]byte   // The ID of a library.
	Text  string           // The path of the text file containing the bases. As a special case, if Text is a directory, it refers to the sglf/sglfv2 files there.
	// Note: the Text field is only relevant to the file system this Library is on.
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
	// While this is costly, in most cases this saves time by not needed to reallocate much.
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
// HashEquals will generally be a faster way of checking equality--this is best used when you need to be sure about library equality (or inequality)
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

// concurrentPath is a type to represent a path, while also being safe for concurrent use.
type concurrentPath struct {
	Lock     sync.RWMutex     // The read/write lock used for concurrency within a path.
	Variants []*KnownVariants // The list of steps, where each entry contains the known variants at that step.
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
func SortLibrary(library *Library) {
	type sortStruct struct { // Temporary struct that groups together the variant and the count for sorting purposes.
		variant *structures.TileVariant
		count   int
	}
	for i := range (*library).Paths {
		(*library).Paths[i].Lock.Lock()
		for _, steplist := range (*library).Paths[i].Variants {
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
		(*library).Paths[i].Lock.Unlock()
	}
}

// TileExists is a function to check if a specific tile exists at a specific path and step in a library.
// Returns the index of the variant and the boolean true, if found--otherwise, returns 0 and false, meaning not found.
func TileExists(path, step int, toCheck *structures.TileVariant, library *Library) (int, bool) {
	(*library).Paths[path].Lock.Lock()
	defer (*library).Paths[path].Lock.Unlock()
	for len((*library).Paths[path].Variants) <= step+toCheck.Length-1 { // Makes enough room so that there are step+1 elements in Paths[path].Variants
		(*library).Paths[path].Variants = append((*library).Paths[path].Variants, nil)
	}
	if len((*library).Paths[path].Variants) > step && (*library).Paths[path].Variants[step] != nil { // Safety to make sure that the KnownVariants struct has been created
		for i, value := range (*library).Paths[path].Variants[step].List {
			if toCheck.Equals(*value) {
				return i, true
			}
		}
		return 0, false
	}
	newKnownVariants := &KnownVariants{make([](*structures.TileVariant), 0, 1), make([]int, 0, 1)}
	(*library).Paths[path].Variants[step] = newKnownVariants
	return 0, false
}

// AddTile is a function to add a tile (without sorting).
// Safe to use without checking existence of the tile beforehand (since the function will do that for you).
// Will return any error encountered.
func AddTile(genomePath, step int, new *structures.TileVariant, bases string, library *Library) error {
	index, ok := TileExists(genomePath, step, new, library)
	if !ok { // Checks if the tile exists already.
		if md5.Sum([]byte(bases)) != new.Hash { // Check to make sure the bases and the tile variant hash do not conflict.
			return ErrInconsistentHash
		}
		new.ReferenceLibrary = library
		(*library).Paths[genomePath].Lock.Lock()
		defer (*library).Paths[genomePath].Lock.Unlock()
		(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
		(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)
		// Added new tile--write the hash and bases to a file.

		info, err := os.Stat(library.Text)
		if err != nil {
			file, err := os.Create(library.Text)
			if err != nil {
				return err
			}
			file.Close()
			info, err = os.Stat(library.Text) // updates the file information.
			if err != nil {
				return err
			}
		}
		if info.IsDir() {
			return ErrCannotAddTile
		}
		new.LookupReference = info.Size()
		file, err := os.OpenFile(library.Text, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	(*library).Paths[genomePath].Lock.RLock()
	defer (*library).Paths[genomePath].Lock.RUnlock()
	(*library).Paths[genomePath].Variants[step].Counts[index]++ // Adds 1 to the count of the tile (since it's already in the library)
	// Nothing to write to disk here--can just return.
	return nil
}

// AddTileUnsafe is a function to add a tile without sorting.
// Unsafe because it doesn't check if the tile is already in the library, unlike AddTile.
// Be careful to check if the tile already exists before using this function to avoid repeats in the library.
// This will NOT write anything to disk, as most functions using this will probably write to disk somewhere else or don't need to write to disk at all.
// If you need to write to disk, you can use AddTile or manually add to disk.
func addTileUnsafe(genomePath, step int, new *structures.TileVariant, library *Library) {
	(*library).Paths[genomePath].Lock.Lock()
	defer (*library).Paths[genomePath].Lock.Unlock()
	(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
	(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)
}

// FindFrequency is a function to find the frequency of a specific tile at a specific path and step.
// A tile that is not found at a specific location has a frequency of 0.
func FindFrequency(path, step int, toFind *structures.TileVariant, library *Library) int {
	if index, ok := TileExists(path, step, toFind, library); ok {
		return (*library).Paths[path].Variants[step].Counts[index]
	}
	return 0
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified through the returned boolean.
func Annotate(path, step int, hash structures.VariantHash, annotation string, library *Library) bool {
	(*library).Paths[path].Lock.Lock()
	defer (*library).Paths[path].Lock.Unlock()
	for _, tile := range (*library).Paths[path].Variants[step].List {
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

var tileBuilder strings.Builder

// bufferedTileRead reads a FastJ file and adds its tiles to the provided library.
// Allows for gzipped FastJ files and regular FastJ files.
// Will return any error encountered.
func bufferedTileRead(fastJFilepath string, library *Library) error {
	var wg sync.WaitGroup
	var baseChannel chan baseInfo
	baseChannel = make(chan baseInfo, 16) // Put information about bases of tiles in here while they need to be processed.
	writeChannel := make(chan bool)
	wg.Add(2)
	go bufferedBaseWrite(library.Text, baseChannel, writeChannel, &wg)
	file := path.Base(fastJFilepath)      // The name of the file.
	splitpath := strings.Split(file, ".") // This is used to make sure the file is in the right format.
	if len(splitpath) != 3 && len(splitpath) != 2 {
		return errors.New("error: Not a valid file " + file) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "fj" || (len(splitpath) == 3 && splitpath[2] != "gz") {
		return errors.New("error: not a valid FastJ file") // Makes sure that the file is a FastJ file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		return errors.New("invalid hex file name") // Makes sure the file title is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0]) + int(pathHex[1]) // Conversion from hex into decimal--this is the path
	var data []byte
	var err error
	if len(splitpath) == 3 {
		data, err = openGZ(fastJFilepath)
	} else {
		data, err = openFile(fastJFilepath)
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
				if commaCounter == 6 { // This is dependent on the location of the length field.
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
			if hexNumber == 811 { // Grows the buffer of tileBuilder to 2^22 bytes if the path is 811, which seems to contain very large (roughly 2.7 million bases per tile) tiles.
				// Not sure of the exact reasons why this is needed, but SGLFv2 construction will fail on paths 811 and 812 without these lines.
				tileBuilder.Grow(4194304)
			}
			for i := range baseData {
				tileBuilder.WriteString(baseData[i])
			}
			bases := tileBuilder.String()
			newTile := &structures.TileVariant{Hash: hashArray, Length: length, Annotation: "", LookupReference: -1, Complete: isComplete(bases), ReferenceLibrary: library}
			if tileIndex, ok := TileExists(hexNumber, step, newTile, library); !ok {
				addTileUnsafe(hexNumber, step, newTile, library)
				baseChannel <- baseInfo{bases, hashArray, newTile}
			} else {
				(*library).Paths[hexNumber].Variants[step].Counts[tileIndex]++ // Increments the count of the tile variant if it is found.
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

// bufferedBaseWrite writes bases and hashes of tiles to the given text file.
// To be used in conjunction with bufferedTileRead.
// Will return any error encountered.
func bufferedBaseWrite(libraryTextFile string, channel chan baseInfo, writeChannel chan bool, group *sync.WaitGroup) error {
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

// writePathToSGLF writes an SGLF for an entire path given a library.
// This assumes that the library has been sorted beforehand.
// Will return any error encountered.
func writePathToSGLF(library *Library, genomePath int, directoryToWriteTo, directoryToGetFrom, textFilename string) error {
	pathFileMap := make(map[string]*os.File, 0)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglf")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		return err
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo, b.String()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriterSize(sglfFile, 4194304) // 4MB buffer

	var fileReader *bufio.Reader
	(*library).Paths[genomePath].Lock.RLock()
	defer (*library).Paths[genomePath].Lock.RUnlock()
	for step := range (*library).Paths[genomePath].Variants {
		if (*library).Paths[genomePath].Variants[step] != nil {
			for i := range (*(*library).Paths[genomePath].Variants[step]).List {
				referenceLibrary, ok := (*(*library).Paths[genomePath].Variants[step]).List[i].ReferenceLibrary.(*Library)
				if !ok {
					return ErrInvalidReferenceLibrary
				}

				file, fileOk := pathFileMap[referenceLibrary.Text]
				if !fileOk {
					info, err := os.Stat(referenceLibrary.Text)
					if err != nil {
						return err
					}
					var fileToOpen string
					if info.IsDir() {
						fileToOpen = path.Join(referenceLibrary.Text, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
					} else {
						fileToOpen = path.Join(referenceLibrary.Text)
					}
					textFile, err := os.OpenFile(fileToOpen, os.O_RDONLY, 0644)
					if err != nil {
						return err
					}
					defer textFile.Close()
					pathFileMap[referenceLibrary.Text] = textFile
					file = textFile
					if fileReader == nil {
						fileReader = bufio.NewReader(textFile)
					}
				}
				_, err = file.Seek((*(*library).Paths[genomePath].Variants[step]).List[i].LookupReference, 0)
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
				bufferedWriter.WriteString(strconv.FormatInt(int64((*(*library).Paths[genomePath].Variants[step]).List[i].Length), 16)) // Tile length
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

// WriteLibraryToSGLF writes the contents of a library to SGLF files.
// Will return any error encountered.
func WriteLibraryToSGLF(library *Library, directoryToWriteTo, directoryToGetFrom, textFile string) error {
	for path := 0; path < structures.Paths; path++ {
		err := writePathToSGLF(library, path, directoryToWriteTo, directoryToGetFrom, textFile)
		if err != nil {
			return err
		}
	}
	err := RemoveIntermediateFile(library)
	if err != nil && err.Error() != ErrRemoveIntermediateDirectory.Error() {
		return err
	}
	return nil
}

// isComplete determines if a set of bases is complete (has no nocalls).
func isComplete(bases string) bool {
	return !strings.ContainsRune(bases, 'n')
}

// writePathToSGLFv2 writes an SGLFv2 for an entire path given a library.
// This assumes that the library has been sorted beforehand.
// Will return any error encountered.
func writePathToSGLFv2(library *Library, genomePath int, directoryToWriteTo string) error {
	pathFileMap := make(map[string]*os.File, 0)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglfv2")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		return err
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo, b.String()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriterSize(sglfFile, 4194304) // 4MB buffer

	var fileReader *bufio.Reader
	bufferedWriter.WriteString("ID:")
	bufferedWriter.WriteString(hex.EncodeToString(library.ID[:]))
	bufferedWriter.WriteString(";Components:")
	if len(library.Components) > 0 {
		bufferedWriter.WriteString(hex.EncodeToString(library.Components[0][:]))
		bufferedWriter.WriteString(",")
		bufferedWriter.WriteString(hex.EncodeToString(library.Components[1][:]))
	}
	bufferedWriter.WriteString("\n")

	(*library).Paths[genomePath].Lock.RLock()
	defer (*library).Paths[genomePath].Lock.RUnlock()
	for step := range (*library).Paths[genomePath].Variants {
		if (*library).Paths[genomePath].Variants[step] != nil {
			for i := range (*(*library).Paths[genomePath].Variants[step]).List {
				referenceLibrary, ok := (*(*library).Paths[genomePath].Variants[step]).List[i].ReferenceLibrary.(*Library)
				if !ok {
					return ErrInvalidReferenceLibrary
				}

				file, fileOk := pathFileMap[referenceLibrary.Text]
				if !fileOk {
					info, err := os.Stat(referenceLibrary.Text)
					if err != nil {
						return err
					}
					var fileToOpen string
					if info.IsDir() {
						fileToOpen = path.Join(referenceLibrary.Text, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
					} else {
						fileToOpen = path.Join(referenceLibrary.Text)
					}
					textFile, err := os.OpenFile(fileToOpen, os.O_RDONLY, 0644)
					if err != nil {
						return err
					}
					defer textFile.Close()
					pathFileMap[referenceLibrary.Text] = textFile
					file = textFile
					if fileReader == nil {
						fileReader = bufio.NewReader(textFile)
					}
				}
				_, err = file.Seek((*(*library).Paths[genomePath].Variants[step]).List[i].LookupReference, 0)
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
				bufferedWriter.WriteString(fmt.Sprintf("%08x", (*(*library).Paths[genomePath].Variants[step]).Counts[i])) // The count of this tile.
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(strconv.FormatInt(int64((*(*library).Paths[genomePath].Variants[step]).List[i].Length), 16)) // Tile length
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

// WriteLibraryToSGLFv2 writes the contents of a library to SGLFv2 files.
// Will return any error encountered.
func WriteLibraryToSGLFv2(library *Library, directoryToWriteTo string) error {
	var emptyID [16]byte
	if library.ID == emptyID { // Ensures that the library will have an ID before writing everything out.
		// In the rare case the library's ID is the emptyID, it's still good to double-check.
		library.AssignID()
	}
	for path := 0; path < structures.Paths; path++ {
		err := writePathToSGLFv2(library, path, directoryToWriteTo)
		if err != nil {
			return err
		}
	}
	err := RemoveIntermediateFile(library)
	if err != nil && err.Error() != ErrRemoveIntermediateDirectory.Error() {
		return err
	}
	return nil
}

// AddLibraryFastJ adds a directory of gzipped FastJ files to a specific library.
// Will return any error encountered.
func AddLibraryFastJ(directory string, library *Library) error {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") { // Checks if a file is a gz file.
			err = bufferedTileRead(path.Join(directory, file.Name()), library)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddPathFromDirectories parses the same path for all genomes, represented by a list of directories, and puts the information in a Library.
// Will return any error encountered.
func AddPathFromDirectories(library *Library, directories []string, genomePath int, gzipped bool) error {
	var filename string
	if gzipped {
		filename = fmt.Sprintf("%04x.fj.gz", genomePath)
	} else {
		filename = fmt.Sprintf("%04x.fj", genomePath)
	}
	for _, directory := range directories {
		err := bufferedTileRead(path.Join(directory, filename), library)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddByDirectories adds information from a list of directories for genomes into a library, but parses by path.
// Will return any error encountered.
func AddByDirectories(library *Library, directories []string, gzipped bool) error {
	for path := 0; path < structures.Paths; path++ {
		err := AddPathFromDirectories(library, directories, path, gzipped)
		if err != nil {
			return err
		}
	}
	return nil
}

// CompileDirectoriesToLibrary creates a new Library based on the directories given, sorts it, and gives it its ID (so this library is ready for use).
// Returns the library pointer and an error, if any (nil if no error was encounted)
func CompileDirectoriesToLibrary(directories []string, libraryTextFile string, gzipped bool) (*Library, error) {
	l := InitializeLibrary(libraryTextFile, nil)
	err := AddByDirectories(&l, directories, gzipped)
	if err != nil {
		return nil, err
	}
	SortLibrary(&l)
	(&l).AssignID()
	return &l, nil
}

// SequentialCompileDirectoriesToLibrary creates a new Library based on the directories given, sorts it, and gives it its ID (so this library is ready for use).
// This adds each directory sequentially (so each genome is done one at a time, rather than doing one path of all genomes all at once)
// Returns the library pointer and an error, if any (nil if no error was encounted)
func SequentialCompileDirectoriesToLibrary(directories []string, libraryTextFile string) (*Library, error) {
	l := InitializeLibrary(libraryTextFile, nil)
	for _, directory := range directories {
		err := AddLibraryFastJ(directory, &l)
		if err != nil {
			return nil, err
		}
	}
	SortLibrary(&l)
	(&l).AssignID()
	return &l, nil
}

// RemoveIntermediateFile removes the file specified by library.Text.
// If it is a directory, it does not remove anything, since then it would remove sglfv2 files.
func RemoveIntermediateFile(library *Library) error {
	info, err := os.Stat(library.Text)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return ErrRemoveIntermediateDirectory
	}
	err = os.Remove(library.Text)
	if err != nil {
		return err
	}
	return nil
}

// InitializeLibrary sets up the basic structure for a library.
// For consistency, it's best to use an absolute path for the text file. Relative paths will still work, but they are not recommended.
func InitializeLibrary(textFile string, componentLibraries [][md5.Size]byte) Library {
	var newLibraryPaths []concurrentPath
	newLibraryPaths = make([]concurrentPath, structures.Paths, structures.Paths)
	for i := range newLibraryPaths {
		var newLock sync.RWMutex
		newLibraryPaths[i] = concurrentPath{newLock, make([]*KnownVariants, 0, 1)} // Lock is copied, but hasn't been used yet, so this is fine.
	}
	return Library{Paths: newLibraryPaths, Text: textFile, Components: componentLibraries}
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
					variantCopy.ReferenceLibrary = destination
					destination.Paths[i].Variants[step].List[j] = &variantCopy
				}
				// copying List and Counts manually instead of copying Variants[step] makes sure the pointers for List and Counts aren't reused.
			}
		}
		source.Paths[i].Lock.RUnlock()
	}
}

// pointVariantsToLibrary points ALL variants of a library to another library.
// Mainly used for merging libraries, since variants need to point to their original libraries.
func pointVariantsToLibrary(library *Library, newLibrary *Library) {
	for i := range library.Paths {
		library.Paths[i].Lock.RLock()
		for step := range library.Paths[i].Variants {
			if library.Paths[i].Variants[step] != nil {
				for _, variant := range library.Paths[i].Variants[step].List {
					variant.ReferenceLibrary = newLibrary
				}
			}
		}
		library.Paths[i].Lock.RUnlock()
	}
}

// MergeLibraries is a function to merge the first library given with the second library.
// This version creates a new library.
// Returns the library pointer and an error, if any (nil if no error was encounted)
func MergeLibraries(libraryToMerge *Library, mainLibrary *Library, textFile string) (*Library, error) {
	allComponents := append([][md5.Size]byte{libraryToMerge.ID, mainLibrary.ID}, libraryToMerge.Components...)
	allComponents = append(allComponents, mainLibrary.Components...)
	newLibrary := InitializeLibrary(textFile, allComponents)
	libraryCopy(&newLibrary, mainLibrary)
	pointVariantsToLibrary(&newLibrary, mainLibrary)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			err := mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j], &newLibrary)
			if err != nil {
				return nil, err
			}
		}
	}
	SortLibrary(&newLibrary)
	return &newLibrary, nil
}

// MergeKnownVariants puts the contents of a KnownVariants at a specific path and step into another library.
// Will return any encountered, if any.
func mergeKnownVariants(genomePath, step int, variantsToMerge *KnownVariants, newLibrary *Library) error {
	if variantsToMerge != nil {
		for i, variant := range variantsToMerge.List {
			if index, ok := TileExists(genomePath, step, variant, newLibrary); !ok {
				addTileUnsafe(genomePath, step, variant, newLibrary)
				newIndex, ok := TileExists(genomePath, step, variant, newLibrary)
				if ok {
					(*newLibrary).Paths[genomePath].Variants[step].Counts[newIndex] += variantsToMerge.Counts[i] - 1
				} else {
					return ErrTileContradiction
				}
			} else {
				(*newLibrary).Paths[genomePath].Variants[step].Counts[index] += variantsToMerge.Counts[i]
			}
		}
	}
	return nil
}

// MergeLibrariesWithoutCreation merges libraries without creating a new one, using the "mainLibrary" instead.
// Returns the library pointer and an error, if any (nil if no error was encounted)
func MergeLibrariesWithoutCreation(libraryToMerge *Library, mainLibrary *Library) (*Library, error) {
	mainLibrary.Components = append(mainLibrary.Components, libraryToMerge.ID, mainLibrary.ID) // This is okay since this involves the mainLibrary's old ID.
	mainLibrary.Components = append(mainLibrary.Components, libraryToMerge.Components...)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			err := mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j], mainLibrary)
			if err != nil {
				return nil, err
			}
		}
	}
	SortLibrary(mainLibrary)
	return mainLibrary, nil
}

// LiftoverMapping is a representation of a liftover from one library to another.
// If a = LiftoverMapping.Mapping[b][c][d], then in path b, step c, variant d in the first library maps to variant a in the second.
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
					index, ok := TileExists(path, step, variant, destination)
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
	textFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	splitpath := strings.Split(filepath, ".")
	if len(splitpath) < 2 || splitpath[1] != "sglfmapping" {
		return nil, [16]byte{}, [16]byte{}, ErrBadLiftover
	}
	info, err := openFile(filepath)
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

// ParseSGLFv2 is a function to put SGLFv2 data back into a library.
// Allows for gzipped SGLFv2 files and regular SGLFv2 files.
// The following bug is believed to be fixed, but just in case it is not and an error occurs here:
// Note: the creation of an SGLFv2 may be incorrect and may put a series of bases before the hash, when the input tile from the FastJ file
// was extremely long (e.g. a couple million bases). This will result in an error here when decoding the hash. To fix this, you can instruct
// the tileBuilder in bufferedTileRead to Grow before constructing the tiles in that path. So far, this error is known to happen for path 811,
// which will create 032b.sglfv2 and 032c.sglfv2 incorrectly.
// It's also possible that the cause is that the two goroutines for reading and writing end too early, leaving some tiles with
// LookupReferences of -1, which results in the incorrect results. This could happen if bufferedTileRead ended first.
// Returns any error encountered, or nil if there's no error.
func ParseSGLFv2(filepath string, library *Library) error {
	file := path.Base(filepath)
	splitpath := strings.Split(file, ".")
	if len(splitpath) != 2 && len(splitpath) != 3 {
		return errors.New("error: Not a valid file") // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "sglfv2" || (len(splitpath) == 3 && splitpath[2] != ".gz") {
		return errors.New("error: not an sglfv2 file") // Makes sure that the file is an SGLFv2 file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		return errors.New("invalid hex file name") // Makes sure the title of the file is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0]) + int(pathHex[1]) // conversion into an integer--this is the path number
	var data []byte
	var err error
	if len(splitpath) == 2 {
		data, err = openFile(filepath)
	} else {
		data, err = openGZ(filepath)
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
			newVariant := structures.TileVariant{Hash: hashArray, Length: int(length), Annotation: "", LookupReference: 27 + int64(len(tileLengthString)) + int64(referenceCounter), Complete: isComplete(fields[2]), ReferenceLibrary: library}
			index, ok := TileExists(hexNumber, int(step), &newVariant, library)
			if !ok {
				addTileUnsafe(hexNumber, int(step), &newVariant, library)
			}
			index, ok = TileExists(hexNumber, int(step), &newVariant, library)
			if ok {
				(*library).Paths[hexNumber].Variants[int(step)].Counts[index] += int(count - 1)
			} else {
				return ErrTileContradiction
			}
			// Adding count-1 instead of count since AddTile will already add 1 to the count of the new tile.
			referenceCounter += len(line) + 1 // +1 to account for the newline.
		} else if line != "" { // This refers to the first line, which contains ID and Component information.
			idSlice := strings.Split(line, ";")
			idString := strings.Split(idSlice[0], ":")[1]
			libraryHash, err := hex.DecodeString(idString)
			if err != nil {
				return err
			}
			var hashArray [md5.Size]byte
			copy(hashArray[:], libraryHash)
			library.ID = hashArray
			components := strings.Split(idSlice[1], ":")[1]
			if components != "" {
				componentStrings := strings.Split(components, ",")
				component1Hash, err := hex.DecodeString(componentStrings[0])
				if err != nil {
					return err
				}
				component2Hash, err := hex.DecodeString(componentStrings[1])
				if err != nil {
					return err
				}
				var component1HashArray [md5.Size]byte
				var component2HashArray [md5.Size]byte
				copy(component1HashArray[:], component1Hash)
				copy(component2HashArray[:], component2Hash)
				library.Components = append(library.Components, component1HashArray, component2HashArray)
			}
		}
	}
	return nil
}

// AddLibrarySGLFv2 adds a directory of SGLFv2 files to a library.
// Library should be initialized with this directory as the Text field, so that text files of bases and directories aren't mixed together.
// Returns any error encountered, or nil if there's no error.
func AddLibrarySGLFv2(library *Library) error {
	info, err := os.Stat(library.Text)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrIncorrectSourceText
	}

	sglfv2Files, err := ioutil.ReadDir(library.Text)
	if err != nil {
		return err
	}
	for _, file := range sglfv2Files {
		if strings.HasSuffix(file.Name(), ".sglfv2") { // Checks if a file is an sglfv2 file.
			err = ParseSGLFv2(path.Join(library.Text, file.Name()), library)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
