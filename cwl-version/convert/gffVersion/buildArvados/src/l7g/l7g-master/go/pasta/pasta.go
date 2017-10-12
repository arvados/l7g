// Package pasta provides primitives for manipulating
// PASTA streams.
//
package pasta

import "os"
import "io"
import "bufio"

import _ "errors"

import _ "fmt"

type ControlMessage struct {
  Type    int
  N       int
  NBytes  int

  Chrom   string
  RefPos  int
  RefLen  int

  Comment string
}


// All valid pasta tokens
//
var Token []byte

// Maps lower case sequence [gcat] reference and alt
// to pasta character
// e.g. reference 'c' and alt 'a': SubMap['c']['a'] = '='
//
var SubMap map[byte]map[byte]byte

// Key is pasta character, value is the lower case
// sequence value of the reference
//
var RefMap map[byte]byte

// Key is the pasta character, value is the implied
// sequence character
//
var AltMap map[byte]byte

var DelMap map[byte]byte
var InsMap map[byte]byte
var IsAltDel map[byte]bool

var RefDelBP map[byte]int
var BPState map[byte]int


// Ref to Alt
//
var gRefBP map[byte]byte
var gAltBP map[byte]byte
var gPastaBPState map[byte]int


const(
  BEG = iota  // 0
  REF = iota
  NOC = iota
  ALT = iota
  MSG = iota
  MSG_REF_NOC = iota
  MSG_CHROM = iota
  MSG_POS = iota
  FIN = iota
  SNP = iota
  SUB = iota
  INS = iota
  DEL = iota
  INDEL = iota
  NOREF = iota

  CHROM = iota
  POS = iota
  COMMENT = iota
)


func init() {
  Token := []byte("acgtnNACGT~?@=:;#&%*+-QSWd!$7EZ'\",_")

  gPastaBPState = make(map[byte]int)

  DelMap = make(map[byte]byte)
  InsMap = make(map[byte]byte)

  DelMap['a'] = '!'
  DelMap['c'] = '$'
  DelMap['g'] = '7'
  DelMap['t'] = 'E'
  DelMap['n'] = 'z'

  IsAltDel = make(map[byte]bool)
  IsAltDel['!'] = true
  IsAltDel['$'] = true
  IsAltDel['7'] = true
  IsAltDel['E'] = true
  IsAltDel['z'] = true


  InsMap['a'] = 'Q'
  InsMap['c'] = 'S'
  InsMap['g'] = 'W'
  InsMap['t'] = 'd'
  InsMap['n'] = 'Z'

  gSub := make(map[byte]map[byte]byte)

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

  // deletion of reference
  //
  gSub['a']['-'] = '!'
  gSub['c']['-'] = '$'
  gSub['t']['-'] = 'E'
  gSub['g']['-'] = '7'
  gSub['n']['-'] = 'z'

  // insertion
  //
  gSub['-'] = make(map[byte]byte)
  gSub['-']['a'] = 'Q'
  gSub['-']['c'] = 'S'
  gSub['-']['g'] = 'W'
  gSub['-']['t'] = 'd'
  gSub['-']['n'] = 'Z'

  gSub['-']['-'] = '.'


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

  //--
  // Alt deleitions

  gRefBP['!'] = 'a'
  gRefBP['$'] = 'c'
  gRefBP['7'] = 'g'
  gRefBP['E'] = 't'
  gRefBP['z'] = 'n'

  // Alt insertions

  gAltBP['Q'] = 'a'
  gAltBP['S'] = 'c'
  gAltBP['W'] = 'g'
  gAltBP['d'] = 't'
  gAltBP['Z'] = 'n'

  //--

  // no-call substitutions

  gRefBP['\''] = 'n'
  gRefBP['"'] = 'n'
  gRefBP[','] = 'n'
  gRefBP['_'] = 'n'

  gAltBP['\''] = 'a'
  gAltBP['"'] = 'c'
  gAltBP[','] = 'g'
  gAltBP['_'] = 't'

  //-

  gRefBP['n'] = 'n'
  gRefBP['N'] = 'n'

  gAltBP['n'] = 'n'
  gAltBP['N'] = 'n'


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

  /*
  gPastaBPState['!'] = INDEL
  gPastaBPState['$'] = INDEL
  gPastaBPState['7'] = INDEL
  gPastaBPState['E'] = INDEL

  gPastaBPState['Q'] = INDEL
  gPastaBPState['S'] = INDEL
  gPastaBPState['W'] = INDEL
  gPastaBPState['d'] = INDEL
  */

  gPastaBPState['!'] = DEL
  gPastaBPState['$'] = DEL
  gPastaBPState['7'] = DEL
  gPastaBPState['E'] = DEL
  gPastaBPState['z'] = DEL

  gPastaBPState['Q'] = INS
  gPastaBPState['S'] = INS
  gPastaBPState['W'] = INS
  gPastaBPState['d'] = INS
  gPastaBPState['Z'] = INS

  BPState = gPastaBPState

  SubMap = gSub
  RefMap = gRefBP
  AltMap = gAltBP

  RefDelBP = make(map[byte]int)
  for i:=0; i<len(Token); i++ {
    if _,ok := RefMap[Token[i]] ; ok {
      RefDelBP[Token[i]] = 1
    } else {
      RefDelBP[Token[i]] = 0
    }
  }

}



type PastaHandle struct {
  Fp *os.File
  Scanner *bufio.Scanner

  //Stream *simplestream.SimpleStream
  //AltStream *simplestream.SimpleStream

  Buf []byte
  Stage []byte
}

