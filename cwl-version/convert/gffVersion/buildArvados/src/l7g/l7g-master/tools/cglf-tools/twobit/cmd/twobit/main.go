// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package main

import (
    "os"
    "github.com/codegangsta/cli"
//    "runtime/pprof"
)

func main() {
    //p, _ := os.Create("twobit.cpuprofile")
    //pprof.StartCPUProfile(p)
    //defer pprof.StopCPUProfile()
    app := cli.NewApp()
    app.Name    = "twobit"
    app.Authors = []cli.Author{cli.Author{Name: "Andrew E. Bruno", Email: "aeb@qnot.org"}}
    app.Usage   = "Read/Write .2bit files"
    app.Version = "0.0.1"
    app.Commands = []cli.Command {
        {
            Name: "convert",
            Usage: "Convert FASTA file to .2bit format.",
            Flags: []cli.Flag{
                &cli.BoolFlag{Name: "to-fasta, f", Usage: "Convert .2bit file to FASTA"},
                &cli.StringFlag{Name: "in, i", Usage: "Input file"},
                &cli.StringFlag{Name: "out, o", Usage: "Output file"},
            },
            Action: func(c *cli.Context) {
                if c.Bool("to-fasta") {
                    ToFasta(c.String("in"), c.String("out"))
                    return
                }

                To2bit(c.String("in"), c.String("out"))
            },
        },
    }

    app.Run(os.Args)
}
