package main

import "fmt"
import "os"
import "io"
import "bufio"

import "github.com/curoverse/l7g/go/pasta"


// Read from an interleaved stream and print out a simplified variant difference format
//
// Each token from the stream should be interleaved and aligned.  Each token can be processed
// two at a time, where the first token is from the first stream and the second is from
// the second stream.  The resulting difference format spits out contigs of ref, non-ref and
// alts where appropriate.
//
// The 'process' callback will be called for every variant line that gets processed.
//
func interleave_to_diff_iface(stream *bufio.Reader, p RefVarPrinter, w io.Writer) error {
  alt0 := []byte{}
  alt1 := []byte{}
  refseq := []byte{}

  ref_start := 0
  ref0_len := 0
  ref1_len := 0

  stream0_pos:=0
  stream1_pos:=0

  info := RefVarInfo{}
  //info := GVCFVarInfo{}
  info.Type = pasta.BEG
  info.MessageType = pasta.BEG
  info.RefSeqFlag = gFullRefSeqFlag
  info.NocSeqFlag = gFullNocSeqFlag
  info.Out = os.Stdout
  info.Chrom = "unk"

  out := bufio.NewWriter(w)

  var bp_anchor_ref byte
  var bp_anchor_prv byte

  curStreamState := pasta.BEG ; _ = curStreamState
  prvStreamState := pasta.BEG ; _ = prvStreamState

  var msg pasta.ControlMessage
  var prev_msg pasta.ControlMessage
  var e error

  var ch1 byte
  var e1 error

  var dbp0 int
  var dbp1 int


  for {
    is_ref0 := false
    is_ref1 := false

    is_noc0 := false
    is_noc1 := false

    message_processed_flag := false

    ch0,e0 := stream.ReadByte()
    for (e0==nil) && ((ch0=='\n') || (ch0==' ') || (ch0=='\r') || (ch0=='\t')) {
      ch0,e0 = stream.ReadByte()
    }
    if e0!=nil { break }

    if ch0=='>' {
      msg,e = pasta.ControlMessageProcess(stream)
      if e!=nil { return fmt.Errorf(fmt.Sprintf("invalid control message %v (%v)", msg, e)) }

      if (msg.Type == pasta.REF) || (msg.Type == pasta.NOC) {
        curStreamState = pasta.MSG_REF_NOC
      } else if msg.Type == pasta.CHROM {
        curStreamState = pasta.MSG_CHROM
      } else if msg.Type == pasta.POS {
        curStreamState = pasta.MSG_POS
      } else {
        //just ignore
        continue
        //return fmt.Errorf("invalid message type")
      }

      message_processed_flag = true
    }

    if !message_processed_flag {
      ch1,e1 = stream.ReadByte()
      for (e1==nil) && ((ch1=='\n') || (ch1==' ') || (ch1=='\r') || (ch1=='\t')) {
        ch1,e1 = stream.ReadByte()
      }
      if e1!=nil { break }

      stream0_pos++
      stream1_pos++

      // special case: nop
      //
      if ch0=='.' && ch1=='.' { continue }

      dbp0 = pasta.RefDelBP[ch0]
      dbp1 = pasta.RefDelBP[ch1]

      if ch0=='a' || ch0=='c' || ch0=='g' || ch0=='t' {
        is_ref0 = true
      } else if ch0=='n' || ch0=='N' || ch0 == 'A' || ch0 == 'C' || ch0 == 'G' || ch0 == 'T' {
        is_noc0 = true
      }

      if ch1=='a' || ch1=='c' || ch1=='g' || ch1=='t' {
        is_ref1 = true
      } else if ch1=='n' || ch1=='N' || ch1 == 'A' || ch1 == 'C' || ch1 == 'G' || ch1 == 'T' {
        is_noc1 = true
      }

      if is_ref0 && is_ref1 {
        curStreamState = pasta.REF
      } else if is_noc0 || is_noc1 {
        curStreamState = pasta.NOC
      } else {
        curStreamState = pasta.ALT
      }

    }

    if curStreamState == pasta.BEG {

      if !is_ref0 || !is_ref1 {
        if bp,ok := pasta.RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        } else if bp, ok := pasta.RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        }
      } else if gFullRefSeqFlag {
        if bp,ok := pasta.RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        } else if bp, ok := pasta.RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        }
      }

      ref0_len+=dbp0
      ref1_len+=dbp1

      if bp_val,ok := pasta.AltMap[ch0] ; ok { alt0 = append(alt0, bp_val) }
      if bp_val,ok := pasta.AltMap[ch1] ; ok { alt1 = append(alt1, bp_val) }

      prvStreamState = curStreamState
      prev_msg = msg

      continue
    }

    if !message_processed_flag {
      if is_ref0 && is_ref1 && ch0!=ch1 {
        return fmt.Errorf(fmt.Sprintf("ERROR: stream position (%d,%d), stream0 token %c (%d), stream1 token %c (%d)",
          stream0_pos, stream1_pos, ch0, ch0, ch1, ch1))
      }
    }

    if (prvStreamState == pasta.REF) && (curStreamState != pasta.REF) {

      info.RefBP = bp_anchor_ref
      e:=p.Print(prvStreamState, ref_start, ref0_len, refseq, nil, out)
      if e!=nil { return e }

      // Save the last ref BP in case the ALT is an indel.
      //
      bp_anchor_prv = '-'
      if len(refseq)>0 { bp_anchor_prv = refseq[len(refseq)-1] }

      ref_start += ref0_len

      ref0_len=0
      ref1_len=0

      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if (prvStreamState == pasta.NOC) && (curStreamState != pasta.NOC) {

      full_noc_flag := gFullNocSeqFlag
      for ii:=0; ii<len(alt0); ii++ { if alt0[ii]!='n' { full_noc_flag = true ; break; } }
      if full_noc_flag { for ii:=0; ii<len(alt1); ii++ { if alt1[ii]!='n' { full_noc_flag = true ; break; } } }

      a0 := string(alt0)
      if len(a0) == 0 { a0 = "-" }

      a1 := string(alt1)
      if len(a1) == 0 { a1 = "-" }

      r := string(refseq)
      if len(r) == 0 { r = "-" }



      info.RefBP = bp_anchor_ref
      info.NocSeqFlag = full_noc_flag
      e:=p.Print(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, out)
      if e!=nil { return e }

      // Save the last ref BP in case the ALT is an indel.
      //
      bp_anchor_prv = '-'
      if len(refseq)>0 { bp_anchor_prv = refseq[len(refseq)-1] }

      ref_start += ref0_len

      ref0_len=0
      ref1_len=0

      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if (prvStreamState == pasta.ALT) && ((curStreamState == pasta.REF) || (curStreamState == pasta.NOC)) {

      a0 := string(alt0)
      if len(a0) == 0 { a0 = "-" }

      a1 := string(alt1)
      if len(a1) == 0 { a1 = "-" }

      r := string(refseq)
      if len(r) == 0 { r = "-" }

      info.RefBP = bp_anchor_prv
      e:=p.Print(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, out)
      if e!=nil { return e }

      ref_start += ref0_len

      ref0_len=0
      ref1_len=0

      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if prvStreamState == pasta.MSG_REF_NOC {

      info.Msg = prev_msg
      info.RefBP = bp_anchor_ref
      e:=p.Print(prvStreamState, ref_start, prev_msg.N, refseq, nil, out)
      if e!=nil { return e }

      ref_start += prev_msg.N

      stream0_pos += prev_msg.N
      stream1_pos += prev_msg.N

      ref0_len=0
      ref1_len=0
      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if prvStreamState == pasta.MSG_CHROM {
      info.Chrom = prev_msg.Chrom
      p.Chrom(prev_msg.Chrom)
    } else if prvStreamState == pasta.MSG_POS {
      ref_start = prev_msg.RefPos
    } else {
      // The current state matches the previous state.
      // Either both the current tokens are non-ref as well as the previous tokens
      // or both the current token and previous tokens are ref.
    }

    if !message_processed_flag {
      if bp_val,ok := pasta.AltMap[ch0] ; ok { alt0 = append(alt0, bp_val) }
      if bp_val,ok := pasta.AltMap[ch1] ; ok { alt1 = append(alt1, bp_val) }

      if !is_ref0 || !is_ref1 {

        if bp,ok := pasta.RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := pasta.RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      } else if gFullRefSeqFlag {

        if bp,ok := pasta.RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := pasta.RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      } else if ref0_len==0 {

        if bp,ok := pasta.RefMap[ch0] ; ok {
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := pasta.RefMap[ch1] ; ok {
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      }

      ref0_len+=dbp0
      ref1_len+=dbp1

    }

    prvStreamState = curStreamState
    prev_msg = msg

  }

  if prvStreamState == pasta.REF {

    info.RefBP = bp_anchor_ref
    e:=p.Print(prvStreamState, ref_start, ref0_len, refseq, [][]byte{alt0, alt1}, out)
    if e!=nil { return e }

  } else if prvStreamState == pasta.NOC {

    full_noc_flag := gFullNocSeqFlag
    for ii:=0; ii<len(alt0); ii++ { if alt0[ii]!='n' { full_noc_flag = true ; break; } }
    if full_noc_flag { for ii:=0; ii<len(alt1); ii++ { if alt1[ii]!='n' { full_noc_flag = true ; break; } } }

    info.NocSeqFlag = full_noc_flag
    info.RefBP = bp_anchor_ref
    e:=p.Print(prvStreamState, ref_start, ref0_len, refseq, [][]byte{alt0, alt1}, out)
    if e!=nil { return e }

  } else if prvStreamState == pasta.ALT {

    a0 := string(alt0)
    if len(a0) == 0 { a0 = "-" }

    a1 := string(alt1)
    if len(a1) == 0 { a1 = "-" }

    r := string(refseq)
    if len(r) == 0 { r = "-" }

    e:=p.Print(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, out)
    if e!=nil { return e }

  } else if prvStreamState == pasta.MSG_REF_NOC {

    info.Msg = prev_msg
    info.RefBP = bp_anchor_ref
    e:=p.Print(prvStreamState, ref_start, prev_msg.N, nil, nil, out)
    if e!=nil { return e }

  } else if prvStreamState == pasta.MSG_CHROM {
    info.Chrom = prev_msg.Chrom
    p.Chrom(prev_msg.Chrom)
  }

  p.PrintEnd(out)

  out.Flush()

  return nil
}
