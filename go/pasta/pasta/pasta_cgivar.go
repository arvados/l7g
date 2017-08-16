package main

import "fmt"
import "strconv"
import "strings"
import "bufio"
//import "io"
//import "bytes"

import "time"

import "github.com/abeconnelly/pasta"

const PASTA_CGIVAR_SOFT_VER = "0.1.0"
const PASTA_CGIVAR_FMT_VER_STR = "2.5"

//var __g_debug bool = false

type CGIRefVar struct {
  Locus int
  Ploidy int
  ChromStr string

  Fields []string

  PrintHeader bool

  Date time.Time

  RefStart int
  RefBufVirtStart int
  RefBufLen int

  RefPos []int
  RefBuf []byte

  Seq [][]byte
  Start int

  PrevLocus int
  PrevStart int
  PrevEnd int
  PrevChomStr string

  LFMod int
  OCounter int
  LCounter int

  CurBeg int

  ChromUpdate bool
  RefPosUpdate bool
  CGIVarRefPos int

  InitState bool
  CGIVarLine string
}

func (g *CGIRefVar) Init() {
  g.Locus = 1
  g.ChromStr = "Unk"
  g.Ploidy=2

  g.Fields = []string{ "locus", "ploidy", "allele",
    "chromosome", "begin", "end", "varType", "reference",
    "alleleSeq", "varScoreVAF", "varScoreEAF", "varFilter",
    "hapLink", "xRef", "alleleFreq", "alternativeCalls" }

  g.PrintHeader = true
  g.Date = time.Now()
  g.Start = 0
  g.Seq = make([][]byte, g.Ploidy)
  for ii:=0; ii<g.Ploidy; ii++ {
    g.Seq[ii] = make([]byte, 0, 1024)
  }

  g.OCounter = 0
  g.LFMod = 50

  g.RefStart = 0
  g.RefBuf = make([]byte, 0, 1024)
  g.RefBufLen = 0
  g.RefBufVirtStart = 0

  g.RefPos = make([]int, 2)
  g.RefPos[0] = 0
  g.RefPos[1] = 0

  g.InitState = true
  g.LCounter = 0
}

func (g *CGIRefVar) Chrom(chrom string) {
  g.ChromStr = chrom
  g.ChromUpdate = true
}

func (g *CGIRefVar) Pos(pos int) {
  g.RefPosUpdate = true
  g.CGIVarRefPos = pos
}

func (g *CGIRefVar) Header(out *bufio.Writer) error {
  var header = []string{}

  header = append(header, "#GENERATED_BY\tpasta tools")
  header = append(header, fmt.Sprintf("#GENERATED_AT\t%d%02d%02d", g.Date.Year(), g.Date.Month(), g.Date.Day()))
  header = append(header, fmt.Sprintf("#SOFTWARE_VERSION\t%s", PASTA_CGIVAR_SOFT_VER))
  header = append(header, fmt.Sprintf("#FORMAT_VERSION\t%s", PASTA_CGIVAR_FMT_VER_STR))
  header = append(header, "#TYPE\tVAR-ANNOTATION")
  header = append(header, "")

  out.WriteString( strings.Join(header, "\n") )
  out.WriteString("\n")
  out.WriteString( ">" + strings.Join(g.Fields, "\t") + "\n" )

  return nil

}

func (g *CGIRefVar) _strip_seqs(refseq []byte, altseq [][]byte) ([]byte, [][]byte) {
  r := []byte{}
  a := [][]byte{}
  for ii:=0; ii<len(altseq); ii++ { a = append(a, []byte{}) }

  for ii:=0; ii<len(refseq); ii++ {
    if refseq[ii]!='-' {
      r = append(r, refseq[ii])
    }
  }

  for ii:=0; ii<len(altseq); ii++ {
    for jj:=0; jj<len(altseq[ii]); jj++ {
      if altseq[ii][jj] != '-' {
        a[ii] = append(a[ii], altseq[ii][jj])
      }
    }
  }

  return r,a
}


