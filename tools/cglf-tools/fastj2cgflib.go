/*

Only meant to be used on a path by path basis.

This only creates the tile library.  The tile library
has only called bases.  No-calls are considered to be
stored with the callset/sample.


* loop through all tiles
* collect all unique tiles in memory
* for each tile position, find most frequent length L_{tilepos}
* For each tilepos
  - loop through all tiles of length L_{tilepos}
  - do a frequency count of each called base
  - take most frequently called base at that position as canonical
* build the canonical tile for each tilepos
* loop through each tile a final time.
  - seed the tile library with the canonical tile
  - if the tile is already there, nothing to do
  - otherwise align the tile against the canonical tile, filling
    in any no-calls with canonicial or 'a', if canonical doesn't exist.

*/


// the magic codes are:
//
// a - canonical tile
// b - simple sub tile
// c - alt data tile (gzip?)
// d - simple indel
// e - simple align
// f - aux align
// g - failed aux align

// The first output column:
// _ - in cache
// * - aux overflow type 1
// ! - aux body seq flag
// ^ - catch all


package main

import "os"
import "fmt"
import "log"
import "crypto/md5"

import "sort"

import "strings"
import "strconv"

import "encoding/binary"

//import "github.com/abeconnelly/autoio"
import "./autoio"
import "github.com/codegangsta/cli"

import "./dlug"

import "./memz"
import _ "math"
import "./twobit"

import "bytes"
import "bufio"
import _ "io/ioutil"
import "compress/gzip"

var VERSION_STR string = "AGPLv3, v0.2.0"

var gProfileFlag bool = false
var gProfileFile string = "fastj2cgflib.cprof"
var gMemProfileFile string = "fastj2cgflib.mprof"
var g_verbose bool = false
var g_output_format string = "gvcf"
var g_tagset_fn string

var g_check_align bool = false

const MAX_GLOB_ALIGN_LEN = 1024

// path step length md5sum
//
type TileLibEntry struct {
  Seq string
  Freq int
  Span int
  Diff []AltEntry

  AuxFlag bool
  InCache bool

  StartTile bool
  EndTile bool
}

//type AuxTileLibEntry struct {
type AuxOverflowLibEntry struct {
  Seq string
  Span int
  Diff []AltEntry

  BaseIndex int
  Type int
  AuxBodySeqFlag bool
  VariantId int
  AuxVariantId int

  InCache bool
  StartTile bool
  EndTile bool

  _z string
}


type CanonSeq struct {
  Seq []byte
  Freq int
  Span int
}

type AltEntry struct {
  Start int
  N     int
  Seq   []byte
}

type SimpleTile struct {
  Path    int
  Step    int
  Variant int

  Seq []byte
}

type TileLibraryPathVerbose struct {
  Path  uint64
  TagSequence [][]byte
  Sequence [][]byte
  Span  []int
  Alt [][]AltEntry
  Lib [][]TileLibEntry

  CanonSeq [][]byte
}

type TileLibPath struct {
  Path uint64
  Tag2BitLenBP uint64
  TagSeq2Bit  []byte

  Seq2BitLenBP  uint64
  Seq2Bit []byte
  Seq2BitBPOffset []uint64

  Span  []byte
  SpanOverflowLen uint64
  SpanOverflow []uint64

  Alt []byte

  VariantInfo []byte
  VariantIdOffset []uint64
  AltOverflow []byte
}

func _diff_str(d []AltEntry) string {
  s := []string{}
  for i:=0; i<len(d); i++ {
    s_ele := fmt.Sprintf("%x+%x{%s}", d[i].Start, d[i].N, d[i].Seq)
    s = append(s, s_ele)
  }
  return strings.Join(s, ";")
}

func _repl(ch byte, n int) []byte {
  b := make([]byte, n)
  for i:=0; i<n; i++ {
    b[i]=ch
  }
  return b
}

func _delta_str(d []memz.Diff) string {
  s := []string{}
  for i:=0; i<len(d); i++ {
    n := 0
    str := []byte{}
    if d[i].Type == memz.DIFF { n = d[i].Len ; str = _repl('X', n) }
    if d[i].Type == memz.GAPA { str = _repl('_', d[i].Len) }
    if d[i].Type == memz.GAPB { n = d[i].Len }


    //s_ele := fmt.Sprintf("%x:%x+%x{%s}", d[i].PosA, d[i].PosB, d[i].Len, t)
    s_ele := fmt.Sprintf("%x+%x{%s}", d[i].PosA, n, str)

    s = append(s, s_ele)
  }
  return strings.Join(s, ";")
}

var g_bytes TileLibPath
var g_var_name_map map[int]string

func simple_text_field(line, field string) (string, int) {

  sep := "\"" + field + "\""
  p:=strings.Index(line, sep)

  n:=len(line)
  if p<0 { return "", p; }
  p+=len(sep)
  if p>=n { return "", -1 }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return "", -2; }
  if line[p] != ':' { return "", -3; }

  p++
  if p==n { return "", -2; }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return "", -4; }
  if line[p] != '"' { return "", -5; }

  p++
  if p==n { return "", -6; }

  p_s := p

  for ; p<n && line[p]!='"'; p++ { }
  if p==n { return "", -7; }

  p_e := p

  return line[p_s:p_e], p_e

}


func simple_bool_field(line, field string) (bool, int) {

  sep := "\"" + field + "\""
  p:=strings.Index(line, sep)

  n:=len(line)
  if p<0 { return false, p; }
  p+=len(sep)
  if p>=n { return false, -1 }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return false, -2; }
  if line[p] != ':' { return false, -3; }

  p++
  if p==n { return false,-2; }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return false, -4; }

  val_start_pos := p
  for ; p<n && line[p]!=',' && line[p]!='}' && line[p]!=' ' && line[p]!='\t' && line[p]!='\n' ; p++ { }
  if p==n { return false, -5; }

  val := line[val_start_pos:p]
  if val=="true" { return true, p }
  if val=="false" { return false, p }
  return false, -6
}

func simple_int_field(line, field string) (int, int) {

  sep := "\"" + field + "\""
  p:=strings.Index(line, sep)

  n:=len(line)
  if p<0 { return 0, p; }

  p+=len(sep)
  if p>=n { return 0, -2 }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return 0, -3; }
  if line[p] != ':' { return 0, -4; }

  p++
  if p==n { return 0, -5; }

  for ; p<n && line[p]==' '; p++ { }
  if p==n { return 0, -6; }

  p_s := p

  for ; p<n && (line[p]!='"') && (line[p]!=','); p++ { }
  if p==n { return 0, -7; }

  p_e := p

  i,e := strconv.ParseInt(line[p_s:p_e], 10, 64)
  if e!=nil { return 0, -8; }

  return int(i), p_e

}

func md5sum2str(md5sum [16]byte) string {
  var str_md5sum [32]byte
  for i:=0; i<16; i++ {
    x := fmt.Sprintf("%02x", md5sum[i])
    str_md5sum[2*i]   = x[0]
    str_md5sum[2*i+1] = x[1]
  }

  return string(str_md5sum[:])
}

// path - step - length - md5sum
//
var tilelib_map map[int]map[int]map[int]map[string]TileLibEntry
var canon_seq_by_tilelen map[int]map[int]map[int]TileLibEntry
var filled_tilelib_map map[int]map[int]map[string]TileLibEntry
var ranked_tilelib_map map[int]map[int][]TileLibEntry

func create_canonical_sequences() {
  canon_seq_by_tilelen = make(map[int]map[int]map[int]TileLibEntry)
  for path,_ := range tilelib_map {
    canon_seq_by_tilelen[path] = make(map[int]map[int]TileLibEntry)
    for step,_ := range tilelib_map[path] {
      canon_seq_by_tilelen[path][step] = make(map[int]TileLibEntry)

      for tilelen,_ := range tilelib_map[path][step] {

        sflag:=false
        eflag:=false
        span := 0
        freq := 0
        canon_freq := make([]int4, tilelen)
        canon_seq := make([]byte, tilelen)
        for m5,_ := range tilelib_map[path][step][tilelen] {
          ent := tilelib_map[path][step][tilelen][m5]

          sflag = ent.StartTile
          eflag = ent.EndTile

          for i:=0; i<tilelen; i++ {
            switch ent.Seq[i] {
            case 'a','A': canon_freq[i].I[0]++
            case 'c','C': canon_freq[i].I[1]++
            case 'g','G': canon_freq[i].I[2]++
            case 't','T': canon_freq[i].I[3]++
            default: canon_freq[i].I[0]++
            }

          }

          freq++
          span = ent.Span

        }

        for i:=0; i<tilelen; i++ {
          acgt := []byte{'a', 'c', 'g', 't'}
          bub_mirror(&(canon_freq[i]), acgt)
          canon_seq[i] = acgt[0]
        }

        canon_seq_by_tilelen[path][step][tilelen] = TileLibEntry{ Seq: string(canon_seq), Freq:freq, Span:span, Diff:nil, StartTile:sflag, EndTile:eflag }

      }
    }
  }

}

