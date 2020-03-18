#!/usr/bin/env python

from __future__ import print_function
from extractor import describe_dna
import collections
import subprocess
import os
import argparse

Window = collections.namedtuple('Window', ['chrom', 'start', 'end'])

# NCBI references:
# https://www.ncbi.nlm.nih.gov/assembly/GCF_000001405.26/
# https://www.ncbi.nlm.nih.gov/assembly/GCF_000001405.13/
# https://www.ncbi.nlm.nih.gov/nucleotide/NC_012920.1

ncbi_prefix = {
    "hg38": {
        "chr1": "NC_000001.11:g.",
        "chr2": "NC_000002.12:g",
        "chr3": "NC_000003.12:g.",
        "chr4": "NC_000004.12:g.",
        "chr5": "NC_000005.10:g.",
        "chr6": "NC_000006.12:g.",
        "chr7": "NC_000007.14:g.",
        "chr8": "NC_000008.11:g.",
        "chr9": "NC_000009.12:g.",
        "chr10": "NC_000010.11:g.",
        "chr11": "NC_000011.10:g.",
        "chr12": "NC_000012.12:g.",
        "chr13": "NC_000013.11:g.",
        "chr14": "NC_000014.9:g.",
        "chr15": "NC_000015.10:g.",
        "chr16": "NC_000016.10:g.",
        "chr17": "NC_000017.11:g.",
        "chr18": "NC_000018.10:g.",
        "chr19": "NC_000019.10:g.",
        "chr20": "NC_000020.11:g.",
        "chr21": "NC_000021.9:g.",
        "chr22": "NC_000022.11:g.",
        "chrX": "NC_000023.11:g.",
        "chrY": "NC_000024.10:g.",
        "chrM": "NC_012920.1:m."},
    "human_g1k_v37": {
        "1": "NC_000001.10:g.",
        "2": "NC_000002.11:g",
        "3": "NC_000003.11:g.",
        "4": "NC_000004.11:g.",
        "5": "NC_000005.9:g.",
        "6": "NC_000006.11:g.",
        "7": "NC_000007.13:g.",
        "8": "NC_000008.10:g.",
        "9": "NC_000009.11:g.",
        "10": "NC_000010.10:g.",
        "11": "NC_000011.9:g.",
        "12": "NC_000012.11:g.",
        "13": "NC_000013.10:g.",
        "14": "NC_000014.8:g.",
        "15": "NC_000015.9:g.",
        "16": "NC_000016.9:g.",
        "17": "NC_000017.10:g.",
        "18": "NC_000018.9:g.",
        "19": "NC_000019.9:g.",
        "20": "NC_000020.10:g.",
        "21": "NC_000021.8:g.",
        "22": "NC_000022.10:g.",
        "X": "NC_000023.10:g.",
        "Y": "NC_000024.9:g.",
        "M": "NC_012920.1:m."}}

def fasta_to_hgvs(ref, sample, seqstart, prefix):
    """Get HGVS using mutalyzer description-extractor."""
    allele = describe_dna(ref, sample)
    for var in allele:
        var.start += seqstart
        var.end += seqstart
    hgvs = "{}{}".format(prefix, allele)
    return hgvs

