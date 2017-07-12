autoio.go
=========

A lightweight (and work in progress) Go library to provide a simple interface to use compressed files
like you use normal files.

Installation
------------

    $ go get github.com/abeconnelly/autoio

Example
-------

```go
package main

import "fmt"
import "github.com/abeconnelly/autoio"

func parrot_file( fn string ) {
  aio,err := autoio.OpenScanner( fn )
  if err!=nil { panic(err) }
  defer aio.Close()

  for aio.Scanner.Scan()  {
    l := aio.Scanner.Text()
    fmt.Println(l)
  }

}

func main() {

  parrot_file( "w.txt" )
  parrot_file( "x.txt.gz" )
  parrot_file( "y.txt.bz2" )
  parrot_file( "-" )

}
```

