package tilelibrary

// This tile library package assumes that any necessary imputation was done beforehand.
// Note: adding tiles to a library at any point will require sorting that library before writing it to a file.

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
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"../structures"
)

// openGZ is a function to open gzipped files and return the corresponding slice of bytes of the data.
// Mostly important for gzipped FastJs, but any gzipped file can be opened too.
func openGZ(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gz.Close()

	data, err := ioutil.ReadAll(gz)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// openFile is a method to get the data of a file and return the corresponding slice of bytes.
// Available mostly as a way to be flexible with files, since with openGZ a gzipped file or a non-gzipped file can be read and used.
func openFile(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// KnownVariants is a struct to hold the known variants in a specific step.
type KnownVariants struct {
	List [](*structures.TileVariant)         // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index after sorting)
	Counts []int // Counts of each variant so far
}

// Library is a type to represent a library of tile variants.
type Library struct {
	Paths []concurrentPath // The paths of the library.
	ID [md5.Size]byte // The ID of a library.
	Text string // The path of the text file containing the bases. As a special case, if Text is a directory, it refers to the sglf/sglfv2 files there.
	// If empty, the library was merged.
	// Note: the Text field is only relevant to the file system this Library is on.
	Length int64 // The total length of text files used by this library or all of its sublibraries.
	Components [][md5.Size]byte // the IDs of the libraries that directly made this library (empty if this library was not made from other)
	sublibraries []*Library // The pointers to the direct sublibraries of this library
}

// Helper function used to make byte slices from unsigned int32s.
// uint32 should be enough space for every integer used in the library.
func uint32ToByteSlice(integer uint32) []byte {
	slice := make([]byte, 4)
	binary.LittleEndian.PutUint32(slice, integer)
	return slice
}

// AssignID assigns a library its ID.
// Current method takes the path, step, variant number, variant count, variant hash, and variant length into account when making the ID.
// TODO: Find another way, if possible, since this takes so much time and memory.
// TODO: Find another way that isn't dependent on order of elements.
func (l *Library) AssignID() {
	var byteList []byte
	byteList = make([]byte, 0, 6000000000) // Allocation of 6 gigabytes, a little over the size of what would be expected of a library with 10 tiles per step.
	for path := range l.Paths {
		byteList = append(byteList, uint32ToByteSlice(uint32(path))...)
		l.Paths[path].Lock.RLock()
		for step, variants := range l.Paths[path].Variants {
			byteList = append(byteList,uint32ToByteSlice(uint32(step))...)
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
func (l Library) Equals(l2 Library) bool {
	for path := range l.Paths {
		l.Paths[path].Lock.RLock()
		l2.Paths[path].Lock.RLock()
		if len(l.Paths[path].Variants) != len(l2.Paths[path].Variants) {
			fmt.Println("error in path lengths", path)
			return false
		}
		for step, stepList := range l.Paths[path].Variants {
			if stepList != nil && l2.Paths[path].Variants[step] != nil {
				if len((*stepList).List) != len(l2.Paths[path].Variants[step].List) {
					fmt.Println("error in number of variants", path, step)
					return false
				}
				for i, variant := range (*stepList).List {
					if !variant.Equals((*l2.Paths[path].Variants[step].List[i])) || l.Paths[path].Variants[step].Counts[i] != l2.Paths[path].Variants[step].Counts[i] {
						fmt.Println("error in variants or counts", path, step)
						return false
					}
				}
			} else if stepList != nil || l2.Paths[path].Variants[step] != nil {
				fmt.Println("one set of variants is nil", path, step)
				return false
			}
		}
		l2.Paths[path].Lock.RUnlock()
		l.Paths[path].Lock.RUnlock()
	}
	return true
}
/*
// A simpler way of checking equality, since two libraries with the same ID are almost certainly equal.
func (l Library) Equals(l2 Library) bool {
	return l.ID=l2.ID
}
*/

// concurrentPath is a type to represent a path, while also being safe for concurrent use.
type concurrentPath struct {
	Lock sync.RWMutex // The read/write lock used for concurrency within a path.
	Variants []*KnownVariants // The list of steps, where each entry contains the known variants at that step.
}

// SortLibrary is a function to sort the library once all initial genomes are done being added.
// This function should only be used once after initial setup of the library, after all tiles have been added, since it sorts everything.
func SortLibrary(library *Library) {
	type sortStruct struct { // Temporary struct that groups together the variant and the count for sorting purposes.
		variant *structures.TileVariant
		count int
	}
	for i := range (*library).Paths {
		(*library).Paths[i].Lock.Lock()
		for _, steplist := range (*library).Paths[i].Variants {
			if steplist != nil {
				var sortStructList []sortStruct
				sortStructList = make([]sortStruct, len((*steplist).List))
				for i:=0; i<len((*steplist).List); i++ {
					sortStructList[i] = sortStruct{(*steplist).List[i], (*steplist).Counts[i]}
				}
				sort.Slice(sortStructList, func(i, j int) bool { return sortStructList[i].count > sortStructList[j].count })
				for j:=0; j<len((*steplist).List); j++ {
					(*steplist).List[j], (*steplist).Counts[j]= sortStructList[j].variant, sortStructList[j].count
				}
			}
		}
		(*library).Paths[i].Lock.Unlock()
	}
}

// TileExists is a function to check if a specific tile exists at a specific path and step in a library.
// Returns the index of the variant, if found--otherwise, returns -1, meaning not found.
func TileExists(path, step int, toCheck *structures.TileVariant, library *Library) int {
	(*library).Paths[path].Lock.Lock()
	defer (*library).Paths[path].Lock.Unlock()
	if len((*library).Paths[path].Variants) > step && (*library).Paths[path].Variants[step] != nil { // Safety to make sure that the KnownVariants struct has been created
		for i, value := range (*library).Paths[path].Variants[step].List {
			if toCheck.Equals(*value) {
				return i
			}
		}
		return -1
	}
	for len((*library).Paths[path].Variants) <= step { // Makes enough room so that there are step+1 elements in Paths[path].Variants
		(*library).Paths[path].Variants = append((*library).Paths[path].Variants, nil)
	}
	newKnownVariants := &KnownVariants{make([](*structures.TileVariant), 0, 1), make([]int, 0, 1)}
	(*library).Paths[path].Variants[step] = newKnownVariants
	return -1
}

// AddTile is a function to add a tile (without sorting).
// Safe to use without checking existence of the tile beforehand (since the function will do that for you).
func AddTile(genomePath, step int, new *structures.TileVariant, library *Library) {
	if index := TileExists(genomePath, step, new, library); index == -1 { // Checks if the tile exists already.
		(*library).Paths[genomePath].Lock.Lock()
		defer (*library).Paths[genomePath].Lock.Unlock()
		(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
		(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)
	} else {
		(*library).Paths[genomePath].Lock.RLock()
		defer (*library).Paths[genomePath].Lock.RUnlock()
		(*library).Paths[genomePath].Variants[step].Counts[index]++ // Adds 1 to the count of the tile (since it's already in the library)
	}
}

// AddTileUnsafe is a function to add a tile without sorting.
// Unsafe because it doesn't check if the tile is already in the library, unlike AddTile.
// Be careful to check if the tile already exists before using this function to avoid repeats in the library.
func addTileUnsafe(genomePath, step int, new *structures.TileVariant, libraryTextFile string, library *Library) {
	(*library).Paths[genomePath].Lock.Lock()
	defer (*library).Paths[genomePath].Lock.Unlock()
	(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
	(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)	
}

// FindFrequency is a function to find the frequency of a specific tile at a specific path and step.
func FindFrequency(path, step int, toFind *structures.TileVariant, library *Library) int {
	if index:= TileExists(path, step, toFind, library); index != -1 {
		return (*library).Paths[path].Variants[step].Counts[index]
	}
	fmt.Println("Variant not found.")
	return 0 // Possibly return an error instead?
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified.
func Annotate(path, step int, hash structures.VariantHash, annotation string, library *Library) {
	(*library).Paths[path].Lock.Lock()
	defer (*library).Paths[path].Lock.Unlock()
	for _, tile := range (*library).Paths[path].Variants[step].List {
		if hash==(*tile).Hash {
			tile.Annotation = annotation
			break
		}
	}
	fmt.Printf("No matching tile found at specified path %v and step %v.\n", path, step) // Information if tile isn't found.
}


// baseInfo is a temporary struct to pass around information about a tile's bases.
type baseInfo struct {
	bases string
	hash structures.VariantHash
	variant *(structures.TileVariant)
}

var tileBuilder strings.Builder
// bufferedTileRead reads a FastJ file and adds its tiles to the provided library.
// Allows for gzipped FastJ files and regular FastJ files.
func bufferedTileRead(fastJFilepath, libraryTextFile string, library *Library) {
	var wg sync.WaitGroup
	var baseChannel chan baseInfo
	baseChannel = make(chan baseInfo, 16) // Put information about bases of tiles in here while they need to be processed.
	writeChannel := make(chan bool)
	wg.Add(2)
	go bufferedBaseWrite(libraryTextFile, baseChannel, writeChannel, &wg)
	file := path.Base(fastJFilepath) // The name of the file.
	splitpath := strings.Split(file, ".") // This is used to make sure the file is in the right format.
	if len(splitpath) != 3 && len(splitpath) != 2 {
		log.Fatal(errors.New("error: Not a valid file "+file)) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "fj" || (len(splitpath)==3 && splitpath[2] != "gz") {
		log.Fatal(errors.New("error: not a valid FastJ file")) // Makes sure that the file is a FastJ file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name")) // Makes sure the file title is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // Conversion from hex into decimal--this is the path
	var data []byte
	if len(splitpath) == 3 {
		data = openGZ(fastJFilepath)
	} else {
		data = openFile(fastJFilepath)
	}
	text := string(data)
	tiles := strings.Split(text, "\n\n") // The divider between two tiles is two newlines.
	
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") { // Makes sure that a "line" starts with the correct character ">"
			stepInHex := line[20:24] // These are the indices where the step is located.
			stepBytes, err := hex.DecodeString(stepInHex)
			if err != nil {
				log.Fatal(err)
			}
			step := 256 * int(stepBytes[0]) + int(stepBytes[1])
			hashString := line[40:72] // These are the indices where the hash is located.
			hash, err := hex.DecodeString(hashString)
			if err != nil {
				log.Fatal(err)
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
					k:=1 // This accounts for the possibility of the length of a tile spanning at least 16 tiles.
					for line[i-k] != ':' { // Goes back a few characters until it knows the string of the tile length.
						k++
					}
					lengthString=line[(i-k+1):i]
					break
				}
			}
			
			length, err := strconv.Atoi(lengthString) // Length is provided here in base 10
			if err != nil {
				log.Fatal(err)
			}
			baseData := strings.Split(line, "\n")[1:] // Data after the first part of the line are the bases of the tile variant.
			tileBuilder.Reset()
			if hexNumber == 811 { // Grows the buffer of tileBuilder to 2^22 bytes if the path is 811, which seems to contain very large (roughly 2.7 million bases per tile) tiles.
				// Not sure of the exact reasons why this is needed, but SGLFv2 construction will fail on paths 811 and 812 without these lines.
				// TODO: find the right limit at which tileBuilder should Grow.
				tileBuilder.Grow(4194304)
			}
			// Test to see if this line takes a long time to run.
			for i := range baseData {
				tileBuilder.WriteString(baseData[i])
			}
			bases := tileBuilder.String()
			newTile := &structures.TileVariant{Hash: hashArray, Length: length, Annotation: "", LookupReference: -1}
			if tileIndex:=TileExists(hexNumber, step, newTile, library); tileIndex==-1 {
				addTileUnsafe(hexNumber, step, newTile, libraryTextFile, library)
				baseChannel <- baseInfo{bases,hashArray, newTile}
			} else {
				(*library).Paths[hexNumber].Variants[step].Counts[tileIndex]++ // Increments the count of the tile variant if it is found.
			}
		}
	}
	wg.Done()
	var nilHash [md5.Size]byte
	baseChannel <- baseInfo{"", nilHash, nil} // Sends a "nil" baseInfo to signal that there are no more tiles.
	<-writeChannel // Make sure the writer is done.
	close(baseChannel)
	splitpath, data, tiles = nil, nil, nil // Clears most things in memory that were used here, to free up memory.
	wg.Wait()
}

// bufferedBaseWrite writes bases and hashes of tiles to the given text file.
// To be used in conjunction with bufferedTileRead and bufferedInfo.
func bufferedBaseWrite(libraryTextFile string, channel chan baseInfo, writeChannel chan bool, group *sync.WaitGroup) {
	err := os.MkdirAll(path.Dir(libraryTextFile), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(libraryTextFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(file)
	for bases := range channel {
		if bases.variant != nil { // Checks if there are any more tiles
			info, err := os.Stat(libraryTextFile)
			if err != nil {
				log.Fatal(err)
			}
			(bases.variant).LookupReference = info.Size()+int64(bufferedWriter.Buffered())
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
}

// writeToTextFile writes the entry of a lookup from a hash to bases for a specific path and step, in a text file.
func writeToTextFile(genomePath, step int, directory, bases, filename string, hash structures.VariantHash) {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	textFile, err := os.OpenFile(path.Join(directory,filename), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fileWriter := bufio.NewWriter(textFile)
	var b strings.Builder
	b.WriteString(hex.EncodeToString(hash[:]))
	b.WriteString(",")
	b.WriteString(bases)
	b.WriteString("\n")
	fileWriter.WriteString(b.String())
	fileWriter.Flush()
	textFile.Close()
}

// writePathToSGLF writes an SGLF for an entire path given a library.
// This assumes that the library has been sorted beforehand.
func writePathToSGLF(library *Library, genomePath, version int, directoryToWriteTo, directoryToGetFrom, textFilename string) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglf")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo,b.String()), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(sglfFile)
	textFile, err := os.OpenFile(path.Join(directoryToGetFrom,textFilename), os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fileReader := bufio.NewReader(textFile)
	(*library).Paths[genomePath].Lock.RLock()
	defer (*library).Paths[genomePath].Lock.RUnlock()
	for step := range (*library).Paths[genomePath].Variants {
		if (*library).Paths[genomePath].Variants[step] != nil {
			for i := range (*(*library).Paths[genomePath].Variants[step]).List {
				textFile.Seek(int64((*(*library).Paths[genomePath].Variants[step]).List[i].LookupReference),0)
				fileReader.Reset(textFile)
				tileString, err := fileReader.ReadString('\n') // This includes the newline at the end.
				if err != nil && err != io.EOF {
					log.Fatal(err)
				}
				stepHex := fmt.Sprintf("%04x", step)
				
				bufferedWriter.WriteString(fmt.Sprintf("%04x", genomePath)) // path
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%02v", version)) // Version
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(stepHex) // Step
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%03x", i)) // Tile variant number
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(fmt.Sprintf("%01x", (*(*library).Paths[genomePath].Variants[step]).List[i].Length)) // Tile length
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(tileString) // Hash and bases of tile.
			}
		}
	}
	bufferedWriter.Flush()
	sglfFile.Close()
	textFile.Close()
}

// WriteLibraryToSGLF writes the contents of a library to SGLF files.
func WriteLibraryToSGLF(library *Library, directoryToWriteTo, directoryToGetFrom, textFile string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLF(library, path, 0, directoryToWriteTo, directoryToGetFrom, textFile)
	}
}

// isComplete determines if a set of bases is complete (has no nocalls).
func isComplete(bases string) bool {
	return !strings.ContainsRune(bases, 'n')
}

// LookupHashAndBases looks up the hash and bases, as a string, of a variant in a specific library.
// This will include a newline at the end, which can be removed.
func (l *Library) LookupHashAndBases(variant *structures.TileVariant, fileMap *map[*Library]*os.File) string {
	if l.Text == "" { // library was merged, and so it must have sublibraries
		if variant.LookupReference < l.sublibraries[0].Length{
			return l.sublibraries[0].LookupHashAndBases(variant, fileMap)
		}
		variantCopy := *variant
		variantCopy.LookupReference -= l.sublibraries[0].Length // reduced as to fit the scale of the second sublibrary.
		return l.sublibraries[1].LookupHashAndBases(&variantCopy, fileMap)
	}
	file := (*fileMap)[l]
	file.Seek(variant.LookupReference, 0)
	bufferedReader := bufio.NewReader(file)
	hashAndBases, err := bufferedReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return hashAndBases
}

// originalLibrary determines which library a variant was originally added from, and returns that library, along with the original reference for that variant in that library.
func originalLibraryAndReference(variant *structures.TileVariant, library *Library) (*Library , int64){
	if library.Text == "" { // library was merged, and so it must have sublibraries
		if variant.LookupReference < library.sublibraries[0].Length{
			return originalLibraryAndReference(variant, library.sublibraries[0])
		}
		variantCopy := *variant
		variantCopy.LookupReference -= library.sublibraries[0].Length // reduced as to fit the scale of the second sublibrary.
		return originalLibraryAndReference(&variantCopy, library.sublibraries[1])
	}
	return library, variant.LookupReference
}

// searchForFiles generates a map of libraries to their corresponding file pointers.
// It does this by checking the tree of libraries one at a time, collecting all of the files for this path.
func searchForFiles(libraryToSearch *Library, genomePath int, fileMap *map[*Library]*os.File) {
	if libraryToSearch.Text=="" {
		searchForFiles(libraryToSearch.sublibraries[0], genomePath, fileMap)
		searchForFiles(libraryToSearch.sublibraries[1], genomePath, fileMap)
	} else {
		info, err := os.Stat(libraryToSearch.Text)
		if err != nil {
			log.Fatal(err)
		}
		var fileToOpen string
		if info.IsDir() {
			fileToOpen = path.Join(libraryToSearch.Text, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
		} else {
			fileToOpen = libraryToSearch.Text // Otherwise, just use the file already given.
		}
		file, err := os.Open(fileToOpen)
		if err != nil {
			log.Fatal(err)
		}
		(*fileMap)[libraryToSearch] = file
	}
}

// writePathToSGLFv2 writes an SGLFv2 for an entire path given a library.
// This assumes that the library has been sorted beforehand.
// TODO: differentiate between writing when the reference is a directory and when it's a text file.
// TODO: Handle merged libraries, which don't necessarily have their own text file of bases and hashes.
// Ideally, don't want to move contents of files around, since the contents of files are very large
// Another scenario: merged library of merged libraries--may need to be able to call this function recursively
// Also, don't want to keep on opening and closing files--figure out how to store them.
// maybe library can keep a list arranged like a heap?
func writePathToSGLFv2(library *Library, genomePath, version int, directoryToWriteTo, directoryToGetFrom, textFilename string) {
	var pathFileMap map[*Library]*os.File
	pathFileMap = make(map[*Library]*os.File, 0)
	searchForFiles(library, genomePath, &pathFileMap)
	for _, file := range pathFileMap {
		defer file.Close()
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglfv2")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	sglfFile, err := os.OpenFile(path.Join(directoryToWriteTo,b.String()), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(sglfFile)
	info, err := os.Stat(path.Join(directoryToGetFrom, textFilename))
	if err != nil {
		log.Fatal(err)
	}
	var fileToOpen string
	if info.IsDir() {
		fileToOpen = path.Join(directoryToGetFrom, fmt.Sprintf("%04x.sglfv2", genomePath)) // If the reference is a directory, point to the corresponding SGLFv2 file as a reference.
	} else {
		fileToOpen = path.Join(directoryToGetFrom, textFilename)
	}
	textFile, err := os.OpenFile(fileToOpen, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fileReader := bufio.NewReader(textFile)
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
				_, err := textFile.Seek((*(*library).Paths[genomePath].Variants[step]).List[i].LookupReference,0)
				if err != nil {
					log.Fatal(err)
				}
				fileReader.Reset(textFile)
				tileString, err := fileReader.ReadString('\n') // This includes the newline at the end.
				if err != nil && err != io.EOF {
					log.Fatal(err)
				}
				stepHex := fmt.Sprintf("%04x", step)
				
				bufferedWriter.WriteString(fmt.Sprintf("%04x", genomePath)) // path
				bufferedWriter.WriteString(".")
				bufferedWriter.WriteString(fmt.Sprintf("%02v", version)) // Version
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
	bufferedWriter.Flush()
	sglfFile.Close()
	textFile.Close()
	
}

// WriteLibraryToSGLFv2 writes the contents of a library to SGLFv2 files.
func WriteLibraryToSGLFv2(library *Library, directoryToWriteTo, directoryToGetFrom, textFile string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLFv2(library, path, 1, directoryToWriteTo, directoryToGetFrom, textFile)
	}
}

// AddLibraryFastJ adds a directory of gzipped FastJ files to a specific library. 
func AddLibraryFastJ(directory, libraryTextFile string, library *Library) {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") { // Checks if a file is a gz file.
			bufferedTileRead(path.Join(directory, file.Name()), libraryTextFile, library)
		}
	}
}

// AddPathFromDirectories parses the same path for all genomes, represented by a list of directories, and puts the information in a Library.
func AddPathFromDirectories(library *Library, directories []string, genomePath int, libraryTextFile string) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%04x",genomePath))
	b.WriteString(".fj.gz")
	for _, directory := range directories {
		bufferedTileRead(path.Join(directory,b.String()),libraryTextFile,library)
	}
}

// AddByDirectories adds information from a list of directories for genomes into a library, but parses by path.
func AddByDirectories(library *Library, directories []string, libraryTextFile string) {
	for path := 0; path < structures.Paths; path++ {
		AddPathFromDirectories(library, directories, path, libraryTextFile)
	}
}

// CompileDirectoriesToLibrary creates a new Library based on the directories given, sorts it, and gives it its ID (so this library is ready for use).
func CompileDirectoriesToLibrary(directories []string, libraryTextFile string) *Library {
	l := InitializeLibrary(libraryTextFile, [][md5.Size]byte{})
	AddByDirectories(&l, directories, libraryTextFile)
	SortLibrary(&l)
	(&l).AssignID()
	info, err := os.Stat(libraryTextFile)
	if err != nil {
		log.Fatal(err)
	}
	l.Length = info.Size()
	return &l
}

// InitializeLibrary sets up the basic structure for a library.
// For consistency, it's best to use an absolute path for the text file.
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
		source.Paths[i].Lock.RLock() // Locked since we're reading from path.Variants when we copy.
		destination.Paths[i].Variants = make([]*KnownVariants, len(source.Paths[i].Variants)) // This is to make sure all elements are copied over.
		copy(destination.Paths[i].Variants, source.Paths[i].Variants)
		source.Paths[i].Lock.RUnlock()
	}
}

// MergeLibraries is a function to merge the first library given with the second library.
// This version creates a new library.
func MergeLibraries(libraryToMerge *Library, mainLibrary *Library) *Library {
	newLibrary := InitializeLibrary("", [][md5.Size]byte{libraryToMerge.ID, mainLibrary.ID})
	libraryCopy(&newLibrary, mainLibrary)
	newLibrary.Length = libraryToMerge.Length + mainLibrary.Length
	newLibrary.sublibraries = []*Library{mainLibrary, libraryToMerge}
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j], &newLibrary, mainLibrary.Length)
		}
	}
	SortLibrary(&newLibrary)
	return &newLibrary
}

// MergeKnownVariants puts the contents of a KnownVariants at a specific path and step into another library.
// Account for CGF files here (try to avoid potential remapping with CGF files)
// Should create a new library and point the old libraries to this library.
func mergeKnownVariants(genomePath, step int, variantsToMerge *KnownVariants, newLibrary *Library, oldLibraryLength int64) {
	originalLibraryStepLength := len((*newLibrary).Paths[genomePath].Variants[step].List)
	newTileCounter := 0
	for i, variant := range variantsToMerge.List {
		if index := TileExists(genomePath, step, variant, newLibrary); index==-1 {
			variantCopy := *variant
			variantCopy.LookupReference += oldLibraryLength // Signals that it's from the libraryToMerge rather than the mainLibrary.
			addTileUnsafe(genomePath, step, &variantCopy, "", newLibrary)
			(*newLibrary).Paths[genomePath].Variants[step].Counts[originalLibraryStepLength+newTileCounter] += variantsToMerge.Counts[i]-1
			newTileCounter++
		} else {
			(*newLibrary).Paths[genomePath].Variants[step].Counts[index] += variantsToMerge.Counts[i]
		}
	}
}

// MergeLibrariesWithoutCreation merges libraries without creating a new one, using the "mainLibrary" instead.
func MergeLibrariesWithoutCreation(text string, libraryToMerge *Library, mainLibrary *Library) *Library {
	mainLibrary.Components = append(mainLibrary.Components, libraryToMerge.ID)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			mergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j], mainLibrary, mainLibrary.Length)
		}
	}
	SortLibrary(mainLibrary)
	return mainLibrary
}

