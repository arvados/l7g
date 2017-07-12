package dlug

//                          0,1,2,3,4,5,6,7,8
var ByteLen   []int = []int{1,2,3,4,5,6,8,9,17}
var PfxBitLen []int = []int{1,2,3,5,5,5,8,8,8}

var BitLen []uint = []uint{7,14,21,27,35,43,56,64,128}
var Pfx []byte = []byte{0,0x80,0xc0,0xe0,0xe8,0xf0,0xf8,0xf9,0xfa,0xff}

func Check(d []byte) bool {
  if len(d)==0 { return false }
  idx := GetDlugIndex(d)
  if idx<0 { return false }
  if idx>= len(ByteLen) { return false }
  if len(d) != ByteLen[idx] { return false }
  return true
}

func CheckCode(d []byte) int {
  if len(d)==0 { return -1 }
  idx := GetDlugIndex(d)
  if idx<0 { return -2 }
  if idx>= len(ByteLen) { return -3 }
  if len(d) != ByteLen[idx] { return -4 }
  return 0
}

// If the byte vectors are equal, we can compare on a byte by byte
// level.
// If they're not, we calculate the first non-zero position
// and compare each byte array from the appropriate position.
//
func Cmp(b0, b1 []byte) int32 {
  if len(b0)==len(b1) {
    if len(b0)==0 { return 0 }

    k := GetDlugIndex(b0)
    if k<0 { return 0 }

    ktst := GetDlugIndex(b1)
    if ktst<0 { return 0 }

    pfx0 := b0[0] & (0xff >> byte(PfxBitLen[k]))
    pfx1 := b1[0] & (0xff >> byte(PfxBitLen[k]))

    if pfx0 < pfx1 { return -1 }
    if pfx0 > pfx1 { return 1 }

    for i:=1; i<len(b0); i++ {
      if b0[i] < b1[i] { return -1 }
      if b0[i] > b1[i] { return 1 }
    }

    return 0
  }

  if len(b0)==0 { return -1 }
  if len(b1)==0 { return  1 }

  idx0 := GetDlugIndex(b0)
  if idx0<0 { return 0 }

  idx1 := GetDlugIndex(b1)
  if idx1<0 { return 0 }

  pfx0 := b0[0] & (0xff >> byte(PfxBitLen[idx0]))
  pfx1 := b1[0] & (0xff >> byte(PfxBitLen[idx1]))

  nz0 := 0
  if pfx0 == 0 {
    n := len(b0)
    for nz0=1; nz0<n; nz0++ { if b0[nz0] > 0 { break } }
  }

  nz1 := 0
  if pfx1 == 0 {
    n := len(b1)
    for nz1=1; nz1<n; nz1++ { if b1[nz1] > 0 { break } }
  }

  remain_len0 := len(b0)-nz0
  remain_len1 := len(b1)-nz1

  if remain_len0 < remain_len1 { return -1 }
  if remain_len0 > remain_len1 { return 1 }

  for i:=0; i<remain_len0; i++ {

    a := b0[nz0+i]
    if nz0+i == 0 { a = (0xff >> byte(PfxBitLen[idx0])) }

    b := b1[nz1+i]
    if nz1+i == 0 { b = (0xff >> byte(PfxBitLen[idx1])) }

    if a < b { return -1 }
    if a > b { return  1 }

  }
  return 0

}

func EqualByte(d []byte, b byte) bool {
  if len(d)==0 { return false }

  if len(d)==1 {
    if (d[0]&(0x80)) != 0 { return false }
    if (d[0]&0x7f) == b { return true }
    return false
  }

  k := GetDlugIndex(d)
  if k<0 { return false }

  if d[0]&byte(0xff << (8-byte(PfxBitLen[k]))) != Pfx[k] { return false }
  n:=len(d)
  if d[n-1]!=b {return false}
  for i:=1; i<(n-1); i++ { if d[i]!=0 { return false } }
  return true
}