func (g *CGIRefVar) PrintAltAlleles(vartype int, ref_start, ref_len int, refseq []byte, altseq [][]byte, out *bufio.Writer) error {

  varscorevaf := ""
  varscoreeaf := ""
  varfilter := ""
  haplink := ""
  xref := ""
  allelefreq := ""
  altcalls := ""

  ref,alt := g._strip_seqs(refseq, altseq)

  for strand:=0; strand<len(alt); strand++ {
    allele_str := fmt.Sprintf("%d", strand+1)

    cur_start := 0
    min_len := len(ref)
    if min_len > len(alt[strand]) { min_len = len(alt[strand]) }

    for min_len>0 {

      noc_pfx_len := 0
      for ii:=0; ii<min_len; ii++ {
        if (alt[strand][cur_start+ii] != 'n') && (alt[strand][cur_start+ii] != 'N') { break }
        noc_pfx_len++
      }

      if noc_pfx_len>0 {

        //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
        out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
          g.Locus, g.Ploidy, allele_str, g.ChromStr,
          ref_start+cur_start, ref_start+cur_start+noc_pfx_len,
          "no-call",
          //refseq, "?",
          "=", "?",
          varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

        cur_start += noc_pfx_len
        min_len -= noc_pfx_len

      }

      ref_pfx_len := 0
      for ii:=0; ii<min_len; ii++ {
        if (ref[cur_start+ii]=='n') || (ref[cur_start+ii]=='N') || (ref[cur_start+ii]!=alt[strand][cur_start+ii]) { break }
        ref_pfx_len++
      }

      if ref_pfx_len>0 {

        //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
        out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
          g.Locus, g.Ploidy, allele_str, g.ChromStr,
          ref_start+cur_start, ref_start+cur_start+ref_pfx_len,
          "ref",
          "=", "=",
          varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

        cur_start += ref_pfx_len
        min_len -= ref_pfx_len

      }

      sub_pfx_len := 0
      for ii:=0; ii<min_len; ii++ {
        if ref[cur_start+ii] == alt[strand][cur_start+ii] { break }
        sub_pfx_len++
      }

      if sub_pfx_len>0 {

        vartype_str := "snp"
        if sub_pfx_len > 1 { vartype_str = "sub" }

        //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
        out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
          g.Locus, g.Ploidy, allele_str, g.ChromStr,
          ref_start+cur_start, ref_start+cur_start+sub_pfx_len,
          vartype_str,
          ref[cur_start:cur_start+sub_pfx_len],
          alt[strand][cur_start:cur_start+sub_pfx_len],
          varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

        cur_start += sub_pfx_len
        min_len -= sub_pfx_len

      }

    }

    if cur_start<len(ref) {

      //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
      out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
        g.Locus, g.Ploidy, allele_str, g.ChromStr,
        ref_start+cur_start, ref_start+cur_start+len(ref[cur_start:]),
        "del",
        ref[cur_start:],
        "",
        varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

    } else if cur_start<len(alt[strand]) {

      //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
      out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
        g.Locus, g.Ploidy, allele_str, g.ChromStr,
        //ref_start+cur_start, ref_start+cur_start+len(alt[strand]),
        ref_start+cur_start, ref_start+cur_start,
        "ins",
        "",
        alt[strand][cur_start:],
        varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

    }

  }

  g.Locus++

  return nil
}

func (g *CGIRefVar) Print(vartype int, ref_start, ref_len int, refseq []byte, altseq [][]byte, out *bufio.Writer) error {

  if g.PrintHeader {
    g.Header(out)
    g.PrintHeader = false
  }

  ploidy := g.Ploidy
  allele_str := "all"
  chrom := g.ChromStr

  varscorevaf := ""
  varscoreeaf := ""
  varfilter := "PASS"
  haplink := ""
  xref := ""
  allelefreq := ""
  altcalls := ""

  vartype_str := "no-ref"

  ref,alt := g._strip_seqs(refseq, altseq) ;  _ = ref ; _ = alt

  if vartype == pasta.NOREF {
    vartype_str = "no-ref"

    //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
    out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
      g.Locus, ploidy, allele_str, chrom,
      ref_start, ref_start+ref_len,
      vartype_str,
      "=", "?",
      varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

    g.Locus++

  } else if vartype == pasta.REF {
    vartype_str = "ref"

    //                           0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
    out.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
      g.Locus, ploidy, allele_str, chrom,
      ref_start, ref_start+ref_len,
      vartype_str,
      "=", "=",
      varscorevaf, varscoreeaf, varfilter, haplink, xref, allelefreq, altcalls))

    g.Locus++

  } else if (vartype==pasta.NOC) || (vartype==pasta.ALT) {

    //snp, sub, ins, del
    //
    g.PrintAltAlleles(vartype, ref_start, ref_len, refseq, altseq, out)
  }

  return nil

}

