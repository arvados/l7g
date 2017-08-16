package gvcf

import "fmt"
import "strconv"
import "strings"
import "bufio"
//import "os"
import "io"
import "bytes"

import "time"

import "github.com/curoverse/l7g/go/pasta"

var VERSION string = "0.1.0"

type GVCFRefVarInfo struct {
  chrom string

  refseq string
  altseq []string
  vartype int
  ref_start int
  ref_len int

  stream_ref_pos int

  right_anchor bool
}

type GVCFRefVar struct {
  Type int
  MessageType int
  RefSeqFlag bool
  NocSeqFlag bool
  Out io.Writer
  Msg pasta.ControlMessage
  RefBP byte
  Allele int

  ChromStr string
  RefPos int

  OCounter int
  LFMod int

  PrintHeader bool
  Reference string
  DataSource string

  PrevRefBPStart byte
  PrevRefBPEnd byte

  VCFVer string

  Date time.Time

  Id string
  Qual string
  Filter string
  Info string
  Format string

  PrevStartRefBase byte
  PrevEndRefBase byte
  PrevStartRefPos int
  PrevRefLen int
  PrevVarType int

  State int

  StateHistory []GVCFRefVarInfo

  StreamRefPos int
}

func (g *GVCFRefVar) Init() {
  g.PrintHeader = true
  g.DataSource = "unknown"
  g.Reference = "unknown"

  g.ChromStr = "Unk"
  g.RefPos = 0
  g.Allele = 2

  g.OCounter = 0
  g.LFMod = 50

  g.VCFVer = "VCFv4.1"
  g.Date = time.Now()

  g.Id = "."
  g.Qual = "."
  g.Filter = ""
  g.Info = ""
  g.Format = "GT"

  g.StreamRefPos = 0

  g.State = pasta.BEG
}

func (g *GVCFRefVar) Chrom(chr string) { g.ChromStr = chr }
func (g *GVCFRefVar) Pos(pos int) { g.RefPos = pos }
func (g *GVCFRefVar) GetRefPos() int { return g.RefPos }
func (g *GVCFRefVar) Header(out *bufio.Writer) error {

  hdr := []string{};
  hdr = append(hdr, fmt.Sprintf("##fileformat=%s", g.VCFVer))
  hdr = append(hdr, fmt.Sprintf("##fileDate=%d%02d%02d", g.Date.Year(), g.Date.Month(), g.Date.Day()))
  hdr = append(hdr, fmt.Sprintf("##source=\"%s\"", g.DataSource))
  hdr = append(hdr, fmt.Sprintf("##reference=\"%s\"", g.Reference))
  hdr = append(hdr, "##FILTER=<ID=NOCALL,Description=\"Some or all of this record had no sequence calls\">")
  hdr = append(hdr, "##FORMAT=<ID=GT,Number=1,Type=String,Description=\"Genotype\">")
  hdr = append(hdr, "##INFO=<ID=END,Number=1,Type=Integer,Description=\"Stop position of the interval\">")
  hdr = append(hdr, "#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT\tSAMPLE")

  out.WriteString( strings.Join(hdr, "\n") + "\n" )

  return nil
}

//---

// 0      1     2   3   4   5    6      7    8      9
// chrom  pos   id  ref alt qual filter info format sample
//
func (g *GVCFRefVar) EmitLine(vartype int, vcf_ref_pos, vcf_ref_len int, vcf_ref_base byte, alt_field string, sample_field string, out *bufio.Writer) error {

  //info_field := fmt.Sprintf("END=%d", vcf_ref_pos+vcf_ref_len)
  info_field := fmt.Sprintf("END=%d", vcf_ref_pos+vcf_ref_len-1)

  //DEBUG
  fmt.Printf("emitline: vartype %v, chrom %s, vcf_ref_pos %v+%v, id %v, vcf_ref_base %v, info_field %v, sample_field %v\n",
    vartype, g.ChromStr, vcf_ref_pos, vcf_ref_len, g.Id, vcf_ref_base, info_field, sample_field)

  return nil
}

func (g *GVCFRefVar) _construct_alt_field() string {
  return ""
}

func (g *GVCFRefVar) _construct_sample_field() string {
  return ""
}

// return reference string, array of alt strings (unique) and the gt string (e.g. "0/0")
//
func (g *GVCFRefVar) _ref_alt_gt_fields(refseq string, altseq []string) (string,[]string,string) {
  local_debug := false
  _allele_n := 0

  _refseq := ""
  if len(refseq)>0 && refseq[0]!='-' {
    _refseq = string(refseq)
  }

  // Find unique altseqs (take out '-' if present)
  //

  gt_idx_str := []string{}
  gt_idx := []int{}

  if local_debug {
    fmt.Printf("# alt_gt_fields, refseq %s, altseq %v\n", refseq, altseq)
  }

  altseq_uniq := []string{}
  _set := make(map[string]int)
  _set[_refseq] = _allele_n
  _allele_n++
  var ts string
  var idx int

  if len(altseq) < g.Allele {
    n := g.Allele - len(altseq)
    for ii:=0; ii<n; ii++ {
      altseq = append(altseq, altseq[0])
    }
  }

  for ii:=0; ii<len(altseq); ii++ {
    if len(altseq[ii])==0 || altseq[ii][0] == '-' {
      ts = ""
    } else {
      ts = string(altseq[ii])
    }

    if _,ok := _set[ts] ; !ok {
      _set[ts] = _allele_n
      _allele_n++
      altseq_uniq = append(altseq_uniq, ts)
    }

    idx,_ = _set[ts]
    gt_idx_str = append(gt_idx_str, fmt.Sprintf("%d", idx))
    gt_idx = append(gt_idx, idx)
  }

  gt_field := strings.Join(gt_idx_str, "/") ; _ = gt_field

  return _refseq, altseq_uniq, gt_field
}