func InterleaveStreams(stream_A, stream_B io.Reader, w io.Writer) error {
  var e0, e1 error
  ref_pos := [2]int{0,0}
  stream_pos := [2]int{0,0} ; _ = stream_pos
  ch_val := [2]byte{0,0}
  dot := [1]byte{'.'}

  stream_a := bufio.NewReader(stream_A)
  stream_b := bufio.NewReader(stream_B)
  out := bufio.NewWriter(w)

  for {

    if ref_pos[0] == ref_pos[1] {

      for {
        ch_val[0],e0 = stream_a.ReadByte()
        if (e0==nil) {
          if (ch_val[0] == '>') {
            msg,e := ControlMessageProcess(stream_a)
            if e!=nil { return e }
            if msg.Type == POS { ref_pos[0] = msg.RefPos }
            continue
          } else if ch_val[0] == '\n' || ch_val[0] == ' ' || ch_val[0] == '\t' {
            continue
          }
        }
        break
      }

      for {
        ch_val[1],e1 = stream_b.ReadByte()
        if (e1==nil) {
          if (ch_val[1] == '>') {
            msg,e := ControlMessageProcess(stream_b)
            if e!=nil { return e }
            if msg.Type == POS { ref_pos[1] = msg.RefPos }
            continue
          } else if ch_val[1] == '\n' || ch_val[1] == ' ' || ch_val[1] == '\t' {
            continue
          }
        }
        break
      }


      stream_pos[0]++
      stream_pos[1]++
    } else if ref_pos[0] < ref_pos[1] {

      for {
        ch_val[0],e0 = stream_a.ReadByte()
        if (e0==nil) && (ch_val[0] == '>') {
          msg,e := ControlMessageProcess(stream_a)
          if e!=nil { return e }
          if msg.Type == POS { ref_pos[0] = msg.RefPos }
          continue
        }
        break
      }

      stream_pos[0]++
    } else if ref_pos[0] > ref_pos[1] {

      for {
        ch_val[1],e1 = stream_b.ReadByte()
        if (e1==nil) && (ch_val[1] == '>') {
          msg,e := ControlMessageProcess(stream_b)
          if e!=nil { return e }
          if msg.Type == POS { ref_pos[1] = msg.RefPos }
          continue
        }
        break
      }

      stream_pos[1]++
    }

    if e0!=nil && e1!=nil { break }

    if (ch_val[0]=='.') && (ch_val[1]=='.') { continue }

    //DEBUG
    /*
    if (ch_val[0] == 'a') && (ch_val[1] == 'Z') {
      fmt.Fprintf(os.Stderr, "bef>>>>> ref_pos %v, stm_pos %v, ch_val %v\n\n")
      os.Stderr.Sync()
    }
    */


    if ref_pos[0] == ref_pos[1] {

      if (ch_val[0]!='Q') && (ch_val[0]!='S') && (ch_val[0]!='W') && (ch_val[0]!='d') && (ch_val[0]!='Z') && (ch_val[0]!='.') && (ch_val[0]!='\n') && (ch_val[0]!=' ') {
        ref_pos[0]++
      }

      if (ch_val[1]!='Q') && (ch_val[1]!='S') && (ch_val[1]!='W') && (ch_val[1]!='d') && (ch_val[1]!='Z') && (ch_val[1]!='.') && (ch_val[1]!='\n') && (ch_val[1]!=' ') {
        ref_pos[1]++
      }

    } else if ref_pos[0] < ref_pos[1] {
      if (ch_val[0]!='Q') && (ch_val[0]!='S') && (ch_val[0]!='W') && (ch_val[0]!='d') && (ch_val[0]!='Z') && (ch_val[0]!='.') && (ch_val[0]!='\n') && (ch_val[0]!=' ') {
        ref_pos[0]++
      }
    } else if ref_pos[0] > ref_pos[1] {

      if (ch_val[1]!='Q') && (ch_val[1]!='S') && (ch_val[1]!='W') && (ch_val[1]!='d') && (ch_val[1]!='Z') && (ch_val[1]!='.') && (ch_val[1]!='\n') && (ch_val[1]!=' ') {
        ref_pos[1]++
      }
    }

    //DEBUG
    /*
    if (ch_val[0] == 'a') && (ch_val[1] == 'Z') {
      fmt.Fprintf(os.Stderr, "aft>>>>> ref_pos %v, stm_pos %v, ch_val %v\n\n")
      os.Stderr.Sync()
    }
    */

    if ref_pos[0]==ref_pos[1] {
      out.WriteByte(ch_val[0])
      out.WriteByte(ch_val[1])
    } else if ref_pos[0] < ref_pos[1] {
      out.WriteByte(ch_val[0])
      out.WriteByte(dot[0])
    } else if ref_pos[0] > ref_pos[1] {
      out.WriteByte(dot[0])
      out.WriteByte(ch_val[1])
    }

  }

  out.Flush()

  return nil
}



/*
func Open(fn string) (p PastaHandle, err error) {
  if fn == "-" {
    p.Fp = os.Stdin
  } else {
    p.Fp,err = os.Open(fn)
  }
  if err!=nil { return }

  p.Reader = bufio.NewReader(p.Fp)
  return p, nil
}

func (p *PastaHandle) Close() {
  p.Fp.Close()
}

func (p *PastaHandle) PeekChar() (byte) {
  if len(p.Stage)==0 { return 0 }
  return p.Stage[0]
}

PASTA_SAUCE := 1024

func (p *PastaHandle) ReadChar() (byte, err) {
  if len(p.Stage)==0 {
    if len(p.Buf)==0 {
      p.Buf = make([]byte, PASTA_SAUCE, PASTA_SAUCE)
    }
    n,e := p.Fp.Read(p.Buf)
    if e!=nil { return 0, e }
    if n==0 { return 0, nil }
    p.Stage = p.Buf[0:n]
  }

  b := p.Stage[0]
  p.Stage = p.Stage[1:]
  return b,nil
}
*/
