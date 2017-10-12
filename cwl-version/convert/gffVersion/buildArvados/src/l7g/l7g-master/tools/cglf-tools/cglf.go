package main

// cglf toolkit

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "log"
import _ "sort"

import "github.com/abeconnelly/autoio"
import "github.com/codegangsta/cli"

import "./dlug"

import "strings"
import "crypto/md5"

import "./twobit"
import "bytes"
import "compress/gzip"

import "io/ioutil"
import "strconv"

var VERSION_STR string = "0.1.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "cglf.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "cglf.mprof"

var g_debug bool = false



type TileLibEntry struct {
  Seq string
  Freq int
  Span int
  Diff []AltEntry
}

type AltEntry struct {
  Start int
  N     int
  Seq   []byte
}

type AltOverflowRec struct {
  BodySeqIndex uint64
  Alt []AltEntry
}

type AltOverflow struct {
  NAltCache uint64
  N uint64
  VariantId []uint64
  Span []uint64
  VariantType []uint64
  VariantIndex []uint64

  AuxBodySeq [][]byte
  AltOverflowRec []AltOverflowRec
  AltData [][]byte
}


func _diff_str(d []AltEntry) string {
  s := []string{}
  for i:=0; i<len(d); i++ {
    s_ele := fmt.Sprintf("%x+%x{%s}", d[i].Start, d[i].N, d[i].Seq)
    s = append(s, s_ele)
  }
  return strings.Join(s, ";")
}




func convert_cache_alt_entry(cachebuf []byte) ([]AltEntry, int) {
  ael := []AltEntry{}

  altnum := uint64(0)
  n := 0
  dn := 0

  altnum,dn = dlug.ConvertUint64(cachebuf[n:])
  n += dn
  if altnum==0 { return nil,n }


  var s,canon_len,alt_len uint64

  for i:=uint64(0); i<altnum; i++ {
    s,dn = dlug.ConvertUint64(cachebuf[n:])
    n += dn

    canon_len,dn = dlug.ConvertUint64(cachebuf[n:])
    n += dn

    alt_len,dn = dlug.ConvertUint64(cachebuf[n:])
    n += dn

    seq_2bit_byte_len := alt_len/4
    if (alt_len%4)!=0 { seq_2bit_byte_len++ }
    seq_2bit := cachebuf[n:n+int(seq_2bit_byte_len)]
    n += int(seq_2bit_byte_len)

    seq := make([]byte, alt_len)
    twobit2seq_sn(seq, seq_2bit, 0, alt_len)

    ael = append(ael, AltEntry{ Start:int(s), N:int(canon_len), Seq: seq })

  }

  return ael,n
}


func (alt_ovf *AltOverflow) init() {
  alt_ovf.VariantId = make([]uint64, 0, 8)
  alt_ovf.Span = make([]uint64, 0, 8)
  alt_ovf.VariantType = make([]uint64, 0, 8)
  alt_ovf.VariantIndex = make([]uint64, 0, 8)

  alt_ovf.AuxBodySeq = make([][]byte, 0, 8)
  alt_ovf.AltOverflowRec = []AltOverflowRec{}
  alt_ovf.AltData = make([][]byte, 0, 8)
}

