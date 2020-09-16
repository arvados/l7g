#!/usr/bin/env python

from __future__ import print_function
import subprocess
import os
import argparse
import re

def parse_band(bandtext):
    """Parse band text to a list of two lists."""
    bandlines = bandtext.split('\n')[:2]
    band = []
    for bandline in bandlines:
        bandstr = bandline[2:-1].split(' ')
        bandsingle = map(int, bandstr)
        band.append(bandsingle)
    return band

def vcfinfo_to_haplotype(vcfinfo):
    """Parse vcf info provided from annotation to haplotype."""
    haplotype = []
    if vcfinfo != '':
        for v in vcfinfo.split(';'):
            triple = v.split('|')
            POS = int(triple[0])
            REF = triple[1]
            ALT = triple[2]
            haplotype.append((POS, REF, ALT))
    return haplotype

def make_vcfblock(haplotypes, chrom):
    """Make vcf block of haplotypes in the form of a list of two haplotypes."""
    vcfblock = ""
    i0 = 0
    i1 = 0
    len0 = len(haplotypes[0])
    len1 = len(haplotypes[1])

    while i0 < len0 or i1 < len1:
        if i1 == len1:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|0\n".format(chrom, haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2])
            i0 += 1
        elif i0 == len0:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t0|1\n".format(chrom, haplotypes[1][i1][0], haplotypes[1][i1][1], haplotypes[1][i1][2])
            i1 += 1
        elif haplotypes[0][i0][0] < haplotypes[1][i1][0]:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|0\n".format(chrom, haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2])
            i0 += 1
        elif haplotypes[0][i0][0] > haplotypes[1][i1][0]:
            vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t0|1\n".format(chrom, haplotypes[1][i1][0], haplotypes[1][i1][1], haplotypes[1][i1][2])
            i1 += 1
        else:
            if haplotypes[0][i0][1] == haplotypes[1][i1][1]:
                if haplotypes[0][i0][2] == haplotypes[1][i1][2]:
                    vcfline = "{}\t{}\t.\t{}\t{}\t.\t.\t.\tGT\t1|1\n".format(chrom, haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2])
                else:
                    vcfline = "{}\t{}\t.\t{}\t{},{}\t.\t.\t.\tGT\t1|2\n".format(chrom, haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[1][i1][2])
            else:
                if haplotypes[0][i0][1].startswith(haplotypes[1][i1][1]):
                    diff = haplotypes[0][i0][1].replace(haplotypes[1][i1][1] ,'', 1)
                    vcfline = "{}\t{}\t.\t{}\t{},{}\t.\t.\t.\tGT\t1|2\n".format(chrom, haplotypes[0][i0][0], haplotypes[0][i0][1], haplotypes[0][i0][2], haplotypes[1][i1][2]+diff)
                elif haplotypes[1][i1][1].startswith(haplotypes[0][i0][1]):
                    diff = haplotypes[1][i1][1].replace(haplotypes[0][i0][1] ,'', 1)
                    vcfline = "{}\t{}\t.\t{}\t{},{}\t.\t.\t.\tGT\t1|2\n".format(chrom, haplotypes[0][i0][0], haplotypes[1][i1][1], haplotypes[0][i0][2]+diff, haplotypes[1][i1][2])
                else:
                    raise Exception("REF {} and REF {} cannot share the same POST {} on chromosome {}".format(haplotypes[0][i0][1], haplotypes[1][i1][1], haplotypes[0][i0][0], chrom))
            i0 += 1
            i1 += 1
        vcfblock += vcfline
    return vcfblock

def get_vcfinfodict(annotation):
    """Load vcfinfodict with key:value pairs of the form shorttileid:vcfinfo."""
    vcfinfodict = {}
    with open(annotation) as f:
        for line in f:
            fields = line.strip().split(',')
            shorttileid = '.'.join(fields[0].split('+')[0].split('.')[2:4])
            vcfinfo = fields[3]
            vcfinfodict[shorttileid] = vcfinfo
    return vcfinfodict

