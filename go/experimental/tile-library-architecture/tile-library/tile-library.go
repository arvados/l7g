package main

// This tile library package assumes that any necessary imputation was done beforehand.

import (
	"bufio"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"../structures"
)

// KnownVariants is a struct to hold the known variants in a specific step.
type KnownVariants struct {
	List [](*structures.TileVariant)         // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index after sorting)
	Counts []int // Counts of each variant so far
}

// VariantLookupTable is a type for looking up the original positions of variants in a list--for now, the table is a list
type VariantLookupTable []int
/*
// Library is a type to represent a library of tile variants.
type Library []concurrentPath
*/

// Library is a type to represent a library of tile variants.
type Library struct {
	Paths []concurrentPath // The paths of the library.
	Text string // The path of the text file containing the bases--if a directory, it refers to the sglf files there.
	Components []string // the paths of the libraries that made this library (empty if this library was not made from other)
}
// TODO: Give some sort of ID to each library for easy reference when convert back from an SGLFv2 file.
// TODO: Keep track of the IDs of the component libraries as well, so they can be referenced.

// Equals checks for equality between two libraries. It does not check similarity in text or components, and tiles are checked by hash.
func (l1 Library) Equals(l2 Library) bool {
	for path := range l1.Paths {
		l1.Paths[path].Lock.RLock()
		l2.Paths[path].Lock.RLock()
		if len(l1.Paths[path].Variants) != len(l2.Paths[path].Variants) {
			return false
		}
		for step, stepList := range l1.Paths[path].Variants {
			if len((*stepList).List) != len(l2.Paths[path].Variants[step].List) {
				return false
			}
			for i, variant := range (*stepList).List {
				if !variant.Equals((*l2.Paths[path].Variants[step].List[i])) || l1.Paths[path].Variants[step].Counts[i] != l2.Paths[path].Variants[step].Counts[i] {
					return false
				}
			}
		}
		l2.Paths[path].Lock.RUnlock()
		l1.Paths[path].Lock.RUnlock()
	}
	return true
}

// TODO: convert SGLFv2 format into libraries for use (for convenience)

// concurrentPath is a type to represent a path, while also being safe for concurrent use.
type concurrentPath struct {
	Lock sync.RWMutex // The read/write lock used for concurrency within a step.
	Variants []*KnownVariants // The list of steps, where each step contains the known variants at that step.
}

// Function to sort the library once all initial genomes are done being added.
// This function should only be used once during initial setup of the library, after all tiles have been added, since it sorts everything.
func sortLibrary(library *Library) {
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
// Returns the index of the variant, if found--otherwise, returns -1.
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
	for len((*library).Paths[path].Variants) <= step {
		(*library).Paths[path].Variants = append((*library).Paths[path].Variants, nil)
	}
	newKnownVariants := &KnownVariants{make([](*structures.TileVariant), 0, 1), make([]int, 0, 1)}
	(*library).Paths[path].Variants[step] = newKnownVariants
	return -1
}

