package main

/*
*/

import "fmt"
import "io"
import "os"
import "runtime"
import "runtime/pprof"

import "bufio"
import "bytes"

import "github.com/curoverse/l7g/go/simplestream"

import "github.com/codegangsta/cli"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "pasta2gff.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "pasta2gff.mprof"

// Ref to Alt
//
var gSub map[byte]map[byte]byte
var gRefBP map[byte]byte
var gAltBP map[byte]byte
var gPastaBPState map[byte]int

func init() {

  gPastaBPState = make(map[byte]int)

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

  gRefBP = make(map[byte]byte)
  gAltBP = make(map[byte]byte)
  gRefBP['a'] = 'a'
  gRefBP['~'] = 'a'
  gRefBP['?'] = 'a'
  gRefBP['@'] = 'a'
  gRefBP['A'] = 'a'

  gAltBP['a'] = 'a'
  gAltBP['~'] = 'c'
  gAltBP['?'] = 'g'
  gAltBP['@'] = 't'
  gAltBP['A'] = 'n'

  //-

  gRefBP['='] = 'c'
  gRefBP['c'] = 'c'
  gRefBP[':'] = 'c'
  gRefBP[';'] = 'c'
  gRefBP['C'] = 'c'

  gAltBP['='] = 'a'
  gAltBP['c'] = 'c'
  gAltBP[':'] = 'g'
  gAltBP[';'] = 't'
  gAltBP['C'] = 'n'

  //-

  gRefBP['#'] = 'g'
  gRefBP['&'] = 'g'
  gRefBP['g'] = 'g'
  gRefBP['%'] = 'g'
  gRefBP['G'] = 'g'

  gAltBP['#'] = 'a'
  gAltBP['&'] = 'c'
  gAltBP['g'] = 'g'
  gAltBP['%'] = 't'
  gAltBP['G'] = 'n'

  //-

  gRefBP['*'] = 't'
  gRefBP['+'] = 't'
  gRefBP['-'] = 't'
  gRefBP['t'] = 't'
  gRefBP['T'] = 't'

  gAltBP['*'] = 'a'
  gAltBP['+'] = 'c'
  gAltBP['-'] = 'g'
  gAltBP['t'] = 't'
  gAltBP['T'] = 'n'

  // alt insertions (will appear in alt but not ref)
  //
  gAltBP['Q'] = 'a'
  gAltBP['S'] = 'c'
  gAltBP['W'] = 'g'
  gAltBP['d'] = 't'

  //-

  // Alt deletetions
  //
  gRefBP['!'] = 'a'
  gRefBP['$'] = 'c'
  gRefBP['7'] = 'g'
  gRefBP['E'] = 't'


  //--
  gPastaBPState['N'] = NOC
  gPastaBPState['n'] = NOC

  gPastaBPState['a'] = REF
  gPastaBPState['~'] = SUB
  gPastaBPState['?'] = SUB
  gPastaBPState['@'] = SUB
  gPastaBPState['A'] = NOC

  //-

  gPastaBPState['='] = SUB
  gPastaBPState['c'] = REF
  gPastaBPState[':'] = SUB
  gPastaBPState[';'] = SUB
  gPastaBPState['C'] = NOC

  //-

  gPastaBPState['#'] = SUB
  gPastaBPState['&'] = SUB
  gPastaBPState['g'] = REF
  gPastaBPState['%'] = SUB
  gPastaBPState['G'] = NOC

  //-

  gPastaBPState['*'] = SUB
  gPastaBPState['+'] = SUB
  gPastaBPState['-'] = SUB
  gPastaBPState['t'] = REF
  gPastaBPState['T'] = NOC

  //-

  gPastaBPState['!'] = INDEL
  gPastaBPState['$'] = INDEL
  gPastaBPState['7'] = INDEL
  gPastaBPState['E'] = INDEL

  gPastaBPState['Q'] = INDEL
  gPastaBPState['S'] = INDEL
  gPastaBPState['W'] = INDEL
  gPastaBPState['d'] = INDEL

}


const(
  REF = iota
  SNP = iota
  SUB = iota
  INDEL = iota
  NOC = iota
  FIN = iota
)

var g_CHROM string = "unk"
var g_INST string = "UNK"

// Position is 0 REFERENCE
// End is INCLUSIVE
//
func emit_ref(bufout *bufio.Writer, s,n int64) {
  bufout.WriteString( fmt.Sprintf("%s\t%s\tREF\t%d\t%d\t.\t+\t.\t.\n", g_CHROM, g_INST, s+1,s+n) )
}

