package memz

import "fmt"

var GAP int

func _score(x,y byte) int {
  if x==y { return 0 }
  return -3
}

var __score [][]int
var Score [][]int

func init() {

  __score = make([][]int, 256)
  for i:=0; i<256; i++ {
    __score[i] = make([]int, 256)
    for j:=0; j<256; j++ {
      __score[i][j] = -3
      if i==j { __score[i][j] = 0 }
    }
  }

  Score = make([][]int, 256)
  for i:=0; i<256; i++ {
    Score[i] = make([]int, 256)
    for j:=0; j<256; j++ {
      Score[i][j] = -3
      if i==j { Score[i][j] = 0 }
    }
  }

  GAP = -2

}

// Calculate the score column from both string
// suffixes.
//
// The implicit matrix can be thought of as
// having the x bytes on the top row and the
// y bytes on the left column with a gap (e.g.
// '-') prefix.
//
// The lower right hand score starts at 0 and
// denotes the base line maximum, regardless
// of what the optimal score is.  All scores
// are backtrack calculated assuming the optimal
// score is 0 in the lower right.
//
// score_col holds the 'last' (left most) score
// column.
//
func da_lin_suffix(x,y []byte, score_col []int) {
  var x_pos int
  var y_pos int
  var mscore int
  var tscore []int

  n := len(x)+2
  m := len(y)+2

  score0 := make([]int, m)
  score1 := make([]int, m)
  for i:=0; i<m; i++ { score0[i] = GAP*((m-1) - i) }

  // We've filled in n-1, start at n-2 and go to 0.
  //
  for col:=n-2; col>0; col-- {
    score1[m-1] = GAP*((n-1)-col)

    x_pos = col-1
    //tscore = __score[x[x_pos]]
    tscore = Score[x[x_pos]]

    // We've filled in the last row entry, start at m-1 and
    // go to 0.
    //
    for row:=m-2; row>=0; row-- {

      y_pos = row-1

      if col>0 && row>0 {
        mscore = tscore[y[y_pos]] + score0[row+1]

        if score0[row]   + GAP > mscore { mscore = score0[row]   + GAP }
        if score1[row+1] + GAP > mscore { mscore = score1[row+1] + GAP }

        score1[row] = mscore
      } else {
        mscore = score0[row+1]
        if score0[row]   + GAP > mscore { mscore = score0[row]   + GAP }
        if score1[row+1] + GAP > mscore { mscore = score1[row+1] + GAP }

        score1[row] = mscore
      }

    }

    score0,score1 = score1,score0
  }

  for i:=0; i<m; i++ {
    score_col[i] = score0[i]
  }

}


func da_lin(x,y []byte, score_col []int) {
  var x_pos int
  var y_pos int
  var tscore []int
  var mscore int

  n := len(x)+1
  m := len(y)+1

  score0 := make([]int, m)
  score1 := make([]int, m)
  for i:=0; i<m; i++ { score0[i] = i*GAP }

  for col:=1; col<n; col++ {
    score1[0] = GAP*col

    x_pos = col-1
    //tscore = __score[x[x_pos]]
    tscore = Score[x[x_pos]]
    for row:=1; row<m; row++ {
      y_pos = row-1

      mscore = tscore[y[y_pos]] + score0[row-1]

      if score0[row]   + GAP > mscore { mscore = score0[row]   + GAP }
      if score1[row-1] + GAP > mscore { mscore = score1[row-1] + GAP }

      score1[row] = mscore
    }

    score0,score1 = score1,score0
  }

  for i:=0; i<m; i++ {
    score_col[i] = score0[i]
  }

}