func fill_sequences() {
  filled_tilelib_map = make(map[int]map[int]map[string]TileLibEntry)
  for path,_ := range tilelib_map {
    filled_tilelib_map[path] = make(map[int]map[string]TileLibEntry)

    for step,_ := range tilelib_map[path] {
      filled_tilelib_map[path][step] = make(map[string]TileLibEntry)

      for tilelen,_ := range tilelib_map[path][step] {
        for _,ent := range tilelib_map[path][step][tilelen] {

          raw_seq := ent.Seq
          canon_seq := canon_seq_by_tilelen[path][step][tilelen].Seq
          filled_seq := make([]byte, tilelen)

          for i:=0; i<tilelen; i++ {
            filled_seq[i] = raw_seq[i]
            if filled_seq[i] == 'n' || filled_seq[i] == 'N' { filled_seq[i] = canon_seq[i] }
          }

          m5_filled := md5sum2str(md5.Sum(filled_seq))

          //DEBUG
          //fmt.Printf("### %x.%x %v %v\n", path, step, ent.StartTile, ent.EndTile)


          if _,ok := filled_tilelib_map[path][step][m5_filled] ; !ok {
            filled_tilelib_map[path][step][m5_filled] = TileLibEntry{ Seq: string(filled_seq), Freq: ent.Freq, Span: ent.Span, Diff: nil, StartTile:ent.StartTile, EndTile:ent.EndTile }
          } else {
            z := filled_tilelib_map[path][step][m5_filled]
            z.Freq += ent.Freq
            filled_tilelib_map[path][step][m5_filled] = z
          }

        }
      }
    }
  }

}

func rank_sequences() {
  ranked_tilelib_map = make(map[int]map[int][]TileLibEntry)
  for path,_ := range filled_tilelib_map {
    ranked_tilelib_map[path] = make(map[int][]TileLibEntry)

    for step,_ := range filled_tilelib_map[path] {
      ranked_tilelib_map[path][step] = make([]TileLibEntry, 0, 8)


      for _,ent := range filled_tilelib_map[path][step] {
        ranked_tilelib_map[path][step] = append(ranked_tilelib_map[path][step], ent)
      }

      sort.Sort(ByTileFreq(ranked_tilelib_map[path][step]))

    }
  }

}

type ByTileFreq []TileLibEntry ;
func (f ByTileFreq) Len() int { return len(f) }
func (f ByTileFreq) Less(i,j int) bool {
  if f[i].Freq == f[j].Freq {
    if f[i].Span == f[j].Span {
      return f[i].Seq < f[j].Seq
    }
    return f[i].Span < f[j].Span
  }
  return f[i].Freq > f[j].Freq
}
func (f ByTileFreq) Swap(i,j int) { f[i],f[j] = f[j],f[i] }

const THRESHOLD = 1000

var aux_overflow_lib map[int]map[int][]AuxOverflowLibEntry

// Sequence A should be the canonical sequence
// Sequence B should be the 'alt' sequence
// The Diff array holds the difference from the global alignment
//

func altentry_from_delta(adelta []memz.Diff, b_seq []byte) []AltEntry {
  aux_alt_diff := make([]AltEntry, 0, 8)
  for ii:=0; ii<len(adelta); ii++ {
    posa := adelta[ii].PosA
    posb := adelta[ii].PosB
    n := adelta[ii].Len

    switch adelta[ii].Type {

    // e.g.
    //    012345678
    // a: fooxvqbar
    // b: fooquxbar
    //
    // Start: 3, N:3, Seq:"qux"
    //
    case memz.DIFF:

      if posb+n > len(b_seq) {
        fmt.Printf("ERROR!!! (diff) %v\nposa %d, posb %d, n %d, b_seq %s len = %d\n", adelta, posa, posb, n, b_seq, len(b_seq))
        os.Exit(1)
      }

      aux_alt_diff = append(aux_alt_diff, AltEntry{ Start: posa, N: n, Seq: b_seq[posb:posb+n] })

    // e.g.
    //    012   345
    // a: foo---bar
    // b: fooquxbar
    //    012345678
    //
    // Start: 3, N: 0, Seq: "qux"
    //
    case memz.GAPA:

      //DEBUG
      if posb+n > len(b_seq) {
        fmt.Printf("ERROR!!! (gapa) %v\nposa %d, posb %d, n %d, b_seq %s len = %d\n", adelta, posa, posb, n, b_seq, len(b_seq))
        os.Exit(1)
      }

      aux_alt_diff = append(aux_alt_diff, AltEntry{ Start: posa, N: 0, Seq: b_seq[posb:posb+n] })

    // e.g.
    //    012345678
    // a: babxuqoof
    // b: bab---oof
    //
    // Start: 3, N: 3, Seq: ""
    //
    case memz.GAPB:
      aux_alt_diff = append(aux_alt_diff, AltEntry{ Start: posa, N: n, Seq: []byte("") })
    }
  }
  return aux_alt_diff
}

// Aux overflow structure should hold all sequences (not just aux)
func print_aux_overflow() {
  for path := range aux_overflow_lib {
    for step := range aux_overflow_lib[path] {
      for i:=0; i<len(aux_overflow_lib[path][step]); i++ {

        ch_code := '_'
        if !aux_overflow_lib[path][step][i].InCache {
          if aux_overflow_lib[path][step][i].Type == 1 {
            ch_code = '*'
          } else {
            if aux_overflow_lib[path][step][i].AuxBodySeqFlag {
              ch_code = '!'
            } else {
              ch_code = '^'
            }
          }
        }

        fmt.Printf("#(%s) %s\n", aux_overflow_lib[path][step][i]._z, _diff_str(aux_overflow_lib[path][step][i].Diff))

        m5 := md5sum2str(md5.Sum([]byte(aux_overflow_lib[path][step][i].Seq)))

        fmt.Printf("%c,%04x.00.%04x.%03x+%d,%s,%s\n", ch_code, path, step, i,
          aux_overflow_lib[path][step][i].Span,
          m5,
          aux_overflow_lib[path][step][i].Seq)

      }
    }
  }
}


type ByAltEntry []AltEntry
func (f ByAltEntry) Len() int { return len(f) }
func (f ByAltEntry) Less(i,j int) bool { return f[i].Start < f[j].Start }
func (f ByAltEntry) Swap(i,j int) { f[i],f[j] = f[j],f[i] }


var g_sanity_flag bool = true
func seq_from_diff(a_diff []AltEntry, seq []byte) []byte {
  o_seq := []byte{}

  //sort.Sort(ByAltEntry(a_diff))

  curpos:=0
  for i:=0; i<len(a_diff); i++ {
    s := a_diff[i].Start
    n := a_diff[i].N
    subseq := []byte(a_diff[i].Seq)

    if curpos < s {
      o_seq = append(o_seq, seq[curpos:s]...)
    }
    curpos=s

    if len(subseq) > 0 {
      o_seq = append(o_seq, subseq...)
    }

    curpos+=n

  }
  if curpos<len(seq) {
    o_seq = append(o_seq, seq[curpos:]...)
  }

  return o_seq
}

