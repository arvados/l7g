#!/usr/bin/env python

from __future__ import print_function
import subprocess
import os
import argparse
import re

def make_bed_path(path, assembly):
    """Make the bed file of a path"""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    pathdec = int(path, 16)

    try:
        indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(path), assemblyindex])
    except subprocess.CalledProcessError:
        raise Exception("No such path as {}".format(path))
    chrom = indexline.split('\t')[0].split(':')[1]
    length = indexline.split('\t')[1]
    offset = indexline.split('\t')[2]

    ps = subprocess.Popen(["bgzip", "-b", offset, "-s", length, assembly],
                      stdout=subprocess.PIPE)
    assemblyline = subprocess.check_output(["tail", "-n", "1"], stdin=ps.stdout)
    ps.wait()
    end = int(assemblyline.split('\t')[1])

    if pathdec == 0:
        start = 0
    else:
        previous_pathdec = pathdec - 1
        previous_path = format(previous_pathdec, '04x')

        indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(previous_path), assemblyindex])
        previous_chrom = indexline.split('\t')[0].split(':')[1]

        if previous_chrom != chrom:
            start = 0
        else:
            previous_length = indexline.split('\t')[1]
            previous_offset = indexline.split('\t')[2]

            ps = subprocess.Popen(["bgzip", "-b", previous_offset, "-s", previous_length, assembly],
                                  stdout=subprocess.PIPE)
            assemblyline = subprocess.check_output(["tail", "-n", "1"], stdin=ps.stdout)
            ps.wait()
            start = int(assemblyline.split('\t')[1])

    print("{}\t{}\t{}".format(chrom, start, end))

def make_bed_steps(path, stepsfile, assembly):
    """Make the bed file of a list of steps in a path"""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    pathdec = int(path, 16)

    try:
        indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(path), assemblyindex])
    except subprocess.CalledProcessError:
        raise Exception("No such path as {}".format(path))
    chrom = indexline.split('\t')[0].split(':')[1]
    length = indexline.split('\t')[1]
    offset = indexline.split('\t')[2]

    with open(stepsfile) as f:
        for line in f:
            linestrip = line.strip()
            step = linestrip.split('+')[0]
            span = int(linestrip.split('+')[1], 16)
            stepdec = int(step, 16)
            spanningtile_stepdec = stepdec + span - 1
            spanningtile_step = format(spanningtile_stepdec, '04x')

            try:
                ps = subprocess.Popen(["bgzip", "-b", offset, "-s", length, assembly],
                                  stdout=subprocess.PIPE)
                assemblyline = subprocess.check_output(["grep", "-P", "^{}\t".format(spanningtile_step)], stdin=ps.stdout)
                ps.wait()
            except subprocess.CalledProcessError:
                raise Exception("No such step as {} with span {} in path {}".format(step, span, path))
            end = int(assemblyline.split('\t')[1])

            # calculate previous tile to derive start position
            if stepdec != 0:
                previous_stepdec = stepdec - 1
                previous_step = format(previous_stepdec, '04x')

                ps = subprocess.Popen(["bgzip", "-b", offset, "-s", length, assembly],
                                      stdout=subprocess.PIPE)
                assemblyline = subprocess.check_output(["grep", "-P", "^{}\t".format(previous_step)], stdin=ps.stdout)
                ps.wait()
                start = int(assemblyline.split('\t')[1])
            elif pathdec == 0:
                start = 0
            else:
                previous_pathdec = pathdec - 1
                previous_path = format(previous_pathdec, '04x')

                indexline = subprocess.check_output(["grep", "-P", ":{}\t".format(previous_path), assemblyindex])
                previous_chrom = indexline.split('\t')[0].split(':')[1]

                if previous_chrom != chrom:
                    start = 0
                else:
                    previous_length = indexline.split('\t')[1]
                    previous_offset = indexline.split('\t')[2]

                    ps = subprocess.Popen(["bgzip", "-b", previous_offset, "-s", previous_length, assembly],
                                          stdout=subprocess.PIPE)
                    assemblyline = subprocess.check_output(["tail", "-n", "1"], stdin=ps.stdout)
                    ps.wait()
                    start = int(assemblyline.split('\t')[1])

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