// Hirschberg's algorithm for dynamic programming.
// Hirschberg's algorithm finds the optimal score
// along with the path in linear space and quadratic
// time.
//
// The implicit DP matrix can be thought of as
// having the x bytes along the columns and y bytes
// along the rows with a gap ('-') prefix character
// attached to each. e.g. for x = 'ox', y = 'fox':
//
//         x0 x1
//       -  o  x
//    -  .  .  .
// y0 f  .  .  .
// y1 o  .  .  .
// y2 x  .  .  .
//
// The Hirschberg algorithm works by finding
// a midpoint where the path must pass through then
// recursively applying the Hirschberg on the upper
// left and lower right sub problem.  The midpoint
// is calculated by keeping the column of scores for
// the upper left and lower right sub problem then
// choosing the appropriate value.
//
// Each midpoint can be stored and constitues only
// an additional linear amount of space to store
// the path.
//
// When the sub problem becomes small enough (for
// example, when |x| or |y| <= 2) vanilla
// dynamic programming can be applied.
//
// HirschbergRecur returns an array and score,
// with even elements of the array being the
// 'x' positions in the implied dynamic programming
// matrix and the odd positions being the implied
// 'y' positions of the implied dynamic programming
// matrix.
//
// e.g.
// x = abc
// y = c
// ipath = [ 0, -1, 1, -1, 2, 0 ]
//
// would imply an alignment of
// abc
// --c
//
func HirschbergRecur(x,y []byte, basec,baser int) ([]int, int) {

  path := []int{}

  n_c := len(x)
  m_r := len(y)

  // The two column score vectors, one
  // more than length of y to include
  // gap.
  //
  y_pfx := make([]int, m_r+1)
  y_sfx := make([]int, m_r+2)

  // Base case, apply vanillla dynamic
  // programming.
  //
  if n_c<=2 || m_r<=2 {
    tpath,sc := Simp_b(x,y,basec,baser)
    path = append(path, tpath...)
    return path,sc
  }

  n2 := n_c/2

  // Find the two score vectors for the upper
  // left and lower right blocks.
  // They need to overlap in one column.
  //

  da_lin(x[:n2],y,y_pfx)
  da_lin_suffix(x[n2:],y,y_sfx)

  // Find the y position with the best cost.
  //
  best_cost := 0
  best_q := 0
  for q:=0; q<=m_r; q++ {
    cost := y_pfx[q] + y_sfx[q+1]
    if q==0 {
      best_cost = cost
      best_q = q
      continue
    }
    if cost > best_cost {
      best_cost = cost
      best_q = q
    }

  }

  //x_path_pos,y_path_pos := n2,best_q+1
  x_path_pos,y_path_pos := n2,best_q

  // left
  //
  tpathl,scl := HirschbergRecur(x[:x_path_pos],y[:y_path_pos], basec,baser)

  // right
  //
  tpathr,scr := HirschbergRecur(x[x_path_pos:],y[y_path_pos:], basec+x_path_pos,baser+y_path_pos)

  path = append(path, tpathl...)
  path = append(path, tpathr...)

  return path,scl+scr

}


func debug_matrix(x,y []byte) {
  // print out debugging matrix
  fmt.Printf("   ")
  for i:=0; i<=len(y); i++ {
    fmt.Printf(" %2d", i)
  }
  fmt.Printf("\n")

  fmt.Printf("     -")
  for i:=0; i<len(y); i++ {
    fmt.Printf("  %c", y[i])
  }
  fmt.Printf("\n")

  fmt.Printf("   -\n")
  for i:=0; i<len(x); i++ {
    fmt.Printf("%2d %c\n", i+1, x[i])
  }
  fmt.Printf("\n")
}