func calc_ranked_diffs() {

  aux_overflow_lib = make(map[int]map[int][]AuxOverflowLibEntry)

  for path,_ := range ranked_tilelib_map {

    aux_overflow_lib[path] = make(map[int][]AuxOverflowLibEntry)

    for step,ranked_tile := range ranked_tilelib_map[path] {
      aux_overflow_lib[path][step] = make([]AuxOverflowLibEntry, 0, 8)

      canon_seq := ranked_tile[0].Seq
      aux_canon_seq := make([][]byte, 0) ; _ = aux_canon_seq
      aux_canon_var_id := make([]int, 0) ; _ = aux_canon_var_id

      // ------------------
      // ADD CANONICAL TILE
      // ------------------
      //
      // aux_overflow_lib should have same length as ranked_tile?
      //
      aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
        AuxOverflowLibEntry{ Seq: ranked_tile[0].Seq,
                             AuxBodySeqFlag: false,
                             Type: 0,
                             Span: ranked_tile[0].Span,
                             VariantId:0,
                             Diff:ranked_tile[0].Diff,
                             StartTile: ranked_tile[0].StartTile,
                             EndTile: ranked_tile[0].EndTile,
                             _z:"a" })

      for tile_idx:=1; tile_idx<len(ranked_tile); tile_idx++ {

        // Simple case where we can do a straight substitution
        //
        if len(ranked_tile[tile_idx].Seq) == len(canon_seq) {
          z := ranked_tile[tile_idx]
          z.Diff = simple_subs_diff([]byte(canon_seq), []byte(ranked_tile[tile_idx].Seq))
          ranked_tile[tile_idx] = z

          // -------------------
          // ADD SIMPLE SUB TILE
          // -------------------
          //
          // the differencece is against the canonical tile
          //
          aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
            AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                                 AuxBodySeqFlag: false,
                                 Type: 0,
                                 Span: ranked_tile[tile_idx].Span,
                                 VariantId:tile_idx,
                                 Diff:z.Diff,
                                 StartTile: ranked_tile[tile_idx].StartTile,
                                 EndTile: ranked_tile[tile_idx].EndTile,
                                 _z:"b" })

          continue
        }


        //-------------------------------
        // The sequence is too long, put it structure for gzip later
        // (AuxBodySeqFlag = true and Type = 1)
        //
        if len(ranked_tile[tile_idx].Seq) >= 1000 {
          //fmt.Printf("SKIPPING %04x.00.%04x.%03x %s\n", path, step, tile_idx, ranked_tile[tile_idx].Seq)

          z := ranked_tile[tile_idx]
          z.AuxFlag = true
          ranked_tile[tile_idx] = z

          // -------------
          // ALT DATA TILE
          // -------------
          //
          // The tile is too long.  Flag it for the 'AltData' structures.
          //
          aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
            AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                                 AuxBodySeqFlag: true,
                                 Type: 1,
                                 Span: ranked_tile[tile_idx].Span,
                                 VariantId:tile_idx,
                                 StartTile: ranked_tile[tile_idx].StartTile,
                                 EndTile: ranked_tile[tile_idx].EndTile,
                                 _z:"c" })
          continue
        }
        //-------------------------------


        alt_entry := simple_indel_diff([]byte(canon_seq), []byte(ranked_tile[tile_idx].Seq))

        // The alt sequence to be stored is small enough so put it in and
        // move on.
        // In other words, we tried a simple alignment (by just constructing
        // a single simple substitution) and got an 'alt' sequence that's less
        // than 200 bp.  It's long but not long enough for us to worry about
        // it.  This heuristic can change in the future, but 200 is a good
        // number as a cutoff to not fiddle with it too much more.  So,
        // even though the alt sequence is longish (but still less than 200bp),
        // treat it like a normal alt against the canonical sequence.
        //
        if len(alt_entry.Seq) < 200 {
          z := ranked_tile[tile_idx]
          z.Diff = append(z.Diff, alt_entry)
          ranked_tile[tile_idx] = z

          // -----------------
          // SIMPLE INDEL TILE
          // -----------------
          //
          // The simple indel (single alt of type indel) was good enough.
          //
          aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
            AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                                 AuxBodySeqFlag: false,
                                 Type: 0,
                                 Span: ranked_tile[tile_idx].Span,
                                 VariantId:tile_idx, Diff:z.Diff,
                                 StartTile: ranked_tile[tile_idx].StartTile,
                                 EndTile: ranked_tile[tile_idx].EndTile,
                                 _z:"d" })
          continue
        }

        // The simple heuristics above have failed, do a global alignment on previous previous
        // auxiliary sequences to try and get the stored diff size lower.  If we don't
        // succeed in getting it lower than a certain threshold, add the whold sequence
        // to the auxiliary sequence structure.
        //
        dy := memz.New()

        // First try global alignment vs. the canonical sequence
        //
        adelta := dy.AlignDelta([]byte(canon_seq), []byte(ranked_tile[tile_idx].Seq))

        bplen := 0
        for xx:=0; xx<len(adelta); xx++ {
          bplen += adelta[xx].Len
        }

        // The simple alignemnt has a sequence that's too large.  Let's try a
        // global alignment against the canonical sequence.  If the global alignment
        // is less than 200 (a heuristic) store it as an alt and move on.
        //
        if bplen<200 {
          a_diff := altentry_from_delta(adelta, []byte(ranked_tile[tile_idx].Seq))
          z := ranked_tile[tile_idx]
          z.Diff = append(z.Diff, a_diff...)
          ranked_tile[tile_idx] = z

          if (g_sanity_flag) {
            tseq := seq_from_diff(a_diff, []byte(canon_seq))
            if string(tseq) != string(ranked_tile[tile_idx].Seq) {
              fmt.Printf("tseq != ranked_tile[%d].Seq\n", tile_idx)
              fmt.Printf("delt: %s\n", _delta_str(adelta))
              fmt.Printf("diff: %s\n", _diff_str(a_diff))
              fmt.Printf("canon_seq\n%s\n", canon_seq)
              fmt.Printf("tseq:\n%s\n", tseq)
              fmt.Printf("ranked_tile[%d].Seq:\n", tile_idx)
              fmt.Printf("%s\n", ranked_tile[tile_idx].Seq)
              os.Exit(1)
            } else {
            }
          }

          // -----------------
          // SIMPLE ALIGN TILE
          // -----------------
          //
          // We've done a global alignment against the canonical tile and got a good enough answer
          //
          aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
            AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                                 AuxBodySeqFlag: false,
                                 Type: 0,
                                 Span: ranked_tile[tile_idx].Span,
                                 VariantId:tile_idx,
                                 Diff:z.Diff,
                                 StartTile: ranked_tile[tile_idx].StartTile,
                                 EndTile: ranked_tile[tile_idx].EndTile,
                                 _z:"e" })
          continue
        }

        // The heuristics up to this point have failed (simple sub., simple indel,
        // global aligment against canonical).  Now we go through each of the
        // 'aux canonical sequences' to see if we get a global aligment that
        // is of reasonable size.
        // If we find a good global alignment, we'll use it, marking the alternative
        // canonincal sequence we used to reference the diffs against.
        // If we haven't found good alignmeents against the alternate canonical sequences,
        // we'll add it to the alternative canonical list in the hopes future sequences
        // will align better to the new addition.
        //

        // Do search against aux canonical sequences up to this point.
        //
        aa:=0
        reclen:=0
        for aa=0; aa<len(aux_canon_seq); aa++ {

          adelta = dy.AlignDelta([]byte(aux_canon_seq[aa]), []byte(ranked_tile[tile_idx].Seq))

          // Get a rough estimate on how big the entry will be
          //
          bplen = 0
          for xx:=0; xx<len(adelta); xx++ {
            bplen += adelta[xx].Len
            reclen++
          }

          // And if it passes our heuristic, use it.
          //
          if bplen < 200 { break }
        }

        // We've found a suitable match for the auxiliary sequence.
        //
        if aa < len(aux_canon_seq) {

          a_diff := altentry_from_delta(adelta, []byte(ranked_tile[tile_idx].Seq))

          z := ranked_tile[tile_idx]
          z.AuxFlag = true
          z.Diff = a_diff
          ranked_tile[tile_idx] = z

          if (g_sanity_flag) {
            tseq := seq_from_diff(a_diff, []byte(aux_canon_seq[aa]))
            if string(tseq) != string(ranked_tile[tile_idx].Seq) {
              fmt.Printf("tseq != ranked_tile[%d].Seq\n", tile_idx)
              fmt.Printf("delt: %s\n", _delta_str(adelta))
              fmt.Printf("diff: %s\n", _diff_str(a_diff))
              fmt.Printf("aux_canon_seq[%d]\n%s\n", aa, aux_canon_seq[aa])
              fmt.Printf("tseq:\n%s\n", tseq)
              fmt.Printf("ranked_tile[%d].Seq:\n", tile_idx)
              fmt.Printf("%s\n", ranked_tile[tile_idx].Seq)
              os.Exit(1)
            } else {
            }
          }

          // --------------
          // AUX ALIGN TILE
          // --------------
          //
          // We've done a global alignment against the auxiliary "canonical" sequences and got a good enough answer
          //
          aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
            AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                                 AuxBodySeqFlag: false,
                                 Type: 0,
                                 Span: ranked_tile[tile_idx].Span,
                                 VariantId:tile_idx, Diff:z.Diff,
                                 BaseIndex:aa+1,
                                 AuxVariantId: aux_canon_var_id[aa],
                                 StartTile: ranked_tile[tile_idx].StartTile,
                                 EndTile: ranked_tile[tile_idx].EndTile,
                                 _z:"f" })

          continue
        }



        // We didn't find a suitable match so add to the auxiliary sequence list
        // (AuxBodySeqFlat = true -> this is an aux body sequence)
        //
        aux_canon_seq = append(aux_canon_seq, []byte(ranked_tile[tile_idx].Seq))
        aux_canon_var_id = append(aux_canon_var_id, tile_idx)

        alt_entry.Seq = []byte(ranked_tile[tile_idx].Seq)

        z := ranked_tile[tile_idx]
        z.AuxFlag = true
        ranked_tile[tile_idx] = z

        // ------------
        // NEW AUX TILE
        // ------------
        //
        // All heuristics have failed, add it to the auxiliary "canonical" tile list in the hopes
        // of helping future tile matches.
        //
        new_ovf_idx := len(aux_canon_seq)
        aux_overflow_lib[path][step] = append(aux_overflow_lib[path][step],
          AuxOverflowLibEntry{ Seq: ranked_tile[tile_idx].Seq,
                               AuxBodySeqFlag: true,
                               Type: 0,
                               BaseIndex: new_ovf_idx,
                               Span: ranked_tile[tile_idx].Span,
                               VariantId: tile_idx,
                               AuxVariantId: tile_idx,
                               StartTile: ranked_tile[tile_idx].StartTile,
                               EndTile: ranked_tile[tile_idx].EndTile,
                               _z:"g" })

      }

    }
  }

}