func peel_alt_overflow_bytes(buf []byte) (AltOverflow, uint64) {

  alt_ovf := AltOverflow{}
  alt_ovf.init()


  var n int

  n_alt_cache,dn := dlug.ConvertUint64(buf[n:])
  n += dn

  alt_ovf.NAltCache = n_alt_cache

  var_num,dn := dlug.ConvertUint64(buf[n:])
  n += dn

  alt_ovf.N = var_num

  if var_num==0 { return AltOverflow{}, uint64(n) }



  variant_id_array := make([]uint64, 0, 8)
  for i:=uint64(0); i<uint64(var_num); i++ {
    vid,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    variant_id_array = append(variant_id_array, uint64(vid))
  }

  alt_ovf.VariantId = variant_id_array


  span_array := make([]uint64, 0, 8)
  for i:=uint64(0); i<uint64(var_num); i++ {
    span,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    span_array = append(span_array, uint64(span))

  }

  alt_ovf.Span = span_array


  variant_type_array := make([]uint64, 0, 8)
  for i:=uint64(0); i<uint64(var_num); i++ {
    typ,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    variant_type_array = append(variant_type_array, uint64(typ))

  }

  alt_ovf.VariantType = variant_type_array


  variant_index_array := make([]uint64, 0, 8)
  for i:=uint64(0); i<uint64(var_num); i++ {
    idx,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    variant_index_array = append(variant_index_array, uint64(idx))

  }

  alt_ovf.VariantIndex = variant_index_array


  alt_overflow_idx_to_variantid := make(map[uint64]uint64)
  alt_overflow_idx_to_span := make(map[uint64]uint64)
  for i:=uint64(0); i<uint64(var_num); i++ {
    if variant_type_array[i]!=0 { continue }
    alt_overflow_idx_to_variantid[ variant_index_array[i] ] = variant_id_array[i]
    alt_overflow_idx_to_span[ variant_index_array[i] ] = variant_index_array[i]
  }


  alt_data_idx_to_variantid := make(map[uint64]uint64)
  for i:=uint64(0); i<uint64(var_num); i++ {
    if variant_type_array[i]!=1 { continue }
    alt_data_idx_to_variantid[ variant_index_array[i] ] = variant_id_array[i]
  }


  aux_body_len := byte2uint64(buf[n:n+8])
  n+=8

  aux_body_seq_offset_bp := make([]uint64, aux_body_len)
  for i:=uint64(0); i<aux_body_len; i++ {
    aux_body_seq_offset_bp[i] = byte2uint64(buf[n:n+8])

    n+=8
  }

  var aux_body_seq_twobit []byte
  if aux_body_len>0 {
    dbp := (aux_body_seq_offset_bp[aux_body_len-1]+3)/4

    aux_body_seq_twobit = buf[n:n+int(dbp)]
    n += int(dbp)
  }


  // Construct AuxBodySeq structure
  //
  aux_body_seq := make([][]byte, aux_body_len)
  for i:=uint64(0); i<aux_body_len; i++ {
    s_bp := uint64(0)
    if i>0 { s_bp = aux_body_seq_offset_bp[i-1] }
    e_bp := aux_body_seq_offset_bp[i]

    aux_body_seq[i] = make([]byte, e_bp-s_bp)
    twobit2seq_sn(aux_body_seq[i], aux_body_seq_twobit, s_bp, e_bp-s_bp)
  }

  alt_ovf.AuxBodySeq = aux_body_seq

  /*
  vid_idx := make(map[int]int)
  for i:=0; i<len(alt_ovf.VariantId); i++ {
    vid_idx[ int(alt_ovf.VariantId[i]) ] = i
  }
  */

  /*
  //DEBUG
  fmt.Printf("#VariantId.:")
  for i:=0; i<len(alt_ovf.VariantId); i++ {
    fmt.Printf(" %d", alt_ovf.VariantId[i])
  }
  fmt.Printf("\n")

  fmt.Printf("#VariantTyp:")
  for i:=0; i<len(alt_ovf.VariantType); i++ {
    fmt.Printf(" %d", alt_ovf.VariantType[i])
  }
  fmt.Printf("\n")

  fmt.Printf("#VariantInd:")
  for i:=0; i<len(alt_ovf.VariantIndex); i++ {
    fmt.Printf(" %d", alt_ovf.VariantIndex[i])
  }
  fmt.Printf("\n")
  */


  altovf_vid_map := make(map[int]int)
  for i:=0; i<len(alt_ovf.VariantId); i++ {
    if alt_ovf.VariantType[i] != 0 { continue }
    altovf_vid_map[ int(alt_ovf.VariantIndex[i]) ] = i
  }

  altdata_vid_map := make(map[int]int) ; _ = altdata_vid_map
  for i:=0; i<len(alt_ovf.VariantId); i++ {
    if alt_ovf.VariantType[i] != 1 { continue }
    altdata_vid_map[ int(alt_ovf.VariantIndex[i]) ] = i
  }


  //
  //----

  alt_overflow_variant_len := byte2uint64(buf[n:n+8])
  n+=8

  alt_overflow_variant_rec_offset := make([]uint64, alt_overflow_variant_len)
  for i:=uint64(0); i<alt_overflow_variant_len; i++ {
    alt_overflow_variant_rec_offset[i] = byte2uint64(buf[n:n+8])
    n+=8
  }

  alt_vid_map := make(map[int]string)

  //experimental
  //alt_ovf.AltOverflowRec = make([]AltOverflowRec, len(alt_ovf.VariantIndex))
  alt_ovf.AltOverflowRec = make([]AltOverflowRec, len(alt_ovf.VariantId))

  for ovf_ind:=uint64(0); ovf_ind<alt_overflow_variant_len; ovf_ind++ {
    body_seq_index,dn := dlug.ConvertUint64(buf[n:]) ; _ = body_seq_index
    n += dn

    nalt,dn := dlug.ConvertUint64(buf[n:]) ; _ = nalt
    n += dn

    aor := AltOverflowRec{ BodySeqIndex: body_seq_index }
    aor.Alt = []AltEntry{}


    for alt_ind:=0; alt_ind<int(nalt); alt_ind++ {
      start_bp,dn := dlug.ConvertUint64(buf[n:]) ; _ = start_bp
      n += dn

      canon_len_bp,dn := dlug.ConvertUint64(buf[n:]) ; _ = canon_len_bp
      n += dn

      alt_len_bp,dn := dlug.ConvertUint64(buf[n:]) ; _ = alt_len_bp
      n += dn

      dn = int((alt_len_bp+3)/4)
      seq_twobit := buf[n:n+dn] ; _ = seq_twobit
      n += dn

      seq := make([]byte, alt_len_bp)
      twobit2seq_sn(seq, seq_twobit, 0, alt_len_bp)

      aor.Alt = append(aor.Alt, AltEntry{ Start: int(start_bp), N:int(canon_len_bp), Seq:seq } )


    }

    vid := alt_ovf.VariantId[ altovf_vid_map[int(ovf_ind)] ]
    alt_ovf.AltOverflowRec[vid] = aor
  }

  alt_data_len := byte2uint64(buf[n:n+8])
  n+=8

  alt_data_offset := make([]uint64, alt_data_len, alt_data_len+1)
  for i:=uint64(0); i<alt_data_len; i++ {
    alt_data_offset[i] = byte2uint64(buf[n:n+8])
    n+=8
  }

  alt_ind:=0 ; _ = alt_ind
  if alt_data_len>0 {
    dn = int(alt_data_offset[alt_data_len-1])

    var bytebuf bytes.Buffer
    bytebuf.Write(buf[n:n+dn])

    gz_reader,e := gzip.NewReader(&bytebuf)
    if e!=nil { panic(e) }
    defer gz_reader.Close()

    raw_bytes,e := ioutil.ReadAll(gz_reader)
    if e!=nil { panic(e) }

    twobit_reader := bytes.NewReader(raw_bytes)

    tb,e := twobit.NewReader(twobit_reader)
    if e!=nil { panic(e) }


    nams := tb.Names()
    for zz:=0; zz<len(nams);zz++ {

      name_parts := strings.Split(nams[zz], ".")
      if len(name_parts) < 4 { panic("...") }
      step,e := strconv.ParseInt(name_parts[2], 16, 64)
      varid,e := strconv.ParseInt(name_parts[3], 16, 64)
      if e!=nil { panic("...") }

      _ = step

      bb,e := tb.Read(nams[zz])
      if e!=nil { panic(e) }
      alt_vid_map[int(varid)] = string(bb)

      aor := AltOverflowRec{ BodySeqIndex: 0 }
      aor.Alt = []AltEntry{}
      aor.Alt = append(aor.Alt, AltEntry{ Start:0, N:-1, Seq: bb })

      alt_ovf.AltOverflowRec[varid] = aor
    }

    n+=dn
  }

  return alt_ovf, uint64(n)
}

func exit_error(msg string, err error ) {
  log.Fatal( fmt.Sprintf("%s, got %v", msg, err) )
}