// Simple dynamic programming.
// Construct an n = |x| by m = |y| matrix (n columns,
// m rows).
//
func Simp_b(x,y []byte, basec, baser int) ([]int, int) {

  path := []int{}

  n_c := len(x)+1
  m_r := len(y)+1

  // Construct matrix and populate first row and column
  //
  M := make([][]int, m_r)
  for i:=0; i<m_r; i++ { M[i] = make([]int, n_c) }
  for i:=0; i<n_c; i++ { M[0][i] = i*GAP }
  for i:=0; i<m_r; i++ { M[i][0] = i*GAP }

  // Fill in matrix
  //
  for r:=1; r<m_r; r++ {
    for c:=1; c<n_c; c++ {
      xpos := c-1
      ypos := r-1
      //s := _score(x[xpos],y[ypos]) + M[r-1][c-1]
      s := Score[x[xpos]][y[ypos]] + M[r-1][c-1]
      if M[r-1][c] + GAP > s { s = M[r-1][c] + GAP }
      if M[r][c-1] + GAP > s { s = M[r][c-1] + GAP }
      M[r][c] = s
    }
  }

  sc := M[m_r-1][n_c-1]


  // Back track to find path.  Assume lower right hand
  // corner entry as starting point.
  //
  cur_c := n_c-1
  cur_r := m_r-1

  for cur_c>0 || cur_r>0 {

    path = append(path, cur_c+basec-1)
    path = append(path, cur_r+baser-1)

    if cur_c==0 { cur_r-- ; continue }
    if cur_r==0 { cur_c-- ; continue }

    //v_0_0 := _score(x[cur_c-1], y[cur_r-1]) + M[cur_r-1][cur_c-1]
    v_0_0 := Score[x[cur_c-1]][y[cur_r-1]] + M[cur_r-1][cur_c-1]
    v_1_0 := GAP + M[cur_r-1][cur_c]
    v_0_1 := GAP + M[cur_r][cur_c-1]

    if v_0_0 == M[cur_r][cur_c] { cur_r-- ; cur_c-- ; continue }
    if v_0_1 == M[cur_r][cur_c] { cur_c-- ; continue }
    if v_1_0 == M[cur_r][cur_c] { cur_r-- ; continue }

    return nil,-1
  }

  z := len(path)
  z2 := z/2
  for i:=0; i<z2; i+=2 {
    path[i],  path[z-i-2] = path[z-i-2],path[i]
    path[i+1],path[z-i-1] = path[z-i-1],path[i+1]
  }

  return path,sc
}

func Simp_b_old(x,y []byte, basec, baser int) ([]int, int) {

  path := []int{}

  n_c := len(x)+1
  m_r := len(y)+1

  // Construct matrix and populate first row and column
  //
  M := make([][]int, m_r)
  for i:=0; i<m_r; i++ { M[i] = make([]int, n_c) }
  for i:=0; i<n_c; i++ { M[0][i] = i*GAP }
  for i:=0; i<m_r; i++ { M[i][0] = i*GAP }

  // Fill in matrix
  //
  for r:=1; r<m_r; r++ {
    for c:=1; c<n_c; c++ {
      xpos := c-1
      ypos := r-1
      s := _score(x[xpos],y[ypos]) + M[r-1][c-1]
      if M[r-1][c] + GAP > s { s = M[r-1][c] + GAP }
      if M[r][c-1] + GAP > s { s = M[r][c-1] + GAP }
      M[r][c] = s
    }
  }

  sc := M[m_r-1][n_c-1]

  // Back track to find path.  Assume lower right hand
  // corner entry as starting point.
  //
  cur_c := n_c-1
  cur_r := m_r-1

  for cur_c>0 || cur_r>0 {

    path = append(path, cur_c+basec-1)
    path = append(path, cur_r+baser-1)

    if cur_c==0 { cur_r-- ; continue }
    if cur_r==0 { cur_c-- ; continue }

    v_0_0 := _score(x[cur_c-1], y[cur_r-1]) + M[cur_r-1][cur_c-1]
    v_1_0 := GAP + M[cur_r-1][cur_c]
    v_0_1 := GAP + M[cur_r][cur_c-1]

    if v_0_0 == M[cur_r][cur_c] { cur_r-- ; cur_c-- ; continue }
    if v_0_1 == M[cur_r][cur_c] { cur_c-- ; continue }
    if v_1_0 == M[cur_r][cur_c] { cur_r-- ; continue }

    return nil,-1
  }

  z := len(path)
  z2 := z/2
  for i:=0; i<z2; i+=2 {
    path[i],  path[z-i-2] = path[z-i-2],path[i]
    path[i+1],path[z-i-1] = path[z-i-1],path[i+1]
  }

  return path,sc
}


