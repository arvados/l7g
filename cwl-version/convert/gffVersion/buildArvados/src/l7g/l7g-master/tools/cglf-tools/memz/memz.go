package memz

import "fmt"
import "math"
import "hash/crc32"
import "../rollsum"

type Memz struct {
  Cache []int
  ScoreMatrix [256][256]int
  Gap int
  GapA int
  GapB int

  N,M int
}

// DIFF is a contig run of mismatched bps
// GAPA is a gap in sequence a (e.g. a(--cc) b(gtcc)
// GAPB is a gap in sequence b (e.g. a(gcat) b(g--t)
// SEQ is sequence (mostly unused and used for a flag)

const (
  DIFF = iota
  GAPA = iota
  GAPB = iota
  SEQ  = iota
)


// examples:
// DIFF : a(tgcc) b(gtcc) -> { PosA:0, PosB:0, Type: DIFF, Len:2 }
// GAPA : a(--cc) b(gtcc) -> { PosA:0, PosB:0, Type: GAPA, Len:2 }
// GAPB : a(tgcc) b(--cc) -> { PosA:0, PosB:0, Type: GAPB, Len:2 }
// SEQ  : a(tgcc) b(tgcc) -> { PosA:0, PosB:0, Type: SEQ, Len:4 }

type Diff struct {
  PosA int
  PosB int
  Type int
  Len int
}


func New() (*Memz) {
  x := Memz{}
  x.Init()
  return &x
}

func (x *Memz) Init() {
  x.Cache = make([]int, 1024*1024)

  x.Gap = -2
  x.GapA = -1
  x.GapB = -1

  for i:=0; i<256; i++ {
    for j:=0; j<256; j++ {
      x.ScoreMatrix[i][j] = -1
      if i==j { x.ScoreMatrix[i][j] = 0 }
    }
  }

}

func (x *Memz) DebugPrintCache() {
  for r:=0; r<x.N; r++ {
    for c:=0; c<x.M; c++ {
      fmt.Printf(" %2d", x.Cache[r*x.M + c])
    }
    fmt.Printf("\n")
  }
  fmt.Printf("\n")

}

func (x *Memz) DebugPrint() {
  fmt.Printf("Score:\n")

  for i:=0; i<256; i++ {
    for j:=0; j<256; j++ {
      fmt.Printf("%2d", x.ScoreMatrix[i][j])
    }
    fmt.Printf("\n")
  }

  fmt.Printf("\n")

  x.DebugPrintCache()
}

func (x *Memz) SetScore(score_matrix map[byte]map[byte]int) {
  for from := range score_matrix {
    for to := range score_matrix[from] {
      x.ScoreMatrix[from][to] = score_matrix[from][to]
    }
  }
}


func (x *Memz) Score(a,b []byte) int {

  // a on rows
  //
  n := len(a)+1

  // b on cols
  //
  m := len(b)+1

  if len(x.Cache) < n*m { x.Cache = make([]int, n*m) }

  x.N = n
  x.M = m

  for c:=1; c<m; c++ { x.Cache[c] = -c }
  for r:=1; r<n; r++ { x.Cache[r*m] = -r }

  for r:=1; r<n; r++ {
    for c:=1; c<m; c++ {
      r0c0 := (r-1)*m + (c-1)
      r0c1 := r0c0+1
      r1c0 := r0c0+m
      r1c1 := r0c0+m+1

      m := x.ScoreMatrix[a[r-1]][b[c-1]] + x.Cache[r0c0]
      if x.Cache[r0c1] + x.GapA > m { m = x.Cache[r0c1] + x.GapA }
      if x.Cache[r1c0] + x.GapB > m { m = x.Cache[r1c0] + x.GapB }
      x.Cache[r1c1] = m

    }
  }

  return x.Cache[n*m-1]

}