func GetDlugIndex(d []byte) int {
  if len(d)==0 { return -1 }
  for i:=0; i<len(ByteLen); i++ {
    if (d[0] & byte(0xff << (8-byte(PfxBitLen[i])))) == Pfx[i] {
      return i
    }
  }
  return -2
}

func GetByteLen(d []byte) int {
  if len(d)==0 { return -1 }
  for i:=0; i<len(ByteLen); i++ {
    if (d[0] & byte(0xff << (8-byte(PfxBitLen[i])))) == Pfx[i] {
      return ByteLen[i]
    }
  }
  return -2
}

func GetDataBitLen(d []byte) int {
  if len(d)==0 { return -1 }
  for i:=0; i<len(ByteLen); i++ {
    if (d[0] & byte(0xff << (8-byte(PfxBitLen[i])))) == Pfx[i] {
      return int(BitLen[i])
    }
  }
  return -2
}

func GetPrefixBitLen(d []byte) int {
  if len(d)==0 { return -1 }
  for i:=0; i<len(ByteLen); i++ {
    if (d[0] & byte(0xff << (8-byte(PfxBitLen[i])))) == Pfx[i] {
      return PfxBitLen[i]
    }
  }
  return -2
}

//-----------------------
// Marshal Byte Functions
//-----------------------

func MarshalByte(b byte) []byte {
  if b<(1<<BitLen[0]) { return []byte{b} }
  return []byte{ 0x80, b }
}

