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
func ParseFastJGenome(filepath string, genome *Genome) {
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
				if commaCounter == 6 {
					lengthString=string(line[i-1])
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
// nil is allowed for the library if the library shouldn't be set yet.
func InitializeGenome(library *tilelibrary.Library) Genome {
	var newPaths [][]Path
	newPaths = make([][]Path, structures.Paths, structures.Paths)
	for i := range newPaths {
		newPaths[i] = make([]Path, 2, 2)
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
// find incomplete tiles:
func WriteGenomeToFile(filename string, g *Genome) {
	f, err := os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bufferedWriter := bufio.NewWriter(f)
	for path := range g.Paths {
		if !(*g.Library).Paths[path].Variants[0].List[g.Paths[path][0][0]].Complete {
			bufferedWriter.WriteString("-2")
		} else {
			bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][0][0])))
		}
		if !(*g.Library).Paths[path].Variants[0].List[g.Paths[path][1][0]].Complete {
			bufferedWriter.WriteString("-2")
		} else {
			bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][0])))
		}
		for i, value := range g.Paths[path][0][1:] {
			bufferedWriter.WriteString(",")
			if !(*g.Library).Paths[path].Variants[i].List[value].Complete {
				bufferedWriter.WriteString("-2")
			} else {
				bufferedWriter.WriteString(strconv.Itoa(int(value)))
			}
			bufferedWriter.WriteString(",")
			if !(*g.Library).Paths[path].Variants[i].List[g.Paths[path][1][i]].Complete {
				bufferedWriter.WriteString("-2")
			} else {
				bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][1][i])))
			}
		}
		bufferedWriter.WriteString("\n")
	}
}

// writes a path of a genome to a numpy file.
func (g *Genome) WriteNumpy(filepath string, path int) {
	npywriter, err := gonpy.NewFileWriter(filepath)
	if err != nil {
		log.Fatal(err)
	}
	sliceOfData := make([]int32, 0, 1)
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
	err = npywriter.WriteInt32(sliceOfData)
	if err != nil {
		log.Fatal(err)
	}
}