type CGLFPath struct {
  Path  uint64
  NStep  uint64

  TagLenBP uint64
  TagSeqTwoBit  []byte

  BodySeqLenBP uint64
  BodySeqOffsetBP []uint64
  BodySeqTwoBit []byte

  Span []byte

  SpanOverflowLen uint64
  SpanOverflow []uint64

  AltCache []byte

  AltOverflowOffset []uint64
  AltOverflow []byte

}

type CGLF struct {
  Magic string
  VersionMajor  uint32
  VersionMinor  uint32
  VersionPatch  uint32

  NPath uint64
  AltStride uint64
  TagStride uint64

  PathOffset []uint64
  Path []CGLFPath

}



type CGLFPath_ struct {
  Path  uint64
  NStep  uint64

  TagLenBP uint64
  TagSeqTwoBit  []byte

  BodySeqLenBP uint64
  BodySeqOffsetBP []uint64
  BodySeqTwoBit []byte

  Span []byte
  SpanOverflowLen uint64
  SpanOverflow []uint64

  AltCache []byte

  AltOverflowOffset []uint64
  AltOverflow []byte

}


var g_tile_lib map[int]map[int][]TileLibEntry
var g_body_seq map[int]string
var g_tag_seq map[int]string
var g_cglf CGLF

type ByAltEntry []AltEntry
func (f ByAltEntry) Len() int { return len(f) }
func (f ByAltEntry) Less(i,j int) bool { return f[i].Start < f[j].Start }
func (f ByAltEntry) Swap(i,j int) { f[i],f[j] = f[j],f[i] }


func construct_body_seq(base_seq string, alt []AltEntry) string {

  //sort.Sort(ByAltEntry(alt))

  ref_pos := 0
  alt_seq := make([]byte, 0, len(base_seq))
  for i:=0; i<len(alt); i++ {

    if alt[i].Start > ref_pos {

      n:=len(base_seq)
      if (ref_pos >= n) || (alt[i].Start>n) {
        es := fmt.Sprintf("nope: ref_pos %d, alt[%d].Start %d, len(base_seq) %d", ref_pos, i, alt[i].Start, n)
        panic(es)
      }

      alt_seq = append(alt_seq, base_seq[ref_pos:alt[i].Start]...)
      ref_pos = alt[i].Start
    }

    if len(alt[i].Seq)>0 {
      alt_seq = append(alt_seq, alt[i].Seq...)
    }

    ref_pos += alt[i].N
  }

  if ref_pos < len(base_seq) {
    alt_seq = append(alt_seq, base_seq[ref_pos:]...)
  }

  return string(alt_seq)

}


func populate_tile_lib_sequences() {

  for path,_ := range g_tile_lib {

    body_seq := g_body_seq
    tag_seq := g_tag_seq

    for step,_ := range g_tile_lib[path] {

      var base_seq string

      if step==0 {
        base_seq = fmt.Sprintf("%s%s", body_seq[step], tag_seq[step])
      } else if step==len(tag_seq) {
        base_seq = fmt.Sprintf("%s%s", tag_seq[step-1], body_seq[step])
      } else {

        if step<1 || step >= len(tag_seq) {
          es := fmt.Sprintf(">>>>> step %d, len(tag_seq) %d\n", step, len(tag_seq))
          fmt.Printf(es)
          panic(es)
        }
        base_seq = fmt.Sprintf("%s%s%s", tag_seq[step-1], body_seq[step], tag_seq[step])
      }


      for var_idx := 0 ; var_idx < len(g_tile_lib[path][step]); var_idx++ {

        var s string

        if var_idx > 0 {

          //DEBUG
          //fmt.Printf("#(A) %x.%x\n", path, step)

          s = construct_body_seq(base_seq, g_tile_lib[path][step][var_idx].Diff)
        } else {
          s = base_seq
        }

        z := g_tile_lib[path][step][var_idx]
        z.Seq = s
        g_tile_lib[path][step][var_idx] = z
      }

    }
  }
}

