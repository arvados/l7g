package genome
/*
	Genome is a package for representing the genome with Go data structures.
*/
import (
	"bufio"
	"compress/gzip"
	"errors"
	"encoding/hex"
	"io/ioutil"
	"path"
	"log"
	"os"
	"strconv"
	"strings"
	"../structures"
	"../tile-library"
	"github.com/kshedden/gonpy" //To be used later for the creation of numpy arrays.
)

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

// Genome is a struct to represent a genome. It contains a pointer to its reference library, which allows for easy tiling.
type Genome struct {
	Paths [][]Path // Paths represents a genome through its paths. Two phases are present here (path and counterpart path).
	Library *tilelibrary.Library // This is the reference library for this Genome.
}

// Path is a type to represent a path, through its steps.
type Path []Step

// Step is a type to represent a step within a path, which can take on a specific tile variant.
type Step int // -1 for a skipped step, and any other integer refers to the tile variant number in the reference library.

// isComplete determines if a set of bases is complete (has no nocalls).
// This is only a helper function for ParseFastJGenome.
func isComplete(bases string) bool {
	return !strings.ContainsRune(bases, 'n')
}

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

// ParseFastJGenome puts the contents of a (gzipped) FastJ into a Genome.
// Works with both gzipped and non-gzipped FastJ files.
func ParseFastJGenome(filepath string, genome *Genome) {
	if genome == nil || genome.Library == nil {
		log.Fatal(errors.New("genome is nil or genome library is nil, cannot parse FastJ"))
	}
	file := path.Base(filepath) // The name of the file.
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
		data = openGZ(filepath)
	} else {
		data = openFile(filepath)
	}
	text := string(data)
	tiles := strings.Split(text, "\n\n") // since the only divider between each tile is two newlines, this works
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") {
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
					k:=1 // This accounts for the possibility of the length of a tile spanning at least 16 tiles.
					for line[i-k] != ':' { // Goes back a few characters until it knows the string of the tile length.
						k++
					}
					lengthString=line[(i-k+1):i]
					break
				}
			}
			
			length, err3 := strconv.Atoi(lengthString)
			if err3 != nil {
				log.Fatal(err3)
			}
			phase, err4 := strconv.Atoi(line[25:28])
			if err4 != nil {
				log.Fatal(err4)
			}
			baseData := strings.Split(line, "\n")[1:]
			var b strings.Builder
			for _, data := range baseData {
				if data != "\n" {
					b.WriteString(data)
				}
			}
			bases := b.String()
			newTile := structures.TileVariant{Hash: hashArray, Length: length, Annotation: "", LookupReference: -1, Complete: isComplete(bases), ReferenceLibrary: genome.Library}
			// In the case of newTile, -1 is a value used to complete the creation of the tile, and has no meaning otherwise.
			if tilelibrary.TileExists(hexNumber, step, &newTile, genome.Library) == -1 {
				log.Fatal(errors.New("this FastJ is not part of the library"))
			}
			for len((genome.Paths)[hexNumber][phase]) <= step+length-1 {
				(genome.Paths)[hexNumber][phase] = append((genome.Paths)[hexNumber][phase], -1) // This adds empty (skipped) steps until we reach the right step number.
			}
			(genome.Paths)[hexNumber][phase][step] = Step(tilelibrary.TileExists(hexNumber, step, &newTile, genome.Library))
		}
	}
	splitpath, data, tiles=  nil, nil, nil // Clears most things in memory that were used here.
}

//CreateGenome puts the contents of a directory of FastJ files into a given Genome.
func CreateGenome(directory string, genome *Genome) {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") {
			ParseFastJGenome(path.Join(directory, file.Name()), genome)
		}
	}
}

// InitializeGenome is a function to initialize a Genome.
// nil is allowed for the library if the library shouldn't be set yet. It can be set manually later.
func InitializeGenome(library *tilelibrary.Library) Genome {
	var newPaths [][]Path
	newPaths = make([][]Path, structures.Paths, structures.Paths)
	for i := range newPaths {
		newPaths[i] = make([]Path, 2, 2)
		newPaths[i][0] = make([]Step, 0)
		newPaths[i][1] = make([]Step, 0)
	}
	return Genome{newPaths, library}
}


/*
skipped tiles: first value is variant, and then -1 (to show tail spanning tiles)
-2 used for incomplete tiles 
put in hash of name of creator file
*/


// WriteGenomeToFile writes a genome to a list format of indices relative to its reference library.
// Will not work if the genome does not have a reference library (nil reference)
func WriteGenomeToFile(filename string, g *Genome) {
	f, err := os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(f)
	for path := range g.Paths {
		if (*g.Library).Paths[path].Variants[0] != nil {
			if g.Paths[path][0][0] >= 0 && !(*g.Library).Paths[path].Variants[0].List[int(g.Paths[path][0][0])].Complete {
				bufferedWriter.WriteString("-2")
			} else {
				bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][0][0])))
			}
			bufferedWriter.WriteString(",")
			if g.Paths[path][1][0] >= 0 && !(*g.Library).Paths[path].Variants[0].List[int(g.Paths[path][1][0])].Complete {
				bufferedWriter.WriteString("-2")
			} else {
				bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][0])))
			}
		} else {
			bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][0][0])))
			bufferedWriter.WriteString(",")
			bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][0])))
		}
		for i, value := range g.Paths[path][0][1:] {
			if (*g.Library).Paths[path].Variants[i+1] != nil {
				bufferedWriter.WriteString(",")
				if value >=0 && !(*g.Library).Paths[path].Variants[i+1].List[int(value)].Complete {
					bufferedWriter.WriteString("-2")
				} else {
					bufferedWriter.WriteString(strconv.Itoa(int(value)))
				}
				bufferedWriter.WriteString(",")
				if g.Paths[path][1][i+1] >= 0 && !(*g.Library).Paths[path].Variants[i+1].List[int(g.Paths[path][1][i+1])].Complete {
					bufferedWriter.WriteString("-2")
				} else {
					bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][i+1])))
				}
			} else {
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(strconv.Itoa(int(value)))
				bufferedWriter.WriteString(",")
				bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][i+1])))
			}
		}
		bufferedWriter.WriteString("\n")
	}
	bufferedWriter.Flush()
}

