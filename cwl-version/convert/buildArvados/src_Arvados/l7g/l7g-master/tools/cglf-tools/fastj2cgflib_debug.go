package main

import "fmt"
import "./dlug"


//DEBUG

func debug_print_alt_overflow_record(buf []byte) {
  n:=0
  n_alt_cache,dn := dlug.ConvertUint64(buf[n:])
  n+=dn

  n_alt_tot,dn := dlug.ConvertUint64(buf[n:])
  n+=dn

  varids := make([]uint64, 0, 8)
  for i:=uint64(0); i<n_alt_tot; i++ {
    varid,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    varids = append(varids, varid)
  }

  spans := make([]uint64, 0, 8)
  for i:=uint64(0); i<n_alt_tot; i++ {
    span,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    spans = append(spans, span)
  }

  vartypes := make([]uint64, 0, 8)
  for i:=uint64(0); i<n_alt_tot; i++ {
    vartype,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    vartypes = append(vartypes, vartype)
  }

  varindexs := make([]uint64, 0, 8)
  for i:=uint64(0); i<n_alt_tot; i++ {
    varindex,dn := dlug.ConvertUint64(buf[n:])
    n+=dn
    varindexs = append(varindexs, varindex)
  }

  aux_body_len := byte2uint64(buf[n:n+8])
  n += 8;

  aux_body_seq_offset_bp := make([]uint64, aux_body_len, aux_body_len+1)
  for i:=uint64(0); i<aux_body_len; i++ {
    aux_body_seq_offset_bp[i] = byte2uint64(buf[n:n+8])
    n+=8
  }


  prev_offset_bp := uint64(0)
  aux_bodyseq := make([][]byte, aux_body_len, aux_body_len+1)
  if (aux_body_len>0) {
    for i:=uint64(0); i<aux_body_len; i++ {
      aux_bodyseq[i] = make([]byte, aux_body_seq_offset_bp[i] - prev_offset_bp)

      twobit2seq_sn(aux_bodyseq[i], buf[n:], prev_offset_bp, aux_body_seq_offset_bp[i] - prev_offset_bp)
      prev_offset_bp = aux_body_seq_offset_bp[i]
    }
    n += int( (aux_body_seq_offset_bp[aux_body_len-1]+3)/4 )
  }


  alt_overflow_variant_len := byte2uint64(buf[n:n+8])
  n+=8

  alt_overflow_variant_rec_offset := make([]uint64, alt_overflow_variant_len, alt_overflow_variant_len+1)
  for i:=uint64(0); i<alt_overflow_variant_len; i++ {
    alt_overflow_variant_rec_offset[i] = byte2uint64(buf[n:n+8])
    n+=8
  }

  alt_overflow_record := []debug_alt_overflow_rec_struct{}
  for i:=uint64(0); i<alt_overflow_variant_len; i++ {
    var dn int
    rec := debug_alt_overflow_rec_struct{}

    rec.BodySeqIndex,dn = dlug.ConvertUint64(buf[n:])
    n += dn

    rec.NAlt,dn = dlug.ConvertUint64(buf[n:])
    n += dn

    rec.Alt = make([]debug_alt_rec_struct, rec.NAlt)
    for j:=uint64(0); j<rec.NAlt; j++ {
      rec.Alt[j].StartBP,dn = dlug.ConvertUint64(buf[n:])
      n += dn

      rec.Alt[j].CanonLenBP,dn = dlug.ConvertUint64(buf[n:])
      n += dn

      rec.Alt[j].AltLenBP,dn = dlug.ConvertUint64(buf[n:])
      n += dn

      if rec.Alt[j].AltLenBP > 0 {

        twobit_byte_len := (rec.Alt[j].AltLenBP+3)/4

        rec.Alt[j].SeqTwoBit = buf[n:n+int(twobit_byte_len)]
        rec.Alt[j].Seq = make([]byte, rec.Alt[j].AltLenBP)

        twobit2seq_sn(rec.Alt[j].Seq, buf[n:], 0, rec.Alt[j].AltLenBP)

        n += int(twobit_byte_len)

      } else {
        rec.Alt[j].SeqTwoBit = make([]byte, 0, 1)
        rec.Alt[j].Seq= make([]byte, 0, 1)
      }

    }

    alt_overflow_record = append(alt_overflow_record, rec)
  }

  alt_data_len := byte2uint64(buf[n:n+8])
  n+=8

  alt_data_offset := make([]uint64, alt_data_len, alt_data_len+1)
  for i:=uint64(0); i<alt_data_len; i++ {
    alt_data_offset[i] = byte2uint64(buf[n:n+8])
    n+=8
  }

  alt_data := make([][]byte, alt_data_len, alt_data_len+1)

  prev_alt_data_offset := uint64(0)

  for i:=uint64(0); i<alt_data_len; i++ {
    dn := alt_data_offset[i]-prev_alt_data_offset
    alt_data[i] = buf[n:n+int(dn)]
    n+=int(dn)

    prev_alt_data_offset = alt_data_offset[i]
  }



  fmt.Printf("NAltCache/N: %d / %d\n", n_alt_cache, n_alt_tot)
  //fmt.Printf("N: %d\n", n_alt_cache)

  fmt.Printf("VariantId: [")
  for i:=uint64(0); i<n_alt_tot; i++ { fmt.Printf(" %d", varids[i]) }
  fmt.Printf("]\n")

  fmt.Printf("Span: [")
  for i:=uint64(0); i<n_alt_tot; i++ { fmt.Printf(" %d", spans[i]) }
  fmt.Printf("]\n")

  fmt.Printf("VariantType: [")
  for i:=uint64(0); i<n_alt_tot; i++ { fmt.Printf(" %d", vartypes[i]) }
  fmt.Printf("]\n")

  fmt.Printf("VariantIndex: [")
  for i:=uint64(0); i<n_alt_tot; i++ { fmt.Printf(" %d", varindexs[i]) }
  fmt.Printf("]\n")
  fmt.Printf("\n")

  fmt.Printf("AuxBodyLen: %d\n", aux_body_len)
  fmt.Printf("AuxBodySeqOffsetBP: [")
  for i:=uint64(0); i<aux_body_len; i++ {
    fmt.Printf(" %d", aux_body_seq_offset_bp[i])
  }
  fmt.Printf("]\n")

  fmt.Printf("AuxBodySeqTwoBit:\n")
  for i:=uint64(0); i<aux_body_len; i++ {
    fmt.Printf("  [%d] %s\n", i, aux_bodyseq[i])
  }
  fmt.Printf("\n")

  fmt.Printf("AltOverflowVariantLen: %d\n", alt_overflow_variant_len)
  fmt.Printf("AltOverflowVariantRecOffset: [")
  for i:=uint64(0); i<alt_overflow_variant_len; i++ { fmt.Printf(" %d", alt_overflow_variant_rec_offset[i]) }
  fmt.Printf("]\n")

  for i:=0; i<len(alt_overflow_record); i++ {
    rec := alt_overflow_record[i]
    fmt.Printf("  [%d] BodySeqIndex: %d\n", i, rec.BodySeqIndex)
    fmt.Printf("  [%d] NAlt: %d\n", i, rec.NAlt)

    for j:=uint64(0); j<rec.NAlt; j++ {
      fmt.Printf("    [%d.%d] StartBP: %d\n", i, j, rec.Alt[j].StartBP)
      fmt.Printf("    [%d.%d] CanonLenBP: %d\n", i, j, rec.Alt[j].CanonLenBP)
      fmt.Printf("    [%d.%d] AltLenBP: %d\n", i, j, rec.Alt[j].AltLenBP)
      fmt.Printf("    [%d.%d] Seq: %s\n", i, j, rec.Alt[j].Seq)
      fmt.Printf("\n")
    }

  }

  fmt.Printf("\n")

  fmt.Printf("AltDataLen: %d\n", alt_data_len)
  fmt.Printf("AltDataOffset: [")
  for i:=uint64(0); i<alt_data_len; i++ { fmt.Printf(" %d", alt_data_offset[i]) }
  fmt.Printf("]\n")

  for i:=0; i<len(alt_data); i++ {
    fmt.Printf("AltData[%d]:\n", i)
    for j:=0; j<len(alt_data[i]); j++ {
      if ((j+1)%60)==0 { fmt.Printf("\n") }
      fmt.Printf(" %2x", alt_data[i][j])
    }
    fmt.Printf("\n")
  }
  fmt.Printf("\n")


  fmt.Printf("\n")
  fmt.Printf("\n")

}

