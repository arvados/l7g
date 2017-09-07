package sloppyjson

/*
 * SloppyJSON:
 *
 * example usage:
 *
 *  s := `{ "keyA" : "valA",
 *          "keyB" : [ "valB.0", "valB.1" ],
 *          "keyC" : { "keyC.0" : "valC.0" },
 *          "keyD" : true,
 *          "keyE" : false,
 *          "keyF" : 3.14159 } `
 *  //...
 *
 *  sj,e := sloppyjson.Loads( s )
 *  if e != nil { panic(e) }
 *
 *  fmt.Printf("  %s %s\n", "keyA", sj.O["keyA"].S )
 *  for i:=0; i<len(sj.O["keyB"].L); i++ { fmt.Printf("  %s [%d] %s\n", "keyB", i, sj.O["keyB"].L[i].S ) }
 *  fmt.Printf("  %s { %s %s }\n", "keyC", "keyC.0", sj.O["keyC"].O["keyC.0"].S )
 *  fmt.Printf("  %s %s\n", "keyD", sj.O["keyD"].Y )
 *  fmt.Printf("  %s %s\n", "keyE", sj.O["keyE"].Y )
 *  fmt.Printf("  %s %f\n", "keyF", sj.O["keyF"].P )
 *
 * Will produce:
 *
 *  keyA valA
 *  keyB [0] valB.0
 *  keyB [1] valB.1
 *  keyC { keyC.0 valC.0 }
 *  keyD true
 *  keyE false
 *  keyF 3.141590
 *
 *
 */


import (
	"fmt"
	"strconv"
	"strings"
)

// S - (S)tring
// L - (L)ist
// O - (O)bject
// P - (P)recision value (float)
// Y - t(Y)pe : (S|L|O|P|true|false|null)
//
type SloppyJSON struct {
  S string
  L []*SloppyJSON
  O map[string]*SloppyJSON
  P float64
  Y string
}

const (
	Whitespaces = " \f\n\r\t\v\u00A0"
)

func ws( indent int ) {
  for i:=0; i<indent; i++ {
    fmt.Printf(" ")
  }
}

func (sjson *SloppyJSON) Printr( indent,dw int ) {
  if sjson.Y == "true" { fmt.Printf("true") ; return }
  if sjson.Y == "false" { fmt.Printf("false") ; return }
  if sjson.Y == "null" { fmt.Printf("null") ; return }

  if sjson.Y == "S" { fmt.Printf("\"%s\"", sjson.S) ; return }
  if sjson.Y == "P" { fmt.Printf("%f", sjson.P) ; return }

  if sjson.Y == "O" {

    fmt.Printf("{\n")

    count:=0
    for k := range sjson.O {
      if count>0 { fmt.Printf(",\n") }
      ws(indent + dw)
      fmt.Printf("\"%s\":", k)
      sjson.O[k].Printr(indent+dw, dw)
      count++
    }
    fmt.Printf("\n")

    ws(indent) ; fmt.Printf("}")

    return
  }

  if sjson.Y == "L" {

    ws(indent) ; fmt.Printf("[\n") ; ws(indent+dw)

    count:=0

    for i:=0; i<len(sjson.L); i++ {
      if count>0 {
        fmt.Printf(",\n")
        ws(indent+dw)
      }
      sjson.L[i].Printr( indent+dw, dw )
      count++
    }
    fmt.Printf("\n")

    ws(indent) ; fmt.Printf("]")

    return
  }

}

func ( sjson *SloppyJSON ) Dump() { sjson.Printr( 0, 2 ) }

func skipspace( dat string, k, n int ) int {
	for ; k<n && strings.IndexByte(Whitespaces, dat[k]) >= 0; k++ {}
	if k == n {
		return -k
	}
	return k
}

func parsefloat( dat string, k, n int ) (*SloppyJSON, int) {

  v := &(SloppyJSON{})
  v.Y = "P"

  pcount := 0
  b:=k
  dig_start := b

  if dat[k] == '-' {
    k = skipspace(dat,k+1,n)
    if k<0 { return nil, k }
    dig_start = k
  }

  // Because JSON can't decode floats like every
  // other sane programming language, we need to put
  // extra checks in.  Here we check for a leading
  // '.' without any digits before it.
  //
  if dat[k] == '.' { return nil, -1 }

  left_digit_count := 0
  right_digit_count := 0
  e_count := 0
  for ; k<n; k++ {
    if dat[k] =='.' {
      pcount++
      if pcount > 1 { return nil, -k }
    } else if ((dat[k]<48) || (dat[k]>57)) &&
              (dat[k]!='e') &&
              (dat[k]!='E') &&
              (dat[k]!='+') &&
              (dat[k]!='-') {
      break
    } else {

      if dat[k] == 'e' || dat[k] == 'E' {

        // Count the number of 'e' or 'E's seen.
        //
        e_count++
        if e_count>1 { return nil, -k }
        continue
      }
      if e_count==0 {

        // And count the number of significant digits
        // seen to the left and right of the decimal point
        //
        if pcount > 0 {
          right_digit_count++
        } else {
          left_digit_count++
        }
      }
    }
  }

  // If we've seen a decimal point and either the left or right
  // significant digit count is zero, that's a parse error.
  //
  if (pcount>0) && ((left_digit_count==0) || (right_digit_count==0)) { return nil, -1 }

  var e error

  v.P,e = strconv.ParseFloat( dat[b:k], 64 )
  if e!=nil { return nil,-k }

  if dat[dig_start] == '.' { return nil, -1 }

  return v, k

}

