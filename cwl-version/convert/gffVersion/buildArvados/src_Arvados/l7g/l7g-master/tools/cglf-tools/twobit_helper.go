package main

import "fmt"

func twobit2seq_offset(seq, seq2bit []byte, offset uint, bp_len uint64) {
  p:=uint64(0)
  for i:=0; i<len(seq2bit); i++ {
    for ; offset<4; offset++ {
      switch (seq2bit[i] & byte(0x3<<(offset*2))) >> byte(offset*2) {
      case 0: seq[p] = 'a'
      case 1: seq[p] = 'c'
      case 2: seq[p] = 'g'
      case 3: seq[p] = 't'
      }
      p++
      if p>=bp_len { break }
    }
    offset=0
    if p>=bp_len { break }
  }
}

func twobit2seq(seq, seq2bit []byte) {
  twobit2seq_offset(seq, seq2bit, 0, uint64(len(seq2bit)*4))
}

func twobit2seq_sn(seq, seq2bit []byte, s_bp, n uint64) {
  s_idx := s_bp/4
  s_offset := s_bp%4

  twobit2seq_offset(seq, seq2bit[s_idx:], uint(s_offset), n)
}


func seq_to_2bit_offset(obuf, seq []byte, obuf_bp_off uint) error {
  obuf_pos := 0
  for i:=0; i<len(seq); i++ {
    if obuf_pos >= len(obuf) {
      return fmt.Errorf("obuf_pos >= len(obuf)")

      /*
      fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
      fmt.Printf("len(obuf) %d, obuf_pos %d, len(seq) %d\n", len(obuf), obuf_pos, len(seq))
      panic(">>>")
      */
    }
    obuf[obuf_pos] &= byte(^(0x03<<(2*obuf_bp_off)))

    switch seq[i] {
      case 'a', 'A':
      case 'c', 'C': obuf[obuf_pos] |= (0x01<<(2*obuf_bp_off))
      case 'g', 'G': obuf[obuf_pos] |= (0x02<<(2*obuf_bp_off))
      case 't', 'T': obuf[obuf_pos] |= (0x03<<(2*obuf_bp_off))
    }

    obuf_bp_off = (obuf_bp_off+1)%4
    if obuf_bp_off==0 { obuf_pos++ }
  }
  return nil
}

func seq_to_2bit(obuf, seq []byte) {
  seq_to_2bit_offset(obuf, seq, uint(0))
  return
}

