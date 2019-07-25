#!/usr/bin/env python

from __future__ import print_function
import sys

def is_header(line):
    """Check if a line is header."""

    return line.startswith('#')

# FIELD index
# CHROM 0, POS 1, REF 3

def main():
    previous_CHROM = ""
    previous_end_POS = 0

    for line in sys.stdin:
        if not is_header(line):
            fields = line.split('\t')
            CHROM = fields[0]
            POS = int(fields[1])
            REF = fields[3]
            if CHROM == previous_CHROM:
                if POS > previous_end_POS:
                    print(line, end='')
                    previous_end_POS = max(previous_end_POS, POS + len(REF) - 1)
            else:
                print(line, end='')
                previous_end_POS = POS + len(REF) - 1
            previous_CHROM = CHROM
        else:
            print(line, end='')

if __name__ == '__main__':
    main()
