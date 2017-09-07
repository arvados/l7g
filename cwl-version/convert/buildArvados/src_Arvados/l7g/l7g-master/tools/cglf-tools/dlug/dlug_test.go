package dlug

import "testing"

import "os"

func TestByte(t *testing.T) {

  marshal_byte_test_len := []int{1, 1,   1,   1,   2,   2 };
  marshal_byte_test_val := []byte{0, 3, 126, 127, 128, 255 };

  for i:=0; i<len(marshal_byte_test_len); i++ {
    x := MarshalByte(marshal_byte_test_val[i])
    n := len(x)-1
    if len(x)!=marshal_byte_test_len[i] || x[n]!=marshal_byte_test_val[i] {
      t.Errorf("MarshalByte could not decode %d, got len(%d)", marshal_byte_test_val[i], len(x))
    }
    if n>1 && x[0] != Pfx[0] {
      t.Errorf("MarshalByte got extraneous higher order bits in return slice (%x)", x[0])
    }
    if !EqualByte(x,marshal_byte_test_val[i]) {
      t.Errorf("EqualByte not equal (value shoudl be %d)!", marshal_byte_test_val[i])
    }

    y := make([]byte, len(x))
    m := FillSliceByte(y, marshal_byte_test_val[i])
    if m!=len(x) { t.Errorf("n (%d) != m (%d)", len(x), m) }
    for j:=0; j<len(y); j++ {
      if x[j] != y[j] { t.Errorf("byte y[%d] (%d) != x[%d] (%d)", j, y[j], j, x[j]) }
    }

  }

  for i:=0; i<len(marshal_byte_test_val); i++ {
    x := MarshalByte(marshal_byte_test_val[i])
    b,n := ConvertByte(x)
    if n!=len(x) {
      t.Errorf("Test %d: Bad return length, got %d expected %d for byte val %x\n", i, n, len(x), marshal_byte_test_val[i])
    }
    if b!=marshal_byte_test_val[i] {
      t.Errorf("Test %d: Bad back conversion, got %x expected %x\n", i, b, marshal_byte_test_val[i])
    }
  }



}

func TestUint32(t *testing.T) {

  marshal_uint32_test_len := []int{   1, 1,   1,   1,   2,   2,     2,     3,       3,       4,         4,         5 };
  marshal_uint32_test_val := []uint32{0, 3, 126, 127, 128, 255, 16383, 16384, 2097151, 2097152, 134217727, 134217728 };

  for i:=0; i<len(marshal_uint32_test_len); i++ {
    x := MarshalUint32(marshal_uint32_test_val[i])
    if !Check(x) { t.Errorf("Failed Check!") }
    n := len(x)-1
    if len(x)!=marshal_uint32_test_len[i] || x[n]!=byte(0xff & marshal_uint32_test_val[i]) {
      t.Errorf("MarshalUint32 could not decode %d, got len(%d) expected len(%d)", marshal_uint32_test_val[i], len(x), marshal_uint32_test_len[i])
    }
    pfx_mask := byte(0xff << (8-byte(PfxBitLen[1])))
    if len(x)==2 && (pfx_mask & x[0]) != Pfx[1] {
      t.Errorf("MarshalUint32 got extraneous higher order bits in return slice (%x) expected %x", x[0], Pfx[1])
    }
    pfx_mask = byte(0xff << (8-byte(PfxBitLen[2])))
    if len(x)==3 && (pfx_mask & x[0]) != Pfx[2] {
      t.Errorf("MarshalUint32 got extraneous higher order bits in return slice (%x) expected %x", x[0], Pfx[2])
    }

    y := make([]byte, len(x))
    m := FillSliceUint32(y, marshal_uint32_test_val[i])
    if m!=len(x) { t.Errorf("n (%d) != m (%d)", len(x), m) }
    for j:=0; j<len(y); j++ {
      if x[j] != y[j] { t.Errorf("byte y[%d] (%d) != x[%d] (%d)", j, y[j], j, x[j]) }
    }

  }

  for i:=0; i<len(marshal_uint32_test_val); i++ {
    x := MarshalUint32(marshal_uint32_test_val[i])
    b,n := ConvertUint32(x)
    if n!=len(x) {
      t.Errorf("Test %d: Bad return length, got %d expected %d for uint32 val %x\n", i, n, len(x), marshal_uint32_test_val[i])
    }
    if b!=marshal_uint32_test_val[i] {
      t.Errorf("Test %d: Bad back conversion, got %x expected %x\n", i, b, marshal_uint32_test_val[i])
    }
  }

}