func (g *GVCFRefVar) _emit_alt_left_anchor(info GVCFRefVarInfo, out *bufio.Writer) {
  local_debug := false

  a_refseq,a_alt,a_gt_field := g._ref_alt_gt_fields(info.refseq, info.altseq)
  _ = a_alt

  a_start := info.ref_start+1
  a_len := info.ref_len

  alt_field := strings.Join(a_alt, ",")

  a_ref_bp := byte('.')
  //if len(a_refseq)>0 { a_ref_bp = a_refseq[0] }

  if local_debug {
    fmt.Printf("#  eala: stream_ref_pos: %d, a_refseq: %s\n", info.stream_ref_pos, a_refseq)
  }

  // experimental
  if info.stream_ref_pos == 0 {
    if len(a_refseq)>0 { a_ref_bp = a_refseq[ len(a_refseq)-1 ] }
  } else {
    if len(a_refseq)>0 { a_ref_bp = a_refseq[0] }
  }



  a_filt_field := "PASS"
  //a_info_field := fmt.Sprintf("END=%d", a_start+a_len)
  a_info_field := fmt.Sprintf("END=%d", a_start+a_len-1)

  //experimental
  if info.stream_ref_pos == 0 {
    a_info_field += fmt.Sprintf(":REF_ANCHOR_AT_END=TRUE")
  }

  if info.vartype == pasta.NOC {
    a_filt_field = "NOCALL"
  }



  //                            0   1   2   3   4   5    6  7   8   9
  out.WriteString( fmt.Sprintf("%s\t%d\t%s\t%c\t%s\t%s\t%s\t%s\t%s\t%s\n",
    g.ChromStr,
    a_start,
    g.Id,
    a_ref_bp,
    alt_field,
    g.Qual,
    a_filt_field,
    a_info_field,
    g.Format,
    a_gt_field) )

}

func (g *GVCFRefVar) _emit_ref_left_anchor(info GVCFRefVarInfo, out *bufio.Writer) {
  a_start := info.ref_start+1
  a_len := info.ref_len
  a_r_seq := info.refseq
  a_ref_bp := byte('.')
  if len(a_r_seq)>0 { a_ref_bp = a_r_seq[0] }
  a_gt_field := "0/0"

  a_filt_field := "PASS"
  //a_info_field := fmt.Sprintf("END=%d", a_start+a_len)
  a_info_field := fmt.Sprintf("END=%d", a_start+a_len-1)

  //                            0   1   2   3   4   5    6  7   8   9
  out.WriteString( fmt.Sprintf("%s\t%d\t%s\t%c\t%s\t%s\t%s\t%s\t%s\t%s\n",
    g.ChromStr,
    a_start,
    g.Id,
    a_ref_bp,
    ".",
    g.Qual,
    a_filt_field,
    a_info_field,
    g.Format,
    a_gt_field) )

}


// _p for 'peel'.  Since GVCF insists on asking for a reference base before
// the event, we need to do contortions to save the relevant information,
// this function being one of them.
//
// This funciton differs from the above in that it allows the reference character
// to be passed in as is added to the appropriate place.
//
func (g *GVCFRefVar) _emit_alt_left_anchor_p(info GVCFRefVarInfo, z byte, out *bufio.Writer, del_start int) {
  local_debug := false

  b_r_seq := info.refseq
  b_refseq,b_alt,b_gt_field := g._ref_alt_gt_fields(b_r_seq, info.altseq)

  _ = b_refseq

  if local_debug { fmt.Printf("#  ealap: stream_ref_pos: %d, z: %c, refseq: %s", info.stream_ref_pos, z, info.refseq) }

  _a := []string{}
  for ii:=0; ii<len(b_alt); ii++ {

    if info.stream_ref_pos == 0 {
      _a = append(_a, fmt.Sprintf("%s%c", b_alt[ii], z))
    } else {
      _a = append(_a, fmt.Sprintf("%c%s", z, b_alt[ii]))
    }
  }
  b_alt_field := strings.Join(_a, ",")

  b_start := info.ref_start + del_start
  b_len := info.ref_len+1
  b_ref_bp := z
  b_filt_field := "PASS"
  //b_info_field := fmt.Sprintf("END=%d", b_start+b_len)
  b_info_field := fmt.Sprintf("END=%d", b_start+b_len-1)

  // In this case, we have no other choice but to put the
  // reference anchor base at the end, violating the
  // VCF spec.  Indicate it with this INFO field in the
  // hopes that whoever downstream runs into this will
  // be able to parse it.
  //
  if info.stream_ref_pos == 0 {
    b_info_field += fmt.Sprintf(":REF_ANCHOR_AT_END=TRUE")
  }


  //                            0   1   2   3   4   5    6  7   8   9
  out.WriteString( fmt.Sprintf("%s\t%d\t%s\t%c\t%s\t%s\t%s\t%s\t%s\t%s\n",
    g.ChromStr,
    b_start,
    g.Id,
    b_ref_bp,
    b_alt_field,
    g.Qual,
    b_filt_field,
    b_info_field,
    g.Format,
    b_gt_field) )

}