// LiftoverMapping is a representation of a liftover from one library to another.
// If a = LiftoverMapping.Mapping[b][c][d], then in path b, step c, variant d in the first library maps to variant a in the second.
type LiftoverMapping struct {
	Mapping [][][]int 
	SourceLibrary *Library // The source library to map from.
	DestinationLibrary *Library // The destination library to map to.
}
// TODO: store mappings.

// CreateMapping creates a liftover mapping from the source library to the destination library.
func CreateMapping(source, destination *Library) LiftoverMapping {
	index := -1
	for i, libraryID := range destination.Components {
		if libraryID == source.ID {
			index = i
			break
		}
	}
	if index == -1 { // Destination was not made from the source--can't guarantee a mapping here.
		log.Fatal(errors.New("source library is not part of the destination library"))
	}
	var mapping [][][]int
	mapping = make([][][]int, structures.Paths, structures.Paths)
	for path := range (*source).Paths {
		(*source).Paths[path].Lock.RLock()
		mapping[path] = make([][]int, len((*source).Paths[path].Variants)) // Number of steps.
		for step, variants := range (*source).Paths[path].Variants {
			for _, variant := range (*variants).List {
				mapping[path][step] = append(mapping[path][step], TileExists(path, step, variant, destination))
			}
		}
		(*source).Paths[path].Lock.RUnlock()
	}
	return LiftoverMapping{mapping, source, destination}
}

