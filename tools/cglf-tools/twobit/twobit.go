// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

// Package twobit implements the 2bit compact randomly-accessible file format
// for storing DNA sequence data.
package twobit

import (
    "fmt"
    "io"
    "bytes"
    "bufio"
    "encoding/binary"
)

// 2bit header
type header struct {
    sig         uint32
    version     uint32
    count       uint32
    reserved    uint32
    byteOrder   binary.ByteOrder
}

// Block represents either blocks of Ns or masked (lower-case) blocks
type Block struct {
    start    int
    count    int
}

// seqRecord stores sequence record from the file index
type seqRecord struct {
    dnaSize      uint32
    nBlocks      []*Block
    mBlocks      []*Block
    reserved     uint32
    sequence     []byte
}

// TwoBit stores the file index and header information of the 2bit file
type twoBit struct {
    reader       io.ReadSeeker
    hdr          header
    index        map[string]int
    records      map[string]*seqRecord
}

type Reader twoBit
type Writer twoBit

func init() {
    NT2BYTES = make([]byte, 256)
    NT2BYTES[BASE_N]    = uint8(0)
    NT2BYTES[BASE_T]    = uint8(0)
    NT2BYTES[BASE_C]    = uint8(1)
    NT2BYTES[BASE_A]    = uint8(2)
    NT2BYTES[BASE_G]    = uint8(3)
    NT2BYTES[BASE_N+32] = uint8(0)
    NT2BYTES[BASE_T+32] = uint8(0)
    NT2BYTES[BASE_C+32] = uint8(1)
    NT2BYTES[BASE_A+32] = uint8(2)
    NT2BYTES[BASE_G+32] = uint8(3)
}

// Return the size in packed bytes of a dna sequence. 4 bases per byte
func packedSize(dnaSize int) (int) {
    return (dnaSize + 3) >> 2
}

// Return length of block
func (b *Block) Length() int {
    return b.start+b.count
}

// Return start of block
func (b *Block) Start() int {
    return b.start
}

// Return count of block
func (b *Block) Count() int {
    return b.count
}

// Return the size in bytes the seqRecord rec will take up in the twobit file
func (rec *seqRecord) size() int {
    size := 16 // dnaSize (4), nBlockCount (4), mBlockCount (4), reserved (4)

    size += 2 * 4 * len(rec.nBlocks) // nBlockStarts, nBlockSizes
    size += 2 * 4 * len(rec.mBlocks) // mBlockStarts, mBlockSizes
    size += len(rec.sequence)   // packedDNA

    return size
}

// Parse the file index of a 2bit file
func (r *Reader) parseIndex() (error) {
    r.index = make(map[string]int)

    for i := 0; i < r.Count(); i++ {
        size := make([]byte, 1)
        _, err := r.reader.Read(size)
        if err != nil {
            return fmt.Errorf("Failed to read file index: %s", err)
        }

        name := make([]byte, size[0])
        _, err = r.reader.Read(name)
        if err != nil {
            return fmt.Errorf("Failed to read file index: %s", err)
        }

        offset := make([]byte, 4)
        _, err = r.reader.Read(offset)
        if err != nil {
            return fmt.Errorf("Failed to read file index: %s", err)
        }

        r.index[string(name)] = int(r.hdr.byteOrder.Uint32(offset))
    }

    return nil
}

/*
// Seek to given offset in file
func (r *Reader) Seek(offset int64) (int64, error) {
    n, err := r.readseeker.Seek(offset, 0)
    r.reader.Reset(r.readseeker)
    return n, err
}
*/

// Parse the header of a 2bit file
func (r *Reader) parseHeader() (error) {
    b := make([]byte, 16)
    _, err := r.reader.Read(b)
    if err != nil {
        return err
    }

    r.hdr.sig = binary.BigEndian.Uint32(b[0:4])
    r.hdr.byteOrder = binary.BigEndian

    if r.hdr.sig != SIG {
        r.hdr.sig = binary.LittleEndian.Uint32(b[0:4])
        r.hdr.byteOrder = binary.LittleEndian
        if r.hdr.sig != SIG {
            return fmt.Errorf("Invalid sig. Not a 2bit file?")
        }
    }

    r.hdr.version = r.hdr.byteOrder.Uint32(b[4:8])
    if r.hdr.version != uint32(0) {
        return fmt.Errorf("Unsupported version %d", r.hdr.version)
    }
    r.hdr.count = r.hdr.byteOrder.Uint32(b[8:12])
    r.hdr.reserved = r.hdr.byteOrder.Uint32(b[12:16])
    if r.hdr.reserved != uint32(0) {
        return fmt.Errorf("Reserved != 0. got %d", r.hdr.reserved)
    }

    return nil
}

