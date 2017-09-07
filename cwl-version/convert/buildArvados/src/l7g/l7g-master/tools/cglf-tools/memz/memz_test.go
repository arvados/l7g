package memz

import "fmt"
import "testing"

//var g_verbose bool = true
var g_verbose bool = false

func Test(t *testing.T) {
  a0 := []byte("tctctcttcctctctttcctctcctctcctctcctctttttttcttttcttctcttttttttttttttgagacagagtttcactctatcgcccaggctggaaagcaatggtgtgatctcagctcactgcaacctctgcctcctgggttcaagcaattctcctgcctcagcctcctgagtagctgggactacaggcatgcactaccatacccagctttttttttttttttttaaaaaaaaaatttttcatagagacagggtttcaccattttggccaggctggtcttgaactcctgacctcaggtgatccgcccaccttggcctcccaaagtgctaagattacaggcataagccactgcgcccagcctggtccccctatttcatttgctcaacagaaacatacaatttgtgagcacccatcacatgtgagaggggcttggacaaacaaggtggaccatcatggtccttgtgagagctcataacgaggaagggaagagggaagaggatgccaattgatgtgtacagggtcctctggagctgacaaatggccttgacaaatactatctccctccatccccgcacccgtt")
  b0 := []byte("tctctcttcctctctttcctctcctctcctctcctctttttttcttttcttctcatttttttttttgagacagagtttcactctatcgcccaggctggaaagcaatggtgtgatctcagctcactgcaacctctgcctcctgggttcaagcaattctcctgcctcagcctcctgagtagctgggactacaggcatgcactaccatacccagctttttttttttttttttaaaaaaaaaatttttcatagagacagggtttcaccattttggccaggctggtcttgaactcctgacctcaggtgatccgcccaccttggcctcccaaagtgctaagattacaggcataagccactgcgcccagcctggtccccctatttcatttgctcaacagaaacatacaatttgtgagcacccatcacatgtgagaggggcttggacaaacaaggtggaccatcatggtccttgtgagagctcataacgaggaagggaagagggaagaggatgccaattgatgtgtacagggtcctctggagctgacaaatggccttgacaaatactatctccctccatccccgcacccgtt")

  a1 := []byte("gggcgggcgggcgggggcagagagtgaaaccgcccccccgccccgcacaaacaagcaccgccgtctgcagcccgaacccgcacccaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaggcactgggaaaatggcggccgag")
  b1 := []byte("gggcgggcgggcgggggcagagagtgaaaccacccccccgccccgcacaaacaagcaccgccgtctgcagcccgaacccgcacccaggccgccacccccggccgcctctttccagccggggagcgagcttcagccgccccccaaaaaaacgaaaagagacgaaaaggcttcgctacggtgctcggttctcccgccggctcggcgagcggtggcggcggtggcggcggcggcggcggcactgggaaaatggcggccgag")

  a2 := []byte("gcatgcat")
  b2 := []byte("catccat")

  a3 := []byte("aaaaaaaaaaaaaaaaaaaaaaaabbbb")
  b3 := []byte("bbbb")

  a4 := []byte("bbbbaaaaaaaaaaaaaaaaaaaaaaaa")
  b4 := []byte("bbbb")

  a5 := []byte("bbbb")
  b5 := []byte("aaaaaaaaaaaaaaaaaaaaaaaabbbb")

  a6 := []byte("bbbb")
  b6 := []byte("bbbbaaaaaaaaaaaaaaaaaaaaaaaa")

  a7 := []byte("bbbb")
  b7 := []byte("aaaaaaaaaaaaaabbbbaaaaaaaaaaaaaaaaaaaaaaaa")

  a8 := []byte("fooquxbar")
  b8 := []byte("fooxzqbar")

  a9 := []byte("foobar")
  b9 := []byte("fooquxbar")

  a10 := []byte("bar")
  b10 := []byte("foobar")

  a11 := []byte("foobar")
  b11 := []byte("bar")

  a12 := []byte("foo")
  b12 := []byte("foobar")

  a13 := []byte("foobar")
  b13 := []byte("foo")

  a14 := []byte("foobarfoobar")
  b14 := []byte("foofoobar")

  a15 := []byte("foofoobar")
  b15 := []byte("foobarfoobar")

  A := make([][]byte, 0, 8)
  B := make([][]byte, 0, 8)

  A = append(A, []byte("foobarfoobazfoo"))
  B = append(B, []byte("foobafooqazfoo"))

  A = append(A, []byte("foobafooqazfoo"))
  B = append(B, []byte("foobarfoobazfoo"))

  A = append(A, []byte("foobarfoobazfoo"))
  B = append(B, []byte("foobafooqqazfoo"))

  A = append(A, []byte("foobafooqqazfoo"))
  B = append(B, []byte("foobarfoobazfoo"))

  for i:=0; i<len(A); i++ {
    __tt(A[i], B[i], t)
  }

  __tt(a0, b0, t)
  __tt(a1, b1, t)
  __tt(a2, b2, t)
  __tt(a3, b3, t)
  __tt(a4, b4, t)
  __tt(a5, b5, t)
  __tt(a6, b6, t)
  __tt(a7, b7, t)

  __tt(a8,b8,t)
  __tt(a9,b9,t)
  __tt(a10,b10,t)
  __tt(a11,b11,t)
  __tt(a12,b12,t)
  __tt(a13,b13,t)
  __tt(a14,b14,t)
  __tt(a15,b15,t)

}


