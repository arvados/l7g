package main // should be changed to package tile-library or package tilelibrary

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"../structures" // try to avoid relative paths.
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
				sort.Slice(sortStructList, func(i, j int) bool { return sortStructList[i].Count > sortStructList[i].Count })
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
func AddTile(path, step int, new structures.TileVariant, bases string, library *Library) {
	if index := TileExists(path, step, new, library); index == -1 { // Checks if the tile exists already.
		(*library)[path][step].List = append((*library)[path][step].List, new)
		(*library)[path][step].Counts = append((*library)[path][step].Counts, 1)
		(*library)[path][step].LookupTable = append((*library)[path][step].LookupTable, len((*library)[path][step].LookupTable))
		//writeToTextFile(path, step, "testing", bases, new.Hash)
		// TODO: implement writing to any directory name
	} else {
		(*library)[path][step].Counts[index]++
	}
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
			annotation, _ := readKeyboard.ReadString('\n')
			tile.Annotation = annotation
			break
		}
	}
	fmt.Printf("No matching tile found at specified path %v and step %v.\n", path, step) // Information if tile isn't found.
}

// writeToTextFile writes the entry of a lookup from a hash to bases for a specific path and step, in a text file.
// To be replaced with using Docker to use Abram's tools to create SGLF files (so more intermediate files aren't necessary)
func writeToTextFile(genomePath, step int, directory, bases string, hash structures.VariantHash) {
	err := os.MkdirAll(path.Join(directory,fmt.Sprintf("%04x", genomePath)), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%04x",step))
	b.WriteString(".txt")
	textFile, err2 := os.OpenFile(path.Join(directory,fmt.Sprintf("%04x", genomePath),b.String()), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}
	bufferedWriter := bufio.NewWriter(textFile)
	
	b.Reset()
	b.WriteString(hex.EncodeToString(hash[:]))
	b.WriteString(",")
	b.WriteString(bases)
	b.WriteString("\n")
	_, err3 := bufferedWriter.WriteString(b.String())
	if err3 != nil {
		log.Fatal(err3)
	}
	bufferedWriter.Flush()

}

// writePathToSGLF writes an SGLF for an entire path given a library.
func writePathToSGLF(library *Library, genomePath, version int, directoryToWriteTo, directoryToGetFrom string) {
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
	textFiles, err3 := ioutil.ReadDir(path.Join(directoryToGetFrom,fmt.Sprintf("%04x", genomePath)))
	if err3 != nil {
		log.Fatal(err3)
	}
	for _, file := range textFiles {
		lines, err4 := os.Open(path.Join(directoryToGetFrom, fmt.Sprintf("%04x", genomePath), file.Name()))
		if err4 != nil {
			log.Fatal(err4)
		}
		scanner := bufio.NewScanner(lines)
		var tiles []string
		tiles = make([]string, 0, 1)
		for scanner.Scan() {
			tiles = append(tiles, scanner.Text())
		}
		for index := range tiles {
			step := strings.Split(file.Name(),".")[0]
			stepInHex, hexErr := hex.DecodeString(step)
			if hexErr != nil {
				log.Fatal(hexErr)
			}
			stepInt := 256 * int(stepInHex[0]) + int(stepInHex[1])
			bufferedWriter.WriteString(fmt.Sprintf("%04x", genomePath)) // path
			bufferedWriter.WriteString(".")
			bufferedWriter.WriteString(fmt.Sprintf("%02v", version)) // Version
			bufferedWriter.WriteString(".")
			bufferedWriter.WriteString(step) // Step
			bufferedWriter.WriteString(".")
			bufferedWriter.WriteString(fmt.Sprintf("%03x", index)) // Tile variant number
			bufferedWriter.WriteString("+")
			bufferedWriter.WriteString(fmt.Sprintf("%01x", (*library)[genomePath][stepInt].List[index].Length)) // Tile length
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(tiles[(*library)[genomePath][stepInt].LookupTable[index]]) // Hash and bases of tile.
			bufferedWriter.WriteString("\n") // New tile.
		}
	}
	bufferedWriter.Flush()
}

// Function to write the contents of a library to SGLF files.
func writeLibraryToSGLF(library *Library, version int, directoryToWriteTo, directoryToGetFrom string) {
	for path := 0; path < structures.Paths; path++ {
		writePathToSGLF(library, path, version, directoryToWriteTo, directoryToGetFrom)
	}
}

// ParseFastJLibrary puts the contents of a (gzipped) FastJ into a Library.
func ParseFastJLibrary(filepath string, library *Library) {
	file := path.Base(filepath) // the name of the file.
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
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // conversion into an integer--this is the path

	data := structures.OpenGZ(filepath)
	text := string(data)
	tiles := strings.Split(text, "\n\n") // since the only divider between each tile is two newlines, this works
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
					lengthString=string(line[i-1]) // account for the possibility of length being at least 16
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
			AddTile(hexNumber, step, newTile, bases, library)
		}
	}
	splitpath, data, tiles=  nil, nil, nil // Clears most things in memory that were used here.

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



