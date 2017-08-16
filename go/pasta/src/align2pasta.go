package main

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "io"
import "bufio"
//import "strings"
//import "strconv"

import "github.com/curoverse/l7g/go/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "align2pasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "align2pasta.mprof"

// Ref to Alt
//
var gSub map[byte]map[byte]byte

func init() {
  gSub = make(map[byte]map[byte]byte)

  gSub['a'] = make(map[byte]byte)
  gSub['a']['a'] = 'a'
  gSub['a']['c'] = '~'
  gSub['a']['g'] = '?'
  gSub['a']['t'] = '@'
  gSub['a']['n'] = 'A'

  gSub['c'] = make(map[byte]byte)
  gSub['c']['a'] = '='
  gSub['c']['c'] = 'c'
  gSub['c']['g'] = ':'
  gSub['c']['t'] = ';'
  gSub['c']['n'] = 'C'

  gSub['g'] = make(map[byte]byte)
  gSub['g']['a'] = '#'
  gSub['g']['c'] = '&'
  gSub['g']['g'] = 'g'
  gSub['g']['t'] = '%'
  gSub['g']['n'] = 'G'

  gSub['t'] = make(map[byte]byte)
  gSub['t']['a'] = '*'
  gSub['t']['c'] = '+'
  gSub['t']['g'] = '-'
  gSub['t']['t'] = 't'
  gSub['t']['n'] = 'T'

  gSub['n'] = make(map[byte]byte)
  gSub['n']['a'] = '\''
  gSub['n']['c'] = '"'
  gSub['n']['g'] = ','
  gSub['n']['t'] = '_'
  gSub['n']['n'] = 'n'

}

func emit_nocall(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch bp {
    case 'A', 'a': aout.WriteByte('A')
    case 'C', 'c': aout.WriteByte('C')
    case 'G', 'g': aout.WriteByte('G')
    case 'T', 't': aout.WriteByte('T')
    case 'N', 'n': aout.WriteByte('n')
    default:
      if bp!='n' && bp!='N'{
        fmt.Printf("!!!! %c ... s%d, n%d\n", bp, start_pos, n)
      panic(bp)
      }

      aout.WriteByte(bp)
    }

  }

  return nil
}

func emit_ref(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch bp {
    case 'a', 'A': aout.WriteByte('a')
    case 'c', 'C': aout.WriteByte('c')
    case 'g', 'G': aout.WriteByte('g')
    case 't', 'T': aout.WriteByte('t')
    case 'n', 'N': aout.WriteByte('n')
    default:
      if bp!='n' && bp!='N'{
        fmt.Printf("!!!! %c ... s%d, n%d\n", bp, start_pos, n)
        panic(bp)
      }
      aout.WriteByte(bp)
    }

  }

  return nil

}