def get_tile_window(path, step, assembly, span, taglen):
    """Derive tile window."""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    pathdec = int(path, 16)
    stepdec = int(step, 16)

    try:
        indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(path), assemblyindex])
    except subprocess.CalledProcessError:
        raise Exception("No such path as {}".format(path))
    chrom = indexline.split('\t')[0].split(':')[1]
    length = indexline.split('\t')[1]
    offset = indexline.split('\t')[2]

    spanningtile_stepdec = stepdec + span - 1
    spanningtile_step = format(spanningtile_stepdec, '04x')

    try:
        ps = subprocess.Popen(["bgzip", "-b", offset, "-s", length, assembly],
                          stdout=subprocess.PIPE)
        assemblyline = subprocess.check_output(["grep", "-P", "^{}\t".format(spanningtile_step)], stdin=ps.stdout)
        ps.wait()
    except subprocess.CalledProcessError:
        raise Exception("No such step as {} with span {} in path {}".format(step, span, path))
    end = int(assemblyline.split('\t')[1]) + 1

    # calculate previous tile to derive start position
    if stepdec != 0:
        previous_stepdec = stepdec - 1
        previous_step = format(previous_stepdec, '04x')

        ps = subprocess.Popen(["bgzip", "-b", offset, "-s", length, assembly],
                              stdout=subprocess.PIPE)
        assemblyline = subprocess.check_output(["grep", "-P", "^{}\t".format(previous_step)], stdin=ps.stdout)
        ps.wait()
        start = int(assemblyline.split('\t')[1]) + 1 - taglen
    elif pathdec == 0:
        start = 1
    else:
        previous_pathdec = pathdec - 1
        previous_path = format(previous_pathdec, '04x')

        indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(previous_path), assemblyindex])
        previous_chrom = indexline.split('\t')[0].split(':')[1]

        if previous_chrom != chrom:
            start = 1
        else:
            previous_length = indexline.split('\t')[1]
            previous_offset = indexline.split('\t')[2]

            ps = subprocess.Popen(["bgzip", "-b", previous_offset, "-s", previous_length, assembly],
                                  stdout=subprocess.PIPE)
            assemblyline = subprocess.check_output(["tail", "-n", "1"], stdin=ps.stdout)
            ps.wait()
            start = int(assemblyline.split('\t')[1]) + 1 - taglen

    return Window(chrom, start, end)

def annotate_tilelib(path, tagver, step, ref, tilelib, tilevars, assembly, taglen):
    """Annotate given tile variants."""
    refname = os.path.basename(ref).split('.')[0]
    if refname not in ncbi_prefix:
        raise Exception("No such reference as {}".format(refname))

    sglf = os.path.join(tilelib, "{}.sglf.gz".format(path))
    sglflines = []
    for tilevar in tilevars:
        try:
            sglfline = subprocess.check_output(["zgrep", "{}.{}.{}.{}+".format(path, tagver, step, tilevar), sglf]).strip()
        except subprocess.CalledProcessError:
            # skip tile variants not found in the library
            continue
        # skip empty tile variants too
        if sglfline.split(',')[-1] != "":
            sglflines.append(sglfline)
    spanset = set([sglfline.split(',')[0].split('+')[1] for sglfline in sglflines])

    # store ref fastas with given span
    windowdict = {}
    reffastadict = {}
    for span in spanset:
        window = get_tile_window(path, step, assembly, int(span), taglen)
        windowdict[span] = window

        rawreffasta = subprocess.check_output(["samtools", "faidx", ref, "{}:{}-{}".format(window.chrom, window.start, window.end)])
        reffasta = ''.join(rawreffasta.split('\n')[1:]).lower()
        reffastadict[span] = reffasta

    # derive HGVS
    for sglfline in sglflines:
        span = sglfline.split(',')[0].split('+')[1]
        window = windowdict[span]
        samplefasta = sglfline.split(',')[2].lower()
        prefix = ncbi_prefix[refname][window.chrom]
        hgvs = fasta_to_hgvs(reffastadict[span], samplefasta, window.start, prefix)
        annotationline = ','.join(sglfline.split(',')[:-1] + [hgvs])
        print(annotationline)

def main():
    parser = argparse.ArgumentParser(description='Output HGVS annotations\
        of tile variants.')
    parser.add_argument('path', metavar='PATH', help='tile path')
    parser.add_argument('step', metavar='STEP', help='tile step')
    parser.add_argument('ref', metavar='REF', help='reference fasta file')
    parser.add_argument('tilelib', metavar='TILELIB', help='tile library directory')
    parser.add_argument('varnum', metavar='VARNUM', type=int, help='the number of tile variants to be annotated,\
        only the first VARNUM variants in the given position are considered')
    parser.add_argument('assembly', metavar='ASSEMBLY', help='assembly file')

    parser.add_argument('--tagver', type=str, default="00",
        help='tag version, default is "00".')
    parser.add_argument('--taglen', type=int, default=24,
        help='tag length, default is 24.')

    args = parser.parse_args()
    tilevars = [format(i, '03x') for i in range(args.varnum)]
    annotate_tilelib(args.path, args.tagver, args.step, args.ref, args.tilelib, tilevars, args.assembly, args.taglen)

if __name__ == '__main__':
    main()