func debug_print_simp_p(x,y []byte, M [][]int) {
  n := len(x)+1
  m := len(y)+1

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==0 { fmt.Printf("     ")
    } else {
      fmt.Printf(" %4d", c)
    }
  }
  fmt.Printf("\n")

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==0 { fmt.Printf(" %4s", "-")
    } else {
      fmt.Printf(" %4c", x[c-1])
    }
  }
  fmt.Printf("\n")

  for r:=0; r<m; r++ {

    if r==0 {
      fmt.Printf(" %4c", '-')
    } else {
      fmt.Printf("%2d", r)
      fmt.Printf("  %c", y[r-1])
    }
    for c:=0; c<n; c++ {

      x_pos := c-1
      y_pos := r-1

      ch := ' '
      if r>0 && M[r-1][c] + GAP == M[r][c] {
        ch = '|'
      }
      if c>0 && M[r][c-1] + GAP == M[r][c] {
        ch = '_'
      }
      if r>0 && c>0 && M[r-1][c-1] + _score(x[x_pos],y[y_pos]) == M[r][c] {
        ch = '\\'
      }

      fmt.Printf(" %c%3d", ch, M[r][c])
    }
    fmt.Printf("\n")

  }
  fmt.Printf("\n")
}

func debug_print_simp(x,y []byte, M [][]int) {
  n := len(x)+1
  m := len(y)+1

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==0 { fmt.Printf("     ")
    } else {
      fmt.Printf(" %4d", c)
    }
  }
  fmt.Printf("\n")

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==0 { fmt.Printf(" %4s", "-")
    } else {
      fmt.Printf(" %4c", x[c-1])
    }
  }
  fmt.Printf("\n")

  for r:=0; r<m; r++ {

    if r==0 {
      fmt.Printf(" %4c", '-')
    } else {
      fmt.Printf("%2d", r)
      fmt.Printf("  %c", y[r-1])
    }
    for c:=0; c<n; c++ {

      mm := 0
      for ii:=0; ii<m; ii++ {
        if ii==0 { mm = M[ii][c] }
        if mm < M[ii][c] { mm=M[ii][c] }
      }

      ch := ' '
      if mm == M[r][c] { ch = '*' }


      fmt.Printf(" %c%3d", ch, M[r][c])
    }
    fmt.Printf("\n")

  }
  fmt.Printf("\n")
}



func debug_print_simp_rev2(x,y []byte, M [][]int) {
  n := len(x)+2
  m := len(y)+2

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==n-1 {
      fmt.Printf(" %4s", "-")
    } else if c==0 {
      fmt.Printf(" %4s", "-")
    } else {
      fmt.Printf(" %4c", x[c-1])
    }
  }
  fmt.Printf("\n")

  for r:=0; r<m; r++ {

    if r==m-1 {
      fmt.Printf(" %4c", '-')
    } else if r==0 {
      fmt.Printf(" %4c", '-')
    } else {
      fmt.Printf(" %4c", y[r-1])
    }
    for c:=0; c<n; c++ {

      mm := 0
      for ii:=0; ii<m; ii++ {
        if ii==0 { mm = M[ii][c] }
        if mm < M[ii][c] { mm=M[ii][c] }
      }

      ch := ' '
      if mm == M[r][c] { ch = '*' }


      fmt.Printf(" %c%3d", ch, M[r][c])
    }
    fmt.Printf("\n")

  }
  fmt.Printf("\n")
}

