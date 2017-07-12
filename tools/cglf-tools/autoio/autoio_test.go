package autoio

import "fmt"
import "testing"

var g_verbose bool = false

func TestOpenLongTextLineFile( t *testing.T ) {
  fn := "./magic/longline.txt"
  first_n_char := 98
  n_char := 100000-1
  n_line:=0

  if g_verbose { fmt.Println(fn) }

  h,err := OpenReadScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.ReadScan() {
    l := h.ReadText()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }


    nc := n_char
    if n_line == 0 { nc = first_n_char }
    if n_line == 101 { nc = 0 }

    if len(l) != nc {
      t.Errorf("did not read %v characters from %s (got %v)", nc, fn, len(l) )
    }

    n_line++
  }

  h.Close()
}

func TestOpenLongTextLineGzFile( t *testing.T ) {
  fn := "./magic/longline.txt.gz"
  first_n_char := 98
  n_char := 100000-1
  n_line:=0

  if g_verbose { fmt.Println(fn) }

  h,err := OpenReadScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.ReadScan() {
    l := h.ReadText()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }


    nc := n_char
    if n_line == 0 { nc = first_n_char }
    if n_line == 101 { nc = 0 }

    if len(l) != nc {
      t.Errorf("did not read %v characters from %s (got %v)", nc, fn, len(l) )
    }

    n_line++
  }

  h.Close()
}

func TestOpenLongTextLineBzip2File( t *testing.T ) {
  fn := "./magic/longline.txt.bz2"
  first_n_char := 98
  n_char := 100000-1
  n_line:=0

  if g_verbose { fmt.Println(fn) }

  h,err := OpenReadScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.ReadScan() {
    l := h.ReadText()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }


    nc := n_char
    if n_line == 0 { nc = first_n_char }
    if n_line == 101 { nc = 0 }

    if len(l) != nc {
      t.Errorf("did not read %v characters from %s (got %v)", nc, fn, len(l) )
    }

    n_line++
  }

  h.Close()
}


func TestOpenTextFile( t *testing.T ) {
  fn := "./magic/test.txt"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}

func TestOpenGzFile( t *testing.T ) {
  fn := "./magic/test.txt.gz"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}

func TestOpenBzip2File( t *testing.T ) {
  fn := "./magic/test.txt.bz2"

  if g_verbose { fmt.Println(fn) }

  h,err := OpenScanner( fn )
  if err != nil { t.Errorf("got error %s", err ) }

  for h.Scanner.Scan() {
    l := h.Scanner.Text()

    if g_verbose { fmt.Printf("  >>> %s\n", l) }

    if l != "test" {
      t.Errorf("did not read test from %s", fn )
    }
  }

  h.Close()
}


