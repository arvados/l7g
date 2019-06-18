package main


import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
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
)

// TileVariant is a struct for representing a tile variant using the hash, length, and any annotation(s) on the variant.
type TileVariant struct {
	Hash       VariantHash // md5sum of the tile variant
	Length     int            // length (span) of the tile
	Annotation string         // any notes about this tile (by default, no comments)
}

// Hash for a tile variant--currently MD5
type VariantHash [md5.Size]byte

// Equals checks for equality of variants based on md5sum.
// This works based on the assumption that no two tiles in the same path and step have the same MD5.
func (t TileVariant) Equals(t2 TileVariant) bool {
	return (t.Hash == t2.Hash)
}

// KnownVariants is a struct to hold the known variants in a specific step.
type KnownVariants struct {
	List []TileVariant          // List to keep track of relative tile ordering (implicitly assigns tile variant numbers by index)
	Counts []int // Counts of each variant so far
	LookupTable VariantLookupTable // The original position of each variant in the List (for reference to text files later)
}

// A type for looking up the original positions of variants in a list--for now, the table is a list
type VariantLookupTable []int

const paths int = 863 // Constant because we know that the number of paths is always the same.

// Genome is a type to represent a genome, through its paths. Two phases are present here (path and counterpart path).
type Genome [][]Path

// Path is a type to represent a path, through its steps.
type Path []Step

// Step is a type to represent a step within a path, which can take on a specific tile variant.
type Step struct {
	Skipped bool // determines if this step has been skipped (due to a spanning tile)
	Variant *TileVariant // the Variant in this step
}

// Library is a type to represent a library of tile variants.
type Library []map[int]*KnownVariants