// AddTile is a function to add a tile (without sorting).
// Safe to use without checking existence of the tile beforehand (since the function will do that for you).
func AddTile(genomePath, step, lookupNumber int, new *structures.TileVariant, libraryTextFile, bases string, library *Library) {
	if index := TileExists(genomePath, step, new, library); index == -1 { // Checks if the tile exists already.
		(*library).Paths[genomePath].Lock.Lock()
		defer (*library).Paths[genomePath].Lock.Unlock()
		(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
		(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)
		//writeToTextFile(genomePath, step, path.Dir(libraryTextFile), bases, path.Base(libraryTextFile), new.Hash)
	} else {
		(*library).Paths[genomePath].Lock.RLock()
		defer (*library).Paths[genomePath].Lock.RUnlock()
		(*library).Paths[genomePath].Variants[step].Counts[index]++ // Adds 1 to the count of the tile (since it's already in the library)
	}
}

// AddTileUnsafe is a function to add a tile without sorting.
// Unsafe because it doesn't check if the tile is already in the library, unlike AddTile.
// Be careful to check if the tile already exists before using this function to avoid repeats in the library.
func AddTileUnsafe(genomePath, step, lookupNumber int, new *structures.TileVariant, libraryTextFile string, library *Library) {
	(*library).Paths[genomePath].Lock.Lock()
	defer (*library).Paths[genomePath].Lock.Unlock()
	(*library).Paths[genomePath].Variants[step].List = append((*library).Paths[genomePath].Variants[step].List, new)
	(*library).Paths[genomePath].Variants[step].Counts = append((*library).Paths[genomePath].Variants[step].Counts, 1)	
}

// FindFrequency is a function to find the frequency of a specific tile at a specific path and step.
func FindFrequency(path, step int, toFind *structures.TileVariant, library *Library) int {
	(*library).Paths[path].Lock.RLock()
	defer (*library).Paths[path].Lock.RUnlock()
	if index:= TileExists(path, step, toFind, library); index != -1 {
		return (*library).Paths[path].Variants[step].Counts[index]
	}
	fmt.Println("Variant not found.")
	return 0 // Possibly return an error instead?
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified.
func Annotate(path, step int, toAnnotate *structures.TileVariant, library *Library) {
	for _, tile := range (*library).Paths[path].Variants[step].List {
		if toAnnotate.Equals(*tile) {
			fmt.Print("Enter annotation: ")
			readKeyboard := bufio.NewReader(os.Stdin)
			annotation, err := readKeyboard.ReadString('\n') // Maybe allow for input for annotation without the use of keyboard input?
			if err!=nil {
				log.Fatal(err)
			}
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
func bufferedTileRead(fastJFilepath, libraryTextFile string, library *Library) {
	var baseChannel chan baseInfo
	baseChannel = make(chan baseInfo, 16) // Put information about bases of tiles in here while they need to be processed.
	go bufferedBaseWrite(libraryTextFile, baseChannel)
	file := path.Base(fastJFilepath) // The name of the file.
	splitpath := strings.Split(file, ".") // This is used to make sure the file is in the right format.
	if len(splitpath) != 3 {
		log.Fatal(errors.New("error: Not a valid gzipped file "+file)) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "fj" || splitpath[2] != "gz" {
		log.Fatal(errors.New("error: not a gzipped FastJ file")) // Makes sure that the file is a FastJ file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name")) // Makes sure the file title is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // Conversion from hex into decimal--this is the path

	data := structures.OpenGZ(fastJFilepath) // potentially create in the string in one pass through the file?
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
			
			length, err := strconv.ParseInt(lengthString, 16, 0)
			if err != nil {
				log.Fatal(err)
			}
			baseData := strings.Split(line, "\n")[1:] // Data after the first part of the line are the bases of the tile variant.
			tileBuilder.Reset()
			for i := range baseData {
				tileBuilder.WriteString(baseData[i])
			}
			bases := tileBuilder.String()
			newTile := &structures.TileVariant{Hash: hashArray, Length: int(length), Annotation: "", LookupReference: -1}
			if tileIndex:=TileExists(hexNumber, step, newTile, library); tileIndex==-1 {
				AddTileUnsafe(hexNumber, step, -1, newTile, libraryTextFile, library) // -1 for unknown lookup reference, set later.
				baseChannel <- baseInfo{bases,hashArray, newTile}
			} else {
				(*library).Paths[hexNumber].Variants[step].Counts[tileIndex]++ // Increments the count of the tile variant if it is found.
			}
		}
	}
	close(baseChannel)
	splitpath, data, tiles = nil, nil, nil // Clears most things in memory that were used here, to free up memory.
}

// bufferedBaseWrite writes bases and hashes of tiles to the given text file.
// To be used in conjunction with bufferedTileRead and bufferedInfo.
func bufferedBaseWrite(libraryTextFile string, channel chan baseInfo) {
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
		info, err := os.Stat(libraryTextFile)
		if err != nil {
			log.Fatal(err)
		}
		(bases.variant).LookupReference = (int(info.Size())+bufferedWriter.Buffered())
		hashString := hex.EncodeToString(bases.hash[:])
		bufferedWriter.WriteString(hashString)
		bufferedWriter.WriteString(",")
		bufferedWriter.WriteString(bases.bases)
		bufferedWriter.WriteString("\n")
	}
	bufferedWriter.Flush()

	file.Close()
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
func WriteLibraryToSGLF(library *Library, version int, directoryToWriteTo, directoryToGetFrom, textFile string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLF(library, path, version, directoryToWriteTo, directoryToGetFrom, textFile)
	}
}

// writePathToSGLFv2 writes an SGLFv2 for an entire path given a library.
// This assumes that the library has been sorted beforehand.
func writePathToSGLFv2(library *Library, genomePath, version int, directoryToWriteTo, directoryToGetFrom, textFilename string) {
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
	textFile, err := os.OpenFile(path.Join(directoryToGetFrom,textFilename), os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fileReader := bufio.NewReader(textFile)
	for step, variants := range (*library).Paths[genomePath].Variants {
		if variants != nil {
			for i, variant := range (*variants).List {
				textFile.Seek(int64(variant.LookupReference),0)
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
				bufferedWriter.WriteString(fmt.Sprintf("%08x", (*variants).Counts[i])) // The count of this tile.
				bufferedWriter.WriteString("+")
				bufferedWriter.WriteString(fmt.Sprintf("%01x", variant.Length)) // Tile length
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(tileString) // Hash and bases of tile.
			}
		}
	}
	bufferedWriter.Flush()
	sglfFile.Close()
	textFile.Close()
}

// WriteLibraryToSGLFv2 writes the contents of a library to SGLFv2 files.
func WriteLibraryToSGLFv2(library *Library, version int, directoryToWriteTo, directoryToGetFrom, textFile string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLFv2(library, path, version, directoryToWriteTo, directoryToGetFrom, textFile)
	}
}
/*
func parseSGLF(filepath string, library *Library) {
	file := path.Base(filepath)
	splitpath := strings.Split(file, ".") 
	if len(splitpath) != 3 {
		log.Fatal(errors.New("error: Not a valid gzipped file")) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "sglf" || splitpath[2] != "gz" {
		log.Fatal(errors.New("error: not a gzipped sglf file")) // Makes sure that the file is an SGLF file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name")) // Makes sure the title of the file is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // conversion into an integer--this is the path number
	data := openGZ(filepath)
	text := string(data)
	tiles := strings.Split(text, "\n")
	var tileData [][]string
	tileData = make([][]string, 0, 1)
	for _, line := range tiles {
		if line != "" {
			tileData = append(tileData, strings.Split(line, ","))
		}
		
	}
	library[hexNumber] = make(map[int]*TilesInStep) // later should be generalized to any library
	for _, oneTile := range tileData {
		lengthInHex, err := strconv.ParseInt(strings.Split(oneTile[0], "+")[1], 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		length := int(lengthInHex)
		newTile := TileCreator(oneTile[2], length, "")
		stepInHex, err2 := hex.DecodeString(strings.Split(oneTile[0], ".")[2])
		if err2 != nil {
			log.Fatal(err2)
		}
		step := 256 * int(stepInHex[0]) + int(stepInHex[1])
		AddTile(hexNumber, step, newTile, library)
	}


	splitpath, data, tiles, tileData =  nil, nil, nil, nil // Clears most things in memory that were used here.
}
*/

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
// Could save space by just putting it in a []*KnownVariants instead of an entire library
func AddPathFromDirectories(library *Library, directories []string, genomePath int, libraryTextFile string) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%04x",genomePath))
	b.WriteString(".fj.gz")
	for _, directory := range directories {
		bufferedTileRead(path.Join(directory,b.String()),libraryTextFile,library)
	}
}

func fjtMakeSGLFFromGenomes(directoryOfGenomes, directoryToWriteTo, tagsetFile string, genomePath int) {
	// commands:
	// samtools faidx ~/keep/by_id/cd9ada494bd979a8bc74e6d59d3e8710+174/tagset.fa.gz 035e.00 | egrep -v '^>' | tr -d '\n' | fold -w 24 > tagset_035e.00
	// docker run -i -v /home/jeremyfchen:/mnt arvados/l7g find /mnt/keep/by_id/6a3b88d7cde57054971eeabe15639cf8+263878/ -name 035e.fj.gz | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g xargs -n1 zcat | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g fjt -C -U | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g fjcsv2sglf /mnt/tagset_035e.00 | bgzip -c > 035e.sglf.gz
	var b strings.Builder

	b.WriteString("samtools faidx ")
	b.WriteString(tagsetFile)
	b.WriteString(fmt.Sprintf(" %04x.00", genomePath))
	b.WriteString(" | egrep -v '^>' | tr -d '\\n' | fold -w 24 > ")
	b.WriteString(fmt.Sprintf("tempTagset_%04x", genomePath))
	createTempTagset:= b.String()
	createTagSetCommandList := strings.Split(createTempTagset, " ")
	b.Reset()
	b.WriteString("docker run -i -v /home/jeremyfchen:/mnt arvados/l7g find ")
	b.WriteString(directoryOfGenomes)
	b.WriteString(fmt.Sprintf(" -name %04x.fj.gz | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g xargs -n1 zcat", genomePath))
	b.WriteString(" | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g fjt -C -U | docker run -i -v /home/jeremyfchen:/mnt arvados/l7g fjcsv2sglf ")
	b.WriteString(fmt.Sprintf("/mnt/tempTagset_%04x | bgzip -c > %04x.sglf.gz", genomePath, genomePath))
	createSGLF := b.String()
	createSGLFCommandList := strings.Split(createSGLF, " ")
	createTempTagsetCmd := exec.Command(createTagSetCommandList[0], createTagSetCommandList[1:]...)
	output, cmdErr := (*createTempTagsetCmd).CombinedOutput()
	if cmdErr != nil {
		fmt.Println(createTempTagset)
		fmt.Println(string(output))
	}
	createSGLFCmd := exec.Command(createSGLFCommandList[0], createSGLFCommandList[1:]...)
	_, cmdErr = (*createSGLFCmd).Output()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	
}

func fjtAddAllFastJs(directoryOfGenomes, directoryToWriteTo, tagsetFile string) {
	for i:=0; i<structures.Paths; i++ {
		fjtMakeSGLFFromGenomes(directoryOfGenomes, directoryToWriteTo, tagsetFile, i)
	}
}
/*
func AddLibrarySGLFbyFJT(library *Library, directoryToWriteTo string, path int) {
	fjtMakeSGLFFromGenomes((*library).Data, directoryToWriteTo, (*library).Tagset, path)
}
*/

// AddByDirectories adds information from a list of directories for genomes into a library, but parses by path.
func AddByDirectories(library *Library, directories []string, libraryTextFile string) {
	for path := 0; path < structures.Paths; path++ {
		AddPathFromDirectories(library, directories, path, libraryTextFile)
	}
}

// InitializeLibrary sets up the basic structure for a library.
// For consistency, it's best to use an absolute path for the text file.
func InitializeLibrary(textFile string, componentLibraries []string) Library {
	var newLibraryPaths []concurrentPath
	newLibraryPaths = make([]concurrentPath, structures.Paths, structures.Paths)
	for i := range newLibraryPaths {
		var newLock sync.RWMutex
		newLibraryPaths[i] = concurrentPath{newLock, make([]*KnownVariants, 0, 1)} // Lock is copied, but hasn't been used yet, so this is fine.
	}
	return Library{newLibraryPaths, textFile, componentLibraries}
}

// Function to copy the contents of the source library into the new library.
// Note that the normal copy function can't be used on the entire library, since locks should not be copied after use.
func libraryCopy(destination, source *Library) {
	for i := range source.Paths {
		source.Paths[i].Lock.RLock() // Locked since we're reading from path.Variants when we copy.
		destination.Paths[i].Variants = make([]*KnownVariants, len(source.Paths[i].Variants)) // This is to make sure all elements are copied over.
		copy(destination.Paths[i].Variants, source.Paths[i].Variants)
		source.Paths[i].Lock.RUnlock()
	}
}

// Function to merge the first library into the second library.
// Need to fix references--some indices are out of bounds.
func mergeLibraries(text string, libraryToMerge *Library, mainLibrary *Library) (*Library, [][][][]int) {
	var listOfReferences [][][][]int
	listOfReferences = make([][][][]int, structures.Paths, structures.Paths)
	for i := range listOfReferences {
		listOfReferences[i] = make([][][]int, 0, 1)
	}
	newLibrary := InitializeLibrary(text, []string{libraryToMerge.Text, mainLibrary.Text})
	libraryCopy(&newLibrary, mainLibrary)
	for i := range (*libraryToMerge).Paths {
		for j := range (*libraryToMerge).Paths[i].Variants {
			listOfReferences[i] = append(listOfReferences[i], MergeKnownVariants(i, j, (*libraryToMerge).Paths[i].Variants[j], &newLibrary))
		}
	}
	type referenceSortStruct struct { // Temporary struct that groups together the variant, count, and references for sorting purposes.
		variant *structures.TileVariant
		count int
		references []int
	}
	for pathNumber := range newLibrary.Paths { // Sorting step.
		newLibrary.Paths[pathNumber].Lock.Lock()
		for step, steplist := range newLibrary.Paths[pathNumber].Variants {
			if steplist != nil {
				var referenceSortStructList []referenceSortStruct
				referenceSortStructList = make([]referenceSortStruct, len((*steplist).List), len((*steplist).List))
				for k:=0; k<len((*steplist).List); k++ {
					newReferences := make([]int, len(listOfReferences[pathNumber][step][k]))
					copy(newReferences, listOfReferences[pathNumber][step][k])
					referenceSortStructList[k] = referenceSortStruct{(*steplist).List[k], (*steplist).Counts[k], newReferences}
				}
				sort.Slice(referenceSortStructList, func(i, j int) bool { return referenceSortStructList[i].count > referenceSortStructList[j].count })
				for l:=0; l<len((*steplist).List); l++ {
					(*steplist).List[l], (*steplist).Counts[l], listOfReferences[pathNumber][step][l]= referenceSortStructList[l].variant, referenceSortStructList[l].count, referenceSortStructList[l].references
				}
			}
		}
		newLibrary.Paths[pathNumber].Lock.Unlock()
	}
	return &newLibrary, listOfReferences
}

// MergeKnownVariants puts the contents of a KnownVariants at a specific path and step into another library.
// Account for CGF files here (try to avoid potential remapping with CGF files)
// Should create a new library and point the old libraries to this library.
func MergeKnownVariants(genomePath, step int, variantsToMerge *KnownVariants, newLibrary *Library) [][]int {
	var references [][]int
	originalLibraryLength := len((*newLibrary).Paths[genomePath].Variants[step].List)
	references = make([][]int, 0, 1)
	for i:=0; i<originalLibraryLength; i++ {
		references = append(references, []int{-1,i})
	}
	newTileCounter := 0
	for i, variant := range variantsToMerge.List {
		if index := TileExists(genomePath, step, variant, newLibrary); index==-1 {
			AddTileUnsafe(genomePath, step, variant.LookupReference, variant, "", newLibrary)
			(*newLibrary).Paths[genomePath].Variants[step].Counts[originalLibraryLength+newTileCounter] += variantsToMerge.Counts[i]-1
			newTileCounter++
			references = append(references, []int{i,-1})
		} else {
			(*newLibrary).Paths[genomePath].Variants[step].Counts[index] += variantsToMerge.Counts[i]
			references[index][0] = i
		}
	}
	return references
}
// TODO: store the references somewhere.


// LiftoverMapping is a representation of a liftover from one library to another.
// If a = LiftoverMapping.Mapping[b][c][d], then in path b, step c, variant d in the first library maps to variant a in the second.
type LiftoverMapping struct {
	Mapping [][][]int 
	SourceLibrary *Library // The source library to map from.
	DestinationLibrary *Library // The destination library to map to.
}
// TODO: store mappings.

// createMapping creates a liftover mapping from the source library to the destination library.
// Other way is to sort destination by reference number according to source, which takes O((m+n)log(m+n)) time, but probably a higher coefficient.
// Can check later which way is faster in practice.
func createMapping(source, destination *Library) LiftoverMapping {
	index := -1
	for i, libraryString := range destination.Components {
		if libraryString == source.Text {
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

// Function to put SGLFv2 data back into a library.
// TODO: check to make sure this and creation of SGLFv2 files works in a test function.
func parseSGLFv2(filepath string, library *Library) {
	file := path.Base(filepath)
	splitpath := strings.Split(file, ".") 
	if len(splitpath) != 3 {
		log.Fatal(errors.New("error: Not a valid gzipped file")) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "sglf" || splitpath[2] != "gz" {
		log.Fatal(errors.New("error: not a gzipped sglf file")) // Makes sure that the file is an SGLF file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name")) // Makes sure the title of the file is four digits of hexadecimal
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // conversion into an integer--this is the path number
	data := structures.OpenGZ(filepath)
	text := string(data)
	tiles := strings.Split(text, "\n")
	referenceCounter := 0
	for _, line := range tiles {
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
		newVariant := structures.TileVariant{hashArray, int(length), "", 28+referenceCounter}
		AddTile(hexNumber, int(step), 28+referenceCounter, &newVariant, filepath, fields[2], library)
		(*library).Paths[hexNumber].Variants[int(step)].Counts[TileExists(hexNumber, int(step), &newVariant, library)] += int(count-1)
		referenceCounter += len(line)
	}
}

// The following main function is only used for testing speed and memory usage of these structures.
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

// Note: may want to take care of assigned consecutive goroutines with different paths, so as not to lock out other goroutines.
func main() {
	log.SetFlags(log.Llongfile)
	flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	var m runtime.MemStats
	fmt.Println("Starting timer...")
	startTime := time.Now()
	//fjtMakeSGLFFromGenomes("/mnt/keep/by_id/6a3b88d7cde57054971eeabe15639cf8+263878/", "l7g/go/experimental/tile-library-architecture", "~/keep/by_id/cd9ada494bd979a8bc74e6d59d3e8710+174/tagset.fa.gz", 862)
	l:=InitializeLibrary("/data-sdc/jc/tile-library/test.txt", []string{})
	/*
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM/035e.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM/035e.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM/035e.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM/035e.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM/035e.fj.gz", "testing/test.txt",&l)
	//readTime := time.Now()
	*/
	/*
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
 	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	*/
	
	AddByDirectories(&l,[]string{"../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM",
	"../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM",
	"../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM",
	"../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM",
	"../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM"},
	"/data-sdc/jc/tile-library/test.txt")
	
	sortLibrary(&l)
	//writePathToSGLF(&l, 862, 0, "testing", "testing", "test.txt")
	//writePathToSGLF(&l, 862, 0, "testing2", "testing2", "test.txt")
	WriteLibraryToSGLF(&l, 0, "/data-sdc/jc/tile-library", "/data-sdc/jc/tile-library", "test.txt")
	finishTime := time.Now()
	runtime.ReadMemStats(&m)
	fmt.Printf("Total time: %v\n", finishTime.Sub(startTime))
	//fmt.Printf("Read time: %v\n", readTime.Sub(startTime))
	//fmt.Printf("Write and sort time: %v\n", finishTime.Sub(readTime))
	fmt.Println(m)
	if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.WriteHeapProfile(f)
        f.Close()
        return
    }
}