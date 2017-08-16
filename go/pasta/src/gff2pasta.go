package main

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "io"
import "bufio"
import "strings"
import "strconv"

import "github.com/curoverse/l7g/go/autoio"
import "github.com/curoverse/l7g/go/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "gff2pasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "gff2pasta.mprof"

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

var g_GFF_REF []byte

func convert(gff_ain *autoio.AutoioHandle, ref_ain *simplestream.SimpleStream, aout *os.File, start_pos int64, allele_num int) error {
  var e error

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for gff_ain.ReadScan() {
    l := gff_ain.ReadText()

    if len(l)==0 || l[0] == '#' { continue }

    gff_parts := strings.Split(l, "\t")
    if len(gff_parts)<9 {return fmt.Errorf("not enough gff parts") }
    chrom   := gff_parts[0] ; _ = chrom
    typ     := gff_parts[2]
    spos1ref,e0 := strconv.ParseInt(gff_parts[3], 10, 64)
    epos1ref,e1 := strconv.ParseInt(gff_parts[4], 10, 64)
    info    := gff_parts[8]

    // REFERENCE CHECK
    //
    gff_ref_seq := make([]byte, 0, 8)
    idx := strings.Index(info, ";ref_allele ")
    if idx>=0 {
      idx += len(";ref_allele ")
      for ; idx<len(info); idx++ {
        if info[idx] == 'a' || info[idx] == 'A' ||
           info[idx] == 'c' || info[idx] == 'C' ||
           info[idx] == 'g' || info[idx] == 'G' ||
           info[idx] == 't' || info[idx] == 'T' ||
           info[idx] == '-' {
          gff_ref_seq = append(gff_ref_seq, info[idx])
        } else {
          break
        }
        if gff_ref_seq[len(gff_ref_seq)-1]=='-' { break }
      }
    }
    g_GFF_REF = gff_ref_seq
    //
    // REFERENCE CHECK

    alt_seq := ""

    del_n := epos1ref-spos1ref+1
    spos0ref := spos1ref-1

    if typ != "REF" {
      info_parts := strings.Split(info, ";")
      alleles_info := info_parts[0]

      if len(alleles_info) < 2 { return fmt.Errorf( fmt.Sprintf("Invalid alleles info (%s)", info) ) }

      alts := strings.Split(alleles_info, " ")
      alt_seqs := strings.Split(alts[1], "/")
      if len(alt_seqs) == 0 {
        alt_seq = alt_seqs[0]
      } else {
        if allele_num < len(alt_seqs) {
          alt_seq = alt_seqs[allele_num]
        } else {
          alt_seq = alt_seqs[0]
        }
        if alt_seq == "-" { alt_seq = "" }
      }
    }

    if e0!=nil { return e0 }
    if e1!=nil { return e1 }

    if start_pos < 0 {
      if gVerboseFlag {
        fmt.Printf("\n{\"comment\":\"initializing start_pos=%d\"}\n", spos0ref)
      }
      start_pos = spos0ref
    }

    if start_pos < spos0ref {
      e = emit_nocall(start_pos, spos0ref-start_pos, ref_ain, bufout)
      if e!=nil { return e }
      start_pos = spos0ref
    }

    if typ=="REF" {
      e = emit_ref(start_pos, del_n, ref_ain, bufout)
      if e!=nil { return e }
      start_pos += del_n
    } else {

      //DEBUG
      //if int(del_n) != len(alt_seq) { fmt.Printf("\n>>>>> [%d] del_n %d, alt_seq %s (%d)\n", start_pos, del_n, alt_seq, len(alt_seq)) }

      e = emit_alt(start_pos, del_n, alt_seq, ref_ain, bufout)
      if e!=nil { return e }
      start_pos += del_n
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
  ref_start = 0
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

  e := convert(&gff_ain, &ref_ain, aout, ref_start, allele)
  if e!=nil && e!=io.EOF { panic(e) }

  aout.Sync()

}

func main() {

  app := cli.NewApp()
  app.Name  = "gff2pasta"
  app.Usage = "gff2pasta"
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