func load_cglf_quick(fp *os.File) {
  var n int
  var e error

  buf := make([]byte, 1024)

  magic := "{\"cglf\":\"bin\""
  for i:=0; i<len(magic); i++ {
    n,e := fp.Read(buf[0:1])
    if n!=1 || e!=nil {
      log.Fatal("NOT a CGLF file")
    }
    if buf[0] != magic[i] {
      log.Fatal("NOT a CGLF file")
    }
  }

  g_cglf.Magic = magic

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading major version number", e) }
  ver_maj := byte2uint32(buf[0:4])

  g_cglf.VersionMajor = ver_maj

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading minor version number", e) }
  ver_min := byte2uint32(buf[0:4])

  g_cglf.VersionMinor = ver_min

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading patch version number", e) }
  ver_pat := byte2uint32(buf[0:4])

  g_cglf.VersionPatch = ver_pat

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading number of paths", e) }
  npath := byte2uint64(buf[0:8])

  g_cglf.NPath = npath

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading alt stride", e) }
  altstride := byte2uint64(buf[0:8])

  g_cglf.AltStride = altstride

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading tag stride", e) }
  tagstride := byte2uint64(buf[0:8])

  g_cglf.TagStride = tagstride

  path_offset := make([]uint64, npath)
  for i:=uint64(0); i<npath; i++ {
    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { log.Fatal( fmt.Sprintf("ERROR reading path %d length, %v", i, e) ) }
    path_offset[i] = byte2uint64(buf[0:8])
  }

  //DEBUG
  //fmt.Printf("# >> npath %d\n", npath)

  g_cglf.PathOffset = path_offset

  g_cglf.Path = make([]CGLFPath, g_cglf.NPath)

  for path_idx:=uint64(0) ; path_idx<npath; path_idx++ {

    _path_idx := int(path_idx)

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    path := byte2uint64(buf[0:8])

    g_cglf.Path[_path_idx].Path = path

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    nstep := byte2uint64(buf[0:8])

    g_cglf.Path[_path_idx].NStep = nstep


    //-----------
    // tag
    //-----------

    //DEBUG
    //fmt.Printf("#>>> path %x, tagstride, %d, nstep %d\n", path, tagstride, nstep)

    ntile := nstep
    taglen_bp := (ntile-1)*tagstride
    tag_seq_2bit := make([]byte, ((taglen_bp+3)/4))

    n,e = fp.Read(tag_seq_2bit)
    if n!=len(tag_seq_2bit) || e!=nil { log.Fatal(e) }

    g_cglf.Path[_path_idx].TagLenBP = taglen_bp
    g_cglf.Path[_path_idx].TagSeqTwoBit = tag_seq_2bit

    //-----------
    // tag
    //-----------


    //-----------
    // seq
    //-----------

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    seq_len_bp := byte2uint64(buf[0:8])

    seq_byte_len := ((seq_len_bp+3)/4)

    // Faster to read the whole byte array then convert in memory?
    //
    seq_bp_offset := make([]uint64, ntile)
    tbuf := make([]byte, ntile*8)
    n,e = fp.Read(tbuf)
    if n!=len(tbuf) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading seq offset (len %d) on path[%d] %d, got %v", ntile, path_idx, path, e)) }

    max_body_byte_len := uint64(0)
    for ii:=uint64(0); ii<ntile; ii++ {
      seq_bp_offset[ii] = byte2uint64(tbuf[8*ii:8*(ii+1)])
      if (ii>0) && ((seq_bp_offset[ii]-seq_bp_offset[ii-1]) > max_body_byte_len) {
        max_body_byte_len = seq_bp_offset[ii]-seq_bp_offset[ii-1]
      }
    }

    seq_2bit := make([]byte, seq_byte_len)
    n,e = fp.Read(seq_2bit)
    if n!=len(seq_2bit) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading seq (len %d) on path[%d] %d, got %v", seq_len_bp, path_idx, path, e)) }

    g_cglf.Path[_path_idx].BodySeqLenBP = seq_len_bp
    g_cglf.Path[_path_idx].BodySeqTwoBit = seq_2bit
    g_cglf.Path[_path_idx].BodySeqOffsetBP = seq_bp_offset

    //-----------
    // seq
    //-----------



    //-----------
    // span
    //-----------

    span_byte := make([]byte, ntile)
    n,e = fp.Read(span_byte)
    if n!=len(span_byte) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading span (len %d) on path[%d] %d, got %v", len(span_byte), path_idx, path, e)) }

    n,e = fp.Read(buf[0:8])
    if n!=8 || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading span overflow length on path[%d] %d, got %v", path_idx, path, e)) }
    span_ovf_len := byte2uint64(buf[0:8])

    span_ovf := []uint64{}
    if span_ovf_len > 0 {
      span_ovf = make([]uint64, 2*span_ovf_len)
      tbuf := make([]byte, 16*span_ovf_len)
      n,e = fp.Read(tbuf)
      if n!=len(tbuf) || e!=nil {
        log.Fatal(fmt.Sprintf("ERROR reading span overflow (len %d) on path[%d] %d, got %v", span_ovf_len, path_idx, path, e))
      }
      for ii:=uint64(0); ii<span_ovf_len; ii++ {
        span_ovf[2*ii] = byte2uint64(tbuf[16*ii  :16*ii+ 8])
        span_ovf[2*ii+1] = byte2uint64(tbuf[16*ii+8:16*ii+16])
      }

    }

    g_cglf.Path[_path_idx].Span = span_byte
    g_cglf.Path[_path_idx].SpanOverflowLen = span_ovf_len
    g_cglf.Path[_path_idx].SpanOverflow = span_ovf

    //-----------
    // span
    //-----------


    //-----------
    // AltCache
    //-----------

    alt_cache_byte := make([]byte, altstride*ntile)
    n,e = fp.Read(alt_cache_byte)
    if n!=len(alt_cache_byte) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt cache (expected %d, got %d) path[%d] %d, got %v", len(alt_cache_byte), n, path_idx, path, e))
    }

    g_cglf.Path[_path_idx].AltCache = alt_cache_byte


    //-----------
    // AltCache
    //-----------

    //-----------
    // AltOverflowOffset
    //-----------

    alt_overflow_offset_bytes := make([]byte, 8*ntile)
    n,e = fp.Read(alt_overflow_offset_bytes)
    if n!=len(alt_overflow_offset_bytes) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt overflow offset (expected %d, got %d) path[%d] %d, got %v",
        len(alt_overflow_offset_bytes), n, path_idx, path, e))
    }

    alt_overflow_offset := make([]uint64, ntile)
    for ii:=uint64(0); ii<ntile; ii++ {
      alt_overflow_offset[ii] = byte2uint64(alt_overflow_offset_bytes[8*ii:8*(ii+1)])
    }

    alt_overflow_byte_len := alt_overflow_offset[ntile-1]

    alt_overflow_bytes := make([]byte, alt_overflow_byte_len)
    n,e = fp.Read(alt_overflow_bytes)
    if n!=len(alt_overflow_bytes) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt overflow (expected %d, got %d) path[%d] %d, got %v",
        len(alt_overflow_bytes), n, path_idx, path, e))
    }

    g_cglf.Path[_path_idx].AltOverflowOffset = alt_overflow_offset
    g_cglf.Path[_path_idx].AltOverflow = alt_overflow_bytes


    //-----------
    // AltOverflowOffset
    //-----------
  }


}

