package main


import (
    "os"
    "bufio"
    "log"
    "fmt"
    "github.com/aebruno/twobit"
)

func IndexInfo(in, out string) {
    if len(in) == 0 {
        log.Fatalln("Please provide an input file (.2bit)")
    }
    if len(out) == 0 {
        log.Fatalln("Please provide an output file (.fa)")
    }

    inFile, err := os.Open(in)
    if err != nil {
        log.Fatal(err)
    }

    defer inFile.Close()

    tb, err := twobit.NewReader(inFile)
    if err != nil {
        log.Fatal(err)
    }

    outFile, err := os.Create(out)
    if err != nil {
        log.Fatal(err)
    }

    defer outFile.Close()

    w := bufio.NewWriter(outFile)

    for name, index := range tb.Indexes() {
        w.WriteString(name)
        w.WriteString("\t")
        idx_str := fmt.Sprintf("%d", index)
        w.WriteString(idx_str)
        w.WriteString("\n")
    }

    w.Flush()

}

func main() {
  IndexInfo("x.2bit", "x.2bit.idx")
}