// Parse the nBlock and mBlock coordinates
func (r *Reader) parseBlockCoords() ([]*Block, error) {
    buf := make([]byte, 4)
    _, err := r.reader.Read(buf)
    if err != nil {
        return nil, fmt.Errorf("Failed to read blockCount: %s", err)
    }

    count := r.hdr.byteOrder.Uint32(buf)

    starts := make([]uint32, count)
    for i := range(starts) {
        _, err := r.reader.Read(buf)
        if err != nil {
            return nil, fmt.Errorf("Failed to block start: %s", err)
        }
        starts[i] = r.hdr.byteOrder.Uint32(buf)
    }

    sizes := make([]uint32, count)
    for i := range(sizes) {
        _, err := r.reader.Read(buf)
        if err != nil {
            return nil, fmt.Errorf("Failed to block size: %s", err)
        }
        sizes[i] = r.hdr.byteOrder.Uint32(buf)
    }

    blocks := make([]*Block, len(starts))

    for i := range(starts) {
        blocks[i] = &Block{start: int(starts[i]), count: int(sizes[i])}
    }

    return blocks, nil
}

// Parse the sequence record information
func (r *Reader) parseRecord(name string, coords bool) (*seqRecord, error) {
    rec := new(seqRecord)

    offset, ok := r.index[name]
    if !ok {
        return nil, fmt.Errorf("Invalid sequence name: %s", name)
    }

    r.reader.Seek(int64(offset), 0)

    buf := make([]byte, 4)
    _, err := r.reader.Read(buf)
    if err != nil {
        return nil, fmt.Errorf("Failed to read dnaSize: %s", err)
    }

    rec.dnaSize = r.hdr.byteOrder.Uint32(buf)

    if coords {
        rec.nBlocks, err = r.parseBlockCoords()
        if err != nil {
            return nil, fmt.Errorf("Failed to read nBlocks: %s", err)
        }

        rec.mBlocks, err = r.parseBlockCoords()
        if err != nil {
            return nil, fmt.Errorf("Failed to read mBlocks: %s", err)
        }

        _, err = r.reader.Read(buf)
        if err != nil {
            return nil, fmt.Errorf("Failed to read reserved: %s", err)
        }

        rec.reserved = r.hdr.byteOrder.Uint32(buf)

        if rec.reserved != uint32(0) {
            return nil, fmt.Errorf("Invalid reserved")
        }
    }

    return rec, nil
}

// Return blocks of Ns in sequence with name
func (r *Reader) NBlocks(name string) ([]*Block, error) {
    rec, err := r.parseRecord(name, true)
    if err != nil {
        return nil, err
    }

    return rec.nBlocks, nil
}

// Read entire sequence.
func (r *Reader) Read(name string) ([]byte, error) {
    return r.ReadRange(name, 0, 0)
}