// right anchor
//
func (g *GVCFRefVar) _emit_alt_right_anchor(info GVCFRefVarInfo, z byte, out *bufio.Writer) {
  local_debug := false

  b_r_seq := info.refseq
  b_refseq,b_alt,b_gt_field := g._ref_alt_gt_fields(b_r_seq, info.altseq)

  if local_debug {
    fmt.Printf("#  eara: z: %c, b_refseq %v, b_alt %v, b_gt_field %v, b_r_seq %v, info.altseq %v\n",
      z, b_refseq, b_alt, b_gt_field, b_r_seq, info.altseq) }

  _ = b_refseq

  _a := []string{}
  for ii:=0; ii<len(b_alt); ii++ {
    _a = append(_a, fmt.Sprintf("%s%c", b_alt[ii], z))
  }
  b_alt_field := strings.Join(_a, ",")

  if local_debug { fmt.Printf("#  eara: stream_pos %d, ref_start %d, z %c..\n", info.stream_ref_pos, info.ref_start, z) }

  b_start := info.ref_start+ + 1
  b_len := info.ref_len+1
  b_ref_bp := z
  b_filt_field := "PASS"
  b_info_field := fmt.Sprintf("END=%d", b_start+b_len-1)

  // In this case, we have no other choice but to put the
  // reference anchor base at the end, violating the
  // VCF spec.  Indicate it with this INFO field in the
  // hopes that whoever downstream runs into this will
  // be able to parse it.
  //
  b_info_field += fmt.Sprintf(":REF_ANCHOR_AT_END=TRUE")

  //                            0   1   2   3   4   5    6  7   8   9
  out.WriteString( fmt.Sprintf("%s\t%d\t%s\t%c\t%s\t%s\t%s\t%s\t%s\t%s\n",
    g.ChromStr,
    b_start,
    g.Id,
    b_ref_bp,
    b_alt_field,
    g.Qual,
    b_filt_field,
    b_info_field,
    g.Format,
    b_gt_field) )

}