func parsefalse( dat string, k,n int ) (*SloppyJSON, int) {

  v := &(SloppyJSON{})
  v.Y = "false"

  if dat[k] != 'f' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'a' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'l' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 's' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'e' { return nil,-k }
  k++ ; //if k==n { return nil,-k }

  return v,k
}

func parsetrue( dat string, k,n int ) (*SloppyJSON, int) {
  v := &(SloppyJSON{})
  v.Y = "true"

  if dat[k] != 't' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'r' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'u' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'e' { return nil,-k }
  k++

  return v,k
}

func parsenull( dat string, k,n int ) (*SloppyJSON, int) {
  v := &(SloppyJSON{})
  v.Y = "null"

  if dat[k] != 'n' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'u' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'l' { return nil,-k }
  k++ ; if k==n { return nil,-k }

  if dat[k] != 'l' { return nil,-k }
  k++

  return v,k
}

func parsesimplestring( dat string, k,n int ) (string, int) {
  b := k
  escape := false
  for ; k<n; k++ {
    if escape { escape = false ; continue }
    if dat[k] == '\\' { escape = true ; continue }
    if dat[k] == '"' { break }
  }

  if k==n { return "",-k }

  return dat[b:k],k+1

}


func parsestring( dat string, k,n int ) (*SloppyJSON, int) {
  escape := false

  tbuf := make([]byte, len(dat))
  tpos := 0

  for ; k<n; k++ {

    if escape {

      // unicode
      //
      if (dat[k] == 'u') || (dat[k] == 'U') {
        if k+5>=n { return nil,-1 }

        r,e := strconv.ParseInt(dat[k+1:k+5], 16, 0)
        if e!=nil { return nil, -1 }
        tstr := string(r)
        for ii:=0; ii<len(tstr); ii++ {
          tbuf[tpos] = tstr[ii]
          tpos++
        }

        k+=4

      } else if dat[k] == '"' { tbuf[tpos] = '"' ; tpos++
      } else if dat[k] == '\\' { tbuf[tpos] = '\\' ; tpos++
      } else if dat[k] == '/' { tbuf[tpos] = '/' ; tpos++
      } else if dat[k] == 'b' { tbuf[tpos] = '\b' ; tpos++
      } else if dat[k] == 'f' { tbuf[tpos] = '\f' ; tpos++
      } else if dat[k] == 'n' { tbuf[tpos] = '\n' ; tpos++
      } else if dat[k] == 'r' { tbuf[tpos] = '\r' ; tpos++
      } else if dat[k] == 't' { tbuf[tpos] = '\t' ; tpos++ }

      escape = false
      continue
    }
    if dat[k] == '\\' {
      escape = true
      continue
    }
    if dat[k] == '"' { break }

    tbuf[tpos] = dat[k]
    tpos++
  }

  if k==n { return nil,-k }

  v:=&(SloppyJSON{})
  v.Y = "S"
  v.S = string(tbuf[0:tpos])

  return v,k+1

}

func parselist( dat string, k int, n int ) ( *SloppyJSON, int) {

  li := &(SloppyJSON{})
  li.Y = "L"
  li.L = make( []*SloppyJSON, 0, 8 )

  k = skipspace(dat,k,n)
  if k<0 { return nil, k }
  if dat[k] == ']' { return li, k+1 }

  var v *SloppyJSON

  if dat[k] == '"' { v,k = parsestring(dat,k+1,n)
  } else if ((dat[k]>='0') && (dat[k]<='9')) ||
             (dat[k]=='.') ||
             (dat[k]=='-') {
    v,k = parsefloat(dat,k,n)
  } else if dat[k] == 't' { v,k = parsetrue(dat,k,n)
  } else if dat[k] == 'f' { v,k = parsefalse(dat,k,n)
  } else if dat[k] == 'n' { v,k = parsenull(dat,k,n)
  } else if dat[k] == '{' { v,k = parseobject(dat,k+1,n)
  } else if dat[k] == '[' { v,k = parselist(dat,k+1,n)
  }

  if k<0 { return nil, k }
  li.L = append(li.L, v)

  for k = skipspace(dat,k,n) ; k<n; k = skipspace(dat,k,n) {
    if k<0 { return nil,k }

    if dat[k] == ']' { break }
    if dat[k] != ',' { return nil, -k }
    k = skipspace(dat,k+1,n)
    if k<0 { return nil, k }

    if dat[k] == '"' { v,k = parsestring(dat,k+1,n)
    } else if ((dat[k]>='0') && (dat[k]<='9')) ||
               (dat[k]=='.') ||
               (dat[k]=='-') {
      v,k = parsefloat(dat,k,n)
    } else if dat[k] == 't' { v,k = parsetrue(dat,k,n)
    } else if dat[k] == 'f' { v,k = parsefalse(dat,k,n)
    } else if dat[k] == 'n' { v,k = parsenull(dat,k,n)
    } else if dat[k] == '{' { v,k = parseobject(dat,k+1,n)
    } else if dat[k] == '[' { v,k = parselist(dat,k+1,n)
    } else { return nil, -k }

    if k<0 { return nil, k }
    li.L = append(li.L, v)

  }
  if k==n { return nil, -k }


  return li,k+1
}

