#!/usr/bin/env python

from __future__ import print_function
import collections
import subprocess
import os
import argparse
import re
import gzip

Window = collections.namedtuple('Window', ['start', 'end'])

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

def make_refdict(ref):
    """Make reference dictionary from reference fasta file."""
    refdict = {}
    chrom = None
    with gzip.open(ref) as f:
        for line in f:
            if line.startswith('>'):
                if chrom != None:
                    refdict[chrom] = ''.join(fasta_list)
                chrom = line.strip().split(' ')[0].replace('>', '')
                fasta_list = []
            else:
                fasta_list.append(line.strip())
    return refdict

def fasta_to_hgvs(ref, sample, seqstart, bashscript):
    """Get HGVS using diff-fasta."""
    return subprocess.check_output([bashscript, ref, sample, str(seqstart-1)])

def get_tile_window(path, step, span, pathassembly, pathstart, taglen):
    """Derive tile window."""
    pathdec = int(path, 16)
    stepdec = int(step, 16)

    spanningtile_stepdec = stepdec + span - 1
    spanningtile_step = format(spanningtile_stepdec, '04x')

    pattern = re.compile(r'^{}\t.*'.format(spanningtile_step), re.MULTILINE)
    match = re.search(pattern, pathassembly)
    if match:
        end = int(match.group().split('\t')[1].strip())
    else:
        raise Exception("No such step as {} with span {} in path {}".format(step, span, path))

    # calculate previous tile to derive start position
    if stepdec != 0:
        previous_stepdec = stepdec - 1
        previous_step = format(previous_stepdec, '04x')
        pattern = re.compile(r'^{}\t.*'.format(previous_step), re.MULTILINE)
        match = re.search(pattern, pathassembly)
        start = int(match.group().split('\t')[1].strip()) + 1 - taglen
    else:
        start = pathstart

    return Window(start, end)

def annotate_tilelib(path, ref, tilelib, assembly, bashscript, taglen, samplefastalimit=None):
    """Annotate given tile variants."""
    pathdec = int(path, 16)
    refname = os.path.basename(ref).split('.')[0]
    if refname not in ncbi_prefix:
        raise Exception("Reference {} not supported".format(refname))
    refdict = make_refdict(ref)

    # prepare pathassembly and pathstart for get_tile_window
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    with open(assemblyindex) as f:
        assemblyindextext = f.read()
    pattern = r'.*:{}\t.*'.format(path)
    match = re.search(pattern, assemblyindextext)
    if match:
        indexline = match.group()
    else:
        raise Exception("No such path as {}".format(path))
    chrom = indexline.split('\t')[0].split(':')[1]
    size = indexline.split('\t')[1]
    offset = indexline.split('\t')[2]
    pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", size, assembly]).strip()
    if pathdec == 0:
        pathstart = 1
    else:
        previous_path = format(pathdec-1, '04x')
        pattern = r'.*:{}\t.*'.format(previous_path)
        match = re.search(pattern, assemblyindextext)
        previous_indexline = match.group()
        previous_chrom = previous_indexline.split('\t')[0].split(':')[1]
        if chrom != previous_chrom:
            pathstart = 1
        else:
            previous_size = previous_indexline.split('\t')[1]
            previous_offset = previous_indexline.split('\t')[2]
            previous_pathassembly = subprocess.check_output(["bgzip", "-b",
                            previous_offset, "-s", previous_size, assembly]).strip()
            lastline = previous_pathassembly.split('\n')[-1]
            pathstart = int(lastline.split('\t')[1].strip()) + 1

    sglf = os.path.join(tilelib, "{}.sglf.gz".format(path))
    with gzip.open(sglf) as f:
        for sglfline in f:
            step = sglfline.split(',')[0].split('.')[2]
            span = int(sglfline.split(',')[0].split('+')[1], 16)
            try:
                window = get_tile_window(path, step, span, pathassembly, pathstart, taglen)
                reffasta = refdict[chrom][window.start-1:window.end]
                samplefasta = sglfline.split(',')[2].strip()
                if samplefastalimit != None and len(samplefasta) >= samplefastalimit:
                    continue
                if reffasta != "" and samplefasta != "":
                    annotationline = ','.join(sglfline.split(',')[:-1])
                    print(annotationline,
                        fasta_to_hgvs(reffasta, samplefasta, window.start, bashscript), sep=',', end='')
            except:
                continue

def main():
    parser = argparse.ArgumentParser(description='Output HGVS annotations\
        of tile variants.')
    parser.add_argument('path', metavar='PATH', help='tile path')
    parser.add_argument('ref', metavar='REF', help='reference fasta file')
    parser.add_argument('tilelib', metavar='TILELIB', help='tile library directory')
    parser.add_argument('assembly', metavar='ASSEMBLY', help='assembly file')
    parser.add_argument('bashscript', metavar='BASHSCRIPT', help='bashscript for diff-fasta')

    parser.add_argument('--taglen', type=int, default=24,
        help='tag length, default is 24.')
    parser.add_argument('--samplefastalimit', type=int,
        help='sample fasta limit, skip if exceeds.')

    args = parser.parse_args()
    annotate_tilelib(args.path, args.ref, args.tilelib, args.assembly, args.bashscript, args.taglen, args.samplefastalimit)

if __name__ == '__main__':
    main()
