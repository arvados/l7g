// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package twobit

const SIG = 0x1A412743

const defaultBufSize = 4096

const BASE_N = 'N'
const BASE_T = 'T'
const BASE_C = 'C'
const BASE_A = 'A'
const BASE_G = 'G'

/*
const BASE_N = 'n'
const BASE_T = 't'
const BASE_C = 'c'
const BASE_A = 'a'
const BASE_G = 'g'
*/

var BYTES2NT = []byte{
    BASE_T,
    BASE_C,
    BASE_A,
    BASE_G,
}

var NT2BYTES = []byte{}