func TestUint64(t *testing.T) {

  marshal_uint64_test_len := []int{
  //0  1    2    3
    1, 1,   1,   1,
  //4  5  6  7  8  9  10
    2, 2, 2, 3, 3, 4, 4,
  //11 12 13 14 15 16 17 18
    5, 5, 6, 6, 8, 8, 9, 9 }

  marshal_uint64_test_val := []uint64{
    0, 3, 126, 127, 128,
    255, 16383, 16384, 2097151,
    2097152, 134217727, 134217728, 34359738367,
    34359738368, 8796093022207, 8796093022208, 72057594037927935, 72057594037927936, (1<<63)+1 }

  for i:=0; i<len(marshal_uint64_test_len); i++ {
    x := MarshalUint64(marshal_uint64_test_val[i])
    z := CheckCode(x)
    idx := GetDlugIndex(x)

    if !Check(x) {
      t.Errorf("Failed Check on trying to convert value %x (test %d) (len(%d)) (check code %d, index %d, expected len(%d)) [%x ...]!",
        marshal_uint64_test_val[i],
        i, len(x), z, idx, ByteLen[idx], x[0] )
    }
    n := len(x)-1
    if len(x)!=marshal_uint64_test_len[i] || x[n]!=byte(0xff & marshal_uint64_test_val[i]) {
      t.Errorf("MarshalUint64 could not decode %d, got len(%d) expected len(%d)", marshal_uint64_test_val[i], len(x), marshal_uint64_test_len[i])
    }
    pfx_mask := byte(0xff << (8-byte(PfxBitLen[1])))
    if len(x)==2 && (pfx_mask & x[0]) != Pfx[1] {
      t.Errorf("MarshalUint64 got extraneous higher order bits in return slice (%x) expected %x", x[0], Pfx[1])
    }
    pfx_mask = byte(0xff << (8-byte(PfxBitLen[2])))
    if len(x)==3 && (pfx_mask & x[0]) != Pfx[2] {
      t.Errorf("MarshalUint64 got extraneous higher order bits in return slice (%x) expected %x", x[0], Pfx[2])
    }

    y := make([]byte, len(x))
    m := FillSliceUint64(y, marshal_uint64_test_val[i])
    if m!=len(x) { t.Errorf("Converting %d, n (%d) != m (%d)", marshal_uint64_test_val[i], len(x), m) }
    for j:=0; j<len(y); j++ {
      if x[j] != y[j] { t.Errorf("Converting %d (%x), byte y[%d] (%d) != x[%d] (%d)", marshal_uint64_test_val[i], marshal_uint64_test_val[i], j, y[j], j, x[j]) }
    }

  }

  for i:=0; i<len(marshal_uint64_test_val); i++ {
    x := MarshalUint64(marshal_uint64_test_val[i])
    b,n := ConvertUint64(x)
    if n!=len(x) {
      t.Errorf("Test %d: Bad return length, got %d expected %d for uint64 val %x\n", i, n, len(x), marshal_uint64_test_val[i])
    }
    if b!=marshal_uint64_test_val[i] {
      t.Errorf("Test %d: Bad back conversion, got %x expected %x\n", i, b, marshal_uint64_test_val[i])
    }
  }

}

