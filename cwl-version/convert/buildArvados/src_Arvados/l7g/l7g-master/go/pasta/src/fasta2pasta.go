package main

/*
  This kind of makes no sense.  Unless the FASTA file is exactly
  aligned with the reference stream then you might get garbage.

  We'll be using it for debugging purposes.
*/

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "bufio"
import "strings"
import "strconv"

import "github.com/curoverse/l7g/go/autoio"
import "github.com/curoverse/l7g/go/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "fasta2pasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "fasta2pasta.mprof"

func convert(fa_ain, ref_ain *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error
  var fa_bp byte

  allele_num := 0

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for ;; {
    if fa_ain.Pos >= fa_ain.N {
      if e:=fa_ain.Refresh() ; e!=nil { return e }
    }

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return e }
    }

    fa_bp = fa_ain.Buf[fa_ain.Pos]
    fa_ain.Pos++

    ref_bp = ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    if fa_bp != 'n' && fa_bp != 'N' && fa_bp == ref_bp {
      switch fa_bp {
      case 'A': aout.WriteByte('a')
      case 'C': aout.WriteByte('c')
      case 'G': aout.WriteByte('g')
      case 'T': aout.WriteByte('t')
      default: aout.WriteByte(fa_bp)
      }
    } else if fa_bp == 'N' || fa_bp == 'n' {
      switch ref_bp {
      case 'A', 'a': aout.WriteByte('!')
      case 'C', 'c': aout.WriteByte('#')
      case 'G', 'g': aout.WriteByte('\'')
      case 'T', 't': aout.WriteByte('4')
      default: aout.WriteByte(ref_bp)
      }
    } else if fa_bp != ref_bp {
      switch fa_bp {
      case 'a', 'A': aout.WriteByte('A')
      case 'c', 'C': aout.WriteByte('C')
      case 'g', 'G': aout.WriteByte('G')
      case 't', 'T': aout.WriteByte('T')
      default: aout.WriteByte(fa_bp)
      }
    } else {
      aout.WriteByte('-')
    }

  }

  return nil
}

func _main(c *cli.Context) {

  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  fasta_ain,err := simplestrea.SimpleStream{}
  fasta_fp := os.Stdin
  if c.String("input") != "-" {
    var e error
    fasta_fp,e = os.Open(c.String("input"))
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", err)
      os.Exit(1)
    }
    defer fasta_fp.Close()
  }
  fasta_ain.Init(fasta_fp)

  ref_ain := simplestream.SimpleStream{}
  ref_fp := os.Stdin
  if c.String("ref-input") != "-" {
    var e error
    ref_fp,e = os.Open(c.String("ref-input"))
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", err)
      os.Exit(1)
    }
    defer ref_fp.Close()
  }
  ref_ain.Init(ref_fp)

  var ref_start int64
  ref_start = -1
  ss := c.Int("ref-start")
  if ss > 0 { ref_start = int64(ss) }

  var seq_start int64
  seq_start = 0 ; _ = seq_start
  ss = c.Int("seq-start")
  if ss > 0 { seq_start = int64(ss) }

  aout := os.Stdout
  if c.String("output") != "-" {
    aout,err = os.Open(c.String("output"))
    if err!=nil {
      fmt.Fprintf(os.Stderr, "%v", err)
      os.Exit(1)
    }
    defer aout.Close()
  }

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

  convert(&fasta_ain, &ref_ain, aout, ref_start)

}

func main() {

  app := cli.NewApp()
  app.Name  = "fasta2pasta"
  app.Usage = "fasta2pasta"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT FASTA",
    },

    cli.StringFlag{
      Name: "ref-input, r",
      Usage: "REF-INPUT FASTA",
    },

    cli.StringFlag{
      Name: "output, o",
      Value: "-",
      Usage: "OUTPUT",
    },

    cli.IntFlag{
      Name: "max-procs, N",
      Value: -1,
      Usage: "MAXPROCS",
    },

    cli.IntFlag{
      Name: "ref-start, S",
      Value: -1,
      Usage: "Start of reference stream (default to start of GFF position)",
    },

    cli.IntFlag{
      Name: "seq-start, s",
      Value: -1,
      Usage: "Start of reference stream (default to start of GFF position)",
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