// (g)VCF lines consist of:
//
// 0      1     2   3   4   5    6      7    8      9
// chrom  pos   id  ref alt qual filter info format sample
//
// Print receives interpreted 'difference stream' lines, one at a time.
//
// We make a simplifying assumption that if there is a nocall region right next to
// an alt call, the alt call gets subsumed into the nocall region.
//
// We print the nocall region with full sequence so that it's recoverable but otherwise it
// looks like a nocall region.
//
// This function got unfortunately quite complicated.  The 'difference stream' that gets
// fed into this function is reporting alternates, nocalls, indels, etc. with only the minimal
// amount of information, not reporting the reference base before or after it.  This means
// we have to keep state in order to report the reference anchor base as the VCF specification
// requires.  In many cases, we have no choice but to violate the VCF spec because we don't
// have information about the reference base that came before the current alternate in question.
// Under this condition, we report the right anchor reference base and put a field in the
// INFO column to indcate that we've done so.
//
// The general scheme is to save state in a structure called `StateHistory`.  We consider
// all pair transitions ( {ALT,NOC,REF} -> {ALT,NOC,REF} ) to determine whether we can emit
// a line and what to emit if we can.
//
// The easy (and hopefully common) case is when there is a REF line followed by a non-REF line.
// In that case, we can emit the REF line before it, peeling off the last reference base to
// use as an anchor for the ALT line if need be (for example, when the ALT line is a straight
// deletion).  Problems arise if there are ALT lines without any REF lines in between them,
// which shouldn't happen if the 'difference stream' is working properly, or, more likely,
// if the difference stream begins on an ALT or NOC.  This special case of the beginning
// stream causes most of the complexity below.
//
// For example, if the first difference line is a deletion, the anchor reference base needs
// to be taken from the following REF line.  If the next REF line has more than one reference
// base, then we can peel it off the beginning, use it as a right anchor in the current
// reported ALT line and promote the REF line to be the base element in the `StateHistory`
// structure.  If the next REF line only has one reference base, then we could end up
// taking a reference base that could be used in subsequent ALT or NOC lines further on
// down the stream.  In order to reduce this cascading domino effect, the ALT line is
// extended in this case.
//
// Since this function is to be used with streams that don't need to start at reference
// position 1, reporting an achor reference base that isn't to the left of the ALT line
// is straight away violating the VCF sepcification.  There's not much we can do since
// we want to report gVCF for arbitrary sequences.  With the VCF specification as stated,
// it's impossible to report arbitrary sequence information with rigid fixed endpoints.
// Instead we make due by occasionally reporting a right anchor point and giving a
// field in the INFO column of `REF_ANCHOR_AT_END=TRUE`.
//
//
func (g *GVCFRefVar) Print(vartype int, ref_start, ref_len int, refseq []byte, altseq [][]byte, out *bufio.Writer) error {
  local_debug := false

  if g.PrintHeader {
    g.Header(out)
    g.PrintHeader = false
  }

  vi := GVCFRefVarInfo{}
  vi.vartype = vartype
  vi.ref_start = ref_start
  vi.ref_len = ref_len
  vi.refseq = string(refseq)
  for ii:=0; ii<len(altseq); ii++ {
    vi.altseq = append(vi.altseq, string(altseq[ii]))
  }
  vi.chrom = g.ChromStr
  vi.stream_ref_pos = g.StreamRefPos

  g.StreamRefPos += ref_len

  g.StateHistory = append(g.StateHistory, vi)

  processing:=false
  if len(g.StateHistory)>1 { processing = true }

  if local_debug {
    fmt.Printf("#\n")
    fmt.Printf("# vartype: %d (REF %d, NOC %d, ALT %d)\n", vartype, pasta.REF, pasta.NOC, pasta.ALT)
    fmt.Printf("# ref_start: %d, ref_len: %d\n", ref_start, ref_len)
    fmt.Printf("# refseq: %s\n", refseq)
    fmt.Printf("# altseq: %s\n", altseq)
    fmt.Printf("# StateHistory: %v\n", g.StateHistory)
  }

  for processing && (len(g.StateHistory)>1) {

    if local_debug { fmt.Printf("#  cp1\n") }

    idx:=1

    if g.StateHistory[idx-1].vartype == pasta.REF  {

      if g.StateHistory[idx].vartype==pasta.REF {

        g._emit_ref_left_anchor(g.StateHistory[idx-1], out)
        g.StateHistory = g.StateHistory[idx:]
        continue

      } else if g.StateHistory[idx].vartype==pasta.NOC {

        b_ref,b_alt,_ := g._ref_alt_gt_fields(g.StateHistory[idx].refseq, g.StateHistory[idx].altseq)

        // b_alt == 0 -> it's a nocall for both reference and alt
        //
        min_alt_len := 0
        if len(b_alt)>0 {
          min_alt_len = len(b_alt[0])
          for ii:=1; ii<len(b_alt); ii++ {
            if min_alt_len > len(b_alt[ii]) {
              min_alt_len = len(b_alt[ii])
            }
          }
        }

        if min_alt_len>0 {

          // The nocall alt will have a reference anchor so we can
          // emit the current reference
          //
          g._emit_ref_left_anchor(g.StateHistory[idx-1], out)
          g.StateHistory = g.StateHistory[idx:]

        } else {

          n := len(g.StateHistory[idx-1].refseq)
          z := g.StateHistory[idx-1].refseq[n-1]

          g.StateHistory[idx-1].refseq = g.StateHistory[idx-1].refseq[0:n-1]
          g.StateHistory[idx-1].ref_len--
          if g.StateHistory[idx-1].ref_len > 0 {
            g._emit_ref_left_anchor(g.StateHistory[idx-1], out)
          }

          g.StateHistory[idx].refseq = string(z) + b_ref
          g.StateHistory[idx].ref_start--
          g.StateHistory[idx].ref_len++
          for ii:=0; ii<len(g.StateHistory[idx].altseq); ii++ {
            if g.StateHistory[idx].altseq[ii] == "-" {
              g.StateHistory[idx].altseq[ii] = string(z)
            } else {
              g.StateHistory[idx].altseq[ii] = string(z) + g.StateHistory[idx].altseq[ii]
            }
          }

          g.StateHistory = g.StateHistory[idx:]
        }

        continue

      } else if g.StateHistory[idx].vartype==pasta.ALT {

        _,b_alt,_ := g._ref_alt_gt_fields(g.StateHistory[idx].refseq, g.StateHistory[idx].altseq)

        min_alt_len := len(b_alt[0])
        for ii:=1; ii<len(b_alt); ii++ {
          if min_alt_len > len(b_alt[ii]) { min_alt_len = len(b_alt[ii]) }
        }

        if (ref_len>0) && (min_alt_len>0) {

          // Not a straight deletion, we can use a reference base
          // as anchor straight out
          //
          g._emit_ref_left_anchor(g.StateHistory[idx-1], out)
          g._emit_alt_left_anchor(g.StateHistory[idx], out)
          g.StateHistory = g.StateHistory[idx+1:]

        } else {

          // The alt is a straight deletion or insertion
          // so we need to peel off a reference base from
          // the previous line and use it as the anchor base.
          //
          n := len(g.StateHistory[idx-1].refseq)
          z := g.StateHistory[idx-1].refseq[n-1]

          g.StateHistory[idx-1].refseq = g.StateHistory[idx-1].refseq[0:n-1]
          g.StateHistory[idx-1].ref_len--
          if g.StateHistory[idx-1].ref_len > 0 {
            g._emit_ref_left_anchor(g.StateHistory[idx-1], out)
          }
          g._emit_alt_left_anchor_p(g.StateHistory[idx], z, out, 0)

          g.StateHistory = g.StateHistory[idx+1:]
        }

      }

    } else if g.StateHistory[idx-1].vartype == pasta.ALT {

      if g.StateHistory[idx].vartype == pasta.REF {

        // In an ALT to REF transition, if the previous
        // ALT line has reference sequence and alternate
        // sequences, we can emit them straight away,
        // keeping the current REF line untouched.
        //
        // If the previous ALT line has no reference
        // sequence, we need to erode the current REF
        // line in order to try to comple with the VCF standard.
        // We'll be violating the standard in this case
        // since the reference base will be at the end of
        // the line instead of the beginning.  We can't
        // do anything except hope to mark it as such
        // with an extra INFO field and hope people using
        // it downstream either don't care to check the refernce
        // base reported or know to look at the INFO field.
        //

        _,a_alt,_ := g._ref_alt_gt_fields(g.StateHistory[idx-1].refseq, g.StateHistory[idx-1].altseq)
        prv_min_alt_len := len(a_alt[0])
        for ii:=1; ii<len(a_alt); ii++ {
          if prv_min_alt_len > len(a_alt[ii]) { prv_min_alt_len = len(a_alt[ii]) }
        }
        prv_ref_len := g.StateHistory[idx-1].ref_len

        if (g.StateHistory[idx-1].stream_ref_pos!=0) {

          // The previous state was ALT with no reference,
          // the current state is REF.  There are two
          // conditions:
          //
          // * The current REF line is more than one base long
          // * The current REF line is one base long
          //
          // In the case of the REF line being more than one
          // base long, we can peel off one of the bases from
          // the front of the reference sequence from the REF
          // line, report the ALT line and continue on.
          //
          // In the case of the REF line being one base long,
          // we subsume the REF line into the ALT and continue
          // on.  The problem with reporting the line straight
          // away is that if there's an ALT line after it,
          // we'll be forcing it to report the reference base
          // at the end.  This way we try to limit the
          // violation of the VCF by creating as few (only one?)
          // ALT lines that require reporting the reference
          // base at the end.
          //

          // This should be the 'normal' case.
          //

          if prv_ref_len>0 {

            // We have reference in the previous alt, emit the line, move on
            //

            g._emit_alt_left_anchor(g.StateHistory[idx-1], out)
            g.StateHistory = g.StateHistory[idx:]
            continue

          } else if len(g.StateHistory[idx].refseq)>1 {

            // Else the previous ALT doesn't have any reference for some reason,
            // emit the previous LAT line with a right reference base anchor
            // and update the REF line as appropriate
            //

            z := g.StateHistory[idx].refseq[0]
            g._emit_alt_right_anchor(g.StateHistory[idx-1], z, out)

            g.StateHistory[idx].refseq = g.StateHistory[idx].refseq[1:]
            g.StateHistory[idx].ref_len--
            g.StateHistory[idx].ref_start++
            g.StateHistory[idx].stream_ref_pos++
            g.StateHistory = g.StateHistory[idx:]
            continue

          } else {

            // Else subsume the reference line into the ALT
            //
            //
            g.StateHistory[idx].ref_start = g.StateHistory[idx-1].ref_start
            g.StateHistory[idx].ref_len += g.StateHistory[idx-1].ref_len
            g.StateHistory[idx].vartype = g.StateHistory[idx-1].vartype

            //g.StateHistory[idx].altseq = g.StateHistory[idx-1].altseq
            g.StateHistory[idx].altseq = g.StateHistory[idx].altseq[0:0]
            for ii:=0; ii<len(g.StateHistory[idx-1].altseq) ; ii++ {
              if g.StateHistory[idx-1].altseq[ii] != "-" {
                g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, g.StateHistory[idx-1].altseq[ii] + g.StateHistory[idx].refseq)
              } else {
                g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, g.StateHistory[idx].refseq)
              }

            }

            g.StateHistory[idx].stream_ref_pos = g.StateHistory[idx-1].stream_ref_pos

            g.StateHistory = g.StateHistory[idx:]
            continue

          }

        } else {

          if local_debug { fmt.Printf("#  ref->alt @ stream_pos 0\n") }

          // The special case that we're beginning our stream
          //

          if len(g.StateHistory[idx].refseq)>1 {

            if local_debug { fmt.Printf("#  ref->alt @ sp0: len(refseq) > 1 (%d)\n", len(g.StateHistory[idx].refseq)) }

            // We have enough reference to pop off the beginning
            // reference base and keep the rest of the REF line.
            //

            z := g.StateHistory[idx].refseq[0]
            g._emit_alt_right_anchor(g.StateHistory[idx-1], z, out)

            g.StateHistory[idx].refseq = g.StateHistory[idx].refseq[1:]
            g.StateHistory[idx].ref_len--
            g.StateHistory[idx].ref_start++
            g.StateHistory[idx].stream_ref_pos++
            g.StateHistory = g.StateHistory[idx:]
            continue

          } else {

            if local_debug { fmt.Printf("#  ref->alt @ sp0: len(refseq) <= 1 (%d)\n", len(g.StateHistory[idx].refseq)) }

            // subsume the reference into the ALT line
            //
            g.StateHistory[idx].altseq = g.StateHistory[idx].altseq[0:0]
            for ii:=0; ii<len(g.StateHistory[idx-1].altseq) ; ii++ {
              if g.StateHistory[idx-1].altseq[ii] != "-" {
                g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, g.StateHistory[idx-1].altseq[ii] + g.StateHistory[idx].refseq)
              } else {
                g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, g.StateHistory[idx].refseq)
              }
            }

            g.StateHistory[idx].stream_ref_pos = g.StateHistory[idx-1].stream_ref_pos

            g.StateHistory[idx].refseq = g.StateHistory[idx-1].refseq + g.StateHistory[idx].refseq
            g.StateHistory[idx].ref_start = g.StateHistory[idx-1].ref_start
            g.StateHistory[idx].ref_len += g.StateHistory[idx-1].ref_len
            g.StateHistory[idx].vartype = g.StateHistory[idx-1].vartype

            g.StateHistory = g.StateHistory[idx:]
            continue

          }

        }

      } else if g.StateHistory[idx].vartype == pasta.ALT {

        if local_debug { fmt.Printf("#  alt->alt\n") }

        // construct the subsumed alternate sequence for each allele
        //
        prv_alt_len := len(g.StateHistory[idx-1].altseq)
        cur_alt_len := len(g.StateHistory[idx].altseq)
        if ((prv_alt_len!=0) && (cur_alt_len!=0) && (cur_alt_len!=prv_alt_len)) {
          return fmt.Errorf(
            fmt.Sprintf("valid alternate sequences lists must have matching lengths (%v != %v) at position %v:%d",
              prv_alt_len, cur_alt_len, g.StateHistory[idx-1].chrom, g.StateHistory[idx-1].ref_start))
        }
        alt_len := prv_alt_len
        if alt_len < cur_alt_len { alt_len = cur_alt_len }
        alt_seqs := []string{}
        for ii:=0 ; ii<alt_len; ii++ { alt_seqs = append(alt_seqs, "") }
        for ii:=0; ii<alt_len; ii++ {
          if (prv_alt_len > 0) && (g.StateHistory[idx-1].altseq[ii] != "-") {
            alt_seqs[ii] = alt_seqs[ii] + g.StateHistory[idx-1].altseq[ii]
          }

          if (cur_alt_len > 0) && (g.StateHistory[idx].altseq[ii] != "-") {
            alt_seqs[ii] = alt_seqs[ii] + g.StateHistory[idx].altseq[ii]
          }

        }

        if local_debug { fmt.Printf("#  alt->alt: alt_seqs %v\n", alt_seqs) }

        // Subsume the previous ALT into the current NOC entry
        //
        g.StateHistory[idx].ref_start = g.StateHistory[idx-1].ref_start
        g.StateHistory[idx].ref_len += g.StateHistory[idx-1].ref_len
        g.StateHistory[idx].stream_ref_pos = g.StateHistory[idx-1].stream_ref_pos

        g.StateHistory[idx].altseq = g.StateHistory[idx].altseq[0:0]
        for ii:=0; ii<len(alt_seqs); ii++ {
          g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, alt_seqs[ii])
        }

        new_ref_seq := g.StateHistory[idx-1].refseq
        if new_ref_seq == "-" { new_ref_seq = "" }
        if g.StateHistory[idx].refseq != "-" {
          new_ref_seq += g.StateHistory[idx].refseq
        }
        if len(new_ref_seq)==0 { new_ref_seq = "-" }
        g.StateHistory[idx].refseq = new_ref_seq

        g.StateHistory = g.StateHistory[idx:]

        continue

      } else if g.StateHistory[idx].vartype == pasta.NOC {

        // construct the subsumed alternate sequence for each allele
        //
        prv_alt_len := len(g.StateHistory[idx-1].altseq)
        cur_alt_len := len(g.StateHistory[idx].altseq)
        if ((prv_alt_len!=0) && (cur_alt_len!=0) && (cur_alt_len!=prv_alt_len)) {
          return fmt.Errorf(
            fmt.Sprintf("valid alternate sequences lists must have matching lengths (%v != %v) at position %v:%d",
              prv_alt_len, cur_alt_len, g.StateHistory[idx-1].chrom, g.StateHistory[idx-1].ref_start) )
        }
        alt_len := prv_alt_len
        if alt_len < cur_alt_len { alt_len = cur_alt_len }
        alt_seqs := []string{}
        for ii:=0 ; ii<alt_len; ii++ { alt_seqs = append(alt_seqs, "") }
        for ii:=0; ii<alt_len; ii++ {
          if (prv_alt_len > 0) && (g.StateHistory[idx-1].altseq[ii] != "-") {
            alt_seqs[ii] = alt_seqs[ii] + g.StateHistory[idx-1].altseq[ii]
          }

          if (cur_alt_len > 0) && (g.StateHistory[idx].altseq[ii] != "-") {
            alt_seqs[ii] = alt_seqs[ii] + g.StateHistory[idx].altseq[ii]
          }

        }

        //if local_debug { fmt.Printf("a_alt: %v, refa: %v, gta: %v\nb_alt: %v, refb: %v, gtb: %v\n", a_alt, refa, gta, b_alt, refb, gtb) }

        // Subsume the previous ALT into the current NOC entry
        //
        g.StateHistory[idx].ref_start = g.StateHistory[idx-1].ref_start
        g.StateHistory[idx].ref_len += g.StateHistory[idx-1].ref_len
        g.StateHistory[idx].stream_ref_pos = g.StateHistory[idx-1].stream_ref_pos

        g.StateHistory[idx].altseq = g.StateHistory[idx].altseq[0:0]
        /*
        for ii:=0; ii<len(a_alt); ii++ {

          if local_debug {
            fmt.Printf("ALT->NOC adding altseq %v\n", string(a_alt[ii]) + string(b_alt[ii]))
          }

          g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, string(a_alt[ii]) + string(b_alt[ii]))
        }
        */

        for ii:=0; ii<len(alt_seqs); ii++ {
          g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, alt_seqs[ii])
        }


        g.StateHistory = g.StateHistory[idx:]

        if local_debug { fmt.Printf("#  StateHistory %v\n", g.StateHistory) }

        continue
      }

    } else if g.StateHistory[idx-1].vartype == pasta.NOC {

      if g.StateHistory[idx].vartype == pasta.REF {

        if local_debug {
          fmt.Printf("# cp (noc->ref)\n")
        }

        g._emit_alt_left_anchor(g.StateHistory[idx-1], out)
        g.StateHistory = g.StateHistory[idx:]
        continue

      } else if g.StateHistory[idx].vartype == pasta.ALT {

        a_seqs := []string{}
        b_seqs := []string{}

        for ii:=0; ii<len(g.StateHistory[idx-1].altseq); ii++ {
          if g.StateHistory[idx-1].altseq[ii] == "-" {
            a_seqs = append(a_seqs, "")
          } else {
            a_seqs = append(a_seqs, g.StateHistory[idx-1].altseq[ii])
          }
        }

        for ii:=0; ii<len(g.StateHistory[idx].altseq); ii++ {
          if g.StateHistory[idx].altseq[ii] == "-" {
            b_seqs = append(b_seqs, "")
          } else {
            b_seqs = append(b_seqs, g.StateHistory[idx].altseq[ii])
          }
        }

        // Subsume the previous ALT into the current NOC entry
        //
        g.StateHistory[idx].ref_start = g.StateHistory[idx-1].ref_start
        g.StateHistory[idx].ref_len += g.StateHistory[idx-1].ref_len
        g.StateHistory[idx].vartype = g.StateHistory[idx-1].vartype
        g.StateHistory[idx].stream_ref_pos = g.StateHistory[idx-1].stream_ref_pos

        ref_b_pos := 0
        ref_b := make([]byte, len(g.StateHistory[idx-1].refseq) + len(g.StateHistory[idx].refseq))
        for ii:=0; ii<len(g.StateHistory[idx-1].refseq); ii++ {
          if (g.StateHistory[idx-1].refseq[ii] != '-') {
            ref_b[ref_b_pos] = g.StateHistory[idx-1].refseq[ii]
            ref_b_pos++
          }
        }

        for ii:=0; ii<len(g.StateHistory[idx].refseq); ii++ {
          if (g.StateHistory[idx].refseq[ii] != '-') {
            ref_b[ref_b_pos] = g.StateHistory[idx].refseq[ii]
            ref_b_pos++
          }
        }
        g.StateHistory[idx].refseq = string(ref_b[:ref_b_pos])


        g.StateHistory[idx].altseq = []string{}
        for ii:=0; ii<len(a_seqs); ii++ {
          g.StateHistory[idx].altseq = append(g.StateHistory[idx].altseq, a_seqs[ii] + b_seqs[ii])
        }

        g.StateHistory = g.StateHistory[idx:]
        continue

      } else if g.StateHistory[idx].vartype == pasta.NOC {

        _,a_alt,_ := g._ref_alt_gt_fields(g.StateHistory[idx-1].refseq, g.StateHistory[idx-1].altseq)

        if len(a_alt)==0 {
          g._emit_alt_left_anchor(g.StateHistory[idx-1], out)
          g.StateHistory = g.StateHistory[idx:]
          continue
        }

        min_length:=len(a_alt[0])
        for ii:=1; ii<len(a_alt); ii++ {
          if min_length < len(a_alt[ii]) { min_length = len(a_alt[ii]) }
        }

        if min_length == 0 {
          g.StateHistory = g.StateHistory[idx:]
          continue
        }

        g._emit_alt_left_anchor(g.StateHistory[idx-1], out)
        g.StateHistory = g.StateHistory[idx:]
        continue

      }

    }

    if len(g.StateHistory) < 2 { processing = false }
  }

  return nil
}