func (x *Memz) Align(a,b []byte) ([]byte, []byte) {
  x.Score(a,b)

  // a on rows
  //
  n:=len(a)+1

  // b on cols
  m:=len(b)+1

  u := make([]byte, 0, n)
  v := make([]byte, 0, m)

  r := n-1
  c := m-1
  for (r>0) && (c>0) {

    r0c0 := (r-1)*m + (c-1)
    r0c1 := r0c0+1
    r1c0 := r0c0+m
    r1c1 := r0c0+m+1

    sc := x.ScoreMatrix[a[r-1]][b[c-1]]
    s0 := x.Cache[r0c0]
    s1 := x.Cache[r0c1]
    s2 := x.Cache[r1c0]

    if (s0 >= s1) && (s0 >= s2) && (sc + x.Cache[r0c0] == x.Cache[r1c1]) {

      u = append(u, a[r-1])
      v = append(v, b[c-1])
      r--
      c--
      continue
    }

    if (s1 >= s0) && (s1 >= s2) && (x.GapA + x.Cache[r0c1] == x.Cache[r1c1]) {

      u = append(u, a[r-1])
      v = append(v, '-')
      r--
      continue
    }

    if (s2 >= s0) && (s2 >= s1) && (x.GapB + x.Cache[r1c0] == x.Cache[r1c1]) {

      u = append(u, '-')
      v = append(v, b[c-1])
      c--
      continue
    }

    if (sc + x.Cache[r0c0] == x.Cache[r1c1]) {
      u = append(u, a[r-1])
      v = append(v, b[c-1])
      r--
      c--
      continue
    }

    if (x.GapA + x.Cache[r0c1] == x.Cache[r1c1]) {
      u = append(u, a[r-1])
      v = append(v, '-')
      r--
      continue
    }

    if (x.GapB + x.Cache[r1c0] == x.Cache[r1c1]) {
      u = append(u, '-')
      v = append(v, b[c-1])
      c--
      continue
    }

    r=-1
    c=-1
    break
  }

  for ; r>0; r-- {
    u = append(u, a[r-1])
    v = append(v, '-')
  }
  for ; c>0; c-- {
    u = append(u, '-')
    v = append(v, b[c-1])
  }

  N := len(u)
  M := len(v)

  for i:=0; i<N/2; i++ { u[i],u[N-i-1] = u[N-i-1],u[i] }
  for i:=0; i<M/2; i++ { v[i],v[M-i-1] = v[M-i-1],v[i] }

  return u, v
}

func (x *Memz) AlignDelta(a,b []byte) ([]Diff) {
  delta := make([]Diff, 0, 16)
  x.Score(a,b)

  curdelta := Diff{PosA:0,PosB:0,Len:0,Type:SEQ}

  // seq a on rows
  //
  n:=len(a)+1

  // seq b on cols
  //
  m:=len(b)+1

  r := n-1
  c := m-1
  for (r>0) && (c>0) {

    r0c0 := (r-1)*m + (c-1)
    r0c1 := r0c0+1
    r1c0 := r0c0+m
    r1c1 := r0c0+m+1

    //  r0c0 (s0)     r0c1 (s1)
    //
    //  r1c0 (s2)     r1c1
    //
    //
    //  -> -B  (gapa)
    //
    //  |
    //  v  A-  (gapb)
    //
    //  \
    //   4 AB  (diff/seq)
    //
    sc := x.ScoreMatrix[a[r-1]][b[c-1]]
    s0 := x.Cache[r0c0]
    s1 := x.Cache[r0c1]
    s2 := x.Cache[r1c0]

    if (s0 >= s1) && (s0 >= s2) && (sc + x.Cache[r0c0] == x.Cache[r1c1]) {
      if a[r-1]!=b[c-1] {
        if curdelta.Type != DIFF {
          if curdelta.Type != SEQ { delta = append(delta, curdelta) }

          curdelta.Type = DIFF
          curdelta.Len = 0
        }
        curdelta.PosA = r-1
        curdelta.PosB = c-1
        curdelta.Len++
      } else {
        if curdelta.Type != SEQ { delta = append(delta, curdelta) }
        curdelta.Type = SEQ
        curdelta.PosA = r-1
        curdelta.PosB = c-1
        curdelta.Len = 0
      }

      r--
      c--
      continue
    }

    // r0c1 is max, crossing vertically which
    // gives a gap in B
    //
    if (s1 >= s0) && (s1 >= s2) && (x.GapB + x.Cache[r0c1] == x.Cache[r1c1]) {
      if curdelta.Type != GAPB {
        if curdelta.Type != SEQ { delta = append(delta, curdelta) }
        //curdelta.Type = GAPA
        curdelta.Type = GAPB
        curdelta.Len = 0
      }
      curdelta.PosA = r-1
      curdelta.PosB = c
      curdelta.Len++
      r--
      continue
    }

    // r1c0 is max, crossing horizontally , which gives us a gap in A
    //
    if (s2 >= s0) && (s2 >= s1) && (x.GapA + x.Cache[r1c0] == x.Cache[r1c1]) {
      //if curdelta.Type != GAPB {
      if curdelta.Type != GAPA {
        if curdelta.Type != SEQ { delta = append(delta, curdelta) }
        //curdelta.Type = GAPB
        curdelta.Type = GAPA
        curdelta.Len = 0
      }
      curdelta.PosA = r
      curdelta.PosB = c-1
      curdelta.Len++
      c--
      continue
    }

    r=-1
    c=-1
    break
  }

  for ; r>0; r-- {
    //if curdelta.Type != GAPA {
    if curdelta.Type != GAPB {
      if curdelta.Type != SEQ { delta = append(delta, curdelta) }
      //curdelta.Type = GAPA
      curdelta.Type = GAPB
      curdelta.Len = 0
    }
    curdelta.PosA = r-1
    curdelta.PosB = c
    curdelta.Len++
  }

  for ; c>0; c-- {
    //if curdelta.Type != GAPB {
    if curdelta.Type != GAPA {
      if curdelta.Type != SEQ { delta = append(delta, curdelta) }
      //curdelta.Type = GAPB
      curdelta.Type = GAPA
      curdelta.Len = 0
    }
    curdelta.PosA = r
    curdelta.PosB = c-1
    curdelta.Len++
  }

  if curdelta.Type != SEQ { delta = append(delta, curdelta) }

  // We did the construction from the bottom up, so reverse to
  // put it in ascending order
  //
  N := len(delta)
  for i:=0; i<N/2; i++ {
    delta[i],delta[N-i-1] = delta[N-i-1], delta[i]
  }

  return delta
}