func debug_print_simp_rev(x,y []byte, M [][]int) {
  n := len(x)+1
  m := len(y)+1

  fmt.Printf("     ")
  for c:=0; c<n; c++ {
    if c==n-1 { fmt.Printf(" %4s", "-")
    } else {
      fmt.Printf(" %4c", x[c])
    }
  }
  fmt.Printf("\n")

  for r:=0; r<m; r++ {

    if r==m-1 {
      fmt.Printf(" %4c", '-')
    } else {
      fmt.Printf(" %4c", y[r])
    }
    for c:=0; c<n; c++ {

      mm := 0
      for ii:=0; ii<m; ii++ {
        if ii==0 { mm = M[ii][c] }
        if mm < M[ii][c] { mm=M[ii][c] }
      }

      ch := ' '
      if mm == M[r][c] { ch = '*' }


      fmt.Printf(" %c%3d", ch, M[r][c])
    }
    fmt.Printf("\n")

  }
  fmt.Printf("\n")
}



func simp(x,y []byte) {
  n := len(x)+1
  m := len(y)+1

  M := make([][]int, m)
  for i:=0; i<m; i++ { M[i] = make([]int, n) }
  for i:=0; i<n; i++ { M[0][i] = i*GAP }
  for i:=0; i<m; i++ { M[i][0] = i*GAP }

  for r:=1; r<m; r++ {
    for c:=1; c<n; c++ {
      xpos := c-1
      ypos := r-1
      s := _score(x[xpos],y[ypos]) + M[r-1][c-1]
      if M[r-1][c] + GAP > s { s = M[r-1][c] + GAP }
      if M[r][c-1] + GAP > s { s = M[r][c-1] + GAP }
      M[r][c] = s
    }
  }

  debug_print_simp(x,y,M)
}

func SimpRev(x,y []byte) {
  n := len(x)+2
  m := len(y)+2

  M := make([][]int, m)
  for i:=0; i<m; i++ { M[i] = make([]int, n) }
  for i:=0; i<n; i++ { M[m-1][n-1-i] = i*GAP }
  for i:=0; i<m; i++ { M[m-1-i][n-1] = i*GAP }

  for r:=m-2; r>=0; r-- {
    for c:=n-2; c>=0; c-- {
      xpos := c-1
      ypos := r-1

      if r>0 && c>0 {
        s := _score(x[xpos],y[ypos]) + M[r+1][c+1]
        if M[r+1][c] + GAP > s { s = M[r+1][c] + GAP }
        if M[r][c+1] + GAP > s { s = M[r][c+1] + GAP }
        M[r][c] = s
      } else {
        s := M[r+1][c+1]
        if M[r+1][c] + GAP > s { s = M[r+1][c] + GAP }
        if M[r][c+1] + GAP > s { s = M[r][c+1] + GAP }
        M[r][c] = s
      }
    }
  }

  debug_print_simp_rev2(x,y,M)
}

// Mostly a wrapper for HirschbergRecur.
// HirschbergRecur returns an array and score,
// with even elements of the array being the
// 'x' positions in the implied dynamic programming
// matrix and the odd positions being the implied
// 'y' positions of the implied dynamic programming
// matrix.
//
// e.g.
// x = abc
// y = c
// ipath = [ 0, -1, 1, -1, 2, 0 ]
//
// would imply an alignment of
// abc
// --c
//
// Hirschberg recur constructs the two strings (with '-'
// as the gap character) and the score.
//
func Hirschberg(x,y []byte) (a,b []byte, sc int) {
  ipath,sc := HirschbergRecur(x,y,0,0)

  if len(ipath)==0 { return nil,nil,-1 }

  prv_x := 0
  prv_y := 0
  for i:=0; i<len(ipath); i+=2 {

    if i==0 {

      if ipath[i]<0 {
        a = append(a,'-')
      } else {
        a = append(a, x[0])
      }

      if ipath[i+1]<0 {
        b = append(b, '-')
      } else {
        b = append(b, y[0])
      }

    } else {

      if (ipath[i]<0) || (ipath[i]==prv_x) {
        a = append(a, '-')
      } else {
        a = append(a, x[ipath[i]])
      }

      if (ipath[i+1]<0) || (ipath[i+1]==prv_y) {
        b = append(b, '-')
      } else {
        b = append(b, y[ipath[i+1]])
      }

    }

    prv_x = ipath[i]
    prv_y = ipath[i+1]
  }

  return a,b, sc
}