// WriteMapping writes a LiftoverMapping to a specified file.
// The format is path/step/source1,destination1;source2,destination2;...
// Use goroutines per step here? Would need a Lock for concurrency issues
func WriteMapping(filename string, mapping LiftoverMapping) {
	textFile, err := os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(textFile)
	bufferedWriter.WriteString(hex.EncodeToString(mapping.SourceLibrary.ID[:]))
	bufferedWriter.WriteString(",")
	bufferedWriter.WriteString(hex.EncodeToString(mapping.DestinationLibrary.ID[:]))
	bufferedWriter.WriteString("\n")
	for path := range mapping.Mapping {
		for step := range mapping.Mapping[path] {
			bufferedWriter.WriteString(fmt.Sprintf("%04x",path))
			bufferedWriter.WriteString("/")
			bufferedWriter.WriteString(fmt.Sprintf("%04x",step))
			bufferedWriter.WriteString("/")
			for index, value := range mapping.Mapping[path][step] {
				bufferedWriter.WriteString(strconv.Itoa(index))
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(strconv.Itoa(value))
				bufferedWriter.WriteString(";")
			}
			bufferedWriter.WriteString("\n")
		}
	}
	bufferedWriter.Flush()
	textFile.Close()
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
func ParseSGLFv2(filepath string, library *Library) {
	file := path.Base(filepath)
	splitpath := strings.Split(file, ".") 
	if len(splitpath) != 2 && len(splitpath) != 3 {
		log.Fatal(errors.New("error: Not a valid file")) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "sglfv2" || (len(splitpath)==3 && splitpath[2] != ".gz") {
		log.Fatal(errors.New("error: not an sglfv2 file")) // Makes sure that the file is an SGLFv2 file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name")) // Makes sure the title of the file is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // conversion into an integer--this is the path number
	var data []byte
	if len(splitpath)== 2 {
		data = openFile(filepath)
	} else {
		data = openGZ(filepath)
	}
	text := string(data)
	tiles := strings.Split(text, "\n")
	referenceCounter := len(tiles[0])+1 // length of the first line plus the newline
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
				log.Fatal(err)
			}
			count, err := strconv.ParseInt(tileCountString, 16, 0)
			if err != nil {
				log.Fatal(err)
			}
			length, err := strconv.ParseInt(tileLengthString, 16, 0)
			if err != nil {
				log.Fatal(err)
			}
			hash, err := hex.DecodeString(hashString)
			if err != nil {
				log.Fatal(err)
			}
			var hashArray [16]byte
			copy(hashArray[:], hash)
			newVariant := structures.TileVariant{Hash: hashArray, Length: int(length), Annotation: "", LookupReference: 27+int64(len(tileLengthString))+int64(referenceCounter)}
			AddTile(hexNumber, int(step), &newVariant, library)
			(*library).Paths[hexNumber].Variants[int(step)].Counts[TileExists(hexNumber, int(step), &newVariant, library)] += int(count-1)
			// Adding count-1 instead of count since AddTile will already add 1 to the count of the new tile.
			referenceCounter += len(line)+1 // +1 to account for the newline.
		} else if line!="" { // This refers to the first line, which contains ID and Component information.
			idSlice := strings.Split(line, ";")
			idString := strings.Split(idSlice[0], ":")[1]
			libraryHash, err := hex.DecodeString(idString)
			if err != nil {
				log.Fatal(err)
			}
			var hashArray[md5.Size]byte
			copy(hashArray[:], libraryHash)
			library.ID=hashArray
			components := strings.Split(idSlice[1], ":")[1]
			if components != "" {
				componentStrings := strings.Split(components, ",")
				component1Hash, err := hex.DecodeString(componentStrings[0])
				if err != nil {
					log.Fatal(err)
				}
				component2Hash, err := hex.DecodeString(componentStrings[1])
				if err != nil {
					log.Fatal(err)
				}
				var component1HashArray [md5.Size]byte
				var component2HashArray [md5.Size]byte
				copy(component1HashArray[:], component1Hash)
				copy(component2HashArray[:], component2Hash)
				library.Components = append(library.Components, component1HashArray, component2HashArray)
			}
		}
	}
}

// AddLibrarySGLFv2 adds a directory of SGLFv2 files to a library.
// Library should be initialized with this directory as the Text field, so that text files of bases and directories aren't mixed together.
func AddLibrarySGLFv2(directory string, library *Library) {
	if directory == library.Text {
		sglfv2Files, err := ioutil.ReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range sglfv2Files {
			if strings.HasSuffix(file.Name(), ".sglfv2") { // Checks if a file is an sglfv2 file.
				if file.Size() > library.Length {
					library.Length = file.Size() // the library takes on the length of the longest SGLFv2 file.
				}
				ParseSGLFv2(path.Join(directory, file.Name()), library)
			}
		}
	} else {
		log.Fatal(errors.New("directory specified is not the same as the library's text field"))
	}
}

