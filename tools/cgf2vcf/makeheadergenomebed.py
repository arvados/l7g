#!/usr/bin/env python

from __future__ import print_function
import collections
import subprocess
import os
import argparse
import sys

ChromEnd = collections.namedtuple('ChromEnd', ['chrom', 'end'])

def get_chromends(assembly):
    """Get the ChromEnd list of a genome based on the assembly file"""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    chromends = []

    with open(assemblyindex) as f:
        previous_chrom = None
        for line in f:
            chrom = line.split('\t')[0].split(':')[1]
            if previous_chrom != None and previous_chrom != chrom:
                fields = previous_line.split('\t')
                length = fields[1]
                offset = fields[2]
                pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
                assemblylines = pathassembly.split('\n')
                end = int(assemblylines[-1].split('\t')[1])
                chromends.append(ChromEnd(previous_chrom, end))
            previous_line = line
            previous_chrom = chrom
        else:
            fields = previous_line.split('\t')
            length = fields[1]
            offset = fields[2]
            pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
            assemblylines = pathassembly.split('\n')
            end = int(assemblylines[-1].split('\t')[1].strip())
            chromends.append(ChromEnd(previous_chrom, end))
    return chromends

def make_header(chromends, sampleid):
    print("##fileformat=VCFv4.2")
    for chromend in chromends:
        print("##contig=<ID={},length={}>".format(chromend.chrom, chromend.end))
    print("#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT\t{}".format(sampleid))

def make_genomebed(chromends):
    for chromend in chromends:
        print("{}\t{}\t{}".format(chromend.chrom, 0, chromend.end))

def main():
    parser = argparse.ArgumentParser(description='Output the vcf header or\
        the whole genome bed.')
    parser.add_argument('assembly', metavar='ASSEMBLY', help='assembly file')

    parser.add_argument('--outtype', choices=['header', 'genomebed'],
        required=True, help='output file type')
    parser.add_argument('--sampleid', help='sample id')
    args = parser.parse_args()

    chromends = get_chromends(args.assembly)
    if args.outtype == 'header':
        if args.sampleid is None:
            raise argparse.ArgumentTypeError('sampleid is required when outtype is header')
        else:
            make_header(chromends, args.sampleid)
    else:
        make_genomebed(chromends)

if __name__ == '__main__':
    main()
