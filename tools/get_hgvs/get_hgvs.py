#!/usr/bin/env python

from __future__ import print_function
from extractor import describe_dna

def fasta_to_hgvs(ref, sample, seqstart, prefix):
    allele = describe_dna(ref, sample)
    for var in allele:
        var.start += seqstart
        var.end += seqstart
    hgvs = "{}{}".format(prefix, allele)
    return hgvs