func __tt(a,b []byte, T *testing.T) {

  mz := New()


  if g_verbose {
    fmt.Printf("\nCHECKING:\n%s\n%s\n\n", a, b)
  }


  sc := mz.Score(a,b)

  x,y := mz.Align(a,b)

  if g_verbose {
    fmt.Printf("%s\n%s\n", x, y)
  }

  check_sc := 0
  if len(x)!=len(y) {
    T.Errorf("aligned sequences don't have same length (%d!=%d):\n%s\n%s\n", len(x), len(y), x, y)
  }
  for i:=0; i<len(x); i++ {
    if x[i] == y[i] { check_sc += mz.ScoreMatrix[x[i]][y[i]] }
    if x[i] == '-' { check_sc += mz.GapA }
    if y[i] == '-' { check_sc += mz.GapB }
    if x[i] != '-' && y[i] != '-' && x[i] != y[i] { check_sc += mz.ScoreMatrix[x[i]][y[i]] }
  }

  if check_sc != sc {
    T.Errorf("scores don't match (%d!=%d)\n", sc, check_sc)
  }

  xx := make([]byte, 0, len(x))
  yy := make([]byte, 0, len(y))
  for i:=0; i<len(x); i++ {
    if x[i]!='-' { xx = append(xx, x[i]) }
  }

  for i:=0; i<len(y); i++ {
    if y[i]!='-' { yy = append(yy, y[i]) }
  }

  if string(xx) != string(a) {
    T.Errorf("ungapped aligned A string doesn't match:\n%s\n%s\n%s\n", x, xx, a)
  }
  if string(yy) != string(b) {
    T.Errorf("ungapped aligned B string doesn't match:\n%s\n%s\n%s\n", y, yy, b)
  }


  d := mz.AlignDelta(a,b)

  xa := make([]byte, 0, 8)
  xb := make([]byte, 0, 8)
  preva:=0
  prevb:=0
  posa:=-1
  posb:=-1
  for i:=0; i<len(d); i++ {

    if g_verbose {
      fmt.Printf("[%d/%d].0 xa: %s\n[%d/%d].0 xb: %s\n\n", i, len(d)-1, xa, i, len(d)-1, xb)
    }

    t := d[i].Type
    posa = d[i].PosA
    posb = d[i].PosB
    n := d[i].Len

    st := "none"
    if t==DIFF {
      st = "DIFF"
    } else if t==GAPA {
      st = "GAPA"
    } else if t==GAPB {
      st = "GAPB"
    } else {
      st = "SEQ"
    }

    if g_verbose {
      fmt.Printf("[%d/%d] Type %d (%s), PosA %d, PosB %d, N %d\n", i, len(d)-1, d[i].Type, st, d[i].PosA, d[i].PosB, d[i].Len)
    }

    if posa>preva { xa = append(xa, a[preva:posa]...) }
    if posb>prevb { xb = append(xb, b[prevb:posb]...) }

    preva = posa
    prevb = posb

    if g_verbose {
      fmt.Printf("[%d/%d].1 xa: %s\n[%d/%d].1 xb: %s\n\n", i, len(d)-1, xa, i, len(d)-1, xb)
    }

    if t == DIFF {
      xa = append(xa, a[posa:posa+n]...)
      xb = append(xb, b[posb:posb+n]...)
      preva += n
      prevb += n
    } else if t == GAPA {
      //for ii:=0; ii<n; ii++ { xb = append(xb, '-') }
      for ii:=0; ii<n; ii++ { xa = append(xa, '-') }
      xb = append(xb, b[posb:posb+n]...)
      prevb+=n
    } else if t == GAPB {
      //for ii:=0; ii<n; ii++ { xa = append(xa, '-') }
      xa = append(xa, a[posa:posa+n]...)
      for ii:=0; ii<n; ii++ { xb = append(xb, '-') }
      preva+=n

    }

    if g_verbose {
      fmt.Printf("[%d/%d].2 xa: %s (%s|%s)\n[%d/%d].2 xb: %s (%s|%s)\n\n",
        i, len(d)-1, xa, x[:len(xa)], x[len(xa):],
        i, len(d)-1, xb, y[:len(xb)], y[len(xb):])
    }

    if string(x[:len(xa)]) != string(xa) {
      T.Errorf("Seq A strings don't match while constructing gapped sequence:\n%s\n%s\n", x, xa)
      //panic("(1)")
    }
    if string(y[:len(xb)]) != string(xb) {
      T.Errorf("Seq B strings don't match while constructing gapped sequence:\n%s\n%s\n", y, xb)
      //panic("(2)")
    }

  }

  x_bp_len := 0
  y_bp_len := 0
  for ii:=0; ii<len(xa); ii++ {
    if xa[ii]!='-' { x_bp_len++ }
  }
  for ii:=0; ii<len(xb); ii++ {
    if xb[ii]!='-' { y_bp_len++ }
  }

  if g_verbose {
    fmt.Printf("x_bp_len %d, preva %d\n", x_bp_len, preva)
    fmt.Printf("y_bp_len %d, prevb %d\n", y_bp_len, prevb)
  }

  if len(a)>x_bp_len { xa = append(xa, a[preva:]...) }
  if len(b)>y_bp_len { xb = append(xb, b[prevb:]...) }

  if g_verbose {
    fmt.Printf("xa: %s\nxb: %s\n", xa, xb)

    fmt.Printf("xa: %s\n x: %s\n", xa, x)
    fmt.Printf("xb: %s\n y: %s\n", xb, y)
  }

  if string(xa) != string(x) {
    T.Errorf("Final A sequences don't match:\n%s\n%s\n", xa, x)
    //panic("(3)")
  }
  if string(xb) != string(y) {
    T.Errorf("Final B sequences don't match:\n%s\n%s\n", xb, y)
    //panic("(4)")
  }

}
