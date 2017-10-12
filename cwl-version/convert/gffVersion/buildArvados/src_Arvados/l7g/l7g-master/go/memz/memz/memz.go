package main

import "os"
import "io/ioutil"
import "fmt"
import "github.com/abeconnelly/memz"

var g_normalize_flag bool

func main() {

  g_normalize_flag = true

  if len(os.Args)<3 {
    fmt.Fprintf(os.Stderr, "provide filenames\n")
    os.Exit(1)
  }

  fn0 := os.Args[1]
  fn1 := os.Args[2]

  s0,e := ioutil.ReadFile(fn0)
  if e!=nil { panic(e) }
  s1,e := ioutil.ReadFile(fn1)
  if e!=nil { panic(e) }

  for i:=0; i<256; i++ {
    memz.Score['n'][i] = 0
    memz.Score[i]['n'] = 0
  }

  X,Y,sc := memz.Hirschberg(s0,s1)

  if g_normalize_flag {
    memz.SeqPairNormalize(X,Y)
  }

  //memz.Simp_b(s0,s1,0,0)

  fmt.Printf("%d\n%s\n%s\n", sc, X, Y)
}