func process_ref_alt_seq(refseq []byte, altseq [][]byte) (string,bool) {
  var type_str string
  noc_flag := false
  indel_flag := false
  n1 := []byte{'n'}

  len_match := true
  for ii:=0; ii<len(altseq); ii++ {
    if len(altseq[ii])!=len(refseq) {
      len_match = false
      break
    }
  }

  if (len(refseq)==1) && len_match {
    for ii:=0; ii<len(altseq); ii++ {
      if altseq[ii][0]=='-' { indel_flag = true; break }
    }
  }

  if len_match && (len(refseq)==1) {
    if indel_flag || (refseq[0]=='-') {
      type_str = "INDEL"
    } else {
      type_str = "SNP"
    }
  } else if len_match {
    type_str = "SUB"

    // In the case:
    // * it's a non 0-length string
    // * the lengths of the altseqs match the refseq
    // * the altseqs are all 'n' (nocall)
    // -> it's a 'true' nocall line
    //
    if len(refseq)>0 {
      noc_flag = true
      for a:=0; a<len(altseq); a++ {
        n := bytes.Count(altseq[a], n1)
        if n!=len(altseq[a]) {
          noc_flag = false
          break
        }
      }
      if noc_flag { type_str = "NOC" }
    }
  } else {
    type_str = "INDEL"
  }

  return type_str, noc_flag
}