func MarshalUint32(u uint32) []byte {
  if u<(1<<BitLen[0]) { return []byte{ byte(u&0xff) } }
  if u<(1<<BitLen[1]) {
    return []byte{ byte(Pfx[1] | byte(0xff & (u>>8))), byte(0xff & u) }
  }
  if u<(1<<BitLen[2]) {
    return []byte{ byte(Pfx[2] | byte(0xff & (u>>16))), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  if u<(1<<BitLen[3]) {
    return []byte{ byte(Pfx[3] | byte(0xff & (u>>24))), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  return []byte{ Pfx[4], byte(0xff & (u>>24)), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
}

func MarshalUint64(u uint64) []byte {
  if u<(1<<BitLen[0]) { return []byte{ byte(u&0xff) } }
  if u<(1<<BitLen[1]) {
    return []byte{ byte(Pfx[1] | byte(0xff & (u>>8))), byte(0xff & u) }
  }
  if u<(1<<BitLen[2]) {
    return []byte{ byte(Pfx[2] | byte(0xff & (u>>16))), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  if u<(1<<BitLen[3]) {
    return []byte{ byte(Pfx[3] | byte(0xff & (u>>24))), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  if u<(1<<BitLen[4]) {
    return []byte{ byte(Pfx[4] | byte(0xff & (u>>32))), byte(0xff & (u>>24)), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  if u<(1<<uint64(BitLen[5])) {
    return []byte{ byte(Pfx[5] | byte(0xff & (u>>40))), byte(0xff & (u>>32)), byte(0xff & (u>>24)), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  if u<(1<<uint64(BitLen[6])) {
    return []byte{ Pfx[6], byte(0xff & (u>>48)), byte(0xff & (u>>40)), byte(0xff & (u>>32)), byte(0xff & (u>>24)), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
  }
  return []byte{ Pfx[7], byte(0xff & (u>>56)), byte(0xff & (u>>48)), byte(0xff & (u>>40)), byte(0xff & (u>>32)), byte(0xff & (u>>24)), byte(0xff & (u>>16)), byte(0xff & (u>>8)), byte(0xff & u) }
}

//---------------------
// Fill Slice Functions
//---------------------

func FillSliceByte(s []byte, b byte) int {
  if len(s) == 0 { return -1 }
  if b<(1<<BitLen[0]) { s[0] = b; return 1 }

  if len(s) < 2 { return -1 }
  s[0] = 0x80
  s[1] = b
  return 2
}

func FillSliceUint32(s []byte, u uint32) int {

  if len(s)==0 { return -1 }
  if u<(1<<BitLen[0]) {
    s[0] = byte(u&0xff)
    return 1
  }

  if len(s)<int(ByteLen[1]) { return -1 }
  if u<(1<<BitLen[1]) {
    s[0] = byte(Pfx[1] | byte(0xff & (u>>8)))
    s[1] = byte(0xff & u)
    return ByteLen[1]
  }

  if len(s)<ByteLen[2] { return -1 }
  if u<(1<<BitLen[2]) {
    s[0] = byte(Pfx[2] | byte(0xff & (u>>16)))
    s[1] = byte(0xff & (u>>8))
    s[2] = byte(0xff & u)
    return ByteLen[2]
  }

  if len(s)<ByteLen[3] { return -1 }
  if u<(1<<BitLen[3]) {
    s[0] = byte(Pfx[3] | byte(0xff & (u>>24)))
    s[1] = byte(0xff & (u>>16))
    s[2] = byte(0xff & (u>>8))
    s[3] = byte(0xff & u)
    return ByteLen[3]
  }

  if len(s)<ByteLen[4] { return -1 }
  s[0] = Pfx[4]
  s[1] = byte(0xff & (u>>24))
  s[2] = byte(0xff & (u>>16))
  s[3] = byte(0xff & (u>>8))
  s[4] = byte(0xff & u)
  return ByteLen[4]

}


func FillSliceUint64(s []byte, u uint64) int {

  if len(s)==0 { return -1 }
  if u<(1<<BitLen[0]) {
    s[0] = byte(u&0xff)
    return 1
  }

  if len(s)<ByteLen[1] { return -1 }
  if u<(1<<BitLen[1]) {
    s[0] = byte(Pfx[1] | byte(0xff & (u>>8)))
    s[1] = byte(0xff & u)
    return ByteLen[1]
  }

  if len(s)<ByteLen[2] { return -1 }
  if u<(1<<BitLen[2]) {
    s[0] = byte(Pfx[2] | byte(0xff & (u>>16)))
    s[1] = byte(0xff & (u>>8))
    s[2] = byte(0xff & u)
    return ByteLen[2]
  }

  if len(s)<ByteLen[3] { return -1 }
  if u<(1<<BitLen[3]) {
    s[0] = byte(Pfx[3] | byte(0xff & (u>>24)))
    s[1] = byte(0xff & (u>>16))
    s[2] = byte(0xff & (u>>8))
    s[3] = byte(0xff & u)
    return ByteLen[3]
  }

  if len(s)<ByteLen[4] { return -1 }
  if u<(1<<BitLen[4]) {
    s[0] = Pfx[4] | byte(0xff & (u>>32))
    s[1] = byte(0xff & (u>>24))
    s[2] = byte(0xff & (u>>16))
    s[3] = byte(0xff & (u>>8))
    s[4] = byte(0xff & u)
    return ByteLen[4]
  }

  if len(s)<ByteLen[5] { return -1 }
  if u<(1<<uint64(BitLen[5])) {
    s[0] = Pfx[5] | byte(0xff & (u>>40))
    s[1] = byte(0xff & (u>>32))
    s[2] = byte(0xff & (u>>24))
    s[3] = byte(0xff & (u>>16))
    s[4] = byte(0xff & (u>>8))
    s[5] = byte(0xff & u)
    return ByteLen[5]
  }

  if len(s)<ByteLen[6] { return -1 }
  if u<(1<<uint64(BitLen[6])) {
    s[0] = Pfx[6]
    s[1] = byte(0xff & (u>>48))
    s[2] = byte(0xff & (u>>40))
    s[3] = byte(0xff & (u>>32))
    s[4] = byte(0xff & (u>>24))
    s[5] = byte(0xff & (u>>16))
    s[6] = byte(0xff & (u>>8))
    s[7] = byte(0xff & u)
    return ByteLen[6]
  }

  if len(s)<ByteLen[7] { return -1 }
  s[0] = Pfx[7]
  s[1] = byte(0xff & (u>>56))
  s[2] = byte(0xff & (u>>48))
  s[3] = byte(0xff & (u>>40))
  s[4] = byte(0xff & (u>>32))
  s[5] = byte(0xff & (u>>24))
  s[6] = byte(0xff & (u>>16))
  s[7] = byte(0xff & (u>>8))
  s[8] = byte(0xff & u)
  return ByteLen[7]

}

//------------------
// Convert Functions
//------------------

func ConvertByte(b []byte) (byte, int) {
  idx := GetDlugIndex(b)
  if idx<0 { return 0,idx }
  if idx==0 { return b[0]&0x7f,1 }

  if len(b) < ByteLen[idx] { return 0,-1 }
  return b[ByteLen[idx]-1], ByteLen[idx]
}

func ConvertUint32(b []byte) (uint32, int) {
  idx := GetDlugIndex(b)
  if idx<0 { return 0,idx }
  if idx==0 { return uint32(b[0]&0x7f),1 }

  if len(b) < ByteLen[idx] { return 0,-1 }
  if idx==1 { return (uint32(b[0]&(^Pfx[1])) << 8) + uint32(b[1]), 2 }
  if idx==2 { return (uint32(b[0]&(^Pfx[2])) << 16) + (uint32(b[1])<<8) + uint32(b[2]), 3 }
  if idx==3 { return (uint32(b[0]&(^Pfx[3])) << 24) + (uint32(b[1])<<16) + (uint32(b[2])<<8) + uint32(b[3]), 4 }

  n := ByteLen[idx]
  return (uint32(b[n-4])<<24) + (uint32(b[n-3])<<16) + (uint32(b[n-2])<<8) + uint32(b[n-1]), n
}

func ConvertUint64(b []byte) (uint64, int) {
  idx := GetDlugIndex(b)
  if idx<0 { return 0,idx }
  if idx==0 { return uint64(b[0]&0x7f),1 }

  if len(b) < ByteLen[idx] { return 0,-1 }
  if idx==1 { return (uint64(b[0]&(^Pfx[1]))<<8) + uint64(b[1]), 2 }
  if idx==2 { return (uint64(b[0]&(^Pfx[2])) << 16) + (uint64(b[1])<<8) + uint64(b[2]), 3 }
  if idx==3 { return (uint64(b[0]&(^Pfx[3])) << 24) + (uint64(b[1])<<16) + (uint64(b[2])<<8) + uint64(b[3]), 4 }
  if idx==4 { return (uint64(b[0]&(^Pfx[4])) << 32) + (uint64(b[1])<<24) + (uint64(b[2])<<16) + (uint64(b[3])<<8) + uint64(b[4]), 5 }
  if idx==5 { return (uint64(b[0]&(^Pfx[5])) << 40) + (uint64(b[1])<<32) + (uint64(b[2])<<24) + (uint64(b[3])<<16) + (uint64(b[4])<<8) + uint64(b[5]), 6 }
  if idx==6 { return (uint64(b[1])<<48) + (uint64(b[2])<<40) + (uint64(b[3])<<32) + (uint64(b[4])<<24) + (uint64(b[5])<<16) + (uint64(b[6])<<8) + uint64(b[7]), 8 }
  if idx==6 { return (uint64(b[1])<<48) + (uint64(b[2])<<40) + (uint64(b[3])<<32) + (uint64(b[4])<<24) + (uint64(b[5])<<16) + (uint64(b[6])<<8) + uint64(b[7]), 8 }

  n := ByteLen[idx]
  return (uint64(b[n-8])<<56) + (uint64(b[n-7])<<48) + (uint64(b[n-6])<<40) + (uint64(b[n-5])<<32) +
         (uint64(b[n-4])<<24) + (uint64(b[n-3])<<16) + (uint64(b[n-2])<<8) + uint64(b[n-1]), n
}