// path step key for the tag prefix of path:step
//
var g_tagset map[int]map[int]string

func load_tagset(scan *autoio.AutoioHandle) error {
  g_tagset = make(map[int]map[int]string)

  line_no := 0
  for scan.ReadScan() {
    l := scan.ReadText();
    line_no++

    if len(l)==0 { continue }
    if l[0]=='\n' { continue }
    if l[0]=='#' { continue }
    if l[0]=='>' { continue }

    parts := strings.Split(l, ",")

    ind:=0
    _path,e := strconv.ParseInt(parts[ind], 16, 64)
    if e!=nil { return e }
    ind++

    _step,e := strconv.ParseInt(parts[ind], 16, 64)
    if e!=nil { return e }
    ind++

    tag := parts[ind]
    ind++

    path := int(_path)
    step := int(_step)

    if _,ok := g_tagset[path] ; !ok { g_tagset[path] = make(map[int]string) }

    g_tagset[path][step] = tag

  }

  return nil
}

func load_fastj(scan *autoio.AutoioHandle) error {
  line_no:=0

  tilelib_map = make(map[int]map[int]map[int]map[string]TileLibEntry)

  cur_seq := make([]byte, 0, 1024)
  tilepath := -1
  tilestep := -1
  nocall := 0

  _ = nocall

  var first_tile bool = true
  var tileid string
  var s_tag string
  var e_tag string
  var md5sum_str string
  var span_len int
  var start_tile_flag bool
  var end_tile_flag bool

  for scan.ReadScan() {
    l := scan.ReadText()
    line_no++
    if len(l)==0 { continue }
    if l[0]=='\n' { continue }

    if l[0]=='>' {

      // store tile sequence
      //
      if !first_tile {
        m5 := md5sum2str( md5.Sum(cur_seq) )
        if m5!=md5sum_str { return fmt.Errorf("md5sums do not match %s != %s (line %d)", m5, md5sum_str, line_no) }

        if _,ok := tilelib_map[tilepath] ; !ok {
          tilelib_map[tilepath] = make(map[int]map[int]map[string]TileLibEntry)
        }
        if _,ok := tilelib_map[tilepath][tilestep] ; !ok {
          tilelib_map[tilepath][tilestep] = make(map[int]map[string]TileLibEntry)
        }
        if _,ok := tilelib_map[tilepath][tilestep][len(cur_seq)] ; !ok {

          if len(cur_seq)==0 { log.Fatal(fmt.Sprintf(">>>> %03x.%04x len(seq)==0", tilepath, tilestep)) }

          tilelib_map[tilepath][tilestep][len(cur_seq)] = make(map[string]TileLibEntry)
        }
        if _,ok := tilelib_map[tilepath][tilestep][len(cur_seq)][m5] ; !ok {
          if len(cur_seq)==0 { log.Fatal(fmt.Sprintf(">>>> %03x.%04x len(seq)==0", tilepath, tilestep)) }

          // We can have no-calls on tags.  This can lead to different sequences on tags
          // depending on tile frequency with no-calls.  We're building the library
          // and making decisions on how to fill in tiles with no-calls.  For tags
          // we really need to force the issue since tags shouldn't be variable.
          // Do that here.
          //

          st_tag := ""
          if tilestep>0 { st_tag = g_tagset[tilepath][tilestep] }

          en_tag := ""
          if _,ok := g_tagset[tilepath][tilestep+span_len] ; ok {
            en_tag = g_tagset[tilepath][tilestep+span_len]
          }

          for x:=0; x<len(st_tag); x++ { cur_seq[x] = st_tag[x] }
          for x:=0; x<len(en_tag); x++ { cur_seq[len(cur_seq)-len(en_tag)+x] = en_tag[x] }

          tilelib_map[tilepath][tilestep][len(cur_seq)][m5] = TileLibEntry{ Seq:string(cur_seq), Freq:1, Span:span_len, StartTile: start_tile_flag, EndTile: end_tile_flag }

          //DEBUG
          //if end_tile_flag { fmt.Printf("####!!!! [%d] %x.%x\n", len(cur_seq), tilepath, tilestep) }

        } else {
          if len(cur_seq)==0 { log.Fatal( fmt.Sprintf(">>>> %03x.%04x len(seq)==0", tilepath, tilestep)) }

          z := tilelib_map[tilepath][tilestep][len(cur_seq)][m5]
          z.Freq++
          tilelib_map[tilepath][tilestep][len(cur_seq)][m5] = z
        }

      }
      first_tile = false

      var pos int =0

      tileid,pos = simple_text_field(l[1:], "tileID")
      if pos<0 { return fmt.Errorf("no tileID found at line %d", line_no) }

      md5sum_str,pos = simple_text_field(l[1:], "md5sum")
      if pos<0 { return fmt.Errorf("no md5sum found at line %d", line_no) }

      span_len,pos = simple_int_field(l[1:], "seedTileLength")
      if pos<0 { return fmt.Errorf("no md5sum found at line %d", line_no) }

      s_tag,pos = simple_text_field(l[1:], "startTag")
      if pos<0 { return fmt.Errorf("no startTag found at line %d", line_no) }
      _ = s_tag

      e_tag,pos = simple_text_field(l[1:], "endTag")
      if pos<0 { return fmt.Errorf("no endTag found at line %d", line_no) }
      _ = e_tag

      start_tile_flag,pos = simple_bool_field(l[1:], "startTile")
      if pos<0 { return fmt.Errorf("no startTile found at line %d", line_no) }
      _ = start_tile_flag

      end_tile_flag,pos = simple_bool_field(l[1:], "endTile")
      if pos<0 { return fmt.Errorf("no endTile found at line %d", line_no) }
      _ = end_tile_flag



      tile_parts := strings.Split(tileid, ".")
      if t,e := strconv.ParseInt(tile_parts[0], 16, 64) ; e==nil {
        tilepath = int(t)
      } else { return e }

      tile_parts = strings.Split(tileid, ".")
      if t,e := strconv.ParseInt(tile_parts[2], 16, 64) ; e==nil {
        tilestep = int(t)
      } else { return e }

      // Header parsed, go on
      //
      cur_seq = cur_seq[0:0]
      continue
    }

    if first_tile { return fmt.Errorf("found body before header (line %d)", line_no) }

    cur_seq = append(cur_seq, l[:]...)

  }

  // store tile sequence
  //
  if !first_tile {
    m5 := md5sum2str( md5.Sum(cur_seq) )
    if m5!=md5sum_str { return fmt.Errorf("md5sums do not match %s != %s (line %d)", m5, md5sum_str, line_no) }

    if _,ok := tilelib_map[tilepath] ; !ok {
      tilelib_map[tilepath] = make(map[int]map[int]map[string]TileLibEntry)
    }
    if _,ok := tilelib_map[tilepath][tilestep] ; !ok {
      tilelib_map[tilepath][tilestep] = make(map[int]map[string]TileLibEntry)
    }
    if _,ok := tilelib_map[tilepath][tilestep][len(cur_seq)] ; !ok {

      if len(cur_seq)==0 { log.Fatal(fmt.Sprintf(">>>> %03x.%04x len(seq)==0", tilepath, tilestep)) }

      tilelib_map[tilepath][tilestep][len(cur_seq)] = make(map[string]TileLibEntry)
    }
    if _,ok := tilelib_map[tilepath][tilestep][len(cur_seq)][m5] ; !ok {
      if len(cur_seq)==0 { log.Fatal( fmt.Sprintf(">>>> %03x.%04x len(seq)==0", tilepath, tilestep)) }
      for x:=0; x<len(s_tag); x++ { cur_seq[x] = s_tag[x] }
      for x:=0; x<len(e_tag); x++ { cur_seq[len(cur_seq)-len(e_tag)+x] = e_tag[x] }
      tilelib_map[tilepath][tilestep][len(cur_seq)][m5] = TileLibEntry{ Seq:string(cur_seq), Freq:1, Span:span_len, StartTile:start_tile_flag, EndTile:end_tile_flag }

      //DEBUG
      //if end_tile_flag { fmt.Printf("####!!!! (2) [%d] %x.%x\n", len(cur_seq), tilepath, tilestep) }


    }
  }

  return nil
}