func emit_alt(start_pos int64, n int64, alt_seq string, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (error) {
  ref_pos := 0

  var ref_bp byte

  for i:=0; i<len(alt_seq); i++ {

    if int64(i)<n {
      if ref_ain.Pos >= ref_ain.N {
        if e:=ref_ain.Refresh() ; e!=nil { return e; }
      }
      ref_bp = ref_ain.Buf[ref_ain.Pos]

      /*
      // REFERENCE CHECK
      //
      lc_ref_bp := ref_bp
      lc_gff_bp := g_GFF_REF[ref_pos]

      if lc_ref_bp == 'A' { lc_ref_bp='a'
      } else if lc_ref_bp == 'C' { lc_ref_bp='c'
      } else if lc_ref_bp == 'G' { lc_ref_bp='g'
      } else if lc_ref_bp == 'T' { lc_ref_bp='t'
      }

      if lc_gff_bp == 'A' { lc_gff_bp='a'
      } else if lc_gff_bp == 'C' { lc_gff_bp='c'
      } else if lc_gff_bp == 'G' { lc_gff_bp='g'
      } else if lc_gff_bp == 'T' { lc_gff_bp='t'
      }

      if lc_ref_bp != lc_gff_bp {
        fmt.Printf("\nREF MISMATCH: GFF reported %c (at %d+%d) but got %c\n", g_GFF_REF[ref_pos], start_pos, ref_pos, ref_bp)
        panic("!!")
      }
      //
      // REFERENCE CHECK
      */


      ref_ain.Pos++
      ref_pos++

      switch ref_bp {
      case 'a', 'A': ref_bp = 'a'
      case 'c', 'C': ref_bp = 'c'
      case 'g', 'G': ref_bp = 'g'
      case 't', 'T': ref_bp = 't'
      case 'n', 'N': ref_bp = 'n'
      default: return fmt.Errorf("invalid character for reference stream ('%c') at %d", ref_bp, ref_pos)
      }


      // It's considered a sub
      //
      if (ref_bp =='n'|| ref_bp=='N') && (alt_seq[i]!='n'&&alt_seq[i]!='N') {
        fmt.Printf("\n\n>>>> start_pos %d, n %d, alt_seq %s, ref_bp %c, alt_seq[%d] %c\n\n", start_pos, n, alt_seq, ref_bp, i, alt_seq[i])
        //panic("whoa!!!")
      }

      switch alt_seq[i] {
      case 'a', 'A': aout.WriteByte( gSub[ref_bp]['a'] )
      case 'c', 'C': aout.WriteByte( gSub[ref_bp]['c'] )
      case 'g', 'G': aout.WriteByte( gSub[ref_bp]['g'] )
      case 't', 'T': aout.WriteByte( gSub[ref_bp]['t'] )
      case 'n', 'N':
        if ref_bp == 'n' || ref_bp == 'N' {
          fmt.Printf("WHOA@! %c (s%d,n%d) [%s]{%d,%c}\n", ref_bp, start_pos, n, alt_seq, i, alt_seq[i])
          panic("-->")
        }
        aout.WriteByte( gSub[ref_bp]['n'] )
      default: return fmt.Errorf("invalid character for alt sequence ('%c') at pos %d", alt_seq[i], i)
      }

    } else {

      // It's considered an insertion
      //
      switch alt_seq[i] {
      case 'a', 'A': aout.WriteByte('Q')
      case 'c', 'C': aout.WriteByte('S')
      case 'g', 'G': aout.WriteByte('W')
      case 't', 'T': aout.WriteByte('d')
      case 'n', 'N': aout.WriteByte('^')
      default: return fmt.Errorf("invalid character for alt sequence ('%c') at pos. %d", alt_seq[i], i)
      }

    }

  }

  for ; int64(ref_pos) < n ; ref_pos++ {
    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return e; }
    }
    ref_bp = ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch ref_bp {
    case 'a', 'A': ref_bp = '!'
    case 'c', 'C': ref_bp = '$'
    case 'g', 'G': ref_bp = '7'
    case 't', 'T': ref_bp = 'E'
    case 'n', 'N': ref_bp = 'z'
    default: return fmt.Errorf("invalid character for reference stream ('%c') at %d", ref_bp, ref_pos)
    }


  }

  return nil

}

var g_convert_debug bool = false