// Read sequence from start to end.
func (r *Reader) ReadRange(name string, start, end int) ([]byte, error) {
    rec, err := r.parseRecord(name, true)
    if err != nil {
        return nil, err
    }

    bases := int(rec.dnaSize)

    // TODO: handle -1 ?
    if start < 0 {
        start = 0
    }

    //TODO: should we error out here?
    if end > bases {
        end = bases
    }

    // TODO: handle -1 ?
    if end == 0 || end < 0 {
        end = bases
    }

    if end <= start {
        return nil, fmt.Errorf("Invalid range: %d-%d", start, end)
    }

    bases = end-start
    size := packedSize(bases)
    if start > 0 {
        shift := packedSize(start)
        if start % 4 != 0 {
            shift--
            size++
        }

        r.reader.Seek(int64(shift), 1)
    }

    dna := make([]byte, size*4)
    chunks := size/defaultBufSize
    if size % defaultBufSize > 0 {
        chunks++
    }

    buf := make([]byte, defaultBufSize)

    i := 0
    for c := 0; c < chunks; c++ {
        sz := defaultBufSize
        if i+defaultBufSize > size {
            sz = size % defaultBufSize
        }
        n, err := r.reader.Read(buf[0:sz])
        if n != sz {
            return nil, fmt.Errorf("Failed to read %d dna bytes: %s", sz, err)
        } else if err != nil && err != io.EOF {
            return nil, fmt.Errorf("Failed to read dna bytes: %s", err)
        }

        for k := 0; k < n; k++ {
            base := buf[k]

            // ~800ms
            /*
            for j := 3; j >= 0; j-- {
                dna[(i*4)+j] = BYTES2NT[int(base & 0x3)]
                base >>= 2
            }
            */

            // ~540ms-650ms
            dna[i*4+3] =  BYTES2NT[((base   ) & 0x3)]
            dna[i*4+2] =  BYTES2NT[((base>>2) & 0x3)]
            dna[i*4+1] =  BYTES2NT[((base>>4) & 0x3)]
            dna[i*4  ] =  BYTES2NT[((base>>6) & 0x3)]

            i++
        }
    }

    seq := dna[(start%4):(start%4)+bases]

    for _, b := range rec.nBlocks {
        if b.Length() < start || b.start > end {
            continue
        }
        idx := b.start-start
        cnt := b.count
        if idx < 0 {
            cnt += idx
            idx = 0
        }
        for i := 0; i < cnt; i++ {
            seq[idx] = BASE_N
            idx++
            if idx >= len(seq) {
                break
            }
        }
    }

    for _, b := range rec.mBlocks {
        if b.Length() < start || b.start > end {
            continue
        }
        idx := b.start-start
        cnt := b.count
        if idx < 0 {
            cnt += idx
            idx = 0
        }

        for i := 0; i < cnt; i++ {
            // Faster lower case.. see: https://groups.google.com/forum/#!topic/golang-nuts/Il2DX4xpW3w
            seq[idx] = seq[idx] + 32 // ('a' - 'A')
            idx++
            if idx >= len(seq) {
                break
            }
        }
    }

    return seq, nil
}

// NewReader returns a new TwoBit file reader which reads from r
func NewReader(r io.ReadSeeker) (*Reader, error) {
    tb := new(Reader)
    tb.reader = r
    err := tb.parseHeader()
    if err != nil {
        return nil, err
    }

    err = tb.parseIndex()
    if err != nil {
        return nil, err
    }

    return tb, nil
}

// Returns the length for sequence with name
func (r *Reader) Length(name string) (int, error) {
    rec, err := r.parseRecord(name, false)
    if err != nil {
        return -1, err
    }

    return int(rec.dnaSize), nil
}

// Returns the length for sequence with name but does not count Ns
func (r *Reader) LengthNoN(name string) (int, error) {
    rec, err := r.parseRecord(name, true)
    if err != nil {
        return -1, err
    }

    n := 0
    for _, b := range rec.nBlocks {
        n += b.count
    }

    return int(rec.dnaSize)-n, nil
}

// Returns the names of sequences in the 2bit file
func (r *Reader) Names() ([]string) {
    names := make([]string, len(r.index))

    i := 0
    for n := range r.index {
        names[i] = n
        i++
    }

    return names
}

// Returns the count of sequences in the 2bit file
func (r *Reader) Count() (int) {
    return int(r.hdr.count)
}

// Returns the version of the 2bit file
func (r *Reader) Version() (int) {
    return int(r.hdr.version)
}

// Unpack array of bytes to DNA string of length sz
func Unpack(raw []byte, sz int) (string) {
    var dna bytes.Buffer
    for _, base := range raw {
        buf := make([]byte, 4)
        for j := 3; j >= 0; j-- {
            buf[j] = BYTES2NT[int(base & 0x3)]
            base >>= 2
        }

        dna.Write(buf)
    }

    return string(dna.Bytes()[0:sz])
}

// Packs DNA sequence string into an array of bytes. 4 bases per byte.
func Pack(s string) ([]byte, error) {
    sz := len(s)
    out := make([]byte, packedSize(sz))

    idx := 0
    for i := range out {
        var b uint8
        for j := 0; j < 4; j++ {
            //val := NT2BYTES['T']
            val := NT2BYTES['T']
            if idx < sz {
                // TODO: add check for unsupported base?
                val = NT2BYTES[s[idx]]
            }
            b <<= 2
            b += val
            idx++
        }
        out[i] = b
    }

    return out, nil
}