// Function to sort the library once all initial genomes are done being added.
// This function should only be used once during initial setup of the library, after all tiles have been added, since it sorts everything.
func sortLibrary(library *Library) {
	type sortStruct struct { // Temporary struct for sorting.
		Variant TileVariant
		Count int
		LookupReference int
	}
	for _, steps := range (*library) {
		for _, steplist := range steps {
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

// TileExists is a function to check if a specific tile exists at a specific path and step.
// Returns the index, if found
func TileExists(path, step int, toCheck TileVariant, library *Library) int {
	if (*library)[path][step] != nil { // safety to make sure that the KnownVariants struct has been created
		for i, value := range (*library)[path][step].List {
			if toCheck.Equals(value) {
				return i
			}
		}
		return -1
	}
	newKnownVariants := &KnownVariants{make([]TileVariant, 0, 1), make([]int, 0, 1), make([]int, 0, 1)}
	(*library)[path][step] = newKnownVariants
	return -1
}

// AddTile is a function to add a tile (without sorting).
func AddTile(path, step int, new TileVariant, library *Library) {
	if index := TileExists(path, step, new, library); index == -1 { // maybe not necessary, but good for safety
		(*library)[path][step].List = append((*library)[path][step].List, new)
		(*library)[path][step].Counts = append((*library)[path][step].Counts, 1)
		(*library)[path][step].LookupTable = append((*library)[path][step].LookupTable, len((*library)[path][step].LookupTable))
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
func FindFrequency(path, step int, toFind TileVariant, library *Library) int {
	if index:= TileExists(path, step, toFind, library); index != -1 {
		return (*library)[path][step].Counts[index]
	}
	fmt.Println("Variant not found.")
	return 0
}

// Annotate is a method to annotate (or re-annotate) a Tile at a specific path and step. If no match is found, the user is notified.
func Annotate(path, step int, toAnnotate TileVariant, library *Library) {
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


// TileCreator is a small function to create a new tile given information about it.
func TileCreator(hash VariantHash, length int, annotation string) TileVariant {
	return TileVariant{hash, length, annotation}
}


// The following files parse gzipped FastJ and SGLF files.
/*
func parseFastJGenome(filepath string, genome *Genome) {
	file := path.Base(filepath) // the name of the file.
	splitpath := strings.Split(file, ".")
	if len(splitpath) != 3 {
		log.Fatal(errors.New("error: Not a valid gzipped file")) // Makes sure that the filepath goes to a valid file
	}
	if splitpath[1] != "fj" || splitpath[2] != "gz" {
		log.Fatal(errors.New("error: not a gzipped FastJ file")) // Makes sure that the file is a FastJ file
	}
	pathHex, hexErr := hex.DecodeString(splitpath[0])
	if len(pathHex) != 2 || hexErr != nil {
		log.Fatal(errors.New("invalid hex file name"))
	}
	hexNumber := 256*int(pathHex[0])+int(pathHex[1]) // conversion into an integer--this is the path
	data := openGZ(filepath)
	text := string(data)
	tiles := strings.Split(text, "\n\n") // since the only divider between each tile is two newlines, this works
	var tileData [][]string
	tileData = make([][]string, 0, 1)
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") {
			tileData = append(tileData, strings.Split(line, "\n")) // within each tile the top information is separated by a newline, and the bases are separated by newlines--need to join the bases together
		}
	}
	for _, oneTile := range tileData {
		tileInfo := strings.Split(oneTile[0], ",")
		tileLength, err := strconv.Atoi(strings.Split(tileInfo[5], ":")[1])
		if err != nil {
			log.Fatal(err)
		}
		tilePathStep := strings.Split(tileInfo[0],":")[1]
		stepInHex, err2 := hex.DecodeString(strings.Split(tilePathStep, ".")[2])
		if err2 != nil {
			log.Fatal(err2)
		}
		step := 256 * int(stepInHex[0]) + int(stepInHex[1])
		phase, err3 := strconv.Atoi(strings.TrimSuffix(strings.Split(tilePathStep, ".")[3], "\""))
		if err3 != nil {
			log.Fatal(err3)
		}
		tileBases := strings.Join(oneTile[1:], "")
		newTile := TileCreator(tileBases, tileLength, "")
		for len(genome[hexNumber][phase]) < step {
			genome[hexNumber][phase] = append(genome[hexNumber][phase], Step{true, nil}) // This adds empty (skipped) steps until we reach the right step number.
		}
		genome[hexNumber][phase] = append(genome[hexNumber][phase],Step{false, &newTile})
	}
	splitpath, data, tiles, tileData =  nil, nil, nil, nil // Clears most things in memory that were used here.
}
*/

func writeToTextFile(genomePath, step int, directory, bases string, hash VariantHash) {
	err := os.MkdirAll(path.Join(directory,fmt.Sprintf("%04x", genomePath)), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%40x",step))
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

}

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
			bufferedWriter.WriteString(fmt.Sprintf("%02v", version))
			bufferedWriter.WriteString(".")
			bufferedWriter.WriteString(step)
			bufferedWriter.WriteString(".")
			bufferedWriter.WriteString(strconv.Itoa(index))
			bufferedWriter.WriteString("+")
			bufferedWriter.WriteString(strconv.Itoa((*library)[genomePath][stepInt].List[index].Length))
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(tiles[(*library)[genomePath][stepInt].LookupTable[index]])
		}
	}
}

func parseFastJLibrary(filepath string, library *Library) {
	file := path.Base(filepath) // the name of the file.
	fmt.Println("Starting path", file)
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

	data := openGZ(filepath)
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
			var hashArray VariantHash
			copy(hashArray[:], hash)
			lengthString := strings.Split(strings.SplitN(line, ",", 7)[5], ":")[1]
			length, err3 := strconv.Atoi(lengthString)
			if err3 != nil {
				log.Fatal(err3)
			}
			newTile := TileCreator(hashArray, length, "")
			AddTile(hexNumber, step, newTile, library)
		}
	}
	splitpath, data, tiles=  nil, nil, nil // Clears most things in memory that were used here.
	fmt.Println("Finished path ", hexNumber)
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

// creates a Genome based on a directory containing all of the needed FastJ files
func createGenome(directory string) Genome {
	var newGenome Genome
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") {
			parseFastJGenome(path.Join(directory, file.Name()), &newGenome)
		}
	}
	return newGenome
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
			parseFastJLibrary(path.Join(directory, file.Name()), library)
		}
	}
}

// Function to open a gzipped file.
func openGZ(filepath string) []byte {
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

// Function to initialize a new Library.
func initializeLibrary() Library {
	var newLibrary Library
	newLibrary = make([]map[int]*KnownVariants, paths, paths)
	for i := range newLibrary {
		newLibrary[i] = make(map[int]*KnownVariants)
	}
	return newLibrary
}


// The following main function is only used for testing speed and memory usage of these structures.

func main() {
	var m runtime.MemStats
	fmt.Println("Starting timer...")
	startTime := time.Now()
	l:=initializeLibrary()
	addLibraryFastJ("../../../../keep/home/tile-library-architecture/Copy of Container output for request su92l-xvhdp-qc9aol66z8oo7ws/hu03E3D2_masterVarBeta-GS000038659-ASM", &l)
	sortLibrary(&l)
	total := time.Since(startTime)
	runtime.ReadMemStats(&m)
	fmt.Printf("Total time: %v\n", total)
	fmt.Println(m)

}