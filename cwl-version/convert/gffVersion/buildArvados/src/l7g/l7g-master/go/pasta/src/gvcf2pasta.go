package main

/* PLACE HOLDER */

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
var gProfileFile string = "gvcf2pasta.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "gvcf2pasta.mprof"

var gCounter int = 0

/*
func emit_nocall(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (int64,error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return 0,e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch bp {
    case 'A', 'a': aout.WriteByte('!')
    case 'C', 'c': aout.WriteByte('#')
    case 'G', 'g': aout.WriteByte('\'')
    case 'T', 't': aout.WriteByte('4')
    default: aout.WriteByte(bp)
    }

    gCounter++

  }

  return start_pos,nil
}
*/

func peel_ref(start_pos int64, n int64, ref_ain *simplestream.SimpleStream) (string,int64,error) {
  refseq := []byte{}
  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return "",0,e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    refseq = append(refseq, bp)
  }

  return string(refseq),start_pos,nil
}

func emit_nocall_ref(start_pos int64, n int64, ref_ain *simplestream.SimpleStream, aout *bufio.Writer) (int64,error) {

  end_pos := start_pos+n
  for ; start_pos < end_pos; start_pos++ {

    if ref_ain.Pos >= ref_ain.N {
      if e:=ref_ain.Refresh() ; e!=nil { return 0,e; }
    }

    bp := ref_ain.Buf[ref_ain.Pos]
    ref_ain.Pos++

    switch bp {
    case 'a', 'A': aout.WriteByte('A')
    case 'c', 'C': aout.WriteByte('C')
    case 'g', 'G': aout.WriteByte('G')
    case 't', 'T': aout.WriteByte('T')
    default: aout.WriteByte(bp)
    }

    gCounter++

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
    case 'a', 'A': aout.WriteByte('a')
    case 'c', 'C': aout.WriteByte('c')
    case 'g', 'G': aout.WriteByte('g')
    case 't', 'T': aout.WriteByte('t')
    default: aout.WriteByte(bp)
    }

    gCounter++
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

func get_field_index(s, field, sep string) (int, error) {
  fmt_parts := strings.Split(s, sep)
  for i:=0; i<len(fmt_parts); i++ {
    if fmt_parts[i] == field {
      return i, nil
    }
  }
  return -1, fmt.Errorf("GT not found")
}

func convert_textvec(s []string) ([]int, error) {
  ivec := []int{}
  for i:=0; i<len(s); i++ {
    _x,e := strconv.Atoi(s[i])
    if e!=nil { return nil, e }
    ivec = append(ivec, _x)
  }
  return ivec, nil
}

func convert(gvcf_ain *autoio.AutoioHandle, ref_ain *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error

  //start_pos := int64(0)
  allele_num := 0 ; _ = allele_num

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  cur_spos := int64(0)

  // All co-ordinates are 0-ref.
  // End is inclusive
  //
  for gvcf_ain.ReadScan() {
    l := gvcf_ain.ReadText()

    if len(l)==0 || l[0] == '#' { continue }

    gvcf_parts := strings.Split(l, "\t")
    if len(gvcf_parts)<9 {return fmt.Errorf("not enough gvcf parts") }
    chrom   := gvcf_parts[0] ; _ = chrom
    spos,e0 := strconv.ParseInt(gvcf_parts[1], 10, 64)
    if e0!=nil { return e0 }
    spos--

    id_str  := gvcf_parts[2] ; _ = id_str
    ref_anch:= gvcf_parts[3] ; _ = ref_anch
    alt_str := gvcf_parts[4] ; _ = alt_str
    qual    := gvcf_parts[5] ; _ = qual
    filt    := gvcf_parts[6] ; _ = filt
    info_str:= gvcf_parts[7] ; _ = info_str
    fmt_str := gvcf_parts[8] ; _ = fmt_str
    samp_str:= gvcf_parts[9] ; _ = samp_str


    // Check for END
    //
    epos := int64(-1)
    info_parts := strings.Split(info_str, ";")
    for i:=0; i<len(info_parts); i++ {
      if strings.HasPrefix(info_parts[i], "END=") {
        end_parts := strings.Split(info_parts[i], "=")

        // End is inclusive
        //
        epos,e = strconv.ParseInt(end_parts[1], 10, 64)
        epos--

        if e!=nil { return e }
        break
      }
    }

    ref_len := int64(len(ref_anch))
    if epos>=0 {
      ref_len = epos - spos + 1
    }

    typ := "NOCALL"
    if filt=="PASS" { typ = "REF" }

    // Catch up to current position
    //
    if (cur_spos >= 0) && ((spos - cur_spos) > 0) {

      //fmt.Printf("\nnocall catchup %d+%d\n", cur_spos, spos-cur_spos)

      emit_nocall_ref(cur_spos, spos-cur_spos, ref_ain, bufout)
    }

    // Update previous end position
    //
    cur_spos = spos + ref_len

    // Process current line
    //
    if typ=="NOCALL" {

      //fmt.Printf("\nnocall ref %d+%d\n", spos, ref_len)

      emit_nocall_ref(spos, ref_len, ref_ain, bufout)

      continue
    }

    refseq,_,e := peel_ref(spos, ref_len, ref_ain)
    if e!=nil { return e }

    gt_idx,er := get_field_index(fmt_str, "GT", ":")
    if er!=nil { return er }

    samp_parts := strings.Split(samp_str, ":")
    if len(samp_parts) <= gt_idx {
      return fmt.Errorf( fmt.Sprintf("%s <-- NO GT FIELD", l) )
    }

    gt_field := samp_parts[gt_idx]

    gt_parts := []string{}
    if strings.Index(gt_field, "/") != -1 {
      gt_parts = strings.Split(gt_field, "/")

      //fmt.Printf("  %s %s (un)\n", gt_parts[0], gt_parts[1])
    } else if strings.Index(gt_field, "|") != -1 {
      gt_parts = strings.Split(gt_field, "|")

      //fmt.Printf("  %s %s (ph)\n", gt_parts[0], gt_parts[1])
    } else {
      gt_parts = append(gt_parts, gt_field)

      //fmt.Printf("  %s\n", gt_field)
    }

    gt_allele_idx,e := convert_textvec(gt_parts)

    //fmt.Printf(">> ref %s\n", refseq)
    alt_fields := strings.Split(alt_str, ",")

    for i:=0; i<len(gt_allele_idx); i++ {
      if gt_allele_idx[i] == 0 {
        //fmt.Printf("> alt%d %s\n", gt_allele_idx[i], refseq)
        //aout.WriteString(refseq)
        bufout.WriteString(refseq)

        gCounter += len(refseq)
      } else if (gt_allele_idx[i]-1) < len(alt_fields) {
        //fmt.Printf("> alt%d %s\n", gt_allele_idx[i],
        //aout.WriteString(alt_fields[gt_allele_idx[i]-1])
        bufout.WriteString(alt_fields[gt_allele_idx[i]-1])

        gCounter += len(alt_fields[gt_allele_idx[i]-1])

      } else {
        return fmt.Errorf( fmt.Sprintf("%s <-- invalid GT field", l) )
      }

      //DEBUG
      // Only first allele for now
      break
      //DEBUG

    }


    //fmt.Printf("chrom %s, spos %d, epos %d\n", chrom, spos, epos)

  }

  //fmt.Printf("\n\ngCounter %d\n", gCounter)

  return nil
}

func _main(c *cli.Context) {

  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  gvcf_ain,err := autoio.OpenReadScannerSimple( c.String("input") ) ; _ = gvcf_ain
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer gvcf_ain.Close()

  ref_ain := simplestream.SimpleStream{}
  ref_fp := os.Stdin
  if c.String("reference") != "-" {
    var e error
    ref_fp,e = os.Open(c.String("reference"))
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

  convert(&gvcf_ain, &ref_ain, aout, ref_start)

}

func main() {

  app := cli.NewApp()
  app.Name  = "gvcf2pasta"
  app.Usage = "gvcf2pasta"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT gVCF",
    },

    cli.StringFlag{
      Name: "reference, r",
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
      Usage: "Start of reference stream (default to start of gVCF position)",
    },

    cli.IntFlag{
      Name: "seq-start, s",
      Value: -1,
      Usage: "Start of reference stream (default to start of gVCF position)",
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