func min(n,m int) int {
  if n<m { return n }
  return m
}

func max(n,m int) int {
  if n<m { return m }
  return n
}

func make_mask(u uint32) uint32 {
  b := uint32(1)
  mask := uint32(1)

  for ; b!=0 ; b = b<<1 {
    if b >= u { break }
    mask = (mask<<1) | 0x1;
  }

  return mask
}

// This is a heuristic alignment.  This does not do
// global alignemtn but does a hybrid.
//
// For two given sequences, split it into 'blocks',
// where the blocks are of variable size and are
// marked by checkpoints in a rolling hash.
//
// The rolling hash checkpoints are marked when
// the low order bits are all zero of a Rabin
// fingerprint.  The number of low order bits
// for the fingerprint is taken to be:
//
//     1 + floor( sqrt( min(length(a), length(b)) ) )
//
// When a sub sequence is marked by it's endpoints
// in the rolling hash, the CRC32 checksum is calculated
// and stored.
//
// The two sequences checksums are then compared and an
// attempt to align matching checksums is made on the assumption
// that both sequences are mostly the same with some minor
// alterations in the middle.
//
// Runs of sub sequence that don't have aligned checksums
// are recursively passed into Memz which a lower threshold
// set to mark when true global alignment should be performed.
//
// If the mismatch run is found to be the length of the sequence,
// this function fails.
//
// If the checksum alignment step fails, this function fails.
//
// **THIS IS NOT GLOBAL ALIGNEMT**.  This is a heuristic.
// This is primarily meant to be run on long strings that are
// relatively similar to each other.  Though this isn't global
// alignemtn, some applications may find the alignment produced
// (if this heuristic succeeds) good enoug.
//
func BlockAlign(a, b []byte) {

  if len(a) < 1000 && len(b) < 1000 {
    fmt.Printf("len(a) %d, len(b) %d, use normal alignment\n", len(a), len(b))
    return
  }

  if len(a) == len(b) {
    n:=len(a)
    i:=0
    for i=0; i<n; i++ { if a[i]!=b[i] { break } }
    if i==n {
      fmt.Printf("match\n")
      return
    }
  }

  rs0 := rollsum.New()
  rs1 := rollsum.New()

  min_ab := min(len(a),len(b))
  block := int(math.Sqrt(float64(min_ab)))
  if block==0 { fmt.Printf("errr\n") }

  a_k := int(len(a)/block) + 1
  b_k := int(len(b)/block) + 1

  a_deci := make([]uint32, 0, a_k)
  b_deci := make([]uint32, 0, b_k)

  pos_a := make([]int, 0, a_k)
  pos_b := make([]int, 0, b_k)

  mask := make_mask(uint32(block))

  prev := 0
  pos := 0
  for i:=0; i<len(a); i++ {
    rs0.Roll(a[i])
    z := rs0.Digest()

    if z&mask == 0 {
      pos = i
      a_deci = append(a_deci, crc32.ChecksumIEEE(a[prev:pos]))
      pos_a = append(pos_a, pos)
    }
    prev=pos
  }
  if pos!=(len(a)-1) {
    a_deci = append(a_deci, crc32.ChecksumIEEE(a[prev:]))
    pos_a = append(pos_a, len(a))
  }

  prev = 0
  pos = 0
  for i:=0; i<len(b); i++ {
    rs1.Roll(b[i])
    z := rs1.Digest()

    if z&mask == 0 {
      pos = i
      b_deci = append(b_deci, crc32.ChecksumIEEE(b[prev:pos]))
      pos_b = append(pos_b, pos)
    }
    prev=pos
  }
  if pos!=(len(b)-1) {
    b_deci = append(b_deci, crc32.ChecksumIEEE(b[prev:]))
    pos_b = append(pos_b, len(b))
  }


  mm := min(len(a_deci), len(b_deci)) ; _ = mm
  MM := max(len(a_deci), len(b_deci)) ; _ = MM
  fmt.Printf("mm %d (%d, %d)\n", mm, len(a_deci), len(b_deci))


  // Find which hash in a matches the hash in b and
  // vice versa.
  // -1 indicates no match
  //
  prev_match := -1
  match_a := make([]int, len(a_deci))
  match_b := make([]int, len(b_deci))

  for i:=0; i<len(a_deci); i++ {

    match_a[i] = -1;
    for bpos:=prev_match+1; bpos<len(b_deci); bpos++ {
      match_b[bpos] = -1

      if a_deci[i] == b_deci[bpos] {
        match_b[bpos] = i
        match_a[i] = bpos
        prev_match = bpos
        break
      }
    }

  }

  //
  posa:=0
  posb:=0

  a_s := uint32(0)
  a_n := uint32(0)
  b_s := uint32(0)
  b_n := uint32(0)

  z := make([][4]uint32, 0, 8)

  for posa<len(match_a) && posb<len(match_b) {

    if (match_a[posa]>=0) && (match_b[posb]>=0) && (a_deci[posa]!=b_deci[posb]) {
      fmt.Printf("ERROR: match_a[%d] %d != match_b[%d] %d\n", posa, match_a[posa], posb, match_b[posb])
      break
    }

    if (match_a[posa] >= 0) && (a_deci[posa] == b_deci[posb]) {
      if a_n > 0 || b_n > 0 {
        z = append(z, [4]uint32{ a_s, a_n, b_s, b_n })
      }

      a_s = uint32(pos_a[posa])
      b_s = uint32(pos_b[posb])
      a_n = 0
      b_n = 0

      posa++
      posb++

      continue
    }


    if match_a[posa] < 0 {
      a_n = uint32(pos_a[posa]) - a_s
      posa++
    }

    if match_b[posb] < 0 {
      b_n = uint32(pos_b[posb]) - b_s
      posb++
    }

  }

  if posa!=len(match_a) {
    a_n = uint32(pos_a[len(pos_a)-1]) - a_s
  }

  if posb!=len(match_b) {
    b_n = uint32(pos_b[len(pos_b)-1]) - b_s
  }

  if a_n > 0 || b_n > 0 {
    z = append(z, [4]uint32{ a_s, a_n, b_s, b_n })
  }




  //DEBUG
  //
  for i:=0; i<MM; i++ {
    pa:=(0) ; da:=uint32(0) ; ma:=(0)
    pb:=(0) ; db:=uint32(0) ; mb:=(0)


    if i < len(pos_a) {
      pa = pos_a[i]
      da = a_deci[i]
      ma = match_a[i]
    }

    if i < len(pos_b) {
      pb = pos_b[i]
      db = b_deci[i]
      mb = match_b[i]
    }

    fmt.Printf("[%d] (%d,%d) %x (>%d) %x (<%d)\n", i, pa, pb, da, ma, db, mb)
    //fmt.Printf("[%d] (%d,%d) %x (>%d) %x (<%d)\n", i, pos_a[i], pos_b[i], a_deci[i], match_a[i], b_deci[i], match_b[i])
  }
  fmt.Printf(">>>>>>>>>>>>>\n")
  //
  //DEBUG

  for i:=0; i<len(z); i++ {
    fmt.Printf("  A(%d+%d) B(%d+%d)\n", z[i][0], z[i][1], z[i][2], z[i][3])

    a_start := z[i][0]
    a_n := z[i][1]

    b_start := z[i][2]
    b_n := z[i][3]

    if a_n == uint32(len(a)) && b_n == uint32(len(b)) {
      fmt.Printf("heuristic failed, bailing out\n")
      return
    }

    BlockAlign(a[a_start:a_start+a_n], b[b_start:b_start+b_n])
  }

  //for i:=0; i<len(match_a); i++ {
  //  fmt.Printf(">> match_a[%d] %d, match_b[%d] %d\n", i, match_a[i], i, match_b[i])
  //}


}

/*
func main() {

  rs := rollsum.New()

  prev_pos := 0
  pos := 0

  _ = x1

  memz([]byte(x0), []byte(x1))

  for i:=0; i<len(x0); i++ {
    rs.Roll(x0[i])
    z := rs.Digest()

    //fmt.Printf("%c", x[i])
    //if z&0x1f == 0 { fmt.Printf("|") }

    ch := '.'
    if z&0x7f == 0 {
      ch = '*'
      pos = i
    }
    fmt.Printf("[%d] %c %x %c (%d)\n", i, x0[i], z, ch, i-prev_pos)
    prev_pos = pos
  }

}
*/
