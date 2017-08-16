package main

import "io"
import "bufio"
import "fmt"
import "strings"

import "github.com/abeconnelly/pasta"

type FASTAInfo struct {

  Allele int

  RefPos int
  RefBuild string
  ChromStr string
  OCounter int
  LFMod int

  Name string

  Out *bufio.Writer
}

func (g *FASTAInfo) Init() {
  g.LFMod = 50
  g.OCounter = 0
  g.ChromStr = "Unk"
  g.RefBuild = "unk"
  g.RefPos = 0
  g.Allele=0
  g.Name = ""
}

func (g *FASTAInfo) WriteFASTAByte(ch byte, out *bufio.Writer) error {
  out.WriteByte(ch)
  g.OCounter++
  if (g.LFMod>0) && (g.OCounter > 0) && ((g.OCounter%g.LFMod)==0) {
    e := out.WriteByte('\n')
    if e!=nil { return e }
  }
  return nil
}

func (g *FASTAInfo) Chrom(chr string) {
  g.ChromStr = chr
}

func (g *FASTAInfo) Pos(pos int) {
  g.RefPos = pos
}

func (g *FASTAInfo) Header(out *bufio.Writer) error {
  out.WriteString(">" + g.Name + "\n")
  return nil
}

func (g *FASTAInfo) Print(vartype int, ref_start, ref_len int, refseq []byte, altseq [][]byte, out *bufio.Writer) error {

  for ii:=0; ii<len(altseq[g.Allele]); ii++ {
    e := g.WriteFASTAByte(altseq[g.Allele][ii], out)
    if e!=nil { return e }
  }

  return nil
}

func (g *FASTAInfo) Stream(pasta_stream *bufio.Reader, out *bufio.Writer) error {
  var ch byte
  var e error
  var msg pasta.ControlMessage
  curStreamState := pasta.BEG ; _ = curStreamState

  for {

    ch,e = pasta_stream.ReadByte()
    for (e==nil) && ((ch=='\n') || (ch==' ') || (ch=='\r') || (ch=='\t')) {
      ch,e = pasta_stream.ReadByte()
    }
    if e!=nil { break }


    if ch=='>' {
      msg,e = pasta.ControlMessageProcess(pasta_stream)
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
    }

    if g.Allele==0 {
      alt_ch,ok := pasta.AltMap[ch]
      if ok {
        g.WriteFASTAByte(_tolch(alt_ch), out)
      }
    } else {
      ref_ch,ok := pasta.RefMap[ch]
      if ok {
        g.WriteFASTAByte(_tolch(ref_ch), out)
      }
    }

  }


  out.Flush()

  if e!=io.EOF { return e }
  return nil
}


func (g *FASTAInfo) PrintEnd(out *bufio.Writer) error {

  if (g.LFMod>0) && ((g.OCounter%g.LFMod)!=0) {
    out.WriteByte('\n')
  }

  out.Flush()
  return nil
}

func (g *FASTAInfo) PastaBegin(out *bufio.Writer) error {
  out.WriteString(fmt.Sprintf(">C{%s}>P{%d}\n", g.ChromStr, g.RefPos))
  return nil
}

func (g *FASTAInfo) PastaRaw(fasta_line string, out *bufio.Writer) error {
  line := strings.Trim(fasta_line, " \n\t\r")
  if line[0]=='>' { return nil }

  for ii:=0; ii<len(fasta_line); ii++ {
    e := g.WriteFASTAByte(_tolch(fasta_line[ii]), out)
    if e!=nil { return e }
  }

  return nil

}

func (g *FASTAInfo) ReadRefByte(ref_stream *bufio.Reader) (byte,error) {
  b,e := ref_stream.ReadByte()
  if e!=nil {
    return b,fmt.Errorf(fmt.Sprintf("ref_stream error: %v", e))
  }
  for b == '\n' || b == ' ' || b == '\t' || b == '\r' {
    b,e = ref_stream.ReadByte()
    if e!=nil {
      return b,fmt.Errorf(fmt.Sprintf("ref_stream error: %v", e))
    }
  }
  return b, nil
}


func (g *FASTAInfo) Pasta(fasta_line string, ref_stream *bufio.Reader, out *bufio.Writer) error {

  if ref_stream==nil { return g.PastaRaw(fasta_line, out) }

  line := strings.Trim(fasta_line, " \n\t\r")
  if line[0]=='>' { return nil }

  for ii:=0; ii<len(fasta_line); ii++ {
    ch := _tolch(fasta_line[ii])

    ref_ch,e := g.ReadRefByte(ref_stream)
    if e!=nil { return e }

    ref_ch = _tolch(ref_ch)

    pasta_ch,ok := pasta.SubMap[ref_ch][ch]
    if !ok {
      return fmt.Errorf(fmt.Sprintf("FASTA Pasta conversion, bad mapping from '%c'->'%c' (%d->%d)", ref_ch, ch, ref_ch, ch))
    }


    e = g.WriteFASTAByte(pasta_ch, out)
    if e!=nil { return e }
  }

  return  nil
}

func (g *FASTAInfo) PastaEnd(out *bufio.Writer) error {

  if (g.LFMod>0) && ((g.OCounter%g.LFMod)!=0) {
    out.WriteByte('\n')
  }


  out.Flush()
  return nil
}