// New Writer
func NewWriter() (*Writer) {
    tb := new(Writer)
    tb.records = make(map[string]*seqRecord)

    return tb
}

func mapMBlocks(seq string) []*Block {
    blocks := make([]*Block, 0)

    n      := len(seq)
    start  := 0
    match  := false
    isLast := false
    for i := 0; i < n; i++ {
        match = seq[i] >= 'a' && seq[i] <= 'z'
        if match {
            if !isLast {
                start = i
            }
        } else {
            if isLast {
                blocks = append(blocks, &Block{start:start, count: i-start})
            }
        }
        isLast = match
    }

    if isLast {
        blocks = append(blocks, &Block{start:start, count: n-start})
    }

    return blocks
}

func mapNBlocks(seq string) []*Block {
    blocks := make([]*Block, 0)

    n      := len(seq)
    start  := 0
    match  := false
    isLast := false
    for i := 0; i < n; i++ {
        match = seq[i] == 'N' || seq[i] == 'n'
        if match {
            if !isLast {
                start = i
            }
        } else {
            if isLast {
                blocks = append(blocks, &Block{start:start, count: i-start})
            }
        }
        isLast = match
    }

    if isLast {
        blocks = append(blocks, &Block{start:start, count: n-start})
    }

    return blocks
}

// Add sequence
func (w *Writer) Add(name, seq string) (error) {
    if len(name) > 255 {
        return fmt.Errorf("Name string cannot be longer than 255 characters")
    }
    rec := new(seqRecord)
    rec.dnaSize = uint32(len(seq))
    rec.nBlocks = mapNBlocks(seq)
    rec.mBlocks = mapMBlocks(seq)

    pack, err := Pack(seq)
    if err != nil {
        return err
    }

    rec.sequence = pack

    w.records[name] = rec

    return nil
}

// Write sequences in 2bit format to out
func (w *Writer) WriteTo(out io.Writer) (error) {
    outbuf := bufio.NewWriter(out)

    buf := make([]byte, 16)
    binary.LittleEndian.PutUint32(buf[0:4], SIG)
    binary.LittleEndian.PutUint32(buf[4:8], uint32(0))
    binary.LittleEndian.PutUint32(buf[8:12], uint32(len(w.records)))
    binary.LittleEndian.PutUint32(buf[12:16], uint32(0))
    _, err := outbuf.Write(buf)
    if err != nil {
        return err
    }

    idxSize := 0
    recSize := 0
    var names []string
    for name, rec := range w.records {
        names = append(names, name)
        idxSize += 5 + len(name)
        recSize += rec.size()
    }

    buf = make([]byte, idxSize)
    offset := 16+idxSize
    idx := 0
    // Write out index
    for _, name := range names {
        buf[idx] = uint8(len(name))
        idx++
        for j := 0; j < len(name); j++ {
            buf[idx] = name[j]
            idx++
        }
        binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(offset))
        offset += w.records[name].size()
        idx += 4
    }

    _, err = outbuf.Write(buf)
    if err != nil {
        return err
    }

    // Write out records
    for _, name := range names {
        rec := w.records[name]
        sz := rec.size()
        buf = make([]byte, sz)

        binary.LittleEndian.PutUint32(buf[0:4], rec.dnaSize)
        binary.LittleEndian.PutUint32(buf[4:8], uint32(len(rec.nBlocks)))
        idx := 8
        for _, b := range rec.nBlocks {
            binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(b.start))
            idx += 4
        }
        for _, b := range rec.nBlocks {
            binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(b.count))
            idx += 4
        }

        binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(len(rec.mBlocks)))
        idx += 4
        for _, b := range rec.mBlocks {
            binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(b.start))
            idx += 4
        }
        for _, b := range rec.mBlocks {
            binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(b.count))
            idx += 4
        }

        // reserved
        binary.LittleEndian.PutUint32(buf[idx:idx+4], uint32(0))
        idx += 4

        copy(buf[idx:sz], rec.sequence[:])

        _, err := outbuf.Write(buf)
        if err != nil {
            return err
        }
    }

    err = outbuf.Flush()
    if err != nil {
        return err
    }

    return nil
}