type int4 struct { I [4]int }
func bub_mirror(a *int4, b []byte) {
  for i:=0; i<len(a.I); i++ {
    for j:=1; j<(len(a.I)-i); j++ {
      if a.I[j] > a.I[j-1] {
        a.I[j],a.I[j-1]=a.I[j-1],a.I[j]
        b[j],b[j-1]=b[j-1],b[j]
      }
    }
  }
}

func seq_equiv(canon_seq []byte, seq string) bool {
  if len(canon_seq) != len(seq) { return false }
  for i:=0; i<len(seq); i++ {
    if seq[i]=='n' || seq[i]=='N' { continue }
    if seq[i]!=canon_seq[i] { return false }
  }
  return true
}

func simple_substitution(canon_seq []byte, seq string) []byte {
  if len(canon_seq)!=len(seq) { return nil }
  b := make([]byte, len(canon_seq))
  for i:=0; i<len(b); i++ {
    b[i] = seq[i]
    if seq[i]=='n' || seq[i]=='N' { b[i] = canon_seq[i] }
  }
  return b
}


// Scan last matching bp and first matching bp.  Give back
// a single indel in the middle.  Scan from the end first
// to favor indels that have an indel at the beginning.
//
// Return and AltEntry that has the start in the canonincal
// sequence, the length of bp sequence in the canonical seq.
// and the replacement sequence.
//
// example:
//
// { 50, 1, {} } -> straight deletion of 1bp at (0ref) bp pos 50
// { 50, 0, {ac} } -> straight insertion of 2bp seq 'ac' at (0ref) bp pos 50
// { 50, 2, {t} } -> indel, replacing 2bp of canon seq. with seq 't' at (0ref) bp pos 50
//
func simple_indel_diff(canon_seq, alt_seq  []byte) AltEntry {
  var s int = 0

  n := len(canon_seq)
  m := len(alt_seq)

  for (n>0) && (m>0) {
    if canon_seq[n-1] != alt_seq[m-1] { break }
    n--
    m--
  }

  for s=0; s < n; s++ {
    if s >= m { break }
    if canon_seq[s] != alt_seq[s] { break }
  }

  del_n := n-s
  del_m := m-s

  if del_m>0 {
    return AltEntry{s,del_n,[]byte(alt_seq[s:s+del_m])}
  }
  return AltEntry{s,del_n,[]byte{}}
}

func simple_subs_diff(canon_seq, alt_seq []byte) []AltEntry {
  a := make([]AltEntry, 0, 3)
  n:=len(canon_seq)

  for i:=0; i<n; i++ {
    b := i
    for ; (i<n) && (canon_seq[i]!=alt_seq[i]) ; i++ { }
    if b<i { a = append(a, AltEntry{b,i-b, alt_seq[b:i]}) }
  }

  return a
}

var final_tile_lib_map map[int]map[int]map[string]TileLibEntry
var tilepos_canon_seq map[int]map[int][]byte
var tilepos_canon_span map[int]map[int]int
var all_canon_seq map[int]map[int]map[int][]byte
var all_canon_freq map[int]map[int]map[int]int
var rank_tile_lib_map map[int]map[int][]TileLibEntry

type CGLFRawBytes struct {

  TagSeq []byte

  BodySeq []byte
  BodyPos uint64

  BodyOffset []uint64

  Span    []uint64



}

func create_alt_cache(cachebuf []byte, ranked_var []TileLibEntry) int {
  n := len(cachebuf)

  tbuf := make([]byte, n)

  for i:=0; i<n; i++ { cachebuf[i] = 0 }

  canon_span := ranked_var[0].Span

  // The 0th position is the canonical tile which, by definition, does
  // not have any alts on it.  Start at 1.
  //
  pos := 0
  for i:=1; i<len(ranked_var); i++ {

    if canon_span != ranked_var[i].Span { return i }
    if ranked_var[i].AuxFlag { return i }

    altnum_byte := dlug.MarshalUint64(uint64(len(ranked_var[i].Diff)))

    // We put the packed bytes into tbuf using sentinal_pos for position.
    // If, after conversion, we find it still would fit into cachebuf,
    // we copy it over.
    //
    // If we spill over, return our current position.
    //
    sentinal_pos := pos

    for j:=0; j<len(altnum_byte); j++ {
      tbuf[sentinal_pos] = altnum_byte[j]
      sentinal_pos++
      if sentinal_pos >= len(tbuf) { return i }
    }

    for j:=0; j<len(ranked_var[i].Diff); j++ {
      start_byte := dlug.MarshalUint64(uint64(ranked_var[i].Diff[j].Start))
      canon_len_byte := dlug.MarshalUint64(uint64(ranked_var[i].Diff[j].N))
      alt_len_byte := dlug.MarshalUint64(uint64(len(ranked_var[i].Diff[j].Seq)))
      alt_2bit_byte_len := len(ranked_var[i].Diff[j].Seq)/4
      if len(ranked_var[i].Diff[j].Seq)%4!=0 { alt_2bit_byte_len++ }

      alt_2bit_seq := make([]byte, alt_2bit_byte_len)
      seq_to_2bit(alt_2bit_seq, []byte(ranked_var[i].Diff[j].Seq))


      for k:=0; k<len(start_byte); k++ {
        tbuf[sentinal_pos] = start_byte[k]
        sentinal_pos++
        if sentinal_pos >= len(tbuf) { return i }
      }

      for k:=0; k<len(canon_len_byte); k++ {
        tbuf[sentinal_pos] = canon_len_byte[k]
        sentinal_pos++
        if sentinal_pos >= len(tbuf) { return i }
      }

      for k:=0; k<len(alt_len_byte); k++ {
        tbuf[sentinal_pos] = alt_len_byte[k]
        sentinal_pos++
        if sentinal_pos >= len(tbuf) { return i }
      }

      for k:=0; k<len(alt_2bit_seq); k++ {
        tbuf[sentinal_pos] = alt_2bit_seq[k]
        sentinal_pos++
        if sentinal_pos >= len(tbuf) { return i }
      }

    }

    if sentinal_pos >= n { return i }
    for j:=pos; j<sentinal_pos; j++ { cachebuf[j] = tbuf[j] }

    pos = sentinal_pos

  }

  // We've processed everything, return the whole length
  //
  return len(ranked_var)
}