func load_cglf(fp *os.File) {
  var n int
  var e error

  buf := make([]byte, 1024)

  magic := "{\"cglf\":\"bin\""
  for i:=0; i<len(magic); i++ {
    n,e := fp.Read(buf[0:1])
    if n!=1 || e!=nil {
      log.Fatal("NOT a CGLF file")
    }
    if buf[0] != magic[i] {
      log.Fatal("NOT a CGLF file")
    }
  }

  g_cglf.Magic = magic

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading major version number", e) }
  ver_maj := byte2uint32(buf[0:4])

  g_cglf.VersionMajor = ver_maj

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading minor version number", e) }
  ver_min := byte2uint32(buf[0:4])

  g_cglf.VersionMinor = ver_min

  n,e = fp.Read(buf[0:4])
  if n<4 || e!=nil { exit_error("ERROR reading patch version number", e) }
  ver_pat := byte2uint32(buf[0:4])

  g_cglf.VersionPatch = ver_pat

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading number of paths", e) }
  npath := byte2uint64(buf[0:8])

  g_cglf.NPath = npath

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading alt stride", e) }
  altstride := byte2uint64(buf[0:8])

  g_cglf.AltStride = altstride

  n,e = fp.Read(buf[0:8])
  if n<8 || e!=nil { exit_error("ERROR reading tag stride", e) }
  tagstride := byte2uint64(buf[0:8])

  g_cglf.TagStride = tagstride

  path_offset := make([]uint64, npath)
  for i:=uint64(0); i<npath; i++ {
    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { log.Fatal( fmt.Sprintf("ERROR reading path %d length, %v", i, e) ) }
    path_offset[i] = byte2uint64(buf[0:8])
  }

  g_cglf.PathOffset = path_offset

  g_cglf.Path = make([]CGLFPath, g_cglf.NPath)

  tag_seq := make([]byte, tagstride)

  g_tile_lib = make(map[int]map[int][]TileLibEntry)

  for path_idx:=uint64(0) ; path_idx<npath; path_idx++ {

    _path_idx := int(path_idx)

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    path := byte2uint64(buf[0:8])

    g_cglf.Path[_path_idx].Path = path

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    nstep := byte2uint64(buf[0:8])

    g_cglf.Path[_path_idx].NStep = nstep


    //-----------
    // tag
    //-----------

    ntile := nstep
    taglen_bp := (ntile-1)*tagstride
    tag_seq_2bit := make([]byte, ((taglen_bp+3)/4))

    n,e = fp.Read(tag_seq_2bit)
    if n!=len(tag_seq_2bit) || e!=nil { log.Fatal(e) }

    g_cglf.Path[_path_idx].TagLenBP = taglen_bp
    g_cglf.Path[_path_idx].TagSeqTwoBit = tag_seq_2bit


    // we can figure out the number of tiles from the tag lengths
    // We don't have a beginning an end tag, so:
    // (number of tiles) = (number of tags) + 1
    //

    g_tag_seq = make(map[int]string)
    tag_seq_map := make(map[int]string)

    for ii:=uint64(0); ii<(ntile-1); ii++ {
      twobit2seq(tag_seq, tag_seq_2bit[6*ii:6*(ii+1)])
      g_tag_seq[int(ii)] = string(tag_seq)
      tag_seq_map[int(ii)] = string(tag_seq)

    }

    //-----------
    // tag
    //-----------


    //-----------
    // seq
    //-----------

    n,e = fp.Read(buf[0:8])
    if n<8 || e!=nil { fmt.Printf("%d %v\n", n, e) }
    seq_len_bp := byte2uint64(buf[0:8])

    seq_byte_len := ((seq_len_bp+3)/4)

    // Faster to read the whole byte array then convert in memory?
    //
    seq_bp_offset := make([]uint64, ntile)
    tbuf := make([]byte, ntile*8)
    n,e = fp.Read(tbuf)
    if n!=len(tbuf) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading seq offset (len %d) on path[%d] %d, got %v", ntile, path_idx, path, e)) }

    max_body_byte_len := uint64(0)
    for ii:=uint64(0); ii<ntile; ii++ {
      seq_bp_offset[ii] = byte2uint64(tbuf[8*ii:8*(ii+1)])
      if (ii>0) && ((seq_bp_offset[ii]-seq_bp_offset[ii-1]) > max_body_byte_len) {
        max_body_byte_len = seq_bp_offset[ii]-seq_bp_offset[ii-1]
      }
    }

    seq_2bit := make([]byte, seq_byte_len)
    n,e = fp.Read(seq_2bit)
    if n!=len(seq_2bit) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading seq (len %d) on path[%d] %d, got %v", seq_len_bp, path_idx, path, e)) }


    g_body_seq = make(map[int]string)
    canon_seq_map := make(map[int]string)

    last_bp_offset := uint64(0)
    for ii:=uint64(0); ii<ntile; ii++ {
      ds := seq_bp_offset[ii]-last_bp_offset
      body_seq := make([]byte, ds)
      twobit2seq_sn(body_seq, seq_2bit, last_bp_offset, ds)
      last_bp_offset = seq_bp_offset[ii]
      g_body_seq[int(ii)] = string(body_seq)
    }

    g_cglf.Path[_path_idx].BodySeqLenBP = seq_len_bp
    g_cglf.Path[_path_idx].BodySeqTwoBit = seq_2bit
    g_cglf.Path[_path_idx].BodySeqOffsetBP = seq_bp_offset


    //-----------
    // seq
    //-----------



    //-----------
    // span
    //-----------

    g_span := make([]uint64, ntile)

    span_byte := make([]byte, ntile)
    n,e = fp.Read(span_byte)
    if n!=len(span_byte) || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading span (len %d) on path[%d] %d, got %v", len(span_byte), path_idx, path, e)) }

    n,e = fp.Read(buf[0:8])
    if n!=8 || e!=nil { log.Fatal(fmt.Sprintf("ERROR reading span overflow length on path[%d] %d, got %v", path_idx, path, e)) }
    span_ovf_len := byte2uint64(buf[0:8])

    span_ovf := []uint64{}
    if span_ovf_len > 0 {
      span_ovf = make([]uint64, 2*span_ovf_len)
      tbuf := make([]byte, 16*span_ovf_len)
      n,e = fp.Read(tbuf)
      if n!=len(tbuf) || e!=nil {
        log.Fatal(fmt.Sprintf("ERROR reading span overflow (len %d) on path[%d] %d, got %v", span_ovf_len, path_idx, path, e))
      }
      for ii:=uint64(0); ii<span_ovf_len; ii++ {
        span_ovf[2*ii] = byte2uint64(tbuf[16*ii  :16*ii+ 8])
        span_ovf[2*ii+1] = byte2uint64(tbuf[16*ii+8:16*ii+16])
      }

    }

    for idx:=0; uint64(idx)<ntile; idx++ {
      g_span[idx] = uint64(span_byte[idx])
      if g_span[idx] == 255 {
        // check overflow
      }
    }

    g_cglf.Path[_path_idx].Span = span_byte
    g_cglf.Path[_path_idx].SpanOverflowLen = span_ovf_len
    g_cglf.Path[_path_idx].SpanOverflow = span_ovf

    //-----------
    // span
    //-----------

    // construct canonical tile sequences explicitely
    //
    for ii:=uint64(0); ii<ntile; ii++ {
      canon_span := g_span[int(ii)]

      pre := ""
      suf := ""
      if ii>0 { pre = tag_seq_map[int(ii)-1] }
      if (ii+canon_span) < ntile { suf = tag_seq_map[int(ii+canon_span-1)] }

      body_seq := g_body_seq[int(ii)]
      canon_seq_map[int(ii)] = pre + string(body_seq) + suf
    }



    //-----------
    // AltCache
    //-----------

    _path := int(path)

    g_tile_lib[_path] = make(map[int][]TileLibEntry)

    alt_cache_byte := make([]byte, altstride*ntile)
    n,e = fp.Read(alt_cache_byte)
    if n!=len(alt_cache_byte) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt cache (expected %d, got %d) path[%d] %d, got %v", len(alt_cache_byte), n, path_idx, path, e))
    }

    idx := 0
    for i:=0; i<len(alt_cache_byte); i+=24 {

      canon_seq := canon_seq_map[int(idx)]

      // add canonical sequence
      //
      g_tile_lib[_path][idx] = append(g_tile_lib[_path][idx], TileLibEntry{ Seq: canon_seq, Freq:1, Span: int(g_span[idx]) })

      n:=0
      for n<24 {
        alt_entry,dn := convert_cache_alt_entry(alt_cache_byte[i+n:i+24])
        n+=dn
        if alt_entry != nil {


          //DEBUG
          //fmt.Printf("#(B) %x.%x\n", path, idx)


          // add cache sequence
          //
          seq := construct_body_seq(canon_seq, alt_entry)
          g_tile_lib[_path][idx] = append(g_tile_lib[_path][idx], TileLibEntry{ Seq: seq, Freq:1, Span: int(g_span[int(idx)]), Diff: alt_entry })

        } else { break }
      }
      idx++
    }

    g_cglf.Path[_path_idx].AltCache = alt_cache_byte


    //-----------
    // AltCache
    //-----------

    //-----------
    // AltOverflowOffset
    //-----------

    alt_overflow_offset_bytes := make([]byte, 8*ntile)
    n,e = fp.Read(alt_overflow_offset_bytes)
    if n!=len(alt_overflow_offset_bytes) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt overflow offset (expected %d, got %d) path[%d] %d, got %v",
        len(alt_overflow_offset_bytes), n, path_idx, path, e))
    }

    alt_overflow_offset := make([]uint64, ntile)
    for ii:=uint64(0); ii<ntile; ii++ {
      alt_overflow_offset[ii] = byte2uint64(alt_overflow_offset_bytes[8*ii:8*(ii+1)])
    }

    alt_overflow_byte_len := alt_overflow_offset[ntile-1]

    alt_overflow_bytes := make([]byte, alt_overflow_byte_len)
    n,e = fp.Read(alt_overflow_bytes)
    if n!=len(alt_overflow_bytes) || e!=nil {
      log.Fatal(fmt.Sprintf("ERROR reading alt overflow (expected %d, got %d) path[%d] %d, got %v",
        len(alt_overflow_bytes), n, path_idx, path, e))
    }

    step_idx:=0
    for pos:=uint64(0); pos<uint64(len(alt_overflow_bytes));  {

      alt_ovf,dn := peel_alt_overflow_bytes(alt_overflow_bytes[pos:])
      pos+=dn

      //DEBUG
      //fmt.Printf("#(D) %x.%x\n", _path, step_idx)

      tla := create_overflow_tilelib_entry(step_idx, int(ntile), canon_seq_map[step_idx], alt_ovf, tag_seq_map)

      g_tile_lib[_path][step_idx] = append(g_tile_lib[_path][step_idx], tla...)

      step_idx++
    }

    g_cglf.Path[_path_idx].AltOverflowOffset = alt_overflow_offset
    g_cglf.Path[_path_idx].AltOverflow = alt_overflow_bytes


    //-----------
    // AltOverflowOffset
    //-----------
  }

}