func (g *GVCFRefVar) PrintEnd(out *bufio.Writer) error {

  idx:=0

  if len(g.StateHistory)==0 { return nil }

  if g.StateHistory[idx].vartype==pasta.REF {
    g._emit_ref_left_anchor(g.StateHistory[idx], out)
  } else if g.StateHistory[idx].vartype==pasta.NOC {
    g._emit_alt_left_anchor(g.StateHistory[idx], out)
  } else if g.StateHistory[idx].vartype==pasta.ALT {
    g._emit_alt_left_anchor(g.StateHistory[idx], out)
  }

  out.Flush()

  return nil
}

//---

func (g *GVCFRefVar) PastaBegin(out *bufio.Writer) error {
  return nil
}

func (g *GVCFRefVar) PastaEnd(out *bufio.Writer) error {

  out.Flush()
  return nil
}

func (g *GVCFRefVar) _parse_info_field_value(info_line string, field string, sep string) (string, error) {
  sa := strings.Split(info_line, sep)
  for ii:=0; ii<len(sa); ii++ {
    fv := strings.Split(sa[ii], "=")
    if len(fv)!=2 { return "", fmt.Errorf("invalud field") }

    if fv[0] == field { return fv[1], nil }
  }
  return "", fmt.Errorf("field not found")
}