// Function to add a directory of gzipped FastJ files to a specific library. 
func addLibraryFastJ(directory string, library *Library) {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") {
			ParseFastJLibrary(path.Join(directory, file.Name()), library)
		}
	}
}

// addPathFromDirectories parses the same path for all genomes, represented by a list of directories, and puts the information in a Library.
// Could save space by just putting it in a []*KnownVariants instead of an entire library
func addPathFromDirectories(library *Library, directories []string, genomePath int) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%04x",genomePath))
	b.WriteString(".fj.gz")
	for _, directory := range directories {
		ParseFastJLibrary(path.Join(directory,b.String()),library)
	}
}

// addByDirectories adds information from a list of directories for genomes into a library, but parses by path.
func addByDirectories(library *Library, directories []string) {
	for path := 0; path < structures.Paths; path++ {
		addPathFromDirectories(library, directories, path)
	}
}



// Function to initialize a new Library.
func initializeLibrary() Library {
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
			mergeKnownVariants(filepathToMerge, i, j, (*libraryToMerge)[i][j], mainLibrary)
		}
	}
}

// Function to merge a KnownVariants at a specific path and step into another library.
// Account for CGF files here (try to avoid potential remapping with CGF files)
func mergeKnownVariants(filepathToMerge string, genomePath, step int, variantsToMerge *KnownVariants, mainLibrary *Library) {
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
		AddTile(genomePath, step, variant, tiles[(*variantsToMerge).LookupTable[i]], mainLibrary)
		(*mainLibrary)[genomePath][step].Counts[TileExists(genomePath, step, variant, mainLibrary)] += (*variantsToMerge).Counts[i]-1
	}
}


// The following main function is only used for testing speed and memory usage of these structures.
// Speed and heap allocation usage: 3-3.5 minutes, 1.5-2.5GB?
// time to make one sglf file for path 24: 1.5 seconds--at this rate would take around 20-22 minutes per 5 genomes, but would probably be less in practice

// time and space to go through 5 genomes by path: 22-23 minutes, 4.5GB
// time and space to go through 5 genomes by directory: 21 minutes, 3.5GB
func main() {
	var m runtime.MemStats
	fmt.Println("Starting timer...")
	startTime := time.Now()
	l:=initializeLibrary()
	addLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", &l)
	addLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM", &l)
	addLibraryFastJ("../../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM", &l)
	//addLibraryFastJ("../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM", &l)
	//addLibraryFastJ("../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM", &l)
	//addByDirectories(&l,[]string{"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000037847-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu01F73B_masterVarBeta-GS000037833-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu02C8E3_masterVarBeta-GS000036653-ASM",
	//"../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu0486D6_masterVarBeta-GS000037846-ASM"})
	
	//writePathToSGLF(&l, 24, 0, "sglf", "testing")
	sortLibrary(&l)
	
	total := time.Since(startTime)
	runtime.ReadMemStats(&m)
	fmt.Printf("Total time: %v\n", total)
	fmt.Println(m)
}