func create_overflow_tilelib_entry(step_idx, ntile int, canon_seq string, alt_ovf AltOverflow, tag_seq_map map[int]string) []TileLibEntry {
  tla := make([]TileLibEntry, alt_ovf.N - alt_ovf.NAltCache, alt_ovf.N - alt_ovf.NAltCache + 1)

  start_tile := alt_ovf.NAltCache

  idx_id_map := make(map[int]int)

  for i:=0; i<len(alt_ovf.VariantId); i++ {
    idx_id_map[ int(alt_ovf.VariantId[i]) ] = i
  }

  for i:=start_tile; i<alt_ovf.N; i++ {
    vid := alt_ovf.VariantId[i] ; _ = vid
    span := int(alt_ovf.Span[i])
    typ := alt_ovf.VariantType[i]
    alt_idx := alt_ovf.VariantIndex[i]

    alt_idx = vid

    if typ==0 {
      alt_rec := alt_ovf.AltOverflowRec[alt_idx]
      aux_body_idx := alt_rec.BodySeqIndex

      pre := ""
      if step_idx > 0 { pre = tag_seq_map[step_idx-1] }

      suf := ""
      if (step_idx+span-1) < ntile {
        suf = tag_seq_map[step_idx+span-1]

        if suf == "" {

          if g_debug {
            fmt.Printf("# DNE step_idx %d, span %d (idx+span-1=%x), ntile %d (%x)\n", step_idx, span, step_idx+span-1, ntile, ntile)
          }
        }
      } else {
        if g_debug {
          fmt.Printf("# SKIPPING step_idx %d + span %d (%x) - 1 (%x) < ntile %d (%x)\n", step_idx, span, step_idx+span-1, ntile, ntile)
        }
      }


      base_seq := canon_seq
      if aux_body_idx > 0 {
        base_seq = pre + string(alt_ovf.AuxBodySeq[aux_body_idx-1]) + suf
      }

      //DEBUG
      //fmt.Printf("#(C) %x\n", vid)
      //fmt.Printf("#%s\n", _diff_str(alt_ovf.AltOverflowRec[alt_idx].Alt))
      //fmt.Printf("#[%d]\n#%s\n", len(base_seq), base_seq)



      seq := construct_body_seq(base_seq, alt_ovf.AltOverflowRec[alt_idx].Alt)
      tla[vid-start_tile] = TileLibEntry{ Seq: seq, Freq:1, Span: int(span), Diff: alt_rec.Alt }

    } else if typ==1 {
      tla[vid-start_tile] = TileLibEntry{ Seq: string(alt_ovf.AltOverflowRec[vid].Alt[0].Seq), Freq:1, Span: int(span) }
    }

  }

  return tla
}

