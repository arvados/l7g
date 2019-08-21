/*Package genome is a package for representing the genome, relative to a tile library, with Go data structures.
Provided are various functions to export/import the data within genomes, along with creating new Genome data structures in memory.
Should be used in conjunction with the tilelibrary and structures packages.*/
package genome

import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"../structures"
	"../tile-library"

	"github.com/kshedden/gonpy"
)

// The following are errors that may be useful to return.

// ErrInvalidGenome is an error for when a file that is expected to be a .genome file is not one.
var ErrInvalidGenome = errors.New("not a valid genome file")

// ErrNoLibraryAttached is an error for when a genome does not have a library in its Library field but needs one.
var ErrNoLibraryAttached = errors.New("genome has no library attached")

// Step is a type to represent a step within a path, which can take on a specific tile variant.
type Step int // -1 for a skipped step, and any other integer refers to the tile variant number in the reference library.

// Path is a type to represent a path of the genome, as a slice of steps.
type Path []Step

/*
Genome is a struct to represent a genome. It contains a pointer to its reference library, which allows for easy tiling.
The genome is split into paths, phases, and steps, and at each step there is potentially a tile variant of bases.
If Paths[a][b][c] = d, then in the genome at path a, phase b, and step c, the variant is variant number d in the Library at path a and step c.
Using the reference library, we can refer to the tile variant at each step using its tile variant number in the library, or using -1 to represent that there's no tile there due to a spanning tile variant.
*/
type Genome struct {
	Paths     [][]Path             // Paths represents a genome through its paths. Two phases are present here (path and counterpart path).
	libraryID [md5.Size]byte       // The ID of the underlying library, for verfication purposes.
	library   *tilelibrary.Library // This is the reference library for this Genome.
}

// AssignLibrary assigns an existing genome to a library. This library must match the ID found in the libraryID field of g.
func (g *Genome) AssignLibrary(library *tilelibrary.Library) error {
	var emptyID [md5.Size]byte
	if library.ID == emptyID {
		library.AssignID()
	}
	if library.ID != g.libraryID {
		return errors.New("genome's library ID and new library ID do not match")
	}
	g.library = library
	return nil
}

// isComplete determines if a set of bases is complete (has no nocalls).
// This is only a helper function for AddFastJ.
func isComplete(bases string) bool {
	return !strings.ContainsRune(bases, 'n')
}

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

// AddFastJ puts the contents of a FastJ into a Genome.
// Works with both gzipped and non-gzipped FastJ files.
func (g *Genome) AddFastJ(filepath string) error {
	if g == nil || g.library == nil {
		return errors.New("genome is nil or genome library is nil, cannot parse FastJ")
	}
	file := path.Base(filepath)           // The name of the file.
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
	pathNumber := 256*int(pathHex[0]) + int(pathHex[1]) // Conversion from hex into decimal--this is the path
	var data []byte
	var err error
	if len(splitpath) == 3 {
		data, err = readAllGZ(filepath)
	} else {
		data, err = ioutil.ReadFile(filepath)
	}
	if err != nil {
		return err
	}
	text := string(data)
	tiles := strings.Split(text, "\n\n") // since the only divider between each tile is two newlines, this works
	for _, line := range tiles {
		if strings.HasPrefix(line, ">") {
			stepInHex := line[20:24]
			stepBytes, err := hex.DecodeString(stepInHex)
			if err != nil {
				return err
			}
			step := 256*int(stepBytes[0]) + int(stepBytes[1])
			hashString := line[40:72]
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

			length, err := strconv.Atoi(lengthString)
			if err != nil {
				return err
			}
			phase, err := strconv.Atoi(line[25:28])
			if err != nil {
				return err
			}
			baseData := strings.Split(line, "\n")[1:]
			var b strings.Builder
			for _, data := range baseData {
				if data != "\n" {
					b.WriteString(data)
				}
			}
			bases := b.String()
			newTile := structures.TileVariant{Hash: hashArray, Length: length, Annotation: "", LookupReference: -1, Complete: isComplete(bases), ReferenceLibrary: g.library}
			// In the case of newTile, -1 is a value used to complete the creation of the tile, and has no meaning otherwise.
			index, ok := g.library.TileExists(pathNumber, step, &newTile)
			if !ok {
				return errors.New("this FastJ is not part of the library")
			}
			for len((g.Paths)[pathNumber][phase]) <= step+length-1 {
				(g.Paths)[pathNumber][phase] = append((g.Paths)[pathNumber][phase], -1) // This adds empty (skipped) steps until we reach the right step number.
			}
			(g.Paths)[pathNumber][phase][step] = Step(index)
		}
	}
	splitpath, data, tiles = nil, nil, nil // Clears most things in memory that were used here.
	return nil
}