func (g *GVCFRefVar) _parameter_index(line string, field string, sep string) (int, error) {
  sa := strings.Split(line, sep)
  for ii:=0; ii<len(sa); ii++ {
    if sa[ii] == field { return ii, nil }
  }
  return -1, fmt.Errorf("field not found")
}

func (g *GVCFRefVar) _get_gt_array(gt_str string, ploidy int) ([]int, error) {
  gt_array := []int{}
  if !strings.ContainsAny(gt_str, "|/") {
    v,e := strconv.Atoi(gt_str)
    if e!=nil { return nil, e }
    gt_array = append(gt_array, v)
    return gt_array, nil
  }

  _sa := strings.Split(gt_str, "/")
  if len(_sa)==1 {
    _sa = strings.Split(gt_str, "|")
  }

  if len(_sa)>ploidy { return nil, fmt.Errorf("invalid GT field") }

  for ii:=0; ii<ploidy; ii++ {
    if ii < len(_sa) {
      v,e := strconv.Atoi(_sa[ii])
      if e!=nil { return nil, e }
      gt_array = append(gt_array, v)
    } else {
      gt_array = append(gt_array, gt_array[ii-1])
    }
  }

  return gt_array, nil
}

func (g *GVCFRefVar) Pasta(gvcf_line string, ref_stream *bufio.Reader, out *bufio.Writer) error {
  var err error
  CHROM_FIELD_POS := 0 ; _ = CHROM_FIELD_POS
  START_FIELD_POS := 1 ; _ = START_FIELD_POS
  ID_FIELD_POS := 2 ; _ = ID_FIELD_POS
  REF_FIELD_POS := 3 ; _ = REF_FIELD_POS
  ALT_FIELD_POS := 4 ; _ = ALT_FIELD_POS
  QUAL_FIELD_POS := 5 ; _ = QUAL_FIELD_POS
  FILTER_FIELD_POS := 6 ; _ = FILTER_FIELD_POS
  INFO_FIELD_POS := 7 ; _ = INFO_FIELD_POS
  FORMAT_FIELD_POS := 8 ; _ = FORMAT_FIELD_POS
  SAMPLE0_FIELD_POS := 9 ; _ = SAMPLE0_FIELD_POS

  // empty line or comment
  //
  if (len(gvcf_line)==0) || (gvcf_line[0]=='#') { return nil }


  line_part := strings.Split(gvcf_line, "\t")

  _start,e := strconv.Atoi(line_part[START_FIELD_POS])
  if e!=nil { return e }

  _end_str,e := g._parse_info_field_value(line_part[INFO_FIELD_POS], "END", ":")
  _end := -1
  if e==nil {
    _end,err = strconv.Atoi(_end_str)
    if err!=nil { return err }
  }
  if _end==-1 { _end = _start+1 }

  ref_anchor_on_left := true
  _,er := g._parse_info_field_value(line_part[INFO_FIELD_POS], "REF_ANCHOR_AT_END", ":")
  if er==nil {
    ref_anchor_on_left = false
  }

  alt_seq := []string{}
  if line_part[ALT_FIELD_POS]!="." {
    alt_seq = strings.Split(line_part[ALT_FIELD_POS], ",")
  }

  gt_samp_idx,e := g._parameter_index(line_part[FORMAT_FIELD_POS], "GT", ":")
  if e!=nil { return e }

  samp_part := strings.Split(line_part[SAMPLE0_FIELD_POS], ":")
  if gt_samp_idx >= len(samp_part) { return fmt.Errorf("GT index overflow") }

  n_allele := 2
  samp_str := samp_part[gt_samp_idx]
  samp_seq_idx,e := g._get_gt_array(samp_str, n_allele)
  if e!=nil { return e }

  ref_anchor_base := line_part[REF_FIELD_POS]
  //refn := _end - _start
  refn := (_end + 1) - _start

  if (samp_seq_idx[0] == samp_seq_idx[1]) && (samp_seq_idx[0] == 0) {

    for ii:=0; ii<refn; ii++ {
      stream_ref_bp,e := ref_stream.ReadByte()
      if e!=nil { return e }
      for stream_ref_bp == '\n' || stream_ref_bp == ' ' || stream_ref_bp == '\t' || stream_ref_bp == '\r' {
        stream_ref_bp,e = ref_stream.ReadByte()
        if e!=nil { return e }
      }


      for a:=0; a<n_allele; a++ {

        if (g.LFMod>0) && (g.OCounter > 0) && ((g.OCounter%g.LFMod)==0) {
          out.WriteByte('\n')
        }
        g.OCounter++

        out.WriteByte(stream_ref_bp)
      }

    }

    return nil
  }

  mM := refn
  for ii:=0; ii<n_allele; ii++ {

    // reference
    //
    if samp_seq_idx[ii]==0 { continue }

    // find maximum of alt sequence lengths
    //
    a_idx := samp_seq_idx[ii]-1
    if mM < len(alt_seq[a_idx]) { mM = len(alt_seq[a_idx]) }
  }



  // Loop through, emitting the appropriate substitution
  // if we have a reference, a deletion if the alt sequence
  // has run out or an insertion if the reference sequence has
  // run out.
  //
  // The reference is 'shifted' to the left, which means there
  // will be (potentially 0-length) substitutions followed by
  // (potentially 0-length) insertions and/or deletions.
  //
  for i:=0; i<mM; i++  {

    // Get the reference base
    //
    var stream_ref_bp byte
    if i<refn {

      stream_ref_bp,e = ref_stream.ReadByte()
      if e!=nil { return e }
      for stream_ref_bp == '\n' || stream_ref_bp == ' ' || stream_ref_bp == '\t' || stream_ref_bp == '\r' {
        stream_ref_bp,e = ref_stream.ReadByte()
        if e!=nil { return e }
      }

    }


    if ref_anchor_on_left {
      if (refn>0) && (i==0) && (stream_ref_bp!=ref_anchor_base[0]) {
        return fmt.Errorf(fmt.Sprintf("stream reference (%c) does not match VCF ref base (%c) at position %d\n", stream_ref_bp, ref_anchor_base[0], _start))
      }
    }
    _ = stream_ref_bp

    // Emit a symbol per alt sequence
    //
    for a:=0; a<n_allele; a++ {

      var bp_ref byte = '-'
      if i<refn {
        bp_ref = stream_ref_bp

        if ref_anchor_on_left {
          if bp_ref != stream_ref_bp {
            return fmt.Errorf( fmt.Sprintf("ref stream to vcf ref mismatch (ref stream %c != vcf ref %c @ %d)", stream_ref_bp, bp_ref, g.RefPos) )
          }
        }

      }

      var bp_alt byte = '-'
      if samp_seq_idx[a]==0 {
        bp_alt = bp_ref
      } else {
        a_idx := samp_seq_idx[a]-1
        if i<len(alt_seq[a_idx]) { bp_alt = alt_seq[a_idx][i] }
      }

      pasta_ch := pasta.SubMap[bp_ref][bp_alt]
      if pasta_ch == 0 { return fmt.Errorf("invalid character") }

      if (g.LFMod>0) && (g.OCounter > 0) && ((g.OCounter%g.LFMod)==0) {
        out.WriteByte('\n')
      }
      g.OCounter++


      out.WriteByte(pasta_ch)
    }

  }

  return nil
}
