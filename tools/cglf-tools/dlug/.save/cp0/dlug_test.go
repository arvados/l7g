package dlug

import "fmt"
import "testing"

func TestGeneric(t *testing.T) {
  fmt.Printf(">>>\n")

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
  }


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

  }

  //                                  0  1    2    3    4    5      6      7        8        9         10         11           12           13             14             15                 16                 17         18
  marshal_uint64_test_len := []int{   1, 1,   1,   1,   2,   2,     2,     3,       3,       4,         4,         5,           5,           6,             6,             8,                 8,                 9,         9 }
  marshal_uint64_test_val := []uint64{0, 3, 126, 127, 128, 255, 16383, 16384, 2097151, 2097152, 134217727, 134217728, 34359738367, 34359738368, 8796093022207, 8796093022208, 72057594037927935, 72057594037927936, (1<<63)+1 }

  for i:=0; i<len(marshal_uint64_test_len); i++ {
    x := MarshalUint64(marshal_uint64_test_val[i])
    z := CheckCode(x)
    idx := GetDlugIndex(x)



    if !Check(x) { t.Errorf("Failed Check on trying to convert value %d (%x) (test %d) (len(%d)) (check code %d, index %d)!", marshal_uint64_test_val[i], marshal_uint64_test_val[i], i, len(x), z, idx) }
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

  }


}
