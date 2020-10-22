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
        if bandline != "[]":
            bandstr = bandline[2:-1].split(' ')
            bandsingle = map(int, bandstr)
        else:
            bandsingle = []
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
            if REF != ALT:
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

def is_path_ignored(path, assembly):
    """
    Check if a path is ignored.
    We ignore a path if the path's first end position is equal to the previous path's last end position,
    i.e., we do not record the path as 'not covered', simply skip to the next path.
    """
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
    length = fields[1]
    offset = fields[2]
    pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
    assemblylines = pathassembly.split('\n')

    if path == "0000":
        is_ignored = False
    else:
        previous_path = format(int(path, 16) - 1, '04x')
        previous_pattern = r'.*:{}\t.*'.format(previous_path)
        match = re.search(previous_pattern, assemblyindextext)
        previous_indexline = match.group()
        previous_fields = previous_indexline.split('\t')
        previous_length = previous_fields[1]
        previous_offset = previous_fields[2]
        previous_pathassembly = subprocess.check_output(["bgzip", "-b", previous_offset, "-s", previous_length, assembly]).strip()
        previous_assemblylines = previous_pathassembly.split('\n')
        firststep_end_path = int(assemblylines[0].split('\t')[1])
        laststep_end_previous_path = int(previous_assemblylines[-1].split('\t')[1])
        is_ignored = (firststep_end_path == laststep_end_previous_path)
    return is_ignored    

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
    pathlen = len(band[0])
    out = {"nocall": "",
           "unannotated": "",
           "homref": ""}
    annotation = os.path.join(annotationlib, path+".csv")
    if pathlen == 0 or not os.path.isfile(annotation):
        out["pathskipped"] = "True"
    else:
        out["pathskipped"] = "False"
        if not is_path_ignored(path, assembly):
            fields = indexline.split('\t')
            chrom = fields[0].split(':')[1]
            length = fields[1]
            offset = fields[2]
            pathassembly = subprocess.check_output(["bgzip", "-b", offset, "-s", length, assembly]).strip()
            assemblylines = pathassembly.split('\n')
            stepdec_set = set([int(line.split('\t')[0], 16) for line in assemblylines])
            missing_stepdec_set = set(range(pathlen)) - stepdec_set
            vcfinfodict = get_vcfinfodict(annotation)

            blockstart_stepdec = None
            for stepdec in range(pathlen):
                step = format(stepdec, '04x')
                # this is when a block starts
                if band[0][stepdec] != -1 and band[1][stepdec] != -1 and stepdec-1 not in missing_stepdec_set:
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
                            elif band[0][stepdec] != -1 and band[1][stepdec] != -1:
                                shorttileid0 = '{}.{}'.format(step, format(band[0][stepdec], '03x'))
                                shorttileid1 = '{}.{}'.format(step, format(band[1][stepdec], '03x'))
                                is_unannotated = shorttileid0 not in vcfinfodict or shorttileid1 not in vcfinfodict
                                if not is_unannotated:
                                    vcfinfo0 = vcfinfodict[shorttileid0]
                                    vcfinfo1 = vcfinfodict[shorttileid1]
                                    haplotype0 = vcfinfo_to_haplotype(vcfinfo0)
                                    haplotype1 = vcfinfo_to_haplotype(vcfinfo1)
                                    haplotypes[0].extend(haplotype0)
                                    haplotypes[1].extend(haplotype1)
            else:
                # reporting the last block
                span = stepdec - blockstart_stepdec + 1
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
    parser.add_argument('--pathskipped', help='output file indicating whether\
        the path is skipped when creating cgf')

    args = parser.parse_args()

    bandtext = subprocess.check_output(["cgft", "-q", "-b", str(int(args.path, 16)), "-i", args.cgf])
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
    if args.pathskipped:
        with open(args.pathskipped, 'w') as f:
            f.write(out["pathskipped"])

if __name__ == '__main__':
    main()