// Add puts the contents of a directory of FastJ files into a given Genome.
func (g *Genome) Add(directory string) error {
	fastJFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range fastJFiles {
		if strings.HasSuffix(file.Name(), ".gz") {
			err = g.AddFastJ(path.Join(directory, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// New is a function to initialize a Genome.
// nil is allowed for the library if the library shouldn't be set yet. It can be set manually later.
func New(library *tilelibrary.Library) *Genome {
	var newPaths [][]Path
	newPaths = make([][]Path, structures.Paths, structures.Paths)
	for i := range newPaths {
		newPaths[i] = make([]Path, 2, 2)
		newPaths[i][0] = make([]Step, 0)
		newPaths[i][1] = make([]Step, 0)
	}
	var emptyID [md5.Size]byte
	if library != nil && library.ID == emptyID {
		library.AssignID()
	}
	return &Genome{newPaths, library.ID, library}
}

// WriteToFile writes a genome to a list format of indices relative to its reference library.
// It alternates between each phase for each step.
// For example, if path 24 had two steps (0 and 1) all with complete tiles, and on phase 0 the step values were 0 and 2, and on phase 1 the step values were 1 and 1,
// the numpy array would be [0 1 2 1], since it writes out the data for step 0 first, and then writes out the values for step 1.
// Will not work if the genome does not have a reference library (nil reference).
func (g *Genome) WriteToFile(filename string) error {
	if g.library == nil {
		return ErrNoLibraryAttached
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriter(f)
	bufferedWriter.WriteString(hex.EncodeToString(g.libraryID[:]))
	bufferedWriter.WriteString("\n")
	for path := range g.Paths {
		if (*g.library).Paths[path].Variants[0] != nil {
			if g.Paths[path][0][0] >= 0 && !(*g.library).Paths[path].Variants[0].List[int(g.Paths[path][0][0])].Complete {
				bufferedWriter.WriteString("-2")
			} else {
				bufferedWriter.WriteString(strconv.Itoa(int(g.Paths[path][0][0])))
			}
			bufferedWriter.WriteString(",")
			if g.Paths[path][1][0] >= 0 && !(*g.library).Paths[path].Variants[0].List[int(g.Paths[path][1][0])].Complete {
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
			if (*g.library).Paths[path].Variants[i+1] != nil {
				bufferedWriter.WriteString(",")
				if value >= 0 && !(*g.library).Paths[path].Variants[i+1].List[int(value)].Complete {
					bufferedWriter.WriteString("-2")
				} else {
					bufferedWriter.WriteString(strconv.Itoa(int(value)))
				}
				bufferedWriter.WriteString(",")
				if g.Paths[path][1][i+1] >= 0 && !(*g.library).Paths[path].Variants[i+1].List[int(g.Paths[path][1][i+1])].Complete {
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
	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

// ReadGenomeFromFile reads a text file containing genome information.
// Current file suffix is .genome (make sure all genomes written to disk have this suffix!)
func ReadGenomeFromFile(filepath string) (*Genome, error) {
	var newPaths [][]Path
	newPaths = make([][]Path, structures.Paths, structures.Paths)
	for i := range newPaths {
		newPaths[i] = make([]Path, 2, 2)
		newPaths[i][0] = make([]Step, 0)
		newPaths[i][1] = make([]Step, 0)
	}
	info, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(info), "\n")
	libraryID, err := hex.DecodeString(lines[0])
	if err != nil {
		return nil, err
	}
	var libraryIDArray [md5.Size]byte
	copy(libraryIDArray[:], libraryID)
	for i, line := range lines[1:] {
		indices := strings.Split(line, ",")
		for j, index := range indices {
			if index != "" {
				indexInt, err := strconv.Atoi(index)
				if err != nil {
					return nil, err
				}
				newPaths[i][j%2] = append(newPaths[i][j%2], Step(indexInt))
			}
		}
	}
	return &Genome{newPaths, libraryIDArray, nil}, nil
}

// WriteNumpy writes the values of a path of a genome to a numpy array.
// It alternates between each phase for each step.
// For example, if path 24 had two steps (0 and 1) all with complete tiles, and on phase 0 the step values were 0 and 2, and on phase 1 the step values were 1 and 1,
// the numpy array would be [0 1 2 1], since it writes out the data for step 0 first, and then writes out the values for step 1.
func (g *Genome) WriteNumpy(directory string, genomePath int) error {
	npywriter, err := gonpy.NewFileWriter(path.Join(directory, fmt.Sprintf("%04x", genomePath)))
	if err != nil {
		return err
	}
	sliceOfData := make([]int32, 0, 1)
	for i, value := range g.Paths[genomePath][0] {
		if value >= 0 && !(*g.library).Paths[genomePath].Variants[i].List[value].Complete {
			sliceOfData = append(sliceOfData, -2)
		} else {
			sliceOfData = append(sliceOfData, int32(value))
		}

		if g.Paths[genomePath][1][i] >= 0 && !(*g.library).Paths[genomePath].Variants[i].List[g.Paths[genomePath][1][i]].Complete {
			sliceOfData = append(sliceOfData, -2)
		} else {
			sliceOfData = append(sliceOfData, int32(g.Paths[genomePath][1][i]))
		}
	}
	err = npywriter.WriteInt32(sliceOfData)
	if err != nil {
		return err
	}
	return nil
}

// ReadGenomePathNumpy reads one path's worth of information from a numpy file.
// This path should be assigned to a path of a genome.
func (g *Genome) ReadGenomePathNumpy(filepath string) error {
	filename := path.Base(filepath)
	pathNumber, err := strconv.ParseInt(filename, 16, 0)
	if err != nil {
		return err
	}
	newPaths := make([]Path, 2)
	newPaths[0] = make([]Step, 0, 1)
	newPaths[1] = make([]Step, 0, 1)
	npyreader, err := gonpy.NewFileReader(filepath)
	if err != nil {
		return err
	}
	pathInfo, err := npyreader.GetInt32()
	if err != nil {
		return err
	}
	for i, index := range pathInfo {
		newPaths[i%2] = append(newPaths[i%2], Step(index))
	}
	g.Paths[int(pathNumber)] = newPaths
	return nil
}

// WriteGenomesPathToNumpy writes multiple genomes' worth of path information to a numpy file.
func WriteGenomesPathToNumpy(genomes []*Genome, filepath string, path int) error {
	if len(genomes) == 0 { // Requires a nonempty list.
		return errors.New("list of genomes is empty--nothing written to numpy")
	}
	npywriter, err := gonpy.NewFileWriter(filepath)
	if err != nil {
		return err
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
			return errors.New("path lengths within each genome are not equal")
		}
	}
	sliceOfData := make([]int32, len(genomes)*lengthOfPath)
	npywriter.Shape = []int{len(genomes), lengthOfPath}
	for _, g := range genomes {
		for i, value := range g.Paths[path][0] {
			if !(*g.library).Paths[path].Variants[i].List[value].Complete {
				sliceOfData = append(sliceOfData, -2)
			} else {
				sliceOfData = append(sliceOfData, int32(value))
			}

			if !(*g.library).Paths[path].Variants[i].List[value].Complete {
				sliceOfData = append(sliceOfData, -2)
			} else {
				sliceOfData = append(sliceOfData, int32(g.Paths[path][1][i]))
			}
		}
	}
	err = npywriter.WriteInt32(sliceOfData)
	if err != nil {
		return err
	}
	return nil
}

// Liftover runs a liftover operation on a genome by changing what library the genome is attached to.
func (g *Genome) Liftover(destination *tilelibrary.Library) error {
	if g.library == nil {
		return ErrNoLibraryAttached
	}
	mapping, err := tilelibrary.CreateMapping(g.library, destination)
	if err != nil {
		return err
	}
	for i, path := range g.Paths {
		for j, phase := range path {
			for step, value := range phase {
				if value >= 0 { // Don't change skipped steps (which have value -1)
					g.Paths[i][j][step] = Step(mapping.Mapping[i][step][value])
				}
			}
		}
	}
	g.libraryID = destination.ID
	g.library = destination // Can set the new reference library, since liftover is complete
	return nil
}