func (g *CGIRefVar) PrintEnd(out *bufio.Writer) error {
  out.Flush()
  return nil
}

//---

func (g *CGIRefVar) PastaBegin(out *bufio.Writer) error {
  g.Seq = make([][]byte, g.Ploidy)
  for ii:=0; ii<g.Ploidy; ii++ {
    g.Seq[ii] = make([]byte, 0, 1024)
  }

  g.Locus = 0

  return nil
}

func (g *CGIRefVar) PastaEnd(out *bufio.Writer) error {

  // Print out the sequence that we can,
  // first filling the the shorter of the
  // two sequences with the 'nop' character
  // '.' (period)
  //
  max_len := len(g.Seq[0])
  min_len := len(g.Seq[0])
  for a:=1; a<len(g.Seq); a++ {
    if max_len < len(g.Seq[a]) { max_len = len(g.Seq[a]) }
    if min_len > len(g.Seq[a]) { min_len = len(g.Seq[a]) }
  }

  for a:=0; a<len(g.Seq); a++ {
    for ii:=min_len; ii<max_len; ii++ {
      g.Seq[a] = append(g.Seq[a], '.')
    }
  }

  for ii:=0; ii<max_len; ii++ {
    for a:=0; a<len(g.Seq); a++ {
      g.WritePastaByte(g.Seq[a][ii], out)
    }
  }

  for a:=0; a<len(g.Seq); a++ { g.Seq[a] = g.Seq[a][max_len:] }

  out.Flush()
  return nil
}

func (g *CGIRefVar) WritePastaByte(pasta_ch byte, out *bufio.Writer) {
  out.WriteByte(pasta_ch)
  g.OCounter++
  if (g.LFMod>0) && (g.OCounter > 0) && ((g.OCounter%g.LFMod)==0) {
    out.WriteByte('\n')
  }
}

func (g *CGIRefVar) RefByte(strand int, ref_stream *bufio.Reader) (byte, error) {

  if strand<=0 {
    ref_bp,e := ref_stream.ReadByte()
    if e!=nil { return ref_bp, e }
    for (ref_bp=='\n') || (ref_bp==' ') || (ref_bp=='\r') || (ref_bp=='\t') {
      ref_bp,e = ref_stream.ReadByte()
      if e!=nil { return ref_bp, e }
    }
    g.RefBuf = append(g.RefBuf, ref_bp)
    g.RefPos[0]++
    return ref_bp,e
  }

  if g.RefPos[1] >= len(g.RefBuf) {
    panic( fmt.Sprintf("g.RefPos[1] %d >= len(g.RefBuf) %d (RefPos:%v, RefBuf:%s, line %d '%s')", g.RefPos[1], len(g.RefBuf), g.RefPos, g.RefBuf, g.LCounter, g.CGIVarLine) )
  }

  ref_bp := g.RefBuf[g.RefPos[1]]
  g.RefPos[1]++

  return ref_bp,nil

}

func (g *CGIRefVar) RefByteReset() {
  g.RefBuf = g.RefBuf[0:0]
  g.RefPos[0] = 0
  g.RefPos[1] = 0
}


