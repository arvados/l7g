package pasta

import "os"
import "io"
import "bufio"

import "errors"

type PastaHandle struct {
  Fp *os.File
  Scanner *bufio.Scanner

  Buf []byte
  Stage []byte
}

func Open(fn string) (p PastaHandle, err error) {
  if fn == "-" {
    p.Fp = os.Stdin
  } else {
    p.Fp,err = os.Open(fn)
  }
  if err!=nil { return }

  p.Reader = bufio.NewReader(p.Fp)
  return p, nil
}

func (p *PastaHandle) Close() {
  p.Fp.Close()
}

func (p *PastaHandle) PeekChar() (byte) {
  if len(p.Stage)==0 return 0
  return p.Stage[0]
}

PASTA_SAUCE := 1024

func (p *PastaHandle) ReadChar() (byte, err) {
  if len(p.Stage)==0 {
    if len(p.Buf)==0 {
      p.Buf = make([]byte, PASTA_SAUCE, PASTA_SAUCE)
    }
    n,e := p.Fp.Read(p.Buf)
    if e!=nil { return 0, e }
    if n==0 { return 0, nil }
    p.Stage = p.Buf[0:n]
  }

  b := p.Stage[0]
  p.Stage = p.Stage[1:]
  return b,nil
}