func emit_csv() {

  for path,_ := range g_tile_lib {
    for step,_ := range g_tile_lib[path] {
      for var_idx := 0 ; var_idx < len(g_tile_lib[path][step]); var_idx++ {

        fmt.Printf("#%s\n", _diff_str(g_tile_lib[path][step][var_idx].Diff))
        fmt.Printf("%04x.00.%04x.%03x+%x,%s\n", path, step, var_idx, g_tile_lib[path][step][var_idx].Span, g_tile_lib[path][step][var_idx].Seq)

      }
    }
  }

}

func show_header_info() {
  npath := g_cglf.NPath

  ver_maj := g_cglf.VersionMajor
  ver_min := g_cglf.VersionMinor
  ver_pat := g_cglf.VersionPatch
  altstride := g_cglf.AltStride
  tagstride := g_cglf.TagStride

  fmt.Printf("Magic: %s\n", g_cglf.Magic)
  fmt.Printf("Version: %d.%d.%d\n", ver_maj, ver_min, ver_pat)
  fmt.Printf("NPath: %d\n", npath)
  fmt.Printf("AltStride: %d\n", altstride)
  fmt.Printf("TagStride: %d\n", tagstride)

  path_offset := g_cglf.PathOffset

  fmt.Printf("PathOffset(%d):", npath)
  for i:=uint64(0); i<npath; i++ { fmt.Printf(" 0x%x", path_offset[i]) }
  fmt.Printf("\n")

  for path_idx:=uint64(0) ; path_idx<npath; path_idx++ {

    path_name := int(g_cglf.Path[path_idx].Path)
    path := g_cglf.Path[path_idx]


    //-----------

    fmt.Printf("------\n")
    fmt.Printf("  Path: %d (%x)\n", path_name, path_name)
    fmt.Printf("  NStep: %d (%x)\n", path.NStep, path.NStep)

    fmt.Printf("  TagLenBP: %d\n", path.TagLenBP)
    fmt.Printf("  TagSeqTwoBit: %d (bytes)\n", len(path.TagSeqTwoBit))

    fmt.Printf("\n")

    fmt.Printf("  BodySeqLenBP: %d\n", path.BodySeqLenBP)
    fmt.Printf("  BodySeqOffsetBP: %d entries (uint64)\n", len(path.BodySeqOffsetBP))
    fmt.Printf("  BodySeqTwoBit: %d (bytes)\n", len(path.BodySeqTwoBit))
    fmt.Printf("\n")

    fmt.Printf("  Span: %d (bytes)\n", len(path.Span))
    fmt.Printf("  SpanOverflowLen: %d\n", path.SpanOverflowLen)
    fmt.Printf("  SpanOverflow: %d entries (uint64)\n", len(path.SpanOverflow))
    fmt.Printf("\n")

    fmt.Printf("  AltCache: %d (bytes), %d entries (%d byte stride)\n", len(path.AltCache), len(path.AltCache)/int(altstride), altstride)
    fmt.Printf("  AltOverflowOffset: %d entries (uint64)\n", len(path.AltOverflowOffset))
    fmt.Printf("  AltOverflow: %d (bytes)\n", len(path.AltOverflow))

    fmt.Printf("\n")

    /*
    var alt_seq_len_bp uint64
    for step := range g_tile_lib[path_name] {
      for alt_idx:=0; alt_idx<len(g_tile_lib[path_name][step]); alt_idx++ {
        for idx := 0; idx<len(g_tile_lib[path_name][step][alt_idx].Diff); idx++ {
          alt_seq_len_bp += uint64(len(g_tile_lib[path_name][alt_idx][alt_idx].Diff[idx].Seq))
        }
      }
    }

    fmt.Printf("  Alt Diff Sequence length: %d BP (~%d bytes)\n", alt_seq_len_bp, alt_seq_len_bp/4)
    */
  }


}

