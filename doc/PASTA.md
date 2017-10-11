PASTA Format
====

## PASTA - a simple verbose stream oriented format for genomic data.

**PASTA is still in the experimental stage**

To facilitate easy conversion between different variant call formats (e.g. VCF, GFF, etc.),
PASTA was invented to serve as an intermediary format that is both simple
to produce and simple to decode from and to various formats.

Though more verbose than other formats, sometimes ease of use
takes precedence.

PASTA format streams genomic sequences, encoding different
information about the base pair emitted with different characters.
Currently, the format is (for a haploid stream):

    a : ref aligned a
    c : ref aligned c
    g : ref aligned g
    t : ref aligned t
    n : ref aligned n

    A : no-call that has reference 'a' at that position
    C : no-call that has reference 'c' at that position
    G : no-call that has reference 'g' at that position
    T : no-call that has reference 't' at that position

    ~ : aligned substitution from ref 'a' to sub 'c' (e.g. SUB 'c')
    ? : aligned substitution from ref 'a' to sub 'g' (e.g. SUB 'g')
    @ : aligned substitution from ref 'a' to sub 't' (e.g. SUB 't')

    = : aligned substitution from ref 'c' to sub 'a' (e.g. SUB 'a')
    : : aligned substitution from ref 'c' to sub 'g' (e.g. SUB 'g')
    ; : aligned substitution from ref 'c' to sub 't' (e.g. SUB 't')

    # : aligned substitution from ref 'g' to sub 'a' (e.g. SUB 'a')
    & : aligned substitution from ref 'g' to sub 'c' (e.g. SUB 'c')
    % : aligned substitution from ref 'g' to sub 't' (e.g. SUB 't')

    * : aligned substitution from ref 't' to sub 'a' (e.g. SUB 'a') (star/asterisk)
    + : aligned substitution from ref 't' to sub 'c' (e.g. SUB 'c') (plus)
    - : aligned substitution from ref 't' to sub 'g' (e.g. SUB 'g') (dash/minus)

    Q : 'a' insertion
    S : 'c' insertion
    W : 'g' insertion
    d : 't' insertion
    Z : 'n' insertion

    ! : deletion of reference 'a'
    $ : deletion of reference 'c'
    7 : deletion of reference 'g'
    E : deletion of reference 't'
    z : deletion of reference 'n' (does this even happen?)

    ' : aligned substitution from ref 'n' to sub 'a' (single quote)
    " : aligned substitution from ref 'n' to sub 'c' (double quote)
    , : aligned substitution from ref 'n' to sub 'g' (comma)
    _ : aligned substitution from ref 'n' to sub 't' (underscore)
    
    > : Beginning of message.  Currently messages of of type ">R{123}" and ">N{456}".

In some sense, a PASTA stream can be thought of as a decorated FASTA stream.

## Example

Consider the following snippet of a PASTA stream:

    gcatGCATgcat?&#dgcat:&*@7$!Egcat

This could produce a FASTA stream:

    gcatnnnngcatgcatgcatgcatgcat

In words:

    gcat - ref
    GCAT - no call with ref 'gcat'
    gcat - ref
    ?=#d - an INDEL 'allele gcat;ref_allele acg'
    gcat - ref
    :&*@7$!E - an INDEL 'allele gcat;ref_allele cgtagcat'
    gcat - ref

## Message

A control message is used to update state.  A control message starts with a `>` character (greater than, ASCII value 62) followed by the message type, typically a one character code, followed by a block starting with `{` and ending with a `}`.

Here is the list of current control messages:

* `>R{\d+}` - a run of reference (e.g. skip `\d+` bases and update current position)
* `>N{\d+}` - a run of no call (e.g. skip `\d+` bases and update current position)
* `>P{\d+}` - update position
* `>C{.*}` - update chromosome name
* `>#{.*}` - comment

In the case of an `R` message, the reference sequence isn't explicitely provided.  In the case of an interleaved stream, `R` and `N` messages are considered homozygous.

For `C` and `#` messages, the message body must not have an end block terminator (`}`).

For example:

* `>R{10}` - a run of reference that is 10 bases long
* `>N{7}` - a run of no-calls that is 7 bases long
* `gcat>R{3}tacg>N{2}acgt` - would translate to `gcat???tacgnnacgt`, where the `??` should be considered reference.
 
## Notes

* INDELs are not explicitely encoded.  By convention an INDEL is a substitution followed by an
  insertion or deletion and can be deduced from the stream provided.
* PASTA is meant to be stream oriented.  This means to convert from a VCF like format, one
  can feed in a VCF like stream along with a (FASTA) reference stream and produce a PASTA
  stream.  So long as the PASTA stream has full information, converting back to the original
  stream (sans annotations) should be possible.  e.g. 'vcf2pasta -i vcf.file -r ref.file | pasta2vcf -i -'
  should produce an equivalent VCF file.
* Since vanilla VCF has only variants, tools should have options to specify whether the PASTA
  stream should produce 'no-calls' or 'ref' for the positions not explicitely called.

## Closing Remarks

By providing a simple, if verbose, intermediate format, conversion between different
schemes should be a matter of producing a PASTA stream that can then be converted
to another format.

The use case PASTA was envisioned for was to convert between the 'VCF-like' formats (VCF,
gVCF, GFF, etc) to FASTJ or FASTA.

## Future work

Right now the stream is simple.  Some extensions to think about, if PASTA proves to be useful,
are:

* binary format