// The aux overflow structures are doing the following
//   - encoding the difference between a base sequence and a 2bit sequence as a series of alt records
//   - encoding a base sequence
//   - putting everything else in a gzipped 2bit 'file'
//
// If the sequence is too long and the alt record would be too big, then the sequence is put into
// a gzipped 2bit file for efficiency.  Some sequences are long but largely redundant.  Rather
// than try to finesse them into a more efficient strucutre, we employ a 'brute force' approach
// of dumping them into a short 2bit file, gzipping it and providing the bytes here.
//
// This also opens the door for other auxiliary records that can hold arbitrary data.
//
// TileLibEntry holds 'simple' differences, one that are different from the canonical sequence.
//   There is a flag to indidcate when the AuxOverflowEntry array should be consulted.
// AuxOverflowLibEntry holds whether it is a 'base sequence', an alt sequence (relative to
//   one of the aux base sequences) or to be put into the AltOverflowRecord portion.
//
// Currently, the AltOverlfow array is as follows:
//
//   AltOverflow []{
//
//      NAltCache     dlug
//      N             dlug
//      VariantId     []dlug
//      Span          []dlug
//      VariantType   []dlug
//      VariantIndex  []dlug
//
//      AuxBodyLen          8byte     // Auxiliary body sequences
//      AuxBodySeqOffsetBP  []8byte   // these positions are referenced in AltOverflowRec
//      AuxBodySeqTwoBit    []byte
//
//
//      AltOverflowVariantLen         8byte     // Number of tile variants in the AltOverflow position
//      AltOverflowVariantRecOffset   []8byte   // Pos k holds byte offset of tile variant (k+1) in AltOverflowRec
//      AltOverflowRec                []AltOverflowRecord // see below
//
//      AltDataLen        8byte    // number of AltData entries
//      AltDataOffset     []8byte  // pos k byte sum of AltData[0:k+1] (i.e. AltDataOffset[0] == len(AltData[0]))
//      AltData           []byte
//    }
//
// With the AltOverflowRecord as:
//
// Type 0:
//
//    AltOverflowRec { // Type 0, 2bit alt record
//      BodySeqIndex    dlug
//      NAlt            dlug
//      Alt[] {
//        StartBP       dlug
//        CanonLenBP    dlug
//        AltLenBP      dlug
//        SeqTwoBit     []byte
//      }
//    }
//
// Type 1: AltData is the payload of the data (twobit gzipped)
//
func create_alt_overflow(ovf_buf []byte, ranked_var []TileLibEntry, overflow_var []AuxOverflowLibEntry, start_tile int, path, v, step int) []byte {

  nalt_cache_bytes := dlug.MarshalUint64(uint64(start_tile))
  n_bytes := dlug.MarshalUint64(uint64(len(ranked_var)))
  vid_array_bytes := make([]byte, 0, 8)
  span_array_bytes := make([]byte, 0, 8)
  variant_type_bytes := make([]byte, 0, 8)
  variant_index_bytes := make([]byte, 0, 8)

  aux_body_len_bytes := make([]byte, 8)
  aux_body_seq_offset_bp_bytes := make([]byte, 0, 8)
  aux_body_seq_twobit_bytes := make([]byte, 0, 8)

  alt_overflow_variant_len_bytes := make([]byte, 8)
  alt_overflow_variant_rec_offset_bytes := make([]byte, 0, 8)
  alt_overflow_rec_bytes := make([]byte, 0, 8)

  alt_data_len_bytes := make([]byte, 8)
  alt_data_offset_bytes := make([]byte, 0, 8)
  alt_data_bytes := make([]byte, 0, 8)


  aux_body_seq := make([]byte, 0, 8)

  buf8 := make([]byte, 8)

  // Stuff in (dummy) values for initial `start_tile` entries
  //
  for i:=0; i<start_tile; i++ {
    vid_bytes := dlug.MarshalUint64(uint64(overflow_var[i].VariantId))
    vid_array_bytes = append(vid_array_bytes, vid_bytes...)

    span_bytes := dlug.MarshalUint64(uint64(overflow_var[i].Span))
    span_array_bytes = append(span_array_bytes, span_bytes...)

    var_type_bytes := dlug.MarshalUint64(uint64(0))
    variant_type_bytes = append(variant_type_bytes, var_type_bytes...)

    var_idx_bytes := dlug.MarshalUint64(uint64(0))
    variant_index_bytes = append(variant_index_bytes, var_idx_bytes...)
  }

  varid_idx_map := make(map[int]int)
  idx:=0
  for i:=0; i<len(overflow_var); i++ {
    varid_idx_map[ overflow_var[i].VariantId ] = 0
    if overflow_var[i].Type!=0 { continue }
    if !overflow_var[i].AuxBodySeqFlag { continue }
    idx++
    varid_idx_map[ overflow_var[i].VariantId ] = idx;
  }


  alt_overflow_index := 0

  cur_body_len_bp := 0
  aux_body_len := 0
  alt_overflow_variant_len := 0
  for i:=start_tile; i<len(overflow_var); i++ {
    if overflow_var[i].Type!=0 { continue }

    use_aux_body_seq := false ; _ = use_aux_body_seq

    // AltOverflow.VariantId[]
    //
    vid_bytes := dlug.MarshalUint64(uint64(overflow_var[i].VariantId))
    vid_array_bytes = append(vid_array_bytes, vid_bytes...)

    // AltOverflow.Span[]
    //
    span_bytes := dlug.MarshalUint64(uint64(overflow_var[i].Span))
    span_array_bytes = append(span_array_bytes, span_bytes...)

    // AltOverflow.VariantType[]
    //
    type_bytes := dlug.MarshalUint64(uint64(overflow_var[i].Type))
    variant_type_bytes = append(variant_type_bytes, type_bytes...)

    // AltOverflow.VariantIndex[]
    //
    // 0 is canonical sequence.  Position in AuxBodySeq structure is value-1.
    //
    idx_bytes := dlug.MarshalUint64(uint64(alt_overflow_index))
    variant_index_bytes = append(variant_index_bytes, idx_bytes...)

    alt_overflow_index++


    // AuxFlag is set, we add it to the 'AuxBodySeq' structures
    //
    if overflow_var[i].AuxBodySeqFlag {
      t := len(overflow_var[i].Seq)

      st:=24
      en:=t-24
      if overflow_var[i].StartTile { st = 0 }
      if overflow_var[i].EndTile { en = t }
      //aux_body_seq = append(aux_body_seq, overflow_var[i].Seq[24:t-24]...)
      aux_body_seq = append(aux_body_seq, overflow_var[i].Seq[st:en]...)

      //cur_body_len_bp += len(overflow_var[i].Seq[24:t-24])
      cur_body_len_bp += len(overflow_var[i].Seq[st:en])
      tobyte64(buf8, uint64(cur_body_len_bp))
      aux_body_seq_offset_bp_bytes = append(aux_body_seq_offset_bp_bytes, buf8...)

      use_aux_body_seq = true
      aux_body_len++
    }

    rec_bytes := make([]byte, 0, 8)

    // Canonical sequence as reference
    //
    body_idx_bytes := dlug.MarshalUint64(uint64(overflow_var[i].BaseIndex))

    // AuxOverflow.AltOverflowRec.BodySeqIndex
    //
    rec_bytes = append(rec_bytes, body_idx_bytes...)

    // AuxOverflow.AltOverflowRec.NAlt
    //
    nalt_bytes := dlug.MarshalUint64(uint64(len(overflow_var[i].Diff)))
    rec_bytes = append(rec_bytes, nalt_bytes...)

    // AuxOverflow.AltOverflowRec.Alt[]
    //
    for aa:=0; aa<len(overflow_var[i].Diff); aa++ {

      // AuxOverflow.AltOverflowRec.Alt.StartBP
      //
      s_bp := dlug.MarshalUint64(uint64(overflow_var[i].Diff[aa].Start))
      rec_bytes = append(rec_bytes, s_bp...)

      // AuxOverflow.AltOverflowRec.Alt.CanonLenBP
      //
      n_bp := dlug.MarshalUint64(uint64(overflow_var[i].Diff[aa].N))
      rec_bytes = append(rec_bytes, n_bp...)

      // AuxOverflow.AltOverflowRec.Alt.AltLenBP
      //
      n_alt_bp := dlug.MarshalUint64(uint64(len(overflow_var[i].Diff[aa].Seq)))
      rec_bytes = append(rec_bytes, n_alt_bp...)

      // AuxOverflow.AltOverflowRec.Alt.SeqTwoBit
      //
      seq_2bit_byte_len := (len(overflow_var[i].Diff[aa].Seq)+3)/4
      seq_2bit := make([]byte, seq_2bit_byte_len)
      seq_to_2bit(seq_2bit, []byte(overflow_var[i].Diff[aa].Seq))
      rec_bytes = append(rec_bytes, seq_2bit...)

    }

    alt_overflow_rec_bytes = append(alt_overflow_rec_bytes, rec_bytes...)


    tobyte64(buf8, uint64(len(alt_overflow_rec_bytes)))
    alt_overflow_variant_rec_offset_bytes = append(alt_overflow_variant_rec_offset_bytes, buf8...)

    alt_overflow_variant_len++
  }


  // Now process Type 1 tiles (gzipped twobit)
  //
  type1_count := 0

  // Collect gzipped 2bit file
  //
  tb := twobit.NewWriter()

  alt_data_index := 0
  for i:=start_tile; i<len(overflow_var); i++ {
    if overflow_var[i].Type!=1 { continue }

    // AltOverflow.VariantId[]
    //
    vid_bytes := dlug.MarshalUint64(uint64(overflow_var[i].VariantId))
    vid_array_bytes = append(vid_array_bytes, vid_bytes...)

    // AltOverflow.Span[]
    //
    span_bytes := dlug.MarshalUint64(uint64(overflow_var[i].Span))
    span_array_bytes = append(span_array_bytes, span_bytes...)

    // AltOverflow.VariantType[]
    //
    type_bytes := dlug.MarshalUint64(uint64(overflow_var[i].Type))
    variant_type_bytes = append(variant_type_bytes, type_bytes...)

    // AltOverflow.VariantIndex[]
    //
    idx_bytes := dlug.MarshalUint64(uint64(alt_data_index))
    variant_index_bytes = append(variant_index_bytes, idx_bytes...)

    // Name of 2bit sequence is PATH.VER.STEP.VARIANTID
    //
    pvst := fmt.Sprintf("%04x.%02x.%04x.%03x", path,v,step, overflow_var[i].VariantId)
    tb.Add(pvst, overflow_var[i].Seq)

    type1_count++
  }

  alt_data_len:=0
  if type1_count>0 {

    var gz_twobit_bytes bytes.Buffer
    gz_twobit_writer := gzip.NewWriter(&gz_twobit_bytes)

    tb.WriteTo(gz_twobit_writer)
    gz_twobit_writer.Flush()
    gz_twobit_writer.Close()

    // We only allow 1 alt data record here.  In the future this might
    // change but we've hard coded it for now.
    //
    b := gz_twobit_bytes.Bytes()
    offset_bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(offset_bytes, uint64(len(b)))
    alt_data_offset_bytes = append(alt_data_offset_bytes, offset_bytes...)
    alt_data_bytes = append(alt_data_bytes, b...)

    alt_data_len++
  }


  binary.LittleEndian.PutUint64(aux_body_len_bytes, uint64(aux_body_len))
  binary.LittleEndian.PutUint64(alt_overflow_variant_len_bytes, uint64(alt_overflow_variant_len))
  binary.LittleEndian.PutUint64(alt_data_len_bytes, uint64(alt_data_len))

  aux_body_seq_twobit_bytes = make([]byte, (len(aux_body_seq)+3)/4)
  seq_to_2bit(aux_body_seq_twobit_bytes, aux_body_seq)

  ovf_buf = append(ovf_buf, nalt_cache_bytes...)
  ovf_buf = append(ovf_buf, n_bytes...)
  ovf_buf = append(ovf_buf, vid_array_bytes...)
  ovf_buf = append(ovf_buf, span_array_bytes...)
  ovf_buf = append(ovf_buf, variant_type_bytes...)
  ovf_buf = append(ovf_buf, variant_index_bytes...)

  ovf_buf = append(ovf_buf, aux_body_len_bytes...)
  ovf_buf = append(ovf_buf, aux_body_seq_offset_bp_bytes...)
  ovf_buf = append(ovf_buf, aux_body_seq_twobit_bytes...)

  ovf_buf = append(ovf_buf, alt_overflow_variant_len_bytes...)
  ovf_buf = append(ovf_buf, alt_overflow_variant_rec_offset_bytes...)
  ovf_buf = append(ovf_buf, alt_overflow_rec_bytes...)

  ovf_buf = append(ovf_buf, alt_data_len_bytes...)
  ovf_buf = append(ovf_buf, alt_data_offset_bytes...)
  ovf_buf = append(ovf_buf, alt_data_bytes...)

  return ovf_buf
}