// ReadGenomeFromFile reads a text file containing genome information.
// Current file suffix is .genome
func ReadGenomeFromFile(filepath string) [][]Path {
	var newPaths [][]Path
	newPaths = make([][]Path, structures.Paths, structures.Paths)
	for i := range newPaths {
		newPaths[i] = make([]Path, 2, 2)
		newPaths[i][0] = make([]Step, 0)
		newPaths[i][1] = make([]Step, 0)
	}
	splitpath := strings.Split(filepath, ".")
	if len(splitpath) != 2 || splitpath[1] != "genome" {
		log.Fatal(errors.New("not a valid genome file"))
	}
	info := openFile(filepath)
	lines := strings.Split(string(info), "\n")
	for i, line := range lines {
		indices := strings.Split(line, ",")
		for j, index := range indices {
			if index != "" {
				indexInt, err := strconv.Atoi(index)
				if err != nil {
					log.Fatal(err)
				}
				newPaths[i][j % 2] = append(newPaths[i][j % 2], Step(indexInt))
			}
		}
	}
	return newPaths
}

// WriteNumpy writes the values of a path of a genome to a numpy array.
// It alternate between each phase for each step.
// writes a path of a genome to a numpy file.
func (g *Genome) WriteNumpy(filepath string, path int) {
	npywriter, err := gonpy.NewFileWriter(filepath)
	if err != nil {
		log.Fatal(err)
	}
	sliceOfData := make([]int32, 0, 1)
	for i, value := range g.Paths[path][0] {
		if value >=0 && !(*g.Library).Paths[path].Variants[i].List[value].Complete {
			sliceOfData = append(sliceOfData, -2)
		} else {
			sliceOfData = append(sliceOfData, int32(value))
		}
		
		if g.Paths[path][1][i] >=0 && !(*g.Library).Paths[path].Variants[i].List[g.Paths[path][1][i]].Complete {
			sliceOfData = append(sliceOfData, -2)
		} else {
			sliceOfData = append(sliceOfData, int32(g.Paths[path][1][i]))
		}
	}
	err = npywriter.WriteInt32(sliceOfData)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadGenomeNumpy reads one path's worth of information from a numpy file.
// This path should be assigned to a path of a genome.
func ReadGenomeNumpy(filepath string) []Path {
	newPaths := make([]Path, 2)
	newPaths[0] = make([]Step, 0, 1)
	newPaths[1] = make([]Step, 0, 1)
	npyreader, err := gonpy.NewFileReader(filepath)
	if err != nil {
		log.Fatal(err)
	}
	pathInfo, err := npyreader.GetInt32()
	if err != nil {
		log.Fatal(err)
	}
	for i, index := range pathInfo {
		newPaths[i%2] = append(newPaths[i%2], Step(index))
	}
	return newPaths
}

// WriteGenomesPathToNumpy writes multiple genomes' worth of path information to a numpy file.
func WriteGenomesPathToNumpy(genomes []*Genome, filepath string, path int) {
	if len(genomes) > 0 { // Requires a nonempty list.
		npywriter, err := gonpy.NewFileWriter(filepath)
		if err != nil {
			log.Fatal(err)
		}
		lengthOfPath := 0
		for _, value := range genomes[0].Paths[path] {
			lengthOfPath += len(value)
		}
		for _, genome := range genomes {
			genomeLengthOfPath := 0
			for _, value := range genome.Paths[path] {
				genomeLengthOfPath += len(value)
			}
			if genomeLengthOfPath != lengthOfPath {
				log.Fatal(errors.New("path lengths within each genome are not equal"))
			}
		}
		sliceOfData := make([]int32, len(genomes) * lengthOfPath)
		npywriter.Shape = []int{lengthOfPath, len(genomes)}
		for _, g := range genomes {
			for i, value := range g.Paths[path][0] {
				if !(*g.Library).Paths[path].Variants[i].List[value].Complete {
					sliceOfData = append(sliceOfData, -2)
				} else {
					sliceOfData = append(sliceOfData, int32(value))
				}
				
				if !(*g.Library).Paths[path].Variants[i].List[value].Complete {
					sliceOfData = append(sliceOfData, -2)
				} else {
					sliceOfData = append(sliceOfData, int32(g.Paths[path][1][i]))
				}
			}
		}
		err = npywriter.WriteInt32(sliceOfData)
		if err != nil {
			log.Fatal(err)
		}
	}
}