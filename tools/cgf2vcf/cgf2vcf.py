#!/usr/bin/env python

from __future__ import print_function
from pyfaidx import Fasta
import pyhgvs
import subprocess
import os
import argparse
import re

# NCBI references:
# https://www.ncbi.nlm.nih.gov/assembly/GCF_000001405.26/
# https://www.ncbi.nlm.nih.gov/assembly/GCF_000001405.13/
# https://www.ncbi.nlm.nih.gov/nucleotide/NC_012920.1

counsyl_prefix = {
    "hg38": {
        "NC_000001.11": "chr1:g.",
        "NC_000002.12": "chr2:g.",
        "NC_000003.12": "chr3:g.",
        "NC_000004.12": "chr4:g.",
        "NC_000005.10": "chr5:g.",
        "NC_000006.12": "chr6:g.",
        "NC_000007.14": "chr7:g.",
        "NC_000008.11": "chr8:g.",
        "NC_000009.12": "chr9:g.",
        "NC_000010.11": "chr10:g.",
        "NC_000011.10": "chr11:g.",
        "NC_000012.12": "chr12:g.",
        "NC_000013.11": "chr13:g.",
        "NC_000014.9": "chr14:g.",
        "NC_000015.10": "chr15:g.",
        "NC_000016.10": "chr16:g.",
        "NC_000017.11": "chr17:g.",
        "NC_000018.10": "chr18:g.",
        "NC_000019.10": "chr19:g.",
        "NC_000020.11": "chr20:g.",
        "NC_000021.9": "chr21:g.",
        "NC_000022.11": "chr22:g.",
        "NC_000023.11": "chrX:g.",
        "NC_000024.10": "chrY:g.",
        "NC_012920.1": "chrM:g."},
    "human_g1k_v37": {
        "NC_000001.10": "1:g.",
        "NC_000002.11": "2:g.",
        "NC_000003.11": "3:g.",
        "NC_000004.11": "4:g.",
        "NC_000005.9": "5:g.",
        "NC_000006.11": "6:g.",
        "NC_000007.13": "7:g.",
        "NC_000008.10": "8:g.",
        "NC_000009.11": "9:g.",
        "NC_000010.10": "10:g.",
        "NC_000011.9": "11:g.",
        "NC_000012.11": "12:g.",
        "NC_000013.10": "13:g.",
        "NC_000014.8": "14:g.",
        "NC_000015.9": "15:g.",
        "NC_000016.9": "16:g.",
        "NC_000017.10": "17:g.",
        "NC_000018.9": "18:g.",
        "NC_000019.9": "19:g.",
        "NC_000020.10": "20:g.",
        "NC_000021.8": "21:g.",
        "NC_000022.10": "22:g.",
        "NC_000023.10": "X:g.",
        "NC_000024.9": "Y:g.",
        "NC_012920.1": "M:g."}}

def parse_band(bandtext):
    """Parse band text to a list of two lists."""
    bandlines = bandtext.split('\n')[:2]
    band = []
    for bandline in bandlines:
        bandstr = bandline.replace('[', '').replace(']', '').strip().split(' ')
        bandsingle = map(int, bandstr)
        band.append(bandsingle)
    return band

def parse_hgvs(hgvs, refname):
    """Parse HGVS to a format admissible to the counsyl hgvs tool."""
    prefix = hgvs.split(':')[0]
    newprefix = counsyl_prefix[refname][prefix]
    mutation = hgvs.split(':')[1].split('.')[1].replace('[', '').replace(']', '')
    mutations = mutation.split(';')
    counsylhgvs_list = ["{}{}".format(newprefix, m) for m in mutations if m != '=']
    return counsylhgvs_list

def make_vcfblock(haplotypes):
    """Make vcf block of haplotypes in the form of a list of two haplotypes."""
    vcfblock = ""
    i0 = 0
    i1 = 0
    len0 = len(haplotypes[0])
    len1 = len(haplotypes[1])

    while i0 < len0 or i1 < len1:
        if i1 == len1:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|0\n".format(haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[0][i0][3])
            i0 += 1
        elif i0 == len0:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t0|1\n".format(haplotypes[1][i1][0], haplotypes[1][i1][1], haplotypes[1][i1][2], haplotypes[1][i1][3])
            i1 += 1
        elif haplotypes[0][i0][1] < haplotypes[1][i1][1]:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|0\n".format(haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[0][i0][3])
            i0 += 1
        elif haplotypes[0][i0][1] > haplotypes[1][i1][1]:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t0|1\n".format(haplotypes[1][i1][0], haplotypes[1][i1][1], haplotypes[1][i1][2], haplotypes[1][i1][3])
            i1 += 1
        else:
            if haplotypes[0][i0][3] == haplotypes[1][i1][3]:
                vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|1\n".format(haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[0][i0][3])
            else:
                vcfline = "{}\t{}\t.\t{}\t{},{}\t.\t.\t.\tGT\t1|2\n".format(haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[0][i0][3], haplotypes[1][i1][3])
            i0 += 1
            i1 += 1
        vcfblock += vcfline
    return vcfblock

