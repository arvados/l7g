// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package twobit

import (
    "testing"
    "bytes"
    "bufio"
    "os"
    "reflect"
    "crypto/md5"
    "fmt"
)

func openTestTwoBit() (*Reader, error) {
    f, err := os.Open("examples/simple.2bit")
    if err != nil {
        return nil, err
    }

    tb, err := NewReader(f)
    if err != nil {
        return nil, err
    }

    return tb, nil
}

func TestHeader(t *testing.T) {
    tb, err := openTestTwoBit()
    if err != nil {
        t.Errorf("%s", err)
    }

    if tb.Count() != 1 {
        t.Errorf("Invalid sequence count: %d != %d", tb.Count(), 1)
    }

    names := map[string]bool{
        "ex1"   : false,
    }

    for name := range tb.index {
        if _, ok := names[name]; !ok {
            t.Errorf("Invalid sequence name: %s", name)
        }
        names[name] = true
    }
    for name, seen := range names {
        if !seen {
            t.Errorf("Sequence name not found in file index: %s", name)
        }
    }
}

func TestNamesLength(t *testing.T) {
    tb, err := openTestTwoBit()
    if err != nil {
        t.Errorf("%s", err)
    }

    names := tb.Names()

    if len(names) != 1 {
        t.Errorf("Invalid length of sequence names: %d != %d", len(names), 1)
    }

    if names[0] != "ex1" {
        t.Errorf("Invalid sequence name: %s != %s", names[0], "ex1")
    }

    sz, err := tb.Length("ex1")
    if err != nil {
        t.Errorf("%s", err)
    }
    if sz != 21 {
        t.Errorf("Invalid length of ex1 sequence: %d != %d", sz, 21)
    }

    sz, err = tb.LengthNoN("ex1")
    if err != nil {
        t.Errorf("%s", err)
    }
    if sz != 15 {
        t.Errorf("Invalid lengthNoN of ex1 sequence: %d != %d", sz, 15)
    }
}

func TestRead(t *testing.T) {
    tb, err := openTestTwoBit()
    if err != nil {
        t.Errorf("%s", err)
    }

    _, err = tb.Read("not-found")
    if err == nil {
        t.Errorf("Found non-existent name")
    }

    regions := map[string][]int {
        "ACTgcctttnnnNantnaCgc": []int{0, 0},
        "ACTgc"                : []int{0, 5},
             "ctttnn"          : []int{5, 11},
                       "tnaCgc": []int{15, 21},
                           "gc": []int{19, 21},
                            "c": []int{20, 21},
    }

    for good, coords := range regions {
        seq, err := tb.ReadRange("ex1", coords[0], coords[1])
        if err != nil {
            t.Errorf("Failed to read sequence: %s", err)
        }

        if string(seq) != good {
            t.Errorf("Invalid sequence: %s != %s", seq, good)
        }
    }
}

func TestPack(t *testing.T) {
    seqs := map[string]string {
        "ACTgcctttnnnNantnaCgc": "ACTGCCTTTTTTTATTTACGC",
        "ACTg"                 : "ACTG",
        "AC"                   : "AC",
        "ACcgcgcgcgcgcg"       : "ACCGCGCGCGCGCG",
    }

    for from, to := range seqs {
        p, err := Pack(from)
        if err != nil {
            t.Errorf("Failed to pack sequence: %s", err)
        }

        b := Unpack(p, len(from))

        if to != b {
            t.Errorf("Invalid sequence packing: %s != %s", to, b)
        }
    }
}

func TestAdd(t *testing.T) {
    tb := NewWriter()

    name := "ex1"
    seq  := "ACTgcctttnnnNantnaCgc"

    err := tb.Add(name, seq)
    if err != nil {
        t.Errorf("Failed to add sequence: %s", err)
    }

    if len(tb.records) != 1 {
        t.Errorf("Invalid sequence count: %d != %d", len(tb.records), 1)
    }

    rec, ok := tb.records[name]
    if !ok {
        t.Errorf("ex1 sequence not found")
    }

    unpacked := Unpack(rec.sequence, len(seq))
    good := "ACTGCCTTTTTTTATTTACGC"
    if unpacked != good {
        t.Errorf("Invalid packed sequence: %s != %s", unpacked, good)
    }

    if len(rec.nBlocks) != 3 {
        t.Errorf("invalid nBlock count: %d != %d", len(rec.nBlocks), 3)
    }

    if len(rec.mBlocks) != 3 {
        t.Errorf("invalid mBlock count: %d != %d", len(rec.mBlocks), 3)
    }

    nBlocks := []*Block{
        &Block{start:9, count:4},
        &Block{start:14, count:1},
        &Block{start:16, count:1},
    }

    mBlocks := []*Block{
        &Block{start:3, count:9},
        &Block{start:13, count:5},
        &Block{start:19, count:2},
    }

    if !reflect.DeepEqual(nBlocks, rec.nBlocks) {
        t.Errorf("invalid nBlocks : %#v != %#v", nBlocks, rec.nBlocks)
    }

    if !reflect.DeepEqual(mBlocks, rec.mBlocks) {
        t.Errorf("invalid mBlocks : %#v != %#v", mBlocks, rec.mBlocks)
    }
}

func TestWrite(t *testing.T) {
    tbw := NewWriter()

    name := "ex1"
    seq  := "ACTgcctttnnnNantnaCgc"

    err := tbw.Add(name, seq)
    if err != nil {
        t.Errorf("Failed to add sequence: %s", err)
    }

    var out bytes.Buffer
    outWriter := bufio.NewWriter(&out)
    tbw.WriteTo(outWriter)
    outWriter.Flush()

    twobitOut := out.Bytes()

    if "e8366f5785d6bf4b34595668d1509cb3" != fmt.Sprintf("%x", md5.Sum(twobitOut)) {
        t.Errorf("Invalid 2bit output. Failed md5sum check")
    }
}
