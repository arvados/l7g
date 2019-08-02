#!/usr/bin/env python

from __future__ import print_function
import argparse
import os
import subprocess

RUNNER_RAM = "50000"
EVAL_TIMEOUT = "2000"

PATHMIN = 0
PATHMAX = 862
SGLFTHRESHOLD = 4000
CHECKNUM = 2
TAGSET = "cd9ada494bd979a8bc74e6d59d3e8710+174/tagset.fa.gz"
REFFA = {"hg19": "ef70506d71ee761b1ec28d67290216a0+1252/hg19.fa.gz",
         "hg38": "ee5b90cf2d5f3573e6d455ab56e15cdf+761/hg38.fa.gz",
         "human_g1k_v37": "5a42cfaddd3a9cfc4fac89b3fe73c6f6+751/human_g1k_v37.fasta.gz"}
AFN = {"hg19": "98c5e71956730c36cc89bb25b99fe58b+965/assembly.00.hg19.fw.gz",
       "hg38": "7deca98a5827e1991bf49a96a0087318+233/assembly.00.hg38.fw.gz",
       "human_g1k_v37": "96fe7d3fdc5b0bd82128131a23117635+269/assembly.00.human_g1k_v37.fw.gz"}

def make_yml_and_run(project_uuid, varstype, inputpdh, ref, chr1, chrM, nchunks):
    yml_text = '''%sdir:
  class: Directory
  location: keep:%s\n''' % (varstype, inputpdh)
    yml_text += 'ref: "%s"\n' % ref
    yml_text += '''reffa:
  class: File
  location: keep:%s
afn:
  class: File
  location: keep:%s
tagset:
  class: File
  location: keep:%s\n''' % (REFFA[ref], AFN[ref], TAGSET)
    chroms_prefix = chr1.replace("1", "")
    chroms_list = ["\"" + chroms_prefix + str(c) + "\"" for c in range(1, 23)]
    chroms_list.append("\"" + chrM + "\"")
    chroms = "[" + ", ".join(chroms_list) + "]"
    yml_text += 'chroms: %s\n' % chroms
    yml_text += '''pathmin: "%d"
pathmax: "%d"
nchunks: "%d"
sglfthreshold: %d
checknum: %d''' % (PATHMIN, PATHMAX, nchunks, SGLFTHRESHOLD, CHECKNUM)

    print("Input yml file:")
    print(yml_text)
    
    yml = "yml/%s_%s.yml" % (varstype, inputpdh)
    with open(yml, 'w') as f:
        f.write(yml_text)
    command = ["arvados-cwl-runner", "--api", "containers",
               "--submit", "--no-wait",
               "--submit-runner-ram", RUNNER_RAM,
               "--eval-timeout", EVAL_TIMEOUT]
    if project_uuid:
        command.extend(["--project-uuid", project_uuid])
    command.extend(["%s2npy-wf.cwl" % varstype, yml])

    print("\nRunning:")
    print(" ".join(command))
    subprocess.check_call(command)
    os.remove(yml)

def main():
    parser = argparse.ArgumentParser(description='Make input yml file and \
        run workflow on arvados to generate npy arrays.')
    parser.add_argument("varstype", choices=['gvcf', 'gff'],
        help='file type of input variant files.')
    parser.add_argument('inputpdh', help='portable data hash of input directory.')
    parser.add_argument('ref', choices=['hg19', 'hg38', 'human_g1k_v37'],
        help='reference name.')
    parser.add_argument('chr1', choices=['chr1', '1'],
        help='chromosome 1 notation, expected to be consistent with chromosome 1-22,X,Y.')
    parser.add_argument('chrM', choices=['chrM', 'M', 'chrMT', 'MT'],
        help='chromosome M notation.')
    parser.add_argument('--project-uuid', help='arvados project-uuid to run workflow.')
    parser.add_argument('--nchunks', type=int, default=15,
        help='number of chunks of tile paths when creating tile library, default is 15.')

    args = parser.parse_args()
    make_yml_and_run(args.project_uuid, args.varstype, args.inputpdh, args.ref, args.chr1, args.chrM, args.nchunks)

if __name__ == '__main__':
    main()