def hgvs_to_haplotype(hgvs, refname, genome):
    """Use counsyl hgvs tool to convert HGVS to haplotype."""
    counsylhgvs_list = parse_hgvs(hgvs, refname)
    haplotype = []
    for counsylhgvs in counsylhgvs_list:
        singlehap = pyhgvs.parse_hgvs_name(counsylhgvs, genome)
        haplotype.append(singlehap)
    return haplotype

def get_vcflines(band, hgvstext, path, ref):
    """Given the HGVS text, get the vcf lines of a band, along with uncalled steps and unannotated steps."""
    refname = os.path.basename(ref).split('.')[0]
    genome = Fasta(ref)
    out = {"uncalled": "",
           "unannotated": ""}

    pathlen = len(band[0])
    blockstart_stepdec = None
    for stepdec in range(pathlen):
        step = format(stepdec, '04x')
        # this is when a block starts
        if band[0][stepdec] != -1 and band[1][stepdec] != -1:
            if blockstart_stepdec != None:
                # reporting previous block
                span = stepdec - blockstart_stepdec
                stepoutput = "{}+{}\n".format(format(blockstart_stepdec, '04x'), span)
                if is_uncalled:
                    out["uncalled"] += stepoutput
                elif is_unannotated:
                    out["unannotated"] += stepoutput
                else:
                    print(vcfblock, end = '')

            is_uncalled = (band[0][stepdec] == -2 or band[1][stepdec] == -2)
            if not is_uncalled:
                # determine whether the tile variants are in the annotated library
                pattern0 = r'{}\..*\.{}\.{}\+.*'.format(path, step, format(band[0][stepdec], '03x'))
                pattern1 = r'{}\..*\.{}\.{}\+.*'.format(path, step, format(band[1][stepdec], '03x'))
                match0 = re.search(pattern0, hgvstext)
                match1 = re.search(pattern1, hgvstext)
                is_unannotated = not (match0 and match1)
                if not is_unannotated:
                    hgvs0 = match0.group().split(',')[-1]
                    hgvs1 = match1.group().split(',')[-1]
                    haplotype0 = hgvs_to_haplotype(hgvs0, refname, genome)
                    haplotype1 = hgvs_to_haplotype(hgvs1, refname, genome)
                    haplotypes = [haplotype0, haplotype1]
                    vcfblock = make_vcfblock(haplotypes)

            blockstart_stepdec = stepdec
        else:
            if not is_uncalled:
                # update whether the block is uncalled
                is_uncalled = (band[0][stepdec] == -2 or band[1][stepdec] == -2)
            if not is_uncalled:
                if not is_unannotated:
                    # update whether the block is unannotated
                    if band[0][stepdec] != -1 or band[1][stepdec] != -1:
                        if band[0][stepdec] != -1:
                            pattern = r'{}\..*\.{}\.{}\+.*'.format(path, step, format(band[0][stepdec], '03x'))
                        else:
                            pattern = r'{}\..*\.{}\.{}\+.*'.format(path, step, format(band[1][stepdec], '03x'))
                        match = re.search(pattern, hgvstext)
                        is_unannotated = not match
                        if not is_unannotated:
                            hgvs = match.group().split(',')[-1]
                            haplotype = hgvs_to_haplotype(hgvs, refname, genome)
                            if band[0][stepdec] != -1:
                                haplotypes = [haplotype, []]
                            else:
                                haplotypes = [[], haplotype]
                            vcfblock = make_vcfblock(haplotypes)
    else:
        # reporting the last block
        span = stepdec - blockstart_stepdec
        stepoutput = "{}+{}\n".format(format(blockstart_stepdec, '04x'), span)
        if is_uncalled:
            out["uncalled"] += stepoutput
        elif is_unannotated:
            out["unannotated"] += stepoutput
        else:
            print(vcfblock, end = '')

    return out

def main():
    parser = argparse.ArgumentParser(description='Output vcf lines of a cgf band\
        in a given path, given an annotated tile library.')
    parser.add_argument('path', metavar='PATH', help='tile path')
    parser.add_argument('ref', metavar='REF', help='reference fasta file')
    parser.add_argument('hgvs', metavar='HGVS', help='HGVS annotation of a tile library')
    parser.add_argument('cgf', metavar='CGF', help='CGF file')

    parser.add_argument('--uncalled', help='output file of uncalled steps')
    parser.add_argument('--unannotated', help='output file of unannotated steps')

    args = parser.parse_args()

    bandtext = subprocess.check_output(["cgft", "-q", "-b", args.path, "-i", args.cgf])
    band = parse_band(bandtext)
    with open(args.hgvs) as f:
        hgvstext = f.read()

    out = get_vcflines(band, hgvstext, args.path, args.ref)
    if args.uncalled:
        with open(args.uncalled, 'w') as f:
            f.write(out["uncalled"])
    if args.unannotated:
        with open(args.unannotated, 'w') as f:
            f.write(out["unannotated"])

if __name__ == '__main__':
    main()