func parseobject( dat string, k int, n int ) (*SloppyJSON, int) {
  var v *SloppyJSON
  var str string

  obj := &(SloppyJSON{})
  obj.Y = "O"
  obj.O = make( map[string]*SloppyJSON )

  k = skipspace(dat,k,n)
  if k<0 { return nil, k }

  if dat[k] == '}' { return obj,k+1 }
  if dat[k] != '"' { return nil, -k }

  str,k = parsesimplestring(dat,k+1,n)
  if k<0 { return nil, k }

  k = skipspace(dat,k,n)
  if k<0 { return nil, k }

  if dat[k] != ':' { return nil, -k }
  k = skipspace(dat,k+1,n)
  if k<0 { return nil, k }

  if dat[k] == '"' { v,k = parsestring(dat,k+1,n)
  } else if ((dat[k]>='0') && (dat[k]<='9')) ||
             (dat[k]=='.') ||
             (dat[k]=='-') {
    v,k = parsefloat(dat,k,n)
  } else if dat[k] == 't' { v,k = parsetrue(dat,k,n)
  } else if dat[k] == 'f' { v,k = parsefalse(dat,k,n)
  } else if dat[k] == 'n' { v,k = parsenull(dat,k,n)
  } else if dat[k] == '{' { v,k = parseobject(dat,k+1,n)
  } else if dat[k] == '[' { v,k = parselist(dat,k+1,n)
  } else { return nil,-k }

  if k<0 { return nil,k }
  obj.O[str] = v

  for k = skipspace(dat,k,n) ; k<n; k = skipspace(dat,k,n) {
    if k<0 { return nil,k }

    if dat[k] == '}' { return obj,k+1 }
    if dat[k] != ',' { return nil,-k }
    k = skipspace(dat,k+1,n)
    if k<0 { return nil,k }

    if dat[k] != '"' { return nil,-k }

    str,k = parsesimplestring(dat,k+1,n)
    if k<0 { return nil,k }

    k = skipspace(dat,k,n)
    if k<0 { return nil,k }

    if dat[k] != ':' { return nil,-k }
    k = skipspace(dat,k+1,n)
    if k<0 { return nil,k }

    if dat[k] == '"' { v,k = parsestring(dat,k+1,n)
    } else if ((dat[k]>='0') && (dat[k]<='9')) ||
               (dat[k]=='.') ||
               (dat[k]=='-') {
      v,k = parsefloat(dat,k,n)
    } else if dat[k] == 't' { v,k = parsetrue(dat,k,n)
    } else if dat[k] == 'f' { v,k = parsefalse(dat,k,n)
    } else if dat[k] == 'n' { v,k = parsenull(dat,k,n)
    } else if dat[k] == '{' { v,k = parseobject(dat,k+1,n)
    } else if dat[k] == '[' { v,k = parselist(dat,k+1,n)
    } else { return nil,-k }
    if k<0 { return nil,k }

    obj.O[str] = v
  }
  if k==n { return nil,-k }

  return obj,k

}

func makeError(dat string, k int) error {
	n := len(dat)
	st := k-10 ; if st<0 { st = 0 }
	en := k+10 ; if en>n { en = n }
	z := fmt.Sprintf("%s(*)%s", dat[st:k], dat[k:en] )
	return fmt.Errorf("Parse error at character %d (%#v)", k, z)
}

func Loads( dat string ) (*SloppyJSON,error) {
  var v *SloppyJSON
  k,n := 0,len(dat)

  if n==0 { return nil, fmt.Errorf("empty string") }

  k=skipspace(dat,k,n)
  if k<0 {
	  return nil, makeError(dat, -k)
  }

  if dat[k] == '[' {
    v,k = parselist( dat, k+1, n )
  } else if dat[k] == '{' {
    v,k = parseobject( dat, k+1 , n )
  } else if dat[k] == '"' {
    v,k = parsestring( dat, k+1 , n )
  } else if (dat[k]>='0' && dat[k]<='9') || dat[k]=='-' {
    v,k = parsefloat(dat, k, n)
  } else if dat[k] == 't' {
    v,k = parsetrue(dat, k, n)
  } else if dat[k] == 'f' {
    v,k = parsefalse(dat, k, n)
  } else if dat[k] == 'n' {
    v,k = parsenull(dat, k, n)
  } else {
    return nil, makeError(dat, k)
  }
  if k<0 {
    return nil, makeError(dat, -k)
  }

  for ; k<n; k++ {
    if strings.IndexByte(Whitespaces, dat[k]) == -1 {
      return nil, makeError(dat, k)
    }
  }
  return v, nil
}
