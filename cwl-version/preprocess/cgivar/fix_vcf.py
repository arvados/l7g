#!/usr/bin/env python

import sys

def is_header(line):
    """Check if a line is header."""

    return line.startswith('#')

def has_END(line):
    """Check if a line has the 'END=' tag."""

    return 'END=' in line

# FIELD index
# CHROM 0, POS 1, REF 3, QUAL 5, INFO 7, FORMAT 8, sample 9

def fix_END(line):

    all_fields = line.split('\t')
    INFO = all_fields[7]
    INFO_fields = INFO.split(';')
    for i, INFO_field in enumerate(INFO_fields):
        if INFO_field.split('=')[0] == 'END':
            INFO_fields[i] = INFO_fields[i].replace('.', '')

    all_fields[7] = ';'.join(INFO_fields)
    line = '\t'.join(all_fields)

    return line

if __name__ == '__main__':
    vcf = sys.argv[1]
    with open(vcf) as g:
        for line in g:
            if is_header(line):
                print line.strip()
            elif has_END(line):
                print fix_END(line).strip()
            else:
                print line.strip()
