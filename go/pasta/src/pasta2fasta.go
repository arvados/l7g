package main

/*
This is more useful
*/

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "bufio"

import "github.com/curoverse/l7g/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "pasta2fasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "pasta2fasta.mprof"

func convert(pa_ain *simplestream.SimpleStream, aout *os.File) error {
  var bp byte

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for ;; {
    if pa_ain.Pos >= pa_ain.N {
      if e:=pa_ain.Refresh() ; e!=nil { return e }
    }

    bp = pa_ain.Buf[pa_ain.Pos]
    pa_ain.Pos++

    switch bp {
    case '=', '#', '*', '\'', 'a', 'Q': bufout.WriteByte('a')
    case '~', '&', '+', '"', 'c', 'S': bufout.WriteByte('c')
    case '?', ':', '-', ',', 'g', 'W': bufout.WriteByte('g')
    case '@', ';', '%', '_', 't', 'd': bufout.WriteByte('t')
    case 'A', 'C', 'G', 'T', 'n', 'N': bufout.WriteByte('n')
    }

  }

  return nil
}

func _main(c *cli.Context) {
  var err error

  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  pasta_ain := simplestream.SimpleStream{}
  pasta_fp := os.Stdin
  if len(c.String("input")) > 0 && c.String("input") != "-" {
    var e error
    pasta_fp,e = os.Open(c.String("input"))
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    defer pasta_fp.Close()
  }
  pasta_ain.Init(pasta_fp)

  aout := os.Stdout
  if c.String("output") != "-" {
    aout,err = os.Open(c.String("output"))
    if err!=nil {
      fmt.Fprintf(os.Stderr, "%v", err)
      os.Exit(1)
    }
    defer aout.Close()
  }


  /*
  aout := os.Stdout
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  */

  if c.Bool( "pprof" ) {
    gProfileFlag = true
    gProfileFile = c.String("pprof-file")
  }

  if c.Bool( "mprof" ) {
    gMemProfileFlag = true
    gMemProfileFile = c.String("mprof-file")
  }

  gVerboseFlag = c.Bool("Verbose")

  if c.Int("max-procs") > 0 {
    runtime.GOMAXPROCS( c.Int("max-procs") )
  }

  if gProfileFlag {
    prof_f,err := os.Create( gProfileFile )
    if err != nil {
      fmt.Fprintf( os.Stderr, "Could not open profile file %s: %v\n", gProfileFile, err )
      os.Exit(2)
    }

    pprof.StartCPUProfile( prof_f )
    defer pprof.StopCPUProfile()
  }

  convert(&pasta_ain, aout)

}

func main() {

  app := cli.NewApp()
  app.Name  = "pasta2fasta"
  app.Usage = "pasta2fasta"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT PASTA",
    },

    cli.StringFlag{
      Name: "output, o",
      Value: "-",
      Usage: "OUTPUT FASTA",
    },

    cli.IntFlag{
      Name: "max-procs, N",
      Value: -1,
      Usage: "MAXPROCS",
    },

    cli.BoolFlag{
      Name: "Verbose, V",
      Usage: "Verbose flag",
    },

    cli.BoolFlag{
      Name: "pprof",
      Usage: "Profile usage",
    },

    cli.StringFlag{
      Name: "pprof-file",
      Value: gProfileFile,
      Usage: "Profile File",
    },

    cli.BoolFlag{
      Name: "mprof",
      Usage: "Profile memory usage",
    },

    cli.StringFlag{
      Name: "mprof-file",
      Value: gMemProfileFile,
      Usage: "Profile Memory File",
    },

  }

  app.Run( os.Args )

  if gMemProfileFlag {
    fmem,err := os.Create( gMemProfileFile )
    if err!=nil { panic(fmem) }
    pprof.WriteHeapProfile(fmem)
    fmem.Close()
  }

}
