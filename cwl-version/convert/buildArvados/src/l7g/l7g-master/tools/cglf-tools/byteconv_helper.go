package main

func byte_low( b byte ) byte {
  if b >= 'A' && b <= 'Z' { return b+32 }
  return b
}


// little endien encoding

func tobyte16(b []byte, u16 uint16) {
  b[0] = byte(u16&0xff)
  b[1] = byte((u16>>8)&0xff)
}

func tobyte32(b []byte, u32 uint32) {
  b[0] = byte(u32&0xff)
  b[1] = byte((u32>>8)&0xff)
  b[2] = byte((u32>>16)&0xff)
  b[3] = byte((u32>>24)&0xff)
}

func tobyte64(b []byte, u64 uint64) {
  b[0] = byte(u64&0xff)
  b[1] = byte((u64>>8)&0xff)
  b[2] = byte((u64>>16)&0xff)
  b[3] = byte((u64>>24)&0xff)
  b[4] = byte((u64>>32)&0xff)
  b[5] = byte((u64>>40)&0xff)
  b[6] = byte((u64>>48)&0xff)
  b[7] = byte((u64>>56)&0xff)
}

func byte2uint32(buf []byte) (u32 uint32) {
  u32 = uint32(buf[0] & 0xff)
  u32 += uint32(buf[1] & 0xff) << 8
  u32 += uint32(buf[2] & 0xff) << 16
  u32 += uint32(buf[3] & 0xff) << 24
  return
}

func byte2uint64(buf []byte) (u64 uint64) {
  u64 = uint64(buf[0] & 0xff)
  u64 += uint64(buf[1] & 0xff) << 8
  u64 += uint64(buf[2] & 0xff) << 16
  u64 += uint64(buf[3] & 0xff) << 24
  u64 += uint64(buf[4] & 0xff) << 32
  u64 += uint64(buf[5] & 0xff) << 40
  u64 += uint64(buf[6] & 0xff) << 48
  u64 += uint64(buf[7] & 0xff) << 56
  return
}
