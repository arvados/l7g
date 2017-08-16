package main

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "bufio"
import "strings"
import "strconv"

import "github.com/curoverse/l7g/autoio"
import "github.com/curoverse/l7g/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "vcf2pasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "vcf2pasta.mprof"

func emit_nocall(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (int64,error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return 0,e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    if bp=='A' { aout.WriteByte('a') }

    switch bp {
    case 'A', 'a': aout.WriteByte('!')
    case 'C', 'c': aout.WriteByte('#')
    case 'G', 'g': aout.WriteByte('\'')
    case 'T', 't': aout.WriteByte('4')
    default: aout.WriteByte(bp)
    }

  }

  return start_pos,nil
}

func emit_ref(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (int64,error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return 0,e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch bp {
    case 'A': aout.WriteByte('a')
    case 'C': aout.WriteByte('c')
    case 'G': aout.WriteByte('g')
    case 'T': aout.WriteByte('t')
    default: aout.WriteByte(bp)
    }

  }

  return start_pos,nil

}

func emit_alt(start_pos int64, n int64, alt_seq string, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (int64,error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {
    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return 0,e; }
    }
    ref_ain.Pos++
  }

  sub_len := n
  if n > int64(len(alt_seq)) { sub_len = int64(len(alt_seq)) }

  for i:=0; i<len(alt_seq); i++ {

    if int64(i)<sub_len {
      switch alt_seq[i] {
      case 'a', 'A': aout.WriteByte('A')
      case 'c', 'C': aout.WriteByte('C')
      case 'g', 'G': aout.WriteByte('G')
      case 't', 'T': aout.WriteByte('T')
      default: aout.WriteByte(alt_seq[i])
      }
    } else {
      switch alt_seq[i] {
      case 'a', 'A': aout.WriteByte('b')
      case 'c', 'C': aout.WriteByte('d')
      case 'g', 'G': aout.WriteByte('h')
      case 't', 'T': aout.WriteByte('u')
      default: aout.WriteByte(alt_seq[i])
      }
    }

    //aout.WriteByte(alt_seq[i])

  }

  return start_pos,nil

}

func convert(gff_ain *autoio.AutoioHandle, ref_ain *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error

  //start_pos := int64(0)
  allele_num := 0

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for gff_ain.ReadScan() {
    l := gff_ain.ReadText()

    if len(l)==0 || l[0] == '#' { continue }

    gff_parts := strings.Split(l, "\t")
    if len(gff_parts)<9 {return fmt.Errorf("not enough gff parts") }
    chrom   := gff_parts[0] ; _ = chrom
    typ     := gff_parts[2]
    spos,e0 := strconv.ParseInt(gff_parts[3], 10, 64)
    epos,e1 := strconv.ParseInt(gff_parts[4], 10, 64)
    info    := gff_parts[8]

    alt_seq := ""

    del_n := epos-spos+1
    spos0ref := spos-1

    if typ != "REF" {
      info_parts := strings.Split(info, ";")
      alleles_info := info_parts[0]

      if len(alleles_info) < 2 { return fmt.Errorf("no") }

      alts := strings.Split(alleles_info, " ")
      alt_seqs := strings.Split(alts[1], "/")
      if len(alt_seqs) == 0 {
        alt_seq = alt_seqs[0]
      } else {
        alt_seq = alt_seqs[allele_num]
        if alt_seq == "-" { alt_seq = "" }
      }
    }

    if e0!=nil { return e0 }
    if e1!=nil { return e1 }

    if start_pos < 0 { start_pos = spos0ref }


    if start_pos < spos0ref {
      start_pos,e = emit_nocall(start_pos, spos0ref-start_pos, ref_ain, bufout)
      if e!=nil { return e }
    }

    if typ=="REF" {
      start_pos,e = emit_ref(start_pos, (spos0ref+del_n)-start_pos, ref_ain, bufout)
      if e!=nil { return e }
    } else {
      start_pos,e = emit_alt(start_pos, (spos0ref+del_n)-start_pos, alt_seq, ref_ain, bufout)
      if e!=nil { return e }
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

  gff_ain,err := autoio.OpenReadScannerSimple( c.String("input") ) ; _ = gff_ain
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer gff_ain.Close()

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

  //aout,err := autoio.CreateWriter( c.String("output") ) ; _ = aout
  aout := os.Stdout
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  //defer func() { aout.Flush() ; aout.Close() }()

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

  convert(&gff_ain, &ref_ain, aout, ref_start)

}

func main() {

  app := cli.NewApp()
  app.Name  = "vcf2pasta"
  app.Usage = "vcf2pasta"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT GFF",
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
