package memz

import "fmt"

// http://genome.sph.umich.edu/wiki/Variant_Normalization
// pseudo code:
//   processing <- true
//   while processing do
//     /* Extend (+) */
//     if alleles end with same nucleotide then
//       truncate rightmost nucleotide on each allele
//     end if
//     /* Delete (-) */
//     if there's an empty allele then
//       extend both alleles 1 nucleotide to the left
//     end if
//   end while
//   while leftmost nucleotide of each allele are the same and all alleles have length 2 or more do
//     truncate leftmost nucleotide of each allele
//   end while
//
// for example, if we had the sequences (pre-aligned):
//
// gactactg
// gact---g
//
// Then the series of operations would be:
//
//    act -(+)-> tact -(-)-> tac -(+)-> ctac -(-)-> cta -(+)-> acta -(-)-> act -(+)-> gact
//    _   -(+)-> t    -(-)-> _   -(+)-> c    -(-)-> _   -(+)-> a    -(-)-> _   -(+)-> g
//
// From a sequence pair-alignment view, it looks vaguely like the following:
//
//  gactactg => gactactg => gactactg => gactactg
//  gact___g => gac___tg => ga___ctg => g___actg
//
// The final step of culling nucleotides from the left does not apply to this example.
//
// Note, alignments like the following may be confusing:
//
//    gcatgcatg
//    g----catg
//
//  In normalized VCF format, this would change to:
//
//    gcatgcatg
//    gcat----g
//
//  This is valid and expected.  It's informative to realize that you're really
//  trying to express the alignment of the following two sequences:
//
//    gcatgcatg
//    gcatg
//
//  Aligning the string 'gcatgcatg' to 'gcatg' could mean
//  "replace the first occurance of 'catg' with gaps" or could also mean "replace the
//  second occurance of 'gcat' with gaps".  The VCF renormalization step chooses the
//  first as the canonical representation.
//

// Note, normalization only goes so far.  The vt folks differentiate between 'nomralization'
// and 'decomposition'.  See https://github.com/atks/vt/issues/16 for a discussion.
// Note that a string like:
//
//    fooquxbar
//    fooxzqbar
//
// can be aligned in two different, equally valid ways:
//
//    fooqux--bar     fooxzq--bar
//    foo--xzqbar     foo--quxbar
//
// for a gap penalty of -2 and a misalignment penalty of -3, both alignemnts
// give -8 and both alignments are normalized.
//



// Move left to the first non gap character
//
func _ldash( seq []byte ) int {
  for n := len(seq)-1; n>0; n-- {
    if seq[n] != '-' { return n }
  }
  return 0
}


func bp_eq( a,b byte ) bool {
  if a==b { return true }
  if (a=='N' || a=='n') && b!='-' { return true }
  if (b=='N' || b=='n') && a!='-' { return true }
  return false
}


func SeqPairNormalize(seq_a,seq_b []byte) error {
  if len(seq_a) != len(seq_b) {
    return fmt.Errorf("Sequence lengths do not match.  Sequences must be of same length")
  }

  b0 := len(seq_a)-1 ; n0 := 1
  b1 := len(seq_b)-1 ; n1 := 1

  for (b0>0) && (b1>0) {

    processing := true
    for processing {
      processing = false

      if (n0>0) && (n1>0) && bp_eq(seq_a[b0+n0-1], seq_b[b1+n1-1]) {
        processing = true
        n0--
        n1--
      }

      if (b0==0) || (b1==0) { break }

      if (n0==0) || (n1==0) {
        processing = true

        b0-- ; n0++
        b1-- ; n1++

        // inefficient, replace with saved left state
        // after we get working.
        //
        // Swap gap character with current character
        l0 := b0
        if seq_a[b0] == '-' {
          l0 = _ldash( seq_a[0:l0] )
        }

        l1 := b1
        if seq_b[b1] == '-' {
          l1 = _ldash( seq_b[0:l1] )
        }

        if bp_eq(seq_a[l0], seq_b[l1]) {
          if seq_a[b0] == '-' {
            seq_a[b0], seq_a[l0] = seq_a[l0], seq_a[b0]
          }
          if seq_b[b1] == '-' {
            seq_b[b1], seq_b[l1] = seq_b[l1], seq_b[b1]
          }
        }
      }

    }

    n0=0
    n1=0

  }

  return nil
}
