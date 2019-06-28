package main // should be changed to package tile-library or package tilelibrary

// This tile library package assumes that any necessary imputation was done beforehand.

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"../structures" // try to avoid relative paths. If possible, move to github.
)

// KnownVariants is a struct to hold the known variants in a specific step.
type KnownVariants struct {
	List [](structures.TileVariant)         // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index after sorting)
	Counts []int // Counts of each variant so far
	LookupTable VariantLookupTable // The original position of each variant in the List (for reference to text files later)
}

// VariantLookupTable is a type for looking up the original positions of variants in a list--for now, the table is a list
type VariantLookupTable []int

// Library is a type to represent a library of tile variants.
// The first slice represents paths, and the second slice represents steps.
type Library [][]*KnownVariants

/*
type Library struct {
	AllVariants [][]*KnownVariants
	Data string // the directory where all the data was collected from
	Tagset string // The tagset for this library
}
*/

// Function to sort the library once all initial genomes are done being added.
// This function should only be used once during initial setup of the library, after all tiles have been added, since it sorts everything.
func sortLibrary(library *Library) {
	type sortStruct struct { // Temporary struct for sorting.
		Variant structures.TileVariant
		Count int
		LookupReference int
	}
	for _, steps := range (*library) {
		for _, steplist := range steps {
			if steplist != nil {
				var sortStructList []sortStruct
				sortStructList = make([]sortStruct, len((*steplist).List))
				for i:=0; i<len((*steplist).List); i++ {
					sortStructList[i] = sortStruct{(*steplist).List[i], (*steplist).Counts[i], (*steplist).LookupTable[i]}
				}
				sort.Slice(sortStructList, func(i, j int) bool { return sortStructList[i].Count > sortStructList[j].Count })
				for j:=0; j<len((*steplist).List); j++ {
					(*steplist).List[j], (*steplist).Counts[j], (*steplist).LookupTable[j] = sortStructList[j].Variant, sortStructList[j].Count, sortStructList[j].LookupReference
				}
			}
		}
	}
}

// TileExists is a function to check if a specific tile exists at a specific path and step in a library.
// Returns the index of the variant, if found--otherwise, returns -1.
func TileExists(path, step int, toCheck structures.TileVariant, library *Library) int {
	if len((*library)[path]) > step && (*library)[path][step] != nil { // Safety to make sure that the KnownVariants struct has been created
		for i, value := range (*library)[path][step].List {
			if toCheck.Equals(value) {
				return i
			}
		}
		return -1
	}
	for len((*library)[path]) <= step {
		(*library)[path] = append((*library)[path], nil)
	}
	newKnownVariants := &KnownVariants{make([](structures.TileVariant), 0, 1), make([]int, 0, 1), make([]int, 0, 1)}
	(*library)[path][step] = newKnownVariants
	return -1
}

// AddTile is a function to add a tile (without sorting).
func AddTile(genomePath, step, lookupNumber int, new structures.TileVariant, libraryTextFile, bases string, library *Library) {
	if index := TileExists(genomePath, step, new, library); index == -1 { // Checks if the tile exists already.
		(*library)[genomePath][step].List = append((*library)[genomePath][step].List, new)
		(*library)[genomePath][step].Counts = append((*library)[genomePath][step].Counts, 1)
		(*library)[genomePath][step].LookupTable = append((*library)[genomePath][step].LookupTable, lookupNumber)
		//writeToTextFile(genomePath, step, path.Dir(libraryTextFile), bases, path.Base(libraryTextFile), new.Hash)
	} else {
		(*library)[genomePath][step].Counts[index]++
	}
}

