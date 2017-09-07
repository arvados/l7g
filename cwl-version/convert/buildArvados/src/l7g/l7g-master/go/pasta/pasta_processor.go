package pasta


import "fmt"
import "os"
import "io"

import "bufio"


type RefVarInfo struct {
  Type int
  MessageType int
  RefSeqFlag bool
  NocSeqFlag bool
  Out io.Writer
  Msg ControlMessage
  RefBP byte

  Chrom string
}


var g_debug bool
var gFullRefSeqFlag bool = true
var gFullNocSeqFlag bool = true


type RefVarProcesser func(int,int,int,[]byte,[][]byte,interface{}) error

// Read from an interleaved stream and print out a simplified variant difference format
//
// Each token from the stream should be interleaved and aligned.  Each token can be processed
// two at a time, where the first token is from the first stream and the second is from
// the second stream.  The resulting difference format spits out contigs of ref, non-ref and
// alts where appropriate.
//
// The 'process' callback will be called for every variant line that gets processed.
//
func InterleaveToDiff(stream *bufio.Reader, process RefVarProcesser) error {
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
  info.Type = BEG
  info.MessageType = BEG
  info.RefSeqFlag = gFullRefSeqFlag
  info.NocSeqFlag = gFullNocSeqFlag

  info.Out = os.Stdout
  info.Chrom = "Unk"
  //info.PrintHeader = true
  //info.Reference = "hg19"

  var bp_anchor_ref byte
  var bp_anchor_prv byte

  if g_debug { fmt.Printf("%v\n", RefDelBP) }

  curStreamState := BEG ; _ = curStreamState
  prvStreamState := BEG ; _ = prvStreamState

  var msg ControlMessage
  var prev_msg ControlMessage
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
      msg,e = ControlMessageProcess(stream)
      if e!=nil { return fmt.Errorf(fmt.Sprintf("invalid control message %v (%v)", msg, e)) }

      if (msg.Type == REF) || (msg.Type == NOC) {
        curStreamState = MSG_REF_NOC
      } else if msg.Type == CHROM {
        curStreamState = MSG_CHROM
      } else if msg.Type == POS {
        curStreamState = MSG_POS
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

      dbp0 = RefDelBP[ch0]
      dbp1 = RefDelBP[ch1]

      if g_debug {
        fmt.Printf("\n")
        fmt.Printf(">>> ch0 %c (%d), ch1 %c (%d), dbp0 +%d, dbp1 +%d, ref0_len %d, ref1_len %d\n", ch0, ch0, ch1, ch1, dbp0, dbp1, ref0_len, ref1_len)
      }

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
        curStreamState = REF
      } else if is_noc0 || is_noc1 {
        curStreamState = NOC
      } else {
        curStreamState = ALT
      }

    }

    if curStreamState == BEG {

      if !is_ref0 || !is_ref1 {
        if bp,ok := RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        } else if bp, ok := RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        }
      } else if gFullRefSeqFlag {
        if bp,ok := RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        } else if bp, ok := RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          bp_anchor_ref = bp
        }
      }

      ref0_len+=dbp0
      ref1_len+=dbp1

      if bp_val,ok := AltMap[ch0] ; ok { alt0 = append(alt0, bp_val) }
      if bp_val,ok := AltMap[ch1] ; ok { alt1 = append(alt1, bp_val) }

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

    if (prvStreamState == REF) && (curStreamState != REF) {

      info.RefBP = bp_anchor_ref
      process(prvStreamState, ref_start, ref0_len, refseq, nil, &info)

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

    } else if (prvStreamState == NOC) && (curStreamState != NOC) {

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
      process(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, &info)

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

    } else if (prvStreamState == ALT) && ((curStreamState == REF) || (curStreamState == NOC)) {

      a0 := string(alt0)
      if len(a0) == 0 { a0 = "-" }

      a1 := string(alt1)
      if len(a1) == 0 { a1 = "-" }

      r := string(refseq)
      if len(r) == 0 { r = "-" }

      info.RefBP = bp_anchor_prv
      process(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, &info)

      ref_start += ref0_len

      ref0_len=0
      ref1_len=0

      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if prvStreamState == MSG_REF_NOC {

      info.Msg = prev_msg
      info.RefBP = bp_anchor_ref
      process(prvStreamState, ref_start, prev_msg.N, refseq, nil, &info)

      ref_start += prev_msg.N
      ref_start += prev_msg.N

      stream0_pos += prev_msg.N
      stream1_pos += prev_msg.N

      ref0_len=0
      ref1_len=0
      alt0 = alt0[0:0]
      alt1 = alt1[0:0]
      refseq = refseq[0:0]

    } else if prvStreamState == MSG_CHROM {
      info.Chrom = prev_msg.Chrom
    } else if prvStreamState == MSG_POS {
      ref_start = prev_msg.RefPos
    } else {
      // The current state matches the previous state.
      // Either both the current tokens are non-ref as well as the previous tokens
      // or both the current token and previous tokens are ref.
    }

    if !message_processed_flag {
      if bp_val,ok := AltMap[ch0] ; ok { alt0 = append(alt0, bp_val) }
      if bp_val,ok := AltMap[ch1] ; ok { alt1 = append(alt1, bp_val) }

      if !is_ref0 || !is_ref1 {

        if bp,ok := RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      } else if gFullRefSeqFlag {

        if bp,ok := RefMap[ch0] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := RefMap[ch1] ; ok {
          refseq = append(refseq, bp)
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      } else if ref0_len==0 {

        if bp,ok := RefMap[ch0] ; ok {
          if ref0_len==0 { bp_anchor_ref = bp }
        } else if bp, ok := RefMap[ch1] ; ok {
          if ref0_len==0 { bp_anchor_ref = bp }
        }
      }

      ref0_len+=dbp0
      ref1_len+=dbp1

    }

    prvStreamState = curStreamState
    prev_msg = msg

  }

  if prvStreamState == REF {

    info.RefBP = bp_anchor_ref
    process(prvStreamState, ref_start, ref0_len, refseq, [][]byte{alt0, alt1}, &info)

  } else if prvStreamState == NOC {
    full_noc_flag := gFullNocSeqFlag
    for ii:=0; ii<len(alt0); ii++ { if alt0[ii]!='n' { full_noc_flag = true ; break; } }
    if full_noc_flag { for ii:=0; ii<len(alt1); ii++ { if alt1[ii]!='n' { full_noc_flag = true ; break; } } }

    info.NocSeqFlag = full_noc_flag
    info.RefBP = bp_anchor_ref
    process(prvStreamState, ref_start, ref0_len, refseq, [][]byte{alt0, alt1}, &info)

  } else if prvStreamState == ALT {

    a0 := string(alt0)
    if len(a0) == 0 { a0 = "-" }

    a1 := string(alt1)
    if len(a1) == 0 { a1 = "-" }

    r := string(refseq)
    if len(r) == 0 { r = "-" }

    process(prvStreamState, ref_start, ref0_len, []byte(r), [][]byte{[]byte(a0), []byte(a1)}, &info)

  } else if prvStreamState == MSG_REF_NOC {

    info.Msg = prev_msg
    info.RefBP = bp_anchor_ref
    process(prvStreamState, ref_start, prev_msg.N, nil, nil, &info)

  } else if prvStreamState == MSG_CHROM {
    info.Chrom = prev_msg.Chrom
  }

  return nil
}