var CGLF_MAGIC_STRING string = "{\"cglf\":\"bin\""
var CGLF_VERSION_MAJOR uint32 = 0
var CGLF_VERSION_MINOR uint32 = 1
var CGLF_VERSION_PATCH uint32 = 0

func write_cglf(ofn string) {

  buf := make([]byte, 1025)

  fo,er := os.Create(ofn)
  if er!=nil { panic(er) }
  defer fo.Close()

  cglf := CGLFRawBytes{}
  cglf.BodyPos = 0
  cglf.TagSeq = make([]byte, 0, 1024)
  cglf.BodySeq = make([]byte, 0, 1024)
  cglf.BodyOffset = make([]uint64, 0, 1024)
  cglf.Span = make([]uint64, 0, 1024)

  var byte_count uint64 = 0

  // WRITE header
  //
  s := fmt.Sprintf("%s", CGLF_MAGIC_STRING)
  fo.Write([]byte(s))
  byte_count += uint64(len(s))

  // WRITE semantic version
  //
  a := make([]byte, 4)
  binary.LittleEndian.PutUint32(a, CGLF_VERSION_MAJOR)
  fo.Write(a)
  byte_count += uint64(len(a))

  binary.LittleEndian.PutUint32(a, CGLF_VERSION_MINOR)
  fo.Write(a)
  byte_count += uint64(len(a))

  binary.LittleEndian.PutUint32(a, CGLF_VERSION_PATCH)
  fo.Write(a)
  byte_count += uint64(len(a))

  a8 := make([]byte, 8)
  //var npath uint64 = 1
  var npath uint64 = uint64(len(ranked_tilelib_map))

  // WRITE NPath
  //
  binary.LittleEndian.PutUint64(a8, npath)
  fo.Write(a8)
  byte_count += uint64(len(a8))

  // WRITE Alternate structure stride
  //
  var altstride uint64 = 24
  binary.LittleEndian.PutUint64(a8, altstride)
  fo.Write(a8)
  byte_count += uint64(len(a8))

  // WRITE Tag stride (24 hardcoded for now)
  //
  var tagstride uint64 = 24
  binary.LittleEndian.PutUint64(a8, tagstride)
  fo.Write(a8)
  byte_count += uint64(len(a8))

  pathoffset := make([]uint64, npath)
  pathoffset[0] = byte_count



  var path_buf []bytes.Buffer
  for tilepath := range ranked_tilelib_map {

    path_buf = append(path_buf, bytes.Buffer{})
    path_buf_idx := len(path_buf)-1
    path_fp := bufio.NewWriter(&(path_buf[path_buf_idx]))


    // WRITE Path
    //
    tobyte64(buf[0:8], uint64(tilepath))
    path_fp.Write(buf[0:8])

    taglen := len(ranked_tilelib_map[tilepath])
    tagseq := make([]byte, taglen*24/3) ; _ = tagseq

    raw_seq := make([]byte, 0, 100000)
    raw_tag_seq := make([]byte, 0, 10000)

    tilestep_sort := make([]int, len(ranked_tilelib_map[tilepath]))
    pos := 0
    for tilestep := range ranked_tilelib_map[tilepath] {
      tilestep_sort[pos] = tilestep
      pos++
    }

    sort.Sort(sort.IntSlice(tilestep_sort))

    canon_seq_pos_bp := make([]uint64, len(tilestep_sort))
    span_bytes := make([]byte, len(tilestep_sort))
    span_overflow := make([]uint64, 0)


    // Calculate the canonical sequence
    //
    ntile := len(tilestep_sort)
    step := 0
    n := 0
    prev_step := -1 ; _ = prev_step
    for i:=0; i<len(tilestep_sort); i++ {
      step = tilestep_sort[i]

      canon_seq := ranked_tilelib_map[tilepath][step][0].Seq
      n = len(canon_seq)

//      if (i>0) && (i<(len(tilestep_sort)-1)) {
//        raw_seq = append(raw_seq, canon_seq[24:n-24]...)
//      } else if i==0 {
//        raw_seq = append(raw_seq, canon_seq[0:n-24]...)
//      } else if i==(len(tilestep_sort)-1) {
//        raw_seq = append(raw_seq, canon_seq[24:n]...)
//      }
      st:=24
      en:=n-24
      if ranked_tilelib_map[tilepath][step][0].StartTile { st = 0 }
      if ranked_tilelib_map[tilepath][step][0].EndTile { en = n }
      raw_seq = append(raw_seq, canon_seq[st:en]...)


      if i>0 {
        raw_tag_seq = append(raw_tag_seq, canon_seq[0:24]...)
      }
      canon_seq_pos_bp[i] = uint64(len(raw_seq))


      canon_tile_ent := ranked_tilelib_map[tilepath][step][0]

      tspan := canon_tile_ent.Span
      if tspan < 255 {
        span_bytes[i] = byte(tspan)
      } else {
        span_bytes[i] = 255

        span_overflow = append(span_overflow, uint64(i))
        span_overflow = append(span_overflow, uint64(tspan))
      }


    }

    tobyte64(buf[0:8], uint64(len(raw_tag_seq)))
    tag_seq_2bit := make([]byte, len(raw_tag_seq)*2/8)
    for i:=0; i<len(raw_tag_seq); i+=24 {
      p:=i*2/8
      if (p+6)>len(tag_seq_2bit) { panic("whoops") }
      if (i+24)>len(raw_tag_seq) { panic("hoops") }
      seq_to_2bit(tag_seq_2bit[p:p+6], raw_tag_seq[i:i+24])
    }

    // WRITE NStep
    //
    var nstep uint64
    nstep = uint64(len(tag_seq_2bit)*4/24)+1
    tobyte64(buf[0:8], nstep)
    path_fp.Write(buf[0:8])


    // WRITE TagSeq2Bit
    //
    path_fp.Write(tag_seq_2bit)


    rr := 0
    if (len(raw_seq)%4) > 0 { rr = 1 }
    body_seq_len_byte := uint64(len(raw_seq)/4 + rr)
    body_seq_len_bp := uint64(len(raw_seq))

    body_seq_2bit := make([]byte, body_seq_len_byte)
    seq_to_2bit(body_seq_2bit, raw_seq)

    // WRITE BodySeqLenBP
    //
    tobyte64(buf[0:8], body_seq_len_bp)
    path_fp.Write(buf[0:8])

    if uint64(len(buf)) < body_seq_len_byte { buf = make([]byte, body_seq_len_byte) }

    // WRITE BodySeq2BitOffset
    //
    for i:=0; i<len(canon_seq_pos_bp); i++ {
      tobyte64(buf[8*i:8*(i+1)], canon_seq_pos_bp[i])
    }
    path_fp.Write(buf[0:8*len(canon_seq_pos_bp)])

    // WRITE BodySeq2Bit
    //
    path_fp.Write(body_seq_2bit)

    // WRITE Span
    //
    path_fp.Write(span_bytes)

    // WRITE SpanOverflowLen
    //
    span_overflow_len := uint64(len(span_overflow)/2)
    tobyte64(buf[0:8], span_overflow_len)
    path_fp.Write(buf[0:8])

    // WRITE SpanOverflow
    //
    if span_overflow_len>0 {
      tbuf := make([]byte, 16*span_overflow_len)
      for ii:=uint64(0); ii<span_overflow_len; ii++ {
        tobyte64(tbuf[16*ii  :16*ii+ 8],  span_overflow[2*ii])
        tobyte64(tbuf[16*ii+8:16*ii+16],  span_overflow[2*ii+1])
      }
      path_fp.Write(tbuf)
    }


    alt_cache_buf := make([]byte, altstride*uint64(ntile))
    alt_overflow_buf := make([]byte, 0, 1024)
    alt_overflow_offset := make([]byte, 8*ntile)

    for i:=0; i<len(tilestep_sort); i++ {
      step := tilestep_sort[i]

      n := create_alt_cache(alt_cache_buf[altstride*uint64(i):altstride*uint64(i+1)], ranked_tilelib_map[tilepath][step])

      for ii:=0; ii<n; ii++ {
        ranked_tilelib_map[tilepath][step][ii].InCache = true
        aux_overflow_lib[tilepath][step][ii].InCache = true
      }

      alt_overflow_buf = create_alt_overflow(alt_overflow_buf, ranked_tilelib_map[tilepath][step], aux_overflow_lib[tilepath][step], n, tilepath, 0, step)
      tobyte64(alt_overflow_offset[8*i:8*(i+1)], uint64(len(alt_overflow_buf)))
    }

    // WRITE Alt(cache)
    //
    path_fp.Write(alt_cache_buf)

    // WRITE AltOverflowOffset
    //
    path_fp.Write(alt_overflow_offset)

    // WRITE AltOverflow
    //
    path_fp.Write(alt_overflow_buf)

    path_fp.Flush()
  }

  // WRITE tile sequence offsets
  //
  tileseqoffset := make([]uint64, len(path_buf), len(path_buf)+1)

  prev_offset := byte_count
  for i:=0; i<len(path_buf); i++ {
    path_bytes := path_buf[i].Bytes()
    tileseqoffset[i] = uint64(len(path_bytes)) + prev_offset
    tobyte64(buf, tileseqoffset[i])
    fo.Write(buf[0:8])

    prev_offset += uint64(len(path_bytes))
  }

  for i:=0; i<len(path_buf); i++ {
    path_bytes := path_buf[i].Bytes()
    fo.Write(path_bytes)
  }

}