// Position is 0 REFERENCE
// End is INCLUSIVE
//
func emit_sub_haploid(bufout *bufio.Writer, s,n int64, typ int, sub, ref []byte) {
  typ_str := ""
  if typ == SNP { typ_str = "SNP"
  } else if typ == SUB { typ_str = "SUB" }

  if len(sub)==0 { sub = []byte{'-'} }
  if len(ref)==0 { ref = []byte{'-'} }

  sub_str := fmt.Sprintf("alleles %s;ref_allele %s", sub, ref)

  bufout.WriteString( fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t.\t+\t.\t%s\n", g_CHROM, g_INST, typ_str, s+1, s+n, sub_str) )
}

// Position is 0 REFERENCE
// End is INCLUSIVE
//
func emit_indel_haploid(bufout *bufio.Writer, s,n int64, sub, ref []byte) {
  if len(sub)==0 { sub = []byte{'-'} }
  if len(ref)==0 { ref = []byte{'-'} }

  indel_str := fmt.Sprintf("alleles %s;ref_allele %s", sub, ref)

  bufout.WriteString( fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t.\t+\t.\t%s\n", g_CHROM, g_INST, "INDEL", s+1, s+n, indel_str) )
}

// Position is 0 REFERENCE
// End is INCLUSIVE
//
func emit_alt(bufout *bufio.Writer, s,n int64, typ int, altA, altB, ref []byte) {
  typ_str := ""
  switch typ {
  case REF: typ_str = "REF"
  case SNP: typ_str = "SNP"
  case SUB: typ_str = "SUB"
  case INDEL: typ_str = "INDEL"
  }
  if len(altA)==0 { altA = []byte{'-'} }
  if len(altB)==0 { altB = []byte{'-'} }
  if len(ref)==0 { ref = []byte{'-'} }

  if bytes.Equal(altA, altB) {
    bufout.WriteString( fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t.\t+\t.\talleles %s;ref_allele %s\n", g_CHROM, g_INST, typ_str, s+1, s+n, altA, ref) )
  } else {
    bufout.WriteString( fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t.\t+\t.\talleles %s/%s;ref_allele %s\n", g_CHROM, g_INST, typ_str, s+1, s+n, altA, altB, ref) )
  }

}

func _getc(s *simplestream.SimpleStream) (byte, error) {
  if s.Pos>=s.N {
    if e=s.Refresh()
    e!=nil { return 0, e }
  }
  ch := s.Buf[s.Pos]
  s.Pos++
  return ch, nil
}

type PastaHaploidState struct {
  ref_seq []byte
  alt_seq []byte

  q_seq []byte
  cur_state int
  next_state int
  ch byte

  ref_pos int
  stream_pos int
}

func (PastaHaploidState *hs) Init() {
  hs.ref_seq = make([]byte, 0, 16)
  hs.alt_seq = make([]byte, 0, 16)
  hs.q_seq = make([]byte, 0, 16)
  hs.cur_state = REF
  hs.next_state = -1
  hs.ref_pos = 0
  hs.stream_pos = 0
  hs.ch = 0
}

func (PastaHaploidState *hs) Update(ch byte) error {
  var ok bool
  var r byte
  hs.next_state,ok = gPastaBPState[bp]
  if !ok  { return fmt.Errorf("Invalid character (%c) at %d", bp, stream_pos) }
  hs.ch = ch

  if r,ok = gRefBp[ch] ; ok {
    hs.ref_seq = append(hs.ref_seq, r)
    hs.ref_pos++
  }

  if r,ok = gAltBP[ch] ; ok {
    hs.alt_seq = append(hs.ref_seq, r)
  }


  return nil
}

func (PastaHaploidState *hs) UpdateEnd() error {
  var ok bool
  var r byte
  ch := hs.ch

  if hs.cur_state == REF {
    if hs.next_state != REF {

    }
  }
}

func min2(a, b int) int {
  if a<b { return a }
  return b
}

func peel(hs0, hs1 PastaHaploidState) {

  m := min2( min2(len(hs0.ref_len), len(hs0.alt_len)), min2(len(hs1.ref_len), len(hs1.alt_len)) )
  for i:=0; i<m; i++ {
    if hs0.ref_seq[i]
  }
}

// Emits are triggered by reading in a new token and seeing a state change.  We must
// keep the currently read token, process it and decide how to interpret the previous tokens.
//
// Both streams REF then either changes, emit REF for previously read stream elements
// Either stream non-REF then turns to REF, emit non-REF for previous streams
//
func convert_diploid(ainA, ainB *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error = nil
  var ok bool
  stream_pos:=-1

  cur_state := REF
  next_state := -1

  var hs0, hs1 PastaHaploidState

  hs0.Init()
  hs1.Init()

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  state:="read-token"


  for {

    if state == "read-token" {

      if hs0.ref_len < hs1.ref_len {
        bp,e := _getc(ainA)
        if e!=nil { hs0.next_state = FIN; break }
        if bp==' ' || bp=='\n' { continue }

        hs0.Update(bp)
        continue
      } else if hs0.ref_len > hs1.ref_len {
        bp,e := _getc(ainA)
        if e!=nil { hs1.next_state = FIN; break }
        if bp==' ' || bp=='\n' { continue }

        hs1.Update(bp)
        continue
      }

      peel(hs0, hs1)

      if next_state0,ok = gPastaBPState[bp0] ; !ok {
        return fmt.Errorf("Invalid character (%c) at %d", bp0, stream_pos)
      }

      if next_state1,ok = gPastaBPState[bp1] ; !ok {
        return fmt.Errorf("Invalid character (%c) at %d", bp1, stream_pos)
      }


    } else {
    }



  }

}

func convert_diploid_old(ainA, ainB *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error = nil
  var ok bool
  stream_pos:=-1

  cur_state := REF
  next_state := -1

  var hs0, hs1 PastaHaploidState

  hs0.Init()
  hs1.Init()

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for ;; {

    if hs0.ref_len < hs1.ref_len {
      bp,e := _getc(ainA)
      if e!=nil { hs0.next_state = FIN; break }
      if bp==' ' || bp=='\n' { continue }

      hs0.Update(bp)
      continue
    } else if ref0_len > ref1_len {
      bp,e := _getc(ainA)
      if e!=nil { hs1.next_state = FIN; break }
      if bp==' ' || bp=='\n' { continue }

      hs1.Update(bp)
      continue
    }

    peel(hs0, hs1)

    if next_state0,ok = gPastaBPState[bp0] ; !ok {
      return fmt.Errorf("Invalid character (%c) at %d", bp0, stream_pos)
    }

    if next_state1,ok = gPastaBPState[bp1] ; !ok {
      return fmt.Errorf("Invalid character (%c) at %d", bp1, stream_pos)
    }

    if next_state0 == REF && next_state1 == REF {
      next_state = REF
    } else if next_state0 == INDEL || next_state1 == INDEL {
      next_state = INDEL
    } else if next_state0 == SUB || next_state1  == SUB {
      next_state = SUB
    } else if next_state0 == NOC || next_state1 == NOC {
      next_state = NOC
    }

    fmt.Printf(">>> next_state %v\n", next_state)

    if cur_state == REF {
      if next_state != REF {

        cur_len := 0
        if len(refa_seq) > len(refb_seq) {
          cur_len = len(refb_seq)
        } else {
          cur_len = len(refa_seq)
        }

        emit_ref(bufout, cur_start, int64(cur_len))

        cur_start += int64(cur_len)
        //cur_len = 1
        cur_state = next_state
        refa_seq = refa_seq[cur_len:]
        refb_seq = refb_seq[cur_len:]
        alt0_seq = alt0_seq[0:0]
        alt1_seq = alt1_seq[0:0]

        refa_len -= int64(cur_len)
        refb_len -= int64(cur_len)
      } else {
        //refa_seq = refa_seq[0:0]
        //refb_seq = refb_seq[0:0]
      }

    } else if cur_state == SUB {
      if next_state == INDEL {
        cur_state = INDEL
      } else if next_state == NOC || next_state == REF {
        //if len(alt0_seq)==1 && len(ref_seq)==1 { cur_state = SNP }
        if len(alt0_seq)==1 && len(refa_seq)==1 && len(refb_seq)==1 { cur_state = SNP }
        //emit_alt(bufout, cur_start, cur_len-1, cur_state, alt0_seq, alt1_seq, ref_seq)

        cur_len := len(refa_seq)
        emit_alt(bufout, cur_start, int64(cur_len), cur_state, alt0_seq, alt1_seq, refa_seq)

        cur_start += int64(cur_len)
        //cur_len = 1
        cur_state = next_state
        refa_seq = refa_seq[0:0]
        refb_seq = refb_seq[0:0]
        alt0_seq = alt0_seq[0:0]
        alt1_seq = alt1_seq[0:0]
      }
    } else if cur_state == INDEL {
      if next_state == INDEL || next_state == SNP || next_state == SUB {
      } else if next_state == REF || next_state == NOC {
        //emit_alt(bufout, cur_start, cur_len-1, INDEL, alt0_seq, alt1_seq, ref_seq)

        cur_len := len(refa_seq)
        emit_alt(bufout, cur_start, int64(cur_len), INDEL, alt0_seq, alt1_seq, refa_seq)

        cur_start += int64(cur_len)
        //cur_len = 1
        refa_seq = refa_seq[0:0]
        refb_seq = refb_seq[0:0]
        alt0_seq = alt0_seq[0:0]
        alt1_seq = alt1_seq[0:0]
        cur_state = next_state
      }
    } else if cur_state == NOC {
      //cur_start += cur_len-1
      //cur_len = 1
      cur_state = next_state
    }

    if r,ok := gRefBP[bp0] ; ok {
      refa_seq = append(refa_seq, r)
      ref0_len++
    }

    if r,ok := gRefBP[bp1] ; ok {
      refb_seq = append(refb_seq, r)
      ref1_len++
    }


    /*
    // assert pasta stream reference equal one another
    //
    if gRefBP[bp0] != gRefBP[bp1] {
      return fmt.Errorf( fmt.Sprintf("ref bases do not match at pos %d (%c != %c)", stream_pos, gRefBP[bp0], gRefBP[bp1]))
    }
    */

    if r,ok := gAltBP[bp0] ; ok {
      alt0_seq = append(alt0_seq, r)
    }

    if r,ok := gAltBP[bp1] ; ok {
      alt1_seq = append(alt1_seq, r)
    }

  }

  //WIP
  /*
  if cur_state == REF {
    if next_state != REF {
      emit_ref(bufout, cur_start, cur_len)
    }
  } else if cur_state == SUB {
    if next_state == INDEL {
      cur_state = INDEL
    } else if next_state == NOC  || next_state == REF || next_state == FIN {
      if len(alt0_seq)==1 && len(ref_seq)==1 { cur_state = SNP }
      emit_alt(bufout, cur_start, cur_len, cur_state, alt0_seq, alt1_seq, ref_seq)
    }
  } else if cur_state == INDEL {
    if next_state == INDEL || next_state == SNP || next_state == SUB {
    } else if next_state == REF || next_state == NOC || next_state == FIN {
      emit_alt(bufout, cur_start, cur_len, INDEL, alt0_seq, alt1_seq, ref_seq)
    }
  } else if cur_state == NOC {
    cur_start += cur_len
    cur_len = 0
    cur_state = next_state
  }
  */


  return e
}

func convert_haploid(ain *simplestream.SimpleStream, aout *os.File, start_pos int64) error {
  var e error = nil
  var ok bool
  stream_pos:=-1

  alt_seq := make([]byte, 0, 1024)
  ref_seq := make([]byte, 0, 1024)

  cur_state := REF
  next_state := -1

  ref_coord := start_pos
  //var cur_len int64 = 0
  var ref_len int64 = 0

  allele_num := 0 ; _ = allele_num

  bufout := bufio.NewWriter(aout)
  defer bufout.Flush()

  for ;; {

    if ain.Pos>=ain.N {

      if e=ain.Refresh()
      e!=nil {
        next_state = FIN
        break
      }

    }

    bp := ain.Buf[ain.Pos]
    ain.Pos++

    if bp == ' ' || bp == '\n' { continue }


    stream_pos++

    next_state,ok = gPastaBPState[bp]
    if !ok {
      return fmt.Errorf("Invalid character (%c) at %d", bp, stream_pos)
    }

    if cur_state == REF {
      if next_state != REF {
        //emit_ref(bufout, ref_coord, int64(len(ref_seq)))
        emit_ref(bufout, ref_coord, ref_len)

        //if int64(len(ref_seq)) != ref_len { panic("cp1") }

        //ref_coord += int64(len(ref_seq))
        ref_coord += ref_len
        cur_state = next_state
        ref_seq = ref_seq[0:0]
        alt_seq = alt_seq[0:0]

        ref_len = 0
      } else {
        ref_seq = ref_seq[0:0]
        alt_seq = alt_seq[0:0]
      }
    } else if cur_state == SUB {
      if next_state == INDEL {
        cur_state = INDEL
      } else if next_state == NOC  || next_state == REF {
        if len(alt_seq)==1 && len(ref_seq)==1 { cur_state = SNP }
        //emit_alt(bufout, ref_coord, int64(len(ref_seq)), cur_state, alt_seq, alt_seq, ref_seq)
        emit_alt(bufout, ref_coord, ref_len, cur_state, alt_seq, alt_seq, ref_seq)

        //if int64(len(ref_seq)) != ref_len { panic("cp2") }

        //ref_coord += int64(len(ref_seq))
        ref_coord += ref_len
        cur_state = next_state
        ref_seq = ref_seq[0:0]
        alt_seq = alt_seq[0:0]

        ref_len = 0
      }
    } else if cur_state == INDEL {
      if next_state == INDEL || next_state == SNP || next_state == SUB {
      } else if next_state == REF || next_state == NOC {
        //emit_alt(bufout, ref_coord, int64(len(ref_seq)), INDEL, alt_seq, alt_seq, ref_seq)
        emit_alt(bufout, ref_coord, ref_len, INDEL, alt_seq, alt_seq, ref_seq)

        //if int64(len(ref_seq)) != ref_len { panic("cp3") }

        //ref_coord += int64(len(ref_seq))
        ref_coord += ref_len
        ref_seq = ref_seq[0:0]
        alt_seq = alt_seq[0:0]
        cur_state = next_state

        ref_len = 0
      }
    } else if cur_state == NOC {
      //ref_coord += int64(len(ref_seq))
      ref_coord += ref_len
      cur_state = next_state
    }

    if r,ok := gRefBP[bp] ; ok {
      ref_seq = append(ref_seq, r)

      ref_len++
    }

    if r,ok := gAltBP[bp] ; ok {
      alt_seq = append(alt_seq, r)
    }

  }

  if cur_state == REF {
    if next_state != REF {
      //emit_ref(bufout, ref_coord, cur_len)
      //emit_ref(bufout, ref_coord, int64(len(ref_seq)))
      emit_ref(bufout, ref_coord, ref_len)
    }
  } else if cur_state == SUB {
    if next_state == INDEL {
      cur_state = INDEL
    } else if next_state == NOC  || next_state == REF || next_state == FIN {
      if len(alt_seq)==1 && len(ref_seq)==1 { cur_state = SNP }
      //emit_alt(bufout, ref_coord, cur_len, cur_state, alt_seq, alt_seq, ref_seq)
      //emit_alt(bufout, ref_coord, int64(len(ref_seq)), cur_state, alt_seq, alt_seq, ref_seq)
      emit_alt(bufout, ref_coord, ref_len, cur_state, alt_seq, alt_seq, ref_seq)
    }
  } else if cur_state == INDEL {
    if next_state == INDEL || next_state == SNP || next_state == SUB {
    } else if next_state == REF || next_state == NOC || next_state == FIN {
      //emit_alt(bufout, ref_coord, cur_len, INDEL, alt_seq, alt_seq, ref_seq)
      //emit_alt(bufout, ref_coord, int64(len(ref_seq)), INDEL, alt_seq, alt_seq, ref_seq)
      emit_alt(bufout, ref_coord, ref_len, INDEL, alt_seq, alt_seq, ref_seq)
    }
  }


  return e
}

func _main(c *cli.Context) {
  var err error

  /*
  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }
  */

  infn_slice := c.StringSlice("input")
  if len(infn_slice)==0 {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  ain_count:=1

  ain := simplestream.SimpleStream{}
  fp := os.Stdin
  //if c.String("input") != "-" {
  if infn_slice[0] != "-" {
    var e error
    //fp ,e = os.Open(c.String("input"))
    fp ,e = os.Open(infn_slice[0])
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    defer fp.Close()
  }
  ain.Init(fp)

  ain2 := simplestream.SimpleStream{}

  if len(infn_slice)>1 {
    ain_count++

    fp2,e := os.Open(infn_slice[1])
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    defer fp2.Close()

    ain2.Init(fp2)
  }

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

  if gProfileFlag {
    prof_f,err := os.Create( gProfileFile )
    if err != nil {
      fmt.Fprintf( os.Stderr, "Could not open profile file %s: %v\n", gProfileFile, err )
      os.Exit(2)
    }

    pprof.StartCPUProfile( prof_f )
    defer pprof.StopCPUProfile()
  }

  //convert(&gff_ain, &ref_ain, aout.Fp, ref_start)
  if ain_count == 1 {
    e := convert_haploid(&ain, aout, ref_start)
    if e!=nil && e!=io.EOF { panic(e) }
  } else {
    e := convert_diploid(&ain, &ain2, aout, ref_start)
    if e!=nil && e!=io.EOF { panic(e) }
  }

}

func main() {

  app := cli.NewApp()
  app.Name  = "pasta2gff"
  app.Usage = "pasta2gff"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    /*
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT",
    },
    */

    cli.StringSliceFlag{
      Name: "input, i",
      Usage: "INPUT",
    },

    cli.StringFlag{
      Name: "ref-input, r",
      Usage: "REF-INPUT",
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
