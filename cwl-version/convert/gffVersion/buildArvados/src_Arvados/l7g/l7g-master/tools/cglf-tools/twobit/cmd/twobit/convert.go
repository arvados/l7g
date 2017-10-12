// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package main

import (
    "os"
    "bufio"
    "log"
    "github.com/aebruno/twobit"
    "github.com/aebruno/gofasta"
)

func ToFasta(in, out string) {
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

    for _, n := range tb.Names() {
        w.WriteString(">")
        w.WriteString(n)
        w.WriteString("\n")
        seq, err := tb.Read(n)
        if err != nil {
            log.Fatal(err)
        }

        n := len(seq)
        cols := 50
        rows := n/cols
        if n % cols > 0 {
            rows++
        }

        for r := 0; r < rows; r++ {
            end := (r*cols)+cols
            if end > n {
                end = n
            }

            w.Write(seq[(r*cols):end])
            w.Write([]byte("\n"))
        }
    }

    w.Flush()
}

func To2bit(in, out string) {
    if len(in) == 0 {
        log.Fatalln("Please provide an input file (.fa)")
    }
    if len(out) == 0 {
        log.Fatalln("Please provide an output file (.2bit)")
    }

    inFile, err := os.Open(in)
    if err != nil {
        log.Fatal(err)
    }

    defer inFile.Close()

    tb := twobit.NewWriter()

    for rec := range gofasta.SimpleParser(inFile) {
        err := tb.Add(rec.Id, rec.Seq)
        if err != nil {
            log.Fatal(err)
        }
    }

    outFile, err := os.Create(out)
    if err != nil {
        log.Fatal(err)
    }

    defer outFile.Close()

    err = tb.WriteTo(outFile)
    if err != nil {
        log.Fatal(err)
    }
}