func convert(ref_ain *simplestream.SimpleStream, seq_ain *simplestream.SimpleStream, fout *os.File, start_pos int64, allele_num int) error {

  aout := bufio.NewWriter(fout)
  if g_convert_debug {
    defer func() { fmt.Printf("\n") ; aout.Flush() }()
  } else {
    defer aout.Flush()
  }

  icount := 0
  ocount := 0

  //DEBUG
  if g_convert_debug { fmt.Printf("start\n") }

  for true {

    //DEBUG
    if g_convert_debug {
      fmt.Printf("ref_ain.Pos %d, ref_ain.N %d\n", ref_ain.Pos, ref_ain.N)
    }

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return e; }
    }
    bp_ref := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    //DEBUG
    if g_convert_debug {
      fmt.Printf("seq_ain.Pos %d, seq_ain.N %d\n", seq_ain.Pos, seq_ain.N)
    }

    if seq_ain.Pos >= seq_ain.N {
      if e:=seq_ain.Refresh() ; e!=nil { return e; }
    }
    bp_seq := seq_ain.Buf[seq_ain.Pos]
    seq_ain.Pos++

    //DEBUG
    if g_convert_debug {
      fmt.Printf("%d,%d bp_ref %c, bp_seq %c\n", icount, ocount, bp_ref, bp_seq)
    }

    if bp_ref == bp_seq && bp_ref == '\n' { break }


    if bp_ref == bp_seq {

      switch bp_ref {
        case 'a', 'A': aout.WriteByte( gSub[bp_ref]['a'] )
        case 'c', 'C': aout.WriteByte( gSub[bp_ref]['c'] )
        case 'g', 'G': aout.WriteByte( gSub[bp_ref]['g'] )
        case 't', 'T': aout.WriteByte( gSub[bp_ref]['t'] )
        case 'n', 'N':
          if bp_ref == 'n' || bp_ref == 'N' {
            return fmt.Errorf("no-call to no-call match: bf_ref:%c i:%d o:%d\n", bp_ref, icount, ocount)
          }
          aout.WriteByte( gSub[bp_ref]['n'] )
        default: return fmt.Errorf("invalid character for alt sequence ('%c') at i:%d o:%d", bp_ref, icount, ocount)
      }

    } else if bp_seq == '-' {

      switch bp_ref {
        case 'a', 'A': aout.WriteByte('!')
        case 'c', 'C': aout.WriteByte('$')
        case 'g', 'G': aout.WriteByte('7')
        case 't', 'T': aout.WriteByte('E')
        case 'n', 'N': aout.WriteByte('z')
        default: return fmt.Errorf("invalid character for alt sequence ('%c') at i:%d o:%d", bp_ref, icount, ocount)
      }

    } else if bp_ref == '-' {

      switch bp_seq {
        case 'a', 'A': aout.WriteByte('Q')
        case 'c', 'C': aout.WriteByte('S')
        case 'g', 'G': aout.WriteByte('W')
        case 't', 'T': aout.WriteByte('d')
        case 'n', 'N': aout.WriteByte('f')
        default: return fmt.Errorf("invalid character for alt sequence ('%c') at i:%d o:%d", bp_ref, icount, ocount)
      }

    } else {
      ref_lc := bp_ref
      if bp_ref <= 70 { ref_lc += 32 }

      seq_lc := bp_seq
      if bp_seq <= 70 { seq_lc += 32 }

      aout.WriteByte( gSub[ref_lc][seq_lc] )
    }

    icount++
    ocount++

  }

  return nil

}

func _main(c *cli.Context) {

  inp_slice := c.StringSlice("input")
  if len(inp_slice)!=2 {
    fmt.Fprintf(os.Stderr, "Input required, exiting\n")
    cli.ShowAppHelp(c)
    os.Exit(1)
  }

  ref_ain := simplestream.SimpleStream{}
  ref_fp := os.Stdin
  if inp_slice[0] != "-" {
    var e error
    ref_fp,e = os.Open(inp_slice[0])
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    defer ref_fp.Close()
  }
  ref_ain.Init(ref_fp)

  seq_ain := simplestream.SimpleStream{}
  seq_fp := os.Stdin
  if (len(inp_slice) >= 2) && (inp_slice[1] != "-") {
    var e error
    seq_fp,e = os.Open(inp_slice[1])
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    defer seq_fp.Close()
  }
  seq_ain.Init(seq_fp)

  var ref_start int64
  ref_start = 0
  ss := c.Int("ref-start")
  if ss > 0 { ref_start = int64(ss) }

  var seq_start int64
  seq_start = 0 ; _ = seq_start
  ss = c.Int("seq-start")
  if ss > 0 { seq_start = int64(ss) }

  aout := os.Stdout
  if c.String("output") != "-" {
    var err error
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

  allele := c.Int("allele")

  if gProfileFlag {
    prof_f,err := os.Create( gProfileFile )
    if err != nil {
      fmt.Fprintf( os.Stderr, "Could not open profile file %s: %v\n", gProfileFile, err )
      os.Exit(2)
    }

    pprof.StartCPUProfile( prof_f )
    defer pprof.StopCPUProfile()
  }

  e := convert(&ref_ain, &seq_ain, aout, ref_start, allele)
  if e!=nil && e!=io.EOF { panic(e) }

  aout.Sync()

}

func main() {

  app := cli.NewApp()
  app.Name  = "align2pasta"
  app.Usage = "align2pasta"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringSliceFlag{
      Name: "input, i",
      Usage: "Input sequence streams",
    },

    cli.StringFlag{
      Name: "ref-input, r",
      Usage: "Reference input (FASTA)",
    },

    cli.StringFlag{
      Name: "output, o",
      Value: "-",
      Usage: "Output",
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

    cli.IntFlag{
      Name: "allele, a",
      Value: 0,
      Usage: "Wich allele to use (default 0)",
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