func TestCmp(t *testing.T) {

  byte_vec := []byte{0, 3, 126, 127, 128, 255 };
  uint32_vec := []uint32{0, 3, 126, 127, 128, 255, 16383, 16384, 2097151, 2097152, 134217727, 134217728 };
  uint64_vec := []uint64{
    0, 3, 126, 127, 128,
    255, 16383, 16384, 2097151,
    2097152, 134217727, 134217728, 34359738367,
    34359738368, 8796093022207, 8796093022208, 72057594037927935, 72057594037927936, (1<<63)+1 }

  _ = uint32_vec
  _ = uint64_vec

  //---

  for i:=0; i<len(byte_vec); i++ {
    xv := byte_vec[i]
    x := MarshalByte(xv)
    for j:=0; j<len(byte_vec); j++ {
      yv := byte_vec[j]
      y := MarshalByte(yv)
      z := Cmp(x, y)
      if xv < yv && z!=-1 { t.Errorf("ERROR: %d < %d but got %d from Cmp", xv, yv, z) }
      if xv>yv && z!=1 { t.Errorf("ERROR: %d > %d but got %d from Cmp", xv, yv, z) }
      if xv==yv && z!=0 { t.Errorf("ERROR: %d == %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(byte_vec); i++ {
    xv := byte_vec[i]
    x := MarshalByte(xv)
    for j:=0; j<len(uint32_vec); j++ {
      yv := uint32_vec[j]
      y := MarshalUint32(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: %d < %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: %d > %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: %d == %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(byte_vec); i++ {
    xv := byte_vec[i]
    x := MarshalByte(xv)
    for j:=0; j<len(uint64_vec); j++ {
      yv := uint64_vec[j]
      y := MarshalUint64(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: %d < %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: %d > %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: %d == %d but got %d from Cmp", xv, yv, z) }
    }
  }


  //---

  for i:=0; i<len(uint32_vec); i++ {
    xv := uint32_vec[i]
    x := MarshalUint32(xv)
    for j:=0; j<len(byte_vec); j++ {
      yv := byte_vec[j]
      y := MarshalByte(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint32) %d < (byte) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint32) %d > (byte) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint32) %d == (byte) %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(uint32_vec); i++ {
    xv := uint32_vec[i]
    x := MarshalUint32(xv)
    for j:=0; j<len(uint32_vec); j++ {
      yv := uint32_vec[j]
      y := MarshalUint32(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint32) %d < (uint32) %d but got %d from Cmp", xv, yv, z) ; os.Exit(1) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint32) %d > (uint32) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint32) %d == (uint32) %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(uint32_vec); i++ {
    xv := uint32_vec[i]
    x := MarshalUint32(xv)
    for j:=0; j<len(uint64_vec); j++ {
      yv := uint64_vec[j]
      y := MarshalUint64(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint32) %d < (uint64) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint32) %d > (uint64) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint32) %d == (uint64) %d but got %d from Cmp", xv, yv, z) }
    }
  }


  //---

  for i:=0; i<len(uint64_vec); i++ {
    xv := uint64_vec[i]
    x := MarshalUint64(xv)
    for j:=0; j<len(byte_vec); j++ {
      yv := byte_vec[j]
      y := MarshalByte(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint64) %d < (byte) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint64) %d > (byte) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint64) %d == (byte) %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(uint64_vec); i++ {
    xv := uint64_vec[i]
    x := MarshalUint64(xv)
    for j:=0; j<len(uint32_vec); j++ {
      yv := uint32_vec[j]
      y := MarshalUint32(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint64) %d < (uint32) %d but got %d from Cmp", xv, yv, z) ; os.Exit(1) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint64) %d > (uint32) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint64) %d == (uint32) %d but got %d from Cmp", xv, yv, z) }
    }
  }

  for i:=0; i<len(uint64_vec); i++ {
    xv := uint64_vec[i]
    x := MarshalUint64(xv)
    for j:=0; j<len(uint64_vec); j++ {
      yv := uint64_vec[j]
      y := MarshalUint64(yv)
      z := Cmp(x, y)
      if uint64(xv) < uint64(yv) && z!=-1 { t.Errorf("ERROR: (uint64) %d < (uint64) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) > uint64(yv) && z!=1 { t.Errorf("ERROR: (uint64) %d > (uint64) %d but got %d from Cmp", xv, yv, z) }
      if uint64(xv) == uint64(yv) && z!=0 { t.Errorf("ERROR: (uint64) %d == (uint64) %d but got %d from Cmp", xv, yv, z) }
    }
  }

}