def get_vcflines(band, annotationlib, path, assembly):
    """Given the HGVS text, get the vcf lines of a band, along with nocall, unannotated, and homref steps."""
    assemblyindex = os.path.splitext(assembly)[0] + ".fwi"
    with open(assemblyindex) as f:
        assemblyindextext = f.read()
    pattern = r'.*:{}\t.*'.format(path)
    match = re.search(pattern, assemblyindextext)
    if match:
        indexline = match.group()
    else:
        raise Exception("No such path as {}".format(path))
    annotation = os.path.join(annotationlib, path+".csv")
    vcfinfodict = get_vcfinfodict(annotation)
    chrom = indexline.split('\t')[0].split(':')[1]
    out = {"nocall": "",
           "unannotated": "",
           "homref": ""}

    pathlen = len(band[0])
    blockstart_stepdec = None
    for stepdec in range(pathlen):
        step = format(stepdec, '04x')
        # this is when a block starts
        if band[0][stepdec] != -1 and band[1][stepdec] != -1:
            if blockstart_stepdec != None:
                # reporting previous block
                span = stepdec - blockstart_stepdec
                stepoutput = "{}+{}\n".format(format(blockstart_stepdec, '04x'), format(span, '01x'))
                if is_nocall:
                    out["nocall"] += stepoutput
                elif is_unannotated:
                    out["unannotated"] += stepoutput
                else:
                    vcfblock = make_vcfblock(haplotypes, chrom)
                    print(vcfblock, end = '')
                    if vcfblock == "":
                        out["homref"] += stepoutput

            is_nocall = (band[0][stepdec] == -2 or band[1][stepdec] == -2)
            if not is_nocall:
                # determine whether the tile variants are in the annotated library
                shorttileid0 = '{}.{}'.format(step, format(band[0][stepdec], '03x'))
                shorttileid1 = '{}.{}'.format(step, format(band[1][stepdec], '03x'))
                is_unannotated = shorttileid0 not in vcfinfodict or shorttileid1 not in vcfinfodict
                if not is_unannotated:
                    vcfinfo0 = vcfinfodict[shorttileid0]
                    vcfinfo1 = vcfinfodict[shorttileid1]
                    haplotype0 = vcfinfo_to_haplotype(vcfinfo0)
                    haplotype1 = vcfinfo_to_haplotype(vcfinfo1)
                    haplotypes = [haplotype0, haplotype1]

            blockstart_stepdec = stepdec
        else:
            if not is_nocall:
                # update whether the block is nocall
                is_nocall = (band[0][stepdec] == -2 or band[1][stepdec] == -2)
            if not is_nocall:
                if not is_unannotated:
                    # update whether the block is unannotated
                    if band[0][stepdec] != -1 or band[1][stepdec] != -1:
                        if band[0][stepdec] != -1:
                            shorttileid = '{}.{}'.format(step, format(band[0][stepdec], '03x'))
                        else:
                            shorttileid = '{}.{}'.format(step, format(band[1][stepdec], '03x'))
                        is_unannotated = shorttileid not in vcfinfodict
                        if not is_unannotated:
                            vcfinfo = vcfinfodict[shorttileid]
                            haplotype = vcfinfo_to_haplotype(vcfinfo)
                            if band[0][stepdec] != -1:
                                haplotypes[0].extend(haplotype)
                            else:
                                haplotypes[1].extend(haplotype)
    else:
        # reporting the last block
        span = stepdec - blockstart_stepdec
        stepoutput = "{}+{}\n".format(format(blockstart_stepdec, '04x'), format(span, '01x'))
        if is_nocall:
            out["nocall"] += stepoutput
        elif is_unannotated:
            out["unannotated"] += stepoutput
        else:
            vcfblock = make_vcfblock(haplotypes, chrom)
            print(vcfblock, end = '')
            if vcfblock == "":
                out["homref"] += stepoutput

    return out

def main():
    parser = argparse.ArgumentParser(description='Output vcf lines of a cgf band\
        in a given path, given an annotated tile library.')
    parser.add_argument('path', metavar='PATH', help='tile path')
    parser.add_argument('assembly', metavar='ASSEMBLY', help='assembly file')
    parser.add_argument('annotationlib', metavar='ANNOTATIONLIB', help='annotation of a tile library')
    parser.add_argument('cgf', metavar='CGF', help='CGF file')

    parser.add_argument('--nocall', help='output file of nocall steps')
    parser.add_argument('--unannotated', help='output file of unannotated steps')
    parser.add_argument('--homref', help='output file of homref steps')

    args = parser.parse_args()

    bandtext = subprocess.check_output(["cgft", "-q", "-b", args.path, "-i", args.cgf])
    band = parse_band(bandtext)

    out = get_vcflines(band, args.annotationlib, args.path, args.assembly)
    if args.nocall:
        with open(args.nocall, 'w') as f:
            f.write(out["nocall"])
    if args.unannotated:
        with open(args.unannotated, 'w') as f:
            f.write(out["unannotated"])
    if args.homref:
        with open(args.homref, 'w') as f:
            f.write(out["homref"])

if __name__ == '__main__':
    main()
