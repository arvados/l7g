package simplestream

import "os"
import "io"

import _ "fmt"

type SimpleStream struct {
  Fp *os.File
  Buf []byte
  Pos int
  N int
  Size int
  IsEOF bool
}

func (s *SimpleStream) Init(fp *os.File) (error) {
  var e error
  s.Fp = fp
  s.Size = 1024*1024

  s.Buf = make([]byte, s.Size)
  s.N,e = s.Fp.Read(s.Buf)
  if s.N==0 && e == io.EOF {
    s.IsEOF = true
    return e
  }
  if e!=nil {return e}

  s.Pos=0
  return nil
}

func (s *SimpleStream) Refresh() error {
  var e error
  s.N,e = s.Fp.Read(s.Buf)
  if s.N==0 && e==io.EOF {
    s.IsEOF=true
    return e
  }
  if e!=nil { panic(e) }
  s.Pos=0

  return nil
}

func (s *SimpleStream) Getc() (byte,error) {
  var e error
  if s.Pos>=s.N {
    if e=s.Refresh()
    e!=nil { return 0, e }
  }
  ch := s.Buf[s.Pos]
  s.Pos++
  return ch, nil
}
