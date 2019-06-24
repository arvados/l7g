package genome
/*
	Genome is a package for representing the genome with Go data structures.
*/
import (
	"errors"
	"encoding/hex"
	"io/ioutil"
	"path"
	"log"
	"strconv"
	"strings"
	"../structures"
)

// Genome is a type to represent a genome, through its paths. Two phases are present here (path and counterpart path).
type Genome [][]Path

// Path is a type to represent a path, through its steps.
type Path []Step

// Step is a type to represent a step within a path, which can take on a specific tile variant.
type Step struct {
	Skipped bool // Determines if this step has been skipped (due to a spanning tile)
	Variant *(structures.TileVariant) // The Variant in this step
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
	data := structures.OpenGZ(filepath)
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
			//bases := b.String()
			newTile := structures.TileCreator(hashArray, length, "")

			for len((*genome)[hexNumber][phase]) < step {
				(*genome)[hexNumber][phase] = append((*genome)[hexNumber][phase], Step{true, nil}) // This adds empty (skipped) steps until we reach the right step number.
			}
			(*genome)[hexNumber][phase] = append((*genome)[hexNumber][phase],Step{false, &newTile})
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
func InitializeGenome() Genome {
	var newGenome Genome
	newGenome = make([][]Path, structures.Paths, structures.Paths)
	for i := range newGenome {
		newGenome[i] = make([]Path, 2, 2)
	}
	return newGenome
}