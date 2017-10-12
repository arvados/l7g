package main

/*
  merge two SGLF libraries on a path-by-path basis.
  Take the source and merge into it the second provided
  sglf, adding tile steps as necessary to it's id.
  This takes the md5sum as reported.  Duplicate md5sums
  in the 'added' sglf are effectively ignored.

  Only the tile path, tile step and span are considered
  from the 'add'ed tile library, the tile variant field
  is ignored.

  There is a hard requirement that TileIDs are in %04x.%02x.%04x.%03x
  format.

  The output merged SGLF is not ordered except to keep the input
  SGLF in the same order with the newly added tiles appended to the
  end.  The new tiles appended at the end will be in the same order
  they were read except for the duplicate tiles that already appeared
  in the original SGLF stream.

  example usage:

    ./merge-tilelib <( zcat lib/012c.sglf.gz ) <( zcat lib-cgi-69/012c.sglf.gz ) | bgzip -c > out/012c.sglf.gz

*/

import "io"
import "os"
import "fmt"
import "bufio"

type SGLFLine struct {
  TileID string
  Md5Sum string
  Seq string
}

func ReadSGLFLine(rdr *bufio.Reader) (SGLFLine, error) {
  x := []byte{}
  sglf_line := SGLFLine{}
  count := 0
  var b byte = 0
  var e error
  for true {
    b,e = rdr.ReadByte()
    if b=='\n' {
      if count==2 {
        sglf_line.Seq = string(x)
        return sglf_line,nil
      }
      return sglf_line, fmt.Errorf("empty line")
    }

    if (count==2) && (e==io.EOF) {
      sglf_line.Seq = string(x)
      return sglf_line,nil
    }

    if e!=nil { return sglf_line, e }

    if b==',' {
      if count==0 {
        sglf_line.TileID = string(x)
      } else if count==1 {
        sglf_line.Md5Sum = string(x)
      } else {
        return sglf_line, fmt.Errorf("bad line")
      }

      x = x[0:0]
      count++
    } else {
      x = append(x, b)
    }

  }

  return sglf_line, fmt.Errorf("error")
}

func emit_merge(src, add []SGLFLine) {
  if len(src)==0 { return }

  m := make(map[string]int)

  for i:=0; i<len(src); i++ {
    m[src[i].Md5Sum] = i
    fmt.Printf("%s,%s,%s\n", src[i].TileID, src[i].Md5Sum, src[i].Seq)
  }

  common_tilepos := src[0].TileID[:12]

  extra_count := 0
  varid := len(src)
  for i:=0; i<len(add) ; i++ {
    if _,ok := m[add[i].Md5Sum] ; ok { continue }
    span := add[i].TileID[17:]
    fmt.Printf("%s.%03x+%s,%s,%s\n", common_tilepos, varid, span, add[i].Md5Sum, add[i].Seq)
    varid++
    extra_count++
  }

}

func main() {

  if len(os.Args) != 3 {
    fmt.Printf("provide two streams\n")
    os.Exit(1)
  }

  src_lib := []SGLFLine{}
  add_lib := []SGLFLine{}

  fn0 := os.Args[1]
  fn1 := os.Args[2]

  fp0,e := os.Open(fn0)
  if e!=nil { panic(e) }
  defer fp0.Close()

  fp1,e := os.Open(fn1)
  if e!=nil { panic(e) }
  defer fp1.Close()

  rdr0 := bufio.NewReader(fp0)
  rdr1 := bufio.NewReader(fp1)

  var err error = nil
  for err==nil {
    sglf_line,e := ReadSGLFLine(rdr0)
    if e!=nil { err=e ; break }

    src_lib = append(src_lib, sglf_line)
  }

  err = nil
  for err==nil {
    sglf_line,e := ReadSGLFLine(rdr1)
    if e!=nil { err=e ; break }

    add_lib = append(add_lib, sglf_line)
  }

  src_beg := 0
  src_n := 0 ; _ = src_n
  add_beg := 0 ; _ = add_beg
  add_n := 0 ; _ = add_n

  debug_count := 0 ; _ = debug_count

  for (src_beg < len(src_lib)) && (add_beg < len(add_lib)) {
    src_n = 0
    add_n = 0

    common_pos := src_lib[src_beg].TileID[:12]
    src_n++

    for ((src_beg + src_n) < len(src_lib)) {
      z := src_lib[src_beg+src_n].TileID[:12]
      if common_pos != z { break }
      src_n++
    }

    for ((add_beg + add_n) < len(add_lib)) {
      z := add_lib[add_beg+add_n].TileID[:12]
      if common_pos != z { break }
      add_n++
    }

    if (src_n>0) {
      emit_merge(src_lib[src_beg:src_beg+src_n], add_lib[add_beg:add_beg+add_n])
    }
    src_beg += src_n
    add_beg += add_n
  }

  /*
  fmt.Printf("# src_beg %d / %d, add_beg %d / %d\n",
    src_beg, len(src_lib),
    add_beg, len(add_lib))
    */

  if (src_beg < len(src_lib)) {
    emit_merge(src_lib[src_beg:], add_lib[0:0])
  }

  if (add_beg < len(add_lib)) {
    emit_merge(add_lib[add_beg:add_beg+1], add_lib[add_beg:])
  }


}
