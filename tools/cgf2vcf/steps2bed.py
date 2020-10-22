#!/usr/bin/env python

from __future__ import print_function
import collections
import subprocess
import os
import argparse
import re

PathInfo = collections.namedtuple("PathInfo", ["chrom", "firstend", "lastend"])

def get_pathinfo(path, assembly, assemblyindextext):
    """Geth PathInfo of a path from the assembly"""
    pattern = r'.*:{}\t.*'.format(path)
    match = re.search(pattern, assemblyindextext)
    if match:
        indexline = match.group()
    else:
        raise Exception("No such path as {}".format(path))
    fields = indexline.split('\t')
    chrom = fields[0].split(':')[1]
    length = fields[1]
    offset = fields[2]
    pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
    assemblylines = pathassembly.split('\n')
    firstend = int(assemblylines[0].split('\t')[1])
    lastend = int(assemblylines[-1].split('\t')[1])
    return PathInfo(chrom, firstend, lastend)

def make_bed_path(path, assembly):
    """Make the bed file of a path"""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    with open(assemblyindex) as f:
        assemblyindextext = f.read()
    pathinfo = get_pathinfo(path, assembly, assemblyindextext)
    if path == "0000":
        start = 0
    else:
        previous_path = format(int(path, 16) - 1, '04x')
        previous_pathinfo = get_pathinfo(previous_path, assembly, assemblyindextext)
        while (previous_pathinfo.lastend >= pathinfo.firstend and previous_pathinfo.chrom == pathinfo.chrom):
            previous_path = format(int(previous_path, 16) - 1, '04x')
            previous_pathinfo = get_pathinfo(previous_path, assembly, assemblyindextext)
        else:
            if previous_pathinfo.chrom != pathinfo.chrom:
                start = 0
            else:
                start = previous_pathinfo.lastend
    print("{}\t{}\t{}".format(pathinfo.chrom, start, pathinfo.lastend))

def make_bed_steps(path, stepsfile, assembly):
    """Make the bed file of a list of steps in a path"""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    with open(assemblyindex) as f:
        assemblyindextext = f.read()
    pattern = r'.*:{}\t.*'.format(path)
    match = re.search(pattern, assemblyindextext)
    if match:
        indexline = match.group()
    else:
        raise Exception("No such path as {}".format(path))
    fields = indexline.split('\t')
    chrom = fields[0].split(':')[1]
    length = fields[1]
    offset = fields[2]
    pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
    assemblylines = pathassembly.split('\n')

    with open(stepsfile) as f:
        for line in f:
            linestrip = line.strip()
            step = linestrip.split('+')[0]
            span = int(linestrip.split('+')[1], 16)
            spanningtile_step = format(int(step, 16) + span - 1, '04x')
            pattern = re.compile(r'^{}\t.*'.format(spanningtile_step), re.MULTILINE)
            match = re.search(pattern, pathassembly)
            if match:
                end = int(match.group().split('\t')[1].strip())
            else:
                raise Exception("No such step as {} with span {} in path {}".format(step, span, path))

            # calculate previous tile to derive start position
            # calculate previous tile when the step is not the first one in the path
            if step != "0000":
                previous_step = format(int(step, 16) - 1, '04x')
                previous_pattern = re.compile(r'^{}\t.*'.format(previous_step), re.MULTILINE)
                previous_match = re.search(previous_pattern, pathassembly)
                start = int(previous_match.group().split('\t')[1].strip())
            elif path == "0000":
                start = 0
            # calculate previous tile when the step is the first one in the path
            else:
                previous_path = format(int(path, 16) - 1, '04x')
                previous_pathinfo = get_pathinfo(previous_path, assembly, assemblyindextext)
                while (previous_pathinfo.lastend >= end and previous_pathinfo.chrom == chrom):
                    previous_path = format(int(previous_path, 16) - 1, '04x')
                    previous_pathinfo = get_pathinfo(previous_path, assembly, assemblyindextext)
                else:
                    if previous_pathinfo.chrom != chrom:
                        start = 0
                    else:
                        start = previous_pathinfo.lastend
            print("{}\t{}\t{}".format(chrom, start, end))

def main():
    parser = argparse.ArgumentParser(description='Output the bed file\
        of a path, or a sub-region of a path.')
    parser.add_argument('path', metavar='PATH', help='tile path')
    parser.add_argument('assembly', metavar='ASSEMBLY', help='assembly file')

    parser.add_argument('--stepsfile', help='steps file indicating a sub-region\
        of a path, each line in the form of "step+span"')

    args = parser.parse_args()
    if args.stepsfile:
        make_bed_steps(args.path, args.stepsfile, args.assembly)
    else:
        make_bed_path(args.path, args.assembly)

if __name__ == '__main__':
    main()