// AddTileUnsafe is a function to add a tile without sorting.
// Unsafe because it doesn't check if the tile is already in the library, unlike AddTile.
func AddTileUnsafe(genomePath, step, lookupNumber int, new structures.TileVariant, libraryTextFile, bases string, library *Library) {
	
	(*library)[genomePath][step].List = append((*library)[genomePath][step].List, new)
	(*library)[genomePath][step].Counts = append((*library)[genomePath][step].Counts, 1)
	(*library)[genomePath][step].LookupTable = append((*library)[genomePath][step].LookupTable, lookupNumber)

}
/*
// AddAndSortTiles takes a list of tiles to put into a path and step, and adds them all at once (and sorts afterwards).
func AddAndSortTiles(path, step int, newTiles []TileVariant, library *Library) {
	for _, tile := range newTiles {
		AddTile(path, step, tile, library)
	}
	toSort := library[path][step]
	sort.Slice(toSort.List, func(i, j int) bool { return toSort.Set[*toSort.List[i].MD5] > toSort.Set[*toSort.List[j].MD5] })
}


// AddPath is a function to add an entire Path to a Library all at once.
func AddPath(pathNumber int, path Path, library *Library) {
	for step, value := range path {
		if !value.Skipped {
			AddTile(pathNumber,step,*value.Variant, library)
		}
	}
}

// AddGenome is a function to add an entire Genome to a Library all at once.
func AddGenome(genome Genome, library *Library) {
	for pathNumber, paths := range genome {
		AddPath(pathNumber, paths[0], library)
		AddPath(pathNumber, paths[1], library)
	}
}
*/

// FindFrequency is a function to find the frequency of a specific tile at a specific path and step.
func FindFrequency(path, step int, toFind structures.TileVariant, library *Library) int {
	if index:= TileExists(path, step, toFind, library); index != -1 {
		return (*library)[path][step].Counts[index]
	}
	fmt.Println("Variant not found.")
	return 0
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified.
func Annotate(path, step int, toAnnotate structures.TileVariant, library *Library) {
	for _, tile := range (*library)[path][step].List {
		if toAnnotate.Equals(tile) {
			fmt.Print("Enter annotation: ")
			readKeyboard := bufio.NewReader(os.Stdin)
			annotation, err := readKeyboard.ReadString('\n')
			if err!=nil {
				log.Fatal(err)
			}
			tile.Annotation = annotation
			break
		}
	}
	fmt.Printf("No matching tile found at specified path %v and step %v.\n", path, step) // Information if tile isn't found.
}


// TODO: make a goroutine do the writing to disk while information relevant to the library structure is added, and add a temporary index number to tiles that have been processed but not written yet.
// fine to read and write concurrently?
// in addition, buffer bases so that there are fewer writes to disk

// Maybe store everything from a path in data, and then write and add everything all at once? Need to check if there are duplicates.


func bufferedTileRead(fastJFilepath, libraryTextFile string, library *Library) {
	var baseChannel chan string
	var startingIndices chan int
	baseChannel = make(chan string, 16) // Put strings of bases in here while they need to be processed.
	startingIndices = make(chan int, 16) // Put pointer indices here while waiting to be processed.

	go bufferedBaseWrite(libraryTextFile, baseChannel, startingIndices)

	file := path.Base(fastJFilepath) // The name of the file.
	splitpath := strings.Split(file, ".")
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

	data := structures.OpenGZ(fastJFilepath)
	text := string(data)
	tiles := strings.Split(text, "\n\n") // The divider between two tiles is two newlines.
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") { // Makes sure that a "line" starts with the correct character ">"
			stepInHex := line[20:24]
			stepBytes, err := hex.DecodeString(stepInHex)
			if err != nil {
				log.Fatal(err)
			}
			step := 256 * int(stepBytes[0]) + int(stepBytes[1])
			hashString := line[40:72]
			hash, err2 := hex.DecodeString(hashString)
			if err2 != nil {
				log.Fatal(err2)
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
					lengthString=string(line[i-1]) // TODO: need to account for the possibility of length being at least 16
					break
				}
			}
			
			length, err3 := strconv.Atoi(lengthString)
			if err3 != nil {
				log.Fatal(err3)
			}
			baseData := strings.Split(line, "\n")[1:]
			var b strings.Builder
			for _, data := range baseData {
				if data != "\n" {
					b.WriteString(data)
				}
			}
			bases := b.String()
			newTile := structures.TileCreator(hashArray, length, "")
			

			if TileExists(hexNumber, step, newTile, library)==-1 {
				baseChannel <- bases
				index := <- startingIndices
				AddTileUnsafe(hexNumber, step, index, newTile, libraryTextFile, bases, library)
			}
			
		}
	}
	close(baseChannel)
	splitpath, data, tiles=  nil, nil, nil // Clears most things in memory that were used here, to free up memory.
}