func write_csv() {

  for path := range ranked_tilelib_map {
    for step := range ranked_tilelib_map[path] {
      for var_idx:=0; var_idx<len(ranked_tilelib_map[path][step]); var_idx++ {

        ch_code := '_'
        if !ranked_tilelib_map[path][step][var_idx].InCache {
          if aux_overflow_lib[path][step][var_idx].AuxBodySeqFlag {
            ch_code = '!'
          } else { ch_code = '^' }
        }
        fmt.Printf("%c,%04x.00.%04x.%03x+%x,%s\n",
          ch_code,
          path, step, var_idx,
          ranked_tilelib_map[path][step][var_idx].Span,
          ranked_tilelib_map[path][step][var_idx].Seq)
      }
    }
  }

}

type FreqM5 struct {
  Freq int
  M5 string
}

type ByFreq []FreqM5
func (f ByFreq) Len() int { return len(f) }
func (f ByFreq) Less(i,j int) bool { return f[i].Freq > f[j].Freq }
func (f ByFreq) Swap(i,j int) { f[i],f[j] = f[j],f[i] }



func _main( c *cli.Context ) {

  g_verbose         = c.Bool("Verbose")
  g_output_format   = c.String("output-format")
  g_tagset_fn       = c.String("tagset")
  //action := c.String("action")
  action := "create"

  if action != "create" {
    fmt.Printf("only the 'create' action is supported now")
    os.Exit(1)
  }

  //ifn := c.String("input") ; _ = ifn
  //ofn := c.String("output")

  if len(c.String("fastj")) == 0 {
    fmt.Fprintf(os.Stderr, "Provide FastJ file\n")
    cli.ShowAppHelp(c)
    os.Exit(1)
  }

  /*
  if len(ofn)==0 {
    fmt.Fprintf(os.Stderr, "Provide output CGLF flie\n")
    cli.ShowAppHelp(c)
    os.Exit(1)
  }
  */

  if len(c.String("tagset"))==0 {
    fmt.Fprintf(os.Stderr, "Provide output tagset flie\n")
    cli.ShowAppHelp(c)
    os.Exit(1)
  }

  tagset_scan,err := autoio.OpenReadScannerSimple( c.String("tagset") )
  if err!=nil { log.Fatal(err) }
  defer tagset_scan.Close()

  e:=load_tagset(&tagset_scan)
  if e!=nil { log.Fatal(e) }

  fastj_scan,err := autoio.OpenReadScannerSimple( c.String("fastj") )
  if err!=nil { log.Fatal(err) }
  defer fastj_scan.Close()

  err = load_fastj(&fastj_scan)
  if err!=nil { log.Fatal(err) }

  create_canonical_sequences()
  fill_sequences()
  rank_sequences()
  calc_ranked_diffs()

  /*
  if action == "create" {
    write_cglf(ofn)
  }
  */

  if g_verbose {
    //write_csv()
    print_aux_overflow()
  }

  return
}

func main() {

  app := cli.NewApp()
  app.Name  = "fastj2cgflib"
  app.Usage = "Create compact genome library (CGLF) from FastJ input stream.  Output format is <code>,<tilepos>,<md5sum>,<sequence>"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "fastj, f",
      Usage: "FastJ input",
    },

    /*
    cli.StringFlag{
      Name: "input, i",
      Usage: "CGLF file (only needed for certain actions)",
    },
    */

    /*
    cli.StringFlag{
      Name: "output, o",
      Value: "lib.cglf",
      Usage: "Output CGLF file (default to 'lib.cglf')",
    },
    */

    cli.StringFlag{
      Name: "tagset, t",
      Usage: "Tagset input (<tilepath>,<tilestep>,<tag>)",
    },

    /*
    cli.StringFlag{
      Name: "action, a",
      Value: "create",
      Usage: "Action (append, create)",
    },
    */

    cli.IntFlag{
      Name: "procs, N",
      Value: -1,
      Usage: "MAXPROCS",
    },

    /*
    cli.StringFlag{
      Name: "output-format, F",
      Value: "gvcf",
      Usage: "Output format: gvcf,compact (defaults to 'gvcf')",
    },
    */

    cli.BoolFlag{
      Name: "Verbose, V",
      Usage: "Verbose flag",
    },

    /*
    cli.BoolFlag{
      Name: "run-tests, T",
      Usage: "Run tests",
    },
    */

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

}