func show_info() {

  npath := g_cglf.NPath
  tag_seq := make([]byte, g_cglf.TagStride)

  ver_maj := g_cglf.VersionMajor
  ver_min := g_cglf.VersionMinor
  ver_pat := g_cglf.VersionPatch
  altstride := g_cglf.AltStride
  tagstride := g_cglf.TagStride

  fmt.Printf("Magic: %s", g_cglf.Magic)
  fmt.Printf("Version: %d.%d.%d\n", ver_maj, ver_min, ver_pat)
  fmt.Printf("NPath: %d\n", npath)
  fmt.Printf("AltStride: %d\n", altstride)
  fmt.Printf("TagStride: %d\n", tagstride)

  path_offset := g_cglf.PathOffset

  fmt.Printf("PathOffset(%d):", npath)
  for i:=uint64(0); i<npath; i++ { fmt.Printf(" %x", path_offset[i]) }
  fmt.Printf("\n")

  for path_idx:=uint64(0) ; path_idx<npath; path_idx++ {

    path := int(g_cglf.Path[path_idx].Path)

    taglen_bp := g_cglf.Path[path_idx].TagLenBP
    seq_len_bp := g_cglf.Path[path_idx].BodySeqLenBP

    ntile := (taglen_bp/tagstride) + 1

    tag_seq_2bit := g_cglf.Path[path_idx].TagSeqTwoBit
    seq_bp_offset := g_cglf.Path[path_idx].BodySeqOffsetBP

    seq_2bit := g_cglf.Path[path_idx].BodySeqTwoBit
    span_byte := g_cglf.Path[path_idx].Span

    span_ovf := g_cglf.Path[path_idx].SpanOverflow
    alt_cache_byte := g_cglf.Path[path_idx].AltCache


    //-----------

    fmt.Printf("  Path: %x\n", path)

    fmt.Printf("  TagLenBP: %d\n", taglen_bp)
    fmt.Printf("  TagSeqTwoBit: %d (bytes)\n", len(g_cglf.Path[path_idx].TagSeqTwoBit))

    fmt.Printf("\n")

    fmt.Printf("  BodySeqLenBP: %d\n", seq_len_bp)
    fmt.Printf("  BodySeqOffsetBP: %d entries (uint64)\n", len(g_cglf.Path[path_idx].BodySeqOffsetBP))
    fmt.Printf("  BodySeqTwoBit: %d (bytes)\n", len(g_cglf.Path[path_idx].BodySeqTwoBit))
    fmt.Printf("\n")

    fmt.Printf("  Span: %d (bytes)\n", len(g_cglf.Path[path_idx].Span))
    fmt.Printf("  SpanOverflowLen: %d\n", g_cglf.Path[path_idx].SpanOverflowLen)
    fmt.Printf("  SpanOverflow: %d entries (uint64)\n", len(g_cglf.Path[path_idx].SpanOverflow))
    fmt.Printf("\n")

    fmt.Printf("  AltCache: %d (bytes), %d entries (%d byte stride)\n", len(g_cglf.Path[path_idx].AltCache), len(g_cglf.Path[path_idx].AltCache)/int(altstride))
    fmt.Printf("  AltOverflowOffset: %d entries (uint64)\n", len(g_cglf.Path[path_idx].AltOverflowOffset))
    fmt.Printf("  AltOverflow: %d (bytes)\n", len(g_cglf.Path[path_idx].AltOverflow))

    fmt.Printf("  (ntile %d)\n", ntile)

    fmt.Printf("\n\n  >> Tag:\n")
    for p:=0; p<len(tag_seq_2bit); p+=6 {
      twobit2seq(tag_seq, tag_seq_2bit[p:p+6])
      fmt.Printf("    [%4x] %02x %02x %02x %02x %02x %02x %s\n",
        p/6,
        tag_seq_2bit[p],
        tag_seq_2bit[p+1],
        tag_seq_2bit[p+2],
        tag_seq_2bit[p+3],
        tag_seq_2bit[p+4],
        tag_seq_2bit[p+5], tag_seq )
    }

    fmt.Printf("\n\n  >> Seq:\n")
    for seq_idx := 0; seq_idx<int(ntile); seq_idx++ {
      bp_s := uint64(0)
      if seq_idx>0 { bp_s = seq_bp_offset[seq_idx-1] }
      bp_e := seq_bp_offset[seq_idx]
      n := bp_e - bp_s

      raw_seq := make([]byte, n)

      twobit2seq_sn(raw_seq, seq_2bit, bp_s, n)
      fmt.Printf("    [%4x %4x+%4x ^%x] %s\n", seq_idx, bp_s, bp_e-bp_s, span_byte[seq_idx], raw_seq[0:n])
    }

    fmt.Printf("\n\n >> SpanOverflow:\n")
    for span_idx:=0; span_idx<len(span_ovf); span_idx++ {
      fmt.Printf("    [%4x] %4x +%2x\n", span_idx, span_ovf[2*span_idx], span_ovf[2*span_idx+1])
    }

    fmt.Printf("\n\n >> AltCache:\n")
    for idx:=0; idx<int(ntile); idx++ {
      fmt.Printf("    [%4x+%x]", idx, altstride)
      for i:=0; i<int(altstride); i++ {
        fmt.Printf(" %02x", alt_cache_byte[altstride*uint64(idx) + uint64(i)])
      }
      fmt.Printf("\n")
    }

    fmt.Printf("\n\n >> AltOveflow:\n")
    fmt.Printf("<todo>")

    fmt.Printf("\n\n >> Alt:\n")
    for idx:=0; idx<int(ntile); idx++ {
      fmt.Printf("    [%4x]", idx)
      for ii:=0; ii<len(g_tile_lib[path][idx]); ii++ {
        fmt.Printf(" ")
        for jj:=0; jj<len(g_tile_lib[path][idx][ii].Diff); jj++ {
          if jj>0 { fmt.Printf(":") }
          d := g_tile_lib[path][idx][ii].Diff[jj]
          fmt.Printf("(%x+%x{%s})", d.Start, d.N, d.Seq)
        }
        fmt.Printf("+%x", g_tile_lib[path][idx][ii].Span)

      }
      fmt.Printf("\n")
    }

    //-----------


  }

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


func show_tile_position(path,step int) {

  for var_idx := 0 ; var_idx < len(g_tile_lib[path][step]); var_idx++ {
    seq := g_tile_lib[path][step][var_idx].Seq
    m5 := md5sum2str(md5.Sum([]byte(seq)))

    fmt.Printf("#%s\n", _diff_str(g_tile_lib[path][step][var_idx].Diff))
    fmt.Printf("%04x.00.%04x.%03x+%x,%s,%s\n", path, step, var_idx, g_tile_lib[path][step][var_idx].Span, m5, g_tile_lib[path][step][var_idx].Seq)

    fmt.Printf("\n")
  }

}


func _main( c *cli.Context ) {

  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  ain,err := os.Open( c.String("input") )
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer ain.Close()

  aout,err := autoio.CreateWriter( c.String("output") ) ; _ = aout
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer func() { aout.Flush() ; aout.Close() }()

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

  //load_cglf(ain)
  //populate_tile_lib_sequences()

  path := 0x2c5
  pos := 0x2e72
  varid := 0x001 ; _ = varid

  switch c.String("action") {
  case "", "info":
    //show_info()
    load_cglf_quick(ain)
    show_header_info()
  case "pos":
    load_cglf_quick(ain)
    show_tile_position(path, pos)
  case "variant":
    load_cglf_quick(ain)
    //show_tile_position(path, pos, varid)
  case "findbytag":
    //?
    load_cglf_quick(ain)
    //find_by_tag(tag)
  case "csv":
    //populate_tile_lib_sequences()
    load_cglf(ain)
    emit_csv()
  }

}

func main() {

  app := cli.NewApp()
  app.Name  = "cglf"
  app.Usage = "cglf"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com}"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT",
    },

    cli.StringFlag{
      Name: "action, a",
      Usage: "ACTION",
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