func bufferedBaseWrite(libraryTextFile string, channel chan string, startingIndices chan int) {
	err := os.MkdirAll(path.Dir(libraryTextFile), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	file, err2 := os.OpenFile(libraryTextFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	bufferedWriter := bufio.NewWriter(file)
	if err2 != nil {
		log.Fatal(err)
	}
	ok := true
	for ok {
		bases, channelOk := <- channel
		ok = channelOk
		if ok {
			info, err3 := os.Stat(libraryTextFile)
			if err3 != nil {
				log.Fatal(err3)
			}
			startingIndices <- (int(info.Size())+bufferedWriter.Buffered())
			hash := md5.Sum([]byte(bases))
			hashString := hex.EncodeToString(hash[:])
			bufferedWriter.WriteString(hashString)
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(bases)
			bufferedWriter.WriteString("\n")
		}
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

	textFile, err2 := os.OpenFile(path.Join(directory,filename), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
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
func writePathToSGLF(library *Library, genomePath, version int, directoryToWriteTo, directoryToGetFrom, textFilename string) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x", genomePath))
	b.WriteString(".sglf")
	err := os.MkdirAll(directoryToWriteTo, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	sglfFile, err2 := os.OpenFile(path.Join(directoryToWriteTo,b.String()), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}
	bufferedWriter := bufio.NewWriter(sglfFile)
	textFile, err2 := os.OpenFile(path.Join(directoryToGetFrom,textFilename), os.O_RDONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}
	for step, variants := range (*library)[genomePath] {
		if variants != nil {
			for i := range (*variants).List {
				textFile.Seek(int64((*variants).LookupTable[i]),0)
				fileReader := bufio.NewReader(textFile)
				tileString, err := fileReader.ReadString('\n')
				if err != nil {
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
				bufferedWriter.WriteString(fmt.Sprintf("%01x", (*library)[genomePath][step].List[i].Length)) // Tile length
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(tileString) // Hash and bases of tile.
			}
		}
		
	}
	bufferedWriter.Flush()
	textFile.Close()
}

// WriteLibraryToSGLF writes the contents of a library to SGLF files.
func WriteLibraryToSGLF(library *Library, version int, directoryToWriteTo, directoryToGetFrom, textFile string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLF(library, path, version, directoryToWriteTo, directoryToGetFrom, textFile)
	}
}

// ParseFastJLibrary puts the contents of a (gzipped) FastJ into a Library.
func ParseFastJLibrary(filepath, libraryTextFile string, library *Library) {
	file := path.Base(filepath) // The name of the file.
	splitpath := strings.Split(file, ".")
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

	data := structures.OpenGZ(filepath)
	text := string(data)
	tiles := strings.Split(text, "\n\n") // The divider between two tiles is two newlines.
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") { // Makes sure that a line starts with the correct character ">"
			stepInHex := line[20:24]
			stepBytes, err := hex.DecodeString(stepInHex)
			if err != nil {
				log.Fatal(err)
			}
			step := 256 * int(stepBytes[0]) + int(stepBytes[1])
			hashString := line[40:72]
			hash, err2 := hex.DecodeString(hashString)
			if err2 != nil {
				log.Fatal(err2)
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
					lengthString=string(line[i-1]) // TODO: need to account for the possibility of length being at least 16
					break
				}
			}
			
			length, err3 := strconv.Atoi(lengthString)
			if err3 != nil {
				log.Fatal(err3)
			}
			baseData := strings.Split(line, "\n")[1:]
			var b strings.Builder
			for _, data := range baseData {
				if data != "\n" {
					b.WriteString(data)
				}
			}
			bases := b.String()
			newTile := structures.TileCreator(hashArray, length, "")
			info, err4 := os.Stat(libraryTextFile)
			var fileLength int
			if err4 != nil {
				fileLength = 0
			} else {
				fileLength = int(info.Size())
			}
			
			AddTile(hexNumber, step, fileLength, newTile, libraryTextFile, bases, library)
		}
	}
	splitpath, data, tiles=  nil, nil, nil // Clears most things in memory that were used here, to free up memory.

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
		if strings.HasSuffix(file.Name(), ".gz") {
			ParseFastJLibrary(path.Join(directory, file.Name()), libraryTextFile, library)
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
		ParseFastJLibrary(path.Join(directory,b.String()),libraryTextFile,library)
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
	_, cmdErr2 := (*createSGLFCmd).Output()
	if cmdErr2 != nil {
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
func InitializeLibrary() Library {
	var newLibrary Library
	newLibrary = make([][]*KnownVariants, structures.Paths, structures.Paths)
	for i := range newLibrary {
		newLibrary[i] = make([]*KnownVariants, 0, 1)
	}
	return newLibrary
}



// Function to merge the first library into the second library.
func mergeLibraries(filepathToMerge string, libraryToMerge *Library, mainLibrary *Library) {
	for i, path := range (*libraryToMerge) {
		for j := range path {
			MergeKnownVariants(filepathToMerge, i, j, (*libraryToMerge)[i][j], mainLibrary)
		}
	}
}

// MergeKnownVariants puts the contents of a KnownVariants at a specific path and step into another library.
// Account for CGF files here (try to avoid potential remapping with CGF files)
func MergeKnownVariants(filepathToMerge string, genomePath, step int, variantsToMerge *KnownVariants, mainLibrary *Library) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x",step))
	b.WriteString(".txt")
	lines, err := os.Open(path.Join(filepathToMerge, fmt.Sprintf("%04x", genomePath), b.String()))
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(lines)
	var tiles []string
	tiles = make([]string, 0, 1)
	for scanner.Scan() {
		tiles = append(tiles, scanner.Text())
	}
	for i, variant := range (*variantsToMerge).List {
		//AddTile(genomePath, step, variant, tiles[(*variantsToMerge).LookupTable[i]], mainLibrary)
		(*mainLibrary)[genomePath][step].Counts[TileExists(genomePath, step, variant, mainLibrary)] += (*variantsToMerge).Counts[i]-1
	}
}


// The following main function is only used for testing speed and memory usage of these structures.
// Speed and heap allocation usage: 3-3.5 minutes, 1.5-2.5GB?
// time to make one sglf file for path 24: 1.5 seconds--at this rate would take around 20-22 minutes per 5 genomes, but would probably be less in practice

// time and space to go through 5 genomes by path: 22-23 minutes, 4.5GB
// time and space to go through 5 genomes by directory: 21 minutes, 3.5GB
// time and space to put 5 genomes in a library and write bases to a file: 32 minutes, 3.5GB
func main() {
	log.SetFlags(log.Llongfile)
	var m runtime.MemStats
	fmt.Println("Starting timer...")
	startTime := time.Now()
	//fjtMakeSGLFFromGenomes("/mnt/keep/by_id/6a3b88d7cde57054971eeabe15639cf8+263878/", "l7g/go/experimental/tile-library-architecture", "~/keep/by_id/cd9ada494bd979a8bc74e6d59d3e8710+174/tagset.fa.gz", 862)
	l:=InitializeLibrary()
	
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM/0017.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM/0017.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM/0017.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM/0017.fj.gz", "testing/test.txt",&l)
	bufferedTileRead("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM/0017.fj.gz", "testing/test.txt",&l)
	
	/*
	ParseFastJLibrary("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM/035e.fj.gz", "testing2/test.txt",&l)
	ParseFastJLibrary("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM/035e.fj.gz", "testing2/test.txt",&l)
	ParseFastJLibrary("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM/035e.fj.gz", "testing2/test.txt",&l)
	ParseFastJLibrary("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM/035e.fj.gz", "testing2/test.txt",&l)
	ParseFastJLibrary("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM/035e.fj.gz", "testing2/test.txt",&l)
	*/
	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
 	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	//AddLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM", "/data-sdc/jc/tile-library/test.txt",&l)
	//addByDirectories(&l,[]string{"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM"})
	sortLibrary(&l)
	writePathToSGLF(&l, 23, 0, "testing", "testing", "test.txt")
	//writePathToSGLF(&l, 862, 0, "testing2", "testing2", "test.txt")
	//WriteLibraryToSGLF(&l, 0, "/data-sdc/jc/tile-library", "/data-sdc/jc/tile-library", "test.txt")
	total := time.Since(startTime)
	runtime.ReadMemStats(&m)
	fmt.Printf("Total time: %v\n", total)
	fmt.Println(m)
}