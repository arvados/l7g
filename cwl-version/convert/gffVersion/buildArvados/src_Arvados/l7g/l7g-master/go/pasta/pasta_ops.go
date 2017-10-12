package pasta

import "fmt"
import "bufio"
import "strconv"


func ControlMessagePrint(msg *ControlMessage, out *bufio.Writer) {

  if msg.Type == REF {
    out.WriteString(fmt.Sprintf(">R{%d}", msg.N))
  } else if msg.Type == POS {
    out.WriteString(fmt.Sprintf(">P{%d}", msg.RefPos))
  } else if msg.Type == NOC {
    out.WriteString(fmt.Sprintf(">P{%d}", msg.N))
  } else if msg.Type == CHROM {
    out.WriteString(fmt.Sprintf(">C{%s}", msg.Chrom))
  } else if msg.Type == COMMENT {
    out.WriteString(fmt.Sprintf(">#{%s}", msg.Comment))
  }

}

func ControlMessageProcess(stream *bufio.Reader) (ControlMessage, error) {
  var msg ControlMessage

  ch,e := stream.ReadByte()
  msg.NBytes++

  if e!=nil { return msg, e }

  if ch=='R' {
    msg.Type = REF
  } else if ch == 'N' {
    msg.Type = NOC
  } else if ch == 'C' {
    msg.Type = CHROM
  } else if ch == 'P' {
    msg.Type = POS
  } else if ch == '#' {
    msg.Type = COMMENT
  } else {
    return msg, fmt.Errorf("Invalid control character %c", ch)
  }

  ch,e = stream.ReadByte()
  msg.NBytes++
  if e!=nil { return msg, e }
  if ch!='{' { return msg, fmt.Errorf("Invalid control block start (expected '{' got %c)", ch) }

  field_str := make([]byte, 0, 32)

  for (e==nil) && (ch!='}') {
    ch,e = stream.ReadByte()
    msg.NBytes++
    if e!=nil { return msg, e }
    field_str = append(field_str, ch)
  }

  n:=len(field_str)

  if (n==0) || (n==1) {
    msg.N = 0
    return msg, nil
  }

  field_str = field_str[:n-1]

  if msg.Type == REF || msg.Type == NOC || msg.Type == POS {
    _i,err := strconv.Atoi(string(field_str))
    if err!=nil { return msg, err }

    if msg.Type == POS {
      msg.RefPos = int(_i)
    } else {
      msg.N = int(_i)
    }
  } else if msg.Type == CHROM {
    msg.Chrom = string(field_str)
  } else if msg.Type == COMMENT {
    msg.Comment = string(field_str)
  }
  return msg, nil

}