type debug_alt_rec_struct struct {
  StartBP uint64
  CanonLenBP uint64
  AltLenBP uint64
  SeqTwoBit []byte
  Seq []byte
}

type debug_alt_overflow_rec_struct struct {
  BodySeqIndex uint64
  NAlt uint64
  Alt []debug_alt_rec_struct
}


/*
// our main structures are
// - ranked_tilelib_map
//   * holds span sequence and diffs
//   * auxflag set to true if we shoudl consult aux_overlfow_lib
// - aux_overlfow_lib
//   * holds span variantid and sequence
//   * AuxBodySeqFlag true if it's an aux canon seq.
//   * AuxBodySeqFlag false means the diffs are stored in the ranked_tilelib_map 
//   * Type 1 is a long sequence meant for gzip'ing
//
func debug_dump() {
  for path := range ranked_tilelib_map {
    for step := range ranked_tilelib_map[path] {

      for idx:=0; idx<len(ranked_tilelib_map[path][step]); idx++ {
        ent := ranked_tilelib_map[path][step][idx]

        if !ent.AuxFlag {
          sa := make([]string, 0, 10)
          for ii:=0; ii<len(ent.Diff); ii++ {
            sa = append(sa, fmt.Sprintf("%x+%x{%s}", ent.Diff[ii].Start, ent.Diff[ii].N, ent.Diff[ii].Seq))
          }
          st := strings.Join(sa, ";")

          if st == "" { st = "." }

          if idx==0 {
            fmt.Printf("%04x.00.%04x.%03x c (+%d) %s %s\n", path, step, idx, ent.Span, st, ent.Seq)
          } else {
            fmt.Printf("%04x.00.%04x.%03x s (+%d) %s %s\n", path,step,idx, ent.Span, st, ent.Seq)
          }

        }

      }

      for idx:=0; idx<len(aux_overflow_lib[path][step]); idx++ {
        ent := aux_overflow_lib[path][step][idx]

        if ent.Type == 0 {

          if ent.AuxBodySeqFlag {
            fmt.Printf("%04x.00.%04x.%03x a+ (+%d) . %s\n", path,step,ent.VariantId, ent.Span, ent.Seq)
          } else {
            idx := ent.VariantId
            ranked_ent := ranked_tilelib_map[path][step][idx]
            sa := make([]string, 0, 10)
            for ii:=0; ii<len(ranked_ent.Diff); ii++ {
              sa = append(sa, fmt.Sprintf("%x+%x{%s}", ranked_ent.Diff[ii].Start, ranked_ent.Diff[ii].N, ranked_ent.Diff[ii].Seq))
            }
            st := strings.Join(sa, ";")

            fmt.Printf("%04x.00.%04x.%03x a (+%d) %d@%s %s\n", path,step,ent.VariantId, ent.Span, ent.AuxVariantId, st, ent.Seq)
          }
        } else if ent.Type == 1 {

          fmt.Printf("%04x.00.%04x.%03x A (+%d) . %s\n", path,step,ent.VariantId, ent.Span, ent.Seq)
        }

      }

    }
  }

}
*/