func (g *CGIRefVar) Pasta(cgivar_line string, ref_stream *bufio.Reader, out *bufio.Writer) error {

  g.LCounter++
  g.CGIVarLine = cgivar_line

  // Skip header, blank space and column definitions
  //
  if len(cgivar_line)==0 || cgivar_line[0] == '#' { return nil }
  if cgivar_line[0]=='>' { return nil }
  if cgivar_line[0]=='\n' { return nil }

  // Parse fields
  //
  fields := strings.Split(cgivar_line, "\t")

  idx:=0
  locus,e := strconv.Atoi(fields[idx]) ; _ = locus
  if e!=nil { return e }

  idx++
  ploidy,e := strconv.Atoi(fields[idx]) ; _ = ploidy
  if e!=nil { return e }

  idx++
  allele := fields[idx] ; _ = allele

  idx++
  chrom := fields[idx] ; _ = chrom

  idx++
  _beg,e := strconv.Atoi(fields[idx]) ; _ = _beg
  if e!=nil { return e }

  idx++
  _end,e := strconv.Atoi(fields[idx]) ; _ = _end
  if e!=nil { return e }

  idx++
  vartype := fields[idx] ; _ = vartype

  if strings.HasPrefix(vartype, "no-call") {
    vartype = "no-call"
  }

  // Default to 'no-call' if we don't know the vartype
  //
  if (vartype != "no-call") &&
     (vartype != "snp") &&
     (vartype != "sub") &&
     (vartype != "ins") &&
     (vartype != "del") &&
     (vartype != "ref") &&
     (vartype != "no-ref") {
    vartype = "no-call"
  }

  idx++
  refseq := fields[idx] ; _ = refseq

  idx++
  alleleseq := fields[idx] ; _ = alleleseq

  dn := _end - _beg

  // Update our knowledge of chromosome name
  // and reference position
  //
  if chrom != g.ChromStr {
    g.ChromStr = chrom
    g.ChromUpdate = true
  }

  if g.PrevLocus != locus {
    if _beg != g.CGIVarRefPos {
      g.CGIVarRefPos = _beg
      g.RefPosUpdate = true
    }
  }

  // Force a control message if we're running
  // for the first time.
  //
  if g.InitState {
    g.ChromUpdate = true
    g.RefPosUpdate = true
    g.InitState = false
  }

  // Print header if there are any new updates
  //
  if g.ChromUpdate {
    out.WriteString( fmt.Sprintf(">C{%s}", g.ChromStr) )
  }

  if g.RefPosUpdate {
    out.WriteString( fmt.Sprintf(">P{%d}", g.CGIVarRefPos) )
  }

  if g.ChromUpdate || g.RefPosUpdate {
    out.WriteByte('\n')
  }

  g.ChromUpdate = false
  g.RefPosUpdate = false
  g.CGIVarRefPos = _end
  g.PrevLocus = locus



  // Get list of sequence indices for updating
  // each of the output pasta streams.
  //

  map_idx := make(map[int]bool)
  for ii:=0; ii<g.Ploidy; ii++ {
    map_idx[ii] = false
  }

  // allele_code holds the allele number
  // (0 for the first allele, 1 for the second, etc.)
  //
  allele_code := -1

  // seq_idx holds an array of the alleles to update.
  // Most of the time this will be a single allele
  // or both alleles.
  //
  seq_idx := []int{}

  if allele=="all" {

    // Update each allele
    //
    for ii:=0; ii<g.Ploidy; ii++ {
      seq_idx = append(seq_idx, ii)
      map_idx[ii] = true
    }
    allele_code = -1

  } else {

    // Otherwise update with the
    // single allele we're updating.
    //
    z,e := strconv.Atoi(allele)
    if e!=nil { return e }
    seq_idx = append(seq_idx, z-1)
    map_idx[z] = true

    allele_code = z-1

  }

  // Fill out previous locus if necessary
  //
  for locus!=g.Locus {

    // We've moved to a new locus, so we need
    // to print out the previous locus's
    // sequence information.
    //

    g.Locus = locus
    g.RefByteReset()

    g.CurBeg = _beg

    // Simple case of single strand, print out PASTA stream
    // without issue
    //
    if g.Ploidy == 1 {
      for ii:=0; ii<len(g.Seq[0]); ii++ {
        g.WritePastaByte(g.Seq[0][ii], out)
      }
      g.Seq[0] = g.Seq[0][0:0]
      break
    }

    // We need to interleave the two 'raw' pasta streams
    // to insert the no-op character ('.') at the appropriate
    // locations.
    // The assumption is that the PASTA stream is ref-aligned
    // at every interleaved character group.
    //
    if g.Ploidy == 2 {

      a_pos := 0
      b_pos := 0

      for (a_pos<len(g.Seq[0])) && (b_pos<len(g.Seq[1])) {
        a_state := pasta.BPState[g.Seq[0][a_pos]]
        b_state := pasta.BPState[g.Seq[1][b_pos]]

        a_pasta_ch := byte('.')
        b_pasta_ch := byte('.')

        if (a_state == pasta.INS) {
          a_pasta_ch = g.Seq[0][a_pos]
          a_pos++
        } else if (b_state == pasta.INS) {
          b_pasta_ch = g.Seq[1][b_pos]
          b_pos++
        } else {
          a_pasta_ch = g.Seq[0][a_pos]
          b_pasta_ch = g.Seq[1][b_pos]
          a_pos++
          b_pos++
        }

        g.WritePastaByte(a_pasta_ch, out)
        g.WritePastaByte(b_pasta_ch, out)
      }


      for ; a_pos<len(g.Seq[0]); a_pos++ {
        g.WritePastaByte(g.Seq[0][a_pos], out)
        g.WritePastaByte('.', out)
      }

      for ; b_pos<len(g.Seq[1]); b_pos++ {
        g.WritePastaByte('.', out)
        g.WritePastaByte(g.Seq[1][b_pos], out)
      }

      g.Seq[0] = g.Seq[0][0:0]
      g.Seq[1] = g.Seq[1][0:0]


      break
    }

    return fmt.Errorf("Ploidy>2 not implemented")

  }

  // Case analysis for each type:
  //   no-ref, ref, no-call, snp, sub, ins, del
  //
  if (vartype == "no-call") || (vartype == "no-ref")  {

    for ii:=0; ii<dn; ii++ {
      ref_bp,e := g.RefByte(allele_code, ref_stream) ; _ = ref_bp
      if e!=nil { return e }

      for a:=0; a<len(seq_idx); a++ {
        idx := seq_idx[a]
        g.Seq[idx] = append(g.Seq[idx], pasta.SubMap[ref_bp]['n'])
      }
    }

  } else if vartype == "ref" {

    for ii:=0; ii<dn; ii++ {
      ref_bp,e := g.RefByte(allele_code, ref_stream)
      if e!=nil { return e }

      for a:=0; a<len(seq_idx); a++ {
        idx := seq_idx[a]
        g.Seq[idx] = append(g.Seq[idx], ref_bp)
      }
    }

  } else if (vartype=="snp") || (vartype == "sub") || (vartype == "ins") || (vartype == "del") {


    // snp, sub, ins and del are all over the place and should be treated
    // as just a general indel.  For example, there are lines labelled
    // as 'snp' that have an extra reference base (so two long) and
    // subs that have an extra insertion or deletion.
    //

    for ii:=0; ii<dn; ii++ {

      ref_bp,e := g.RefByte(allele_code, ref_stream)
      if e!=nil { return e }

      if ii<len(alleleseq) {

        pasta_ch,ok := pasta.SubMap[ref_bp][_tolch(alleleseq[ii])]
        if !ok { return fmt.Errorf(fmt.Sprintf("bad sub map from [%c,%c] (%d,%d)", ref_bp, alleleseq[ii], ref_bp, alleleseq[ii]))}

        for a:=0; a<len(seq_idx); a++ {
          idx := seq_idx[a]

          g.Seq[idx] = append(g.Seq[idx], pasta_ch)

        }

      } else {

        pasta_ch,ok := pasta.DelMap[_tolch(ref_bp)]
        if !ok { panic("cp sub-del") }

        for a:=0; a<len(seq_idx); a++ {
          idx := seq_idx[a]
          g.Seq[idx] = append(g.Seq[idx], pasta_ch)
        }

      }

    }

    for ii:=dn; ii<len(alleleseq); ii++ {

      pasta_ch,ok := pasta.InsMap[_tolch(alleleseq[ii])]
      if !ok { panic("cp sub-ins") }

      for a:=0; a<len(seq_idx); a++ {
        idx := seq_idx[a]
        g.Seq[idx] = append(g.Seq[idx], pasta_ch)
      }
    }

  } else {
    return fmt.Errorf("invalid vartype '%s' at line %d ('%s')", vartype, g.LCounter, g.CGIVarLine)
  }

  return nil
}
