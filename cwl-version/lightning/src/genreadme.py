#!/usr/bin/env python

from __future__ import print_function
import os
import sys

text = """h1. Data release readme

Data release candidate: {}
Description: This dataset contains {} human genomes ({}) encoded using the Lightning tiling system for the AI4AD project. It is published at {}. 

Collection contains:
* library_full/ -- Full Tiled Data Set
** matrix.0000.npy, matrix.0001.npy, matrix.0002.npy, ... -- tile variant# for each (sample, tag)
** chunk-tag-offset.csv -- tag offset for each matrix.NNNN.npy file
** samples.csv --  sample ID for each row of matrix.NNNN.npy
* library_filtered/ -- Filtered Tiled Data Set (filtered using chi-square test between tile variants and AD phenotype)
** onehot.npy -- one-hot representation of tiled data filtered by p-value
** onehot-columns.npy -- tag, variant, het/hom, p-value for each column of onehot.npy
** samples.csv -- sample ID for each row of onehot.npy
* GRCh38.86_library_annotation/ -- Annotations for Tiled Data Set
** GRCh38.86_library_snpeff_dbsnp_gnomad.vcf.gz -- annotations for each genomic variant found in tiled dataset
** GRCh38.86_library_snpeff_dbsnp_gnomad.vcf.gz.tbi -- index for annotations vcf
** GRCh38.86_library_summary.txt -- % of variants in each chromosome that were found in gnomad
** hg38.fa.gz.bed -- position of tile set in reference genome

Tiling Background:

Tiling abstracts a called genome by partitioning it into overlapping variable length shorter sequences, known as tiles. A tile is a genomic sequence that is braced on either side by 24 base (24-mer) "tags".

A tile sequence must be at least 248 base pairs long where each tile is labeled with a "position" according to the number of tiles before it. One tile position can have multiple tile variants, one for each sequence observed at that position. When a variation occurs on a tag, we allow tile variants to span multiple steps where the tags would normally end. These tiles that span multiple steps are known as "spanning tiles"

Our choice of tags ("tag-set") partition the human reference genome into 10,655,006 tiles, composed of 3.1 billion bases (with an average of around 315 bases per tile). The set of all positions and tile variants are stored in is what we call the tile library. An individual's genome can then be easily represented as an array of tag sets that reference tiles in the tile library. Each position in the array corresponds to a tile position and points to the tile variant observed at that position for that individual.

To create the tiled genomes, we use Lightning, a system that allows for efficient access to large scale population genomic data with a focus on clinical and research use. The Lightning system is a combination of a conceptual way to think about genomes (genomic tiling), the internal representation of genomes for efficient access, and the software that manages access to the data.

h2. Read me for library_full

Directory:  library_full/

Files:

* matrix.XXX.npy:  numpy-encoded matrix with one row per genome, and a pair of columns per tag / tile position (one for each allele). Each matrix element is an integer. For easier loading, the numpy matrix is broken into chunks. : 
** -1 indicates a "low quality" tile variant containing no-calls.
** 0 indicates the tag for this tile was not found, i.e., this part of the genome is covered by a spanning tile in an earlier (leftward) column.
** Tile variants can span multiple tile positions  if a tag is not found and are known as spanning tiles. 
** 1 indicates the most common high quality variant of this tile in this dataset; 2 indicates the 2nd most common; etc.

* chunk-tag-offset.csv - common separated text file that indicates tag offset for each matrix.NNNN.npy file
** Columns are file name and offset

* samples.csv: mapping from numpy file (matrix.npy) and row number to input ID for each tiled genome
** Columns are row number, genome ID (usually taken from tile name of gvcf/vcf, and name of npy output
        - Example: 0,"A-WCAP-WC000711-BL-COL-39141BL1","matrix.npy"


h2. Read me for library_filtered

Directory: library_filtered/
Files:
* onehot.npy -- 
**  The tile variants have been filtered using a chi2 filter between each separate tile variant and the AD phenotype. Only tile positions with 90% coverage are included (i.e. 90% of the tile variants in a tile position do not contain no-calls).  
** Contains the positions of the non-zeros elements of the filtered sparse matrix.: two rows: 1) row position 2) column position
** This sparse numpy-encoded matrix has one row per genome, and a pair of columns per tile variant. One column represents the heterozygous tile variant (i.e. tile variant found in 1 allele) and one for homozygous tile variant (i.e. tile variant found in 2 alleles). Each matrix element is an integer with a 1 indicating the tile variant is present in that form and a 0 indicating the tile variant is not present in that format.
** Can create a sparse matrix with the following commands in python:

import numpy as np
from scipy.sparse import csr_matrix

Xrc = np.load('onehot.npy')
data = np.ones(Xrc[0,:].shape)
row_ind = Xrc[0,:]
col_ind = Xrc[1,:]
filtered = csr_matrix((data, (row_ind, col_ind)))
    
* onehot-columns.npy -
numpy file containing information corresponding to each column of the one-hot matrix representation of the filtered data.
Columns are as follows: tag, tile variant, zygosity with heterozygous = 0 and homozygous = 1, p-value * 1e6 for each column of onehot.npy
* samples.csv -mapping from numpy file (matrix.npy) and row number to input ID for each tiled genome
** Columns are row number, genome ID (usually taken from tile name of gvcf/vcf, and name of npy output
    - Example: 0,"A-WCAP-WC000711-BL-COL-39141BL1","matrix.npy"


h2. Read me for annotations

Directory: GRCh38.86_library_annotation/

Files:
* GRCh38.86_library_snpeff_dbsnp_gnomad.vcf.gz
** gzipped vcf of each genomic variant found in tile variants containing frequencies and other annotation details (gene, predicted effects, etc) from dbsnp and gnmad.  
** ID contains both HGVS and rsID (if found) and INFO contains tile variant (TV: tileposition-tilevariant) as well as the other annotations. All tiles variants contains that genomic variant are listed in the TV field. 
- Example: 
- #CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO
- chr9	45079	chr9:g.45080del;rs55984476	TC	T	.	.	TV=,5649728-1,;ANN=T|intergenic_region|MODIFIER|FAM138C-PGM5P3-AS1|ENSG00000218839-ENSG00000277631|intergenic_region|ENSG00000218839-ENSG00000277631|||n.45080delC||||||;AC=129535;AN=129536;AF=0.999992;AF_afr=0.999966;AF_amr=1;AF_asj=1;AF_eas=1;AF_fin=1;AF_nfe=1;AF_oth=1
** In this annotation file, for simplicity the name of the chromosome is used instead of the proper HGVS annotation for the reference and chromosome. If you want to search the HGVS annotation you will need to replace it. 
        - Example: chr3:g.36130213T>A -> NC_000003.12:g.36130213T>A
        - Example: chr10:g.13511587G>A -> NC_000010.12:g.13511587G>A

* GRCh38.86_library_snpeff_dbsnp_gnomad.vcf.gz.tbi
** index file for vcf of annotations

* GRCh38.86_library_summary.txt 
** text file containing % of variants in each chromosome that were found in gnomad
* GRCh38.86_reference_tiles.bed
** bed file containing tile locations on GRCh38 for reference. 
** The columns are as follows:
** 1) Chromosome
** 2) Tile start (including tag)
** 3) Tile end (including tag)
** 4) Tag #
** 5) Coverage (this gives a score 0-1000 of how many times this tile is placed in a set of genomes, 1000 means the tag is found in every genome of the set. 0 indicates the tag is not found in any of the genomes.  Tag may not be placed due to variants or no-calls existing on the tag. 
** 6) Strand (always ., included so that our bed file maintains the bed standard format)
** 7) Tile start (not including tag)
** 8) Tile end (not including tag
- Example: 
M 0 467 10654109 870 . 0 443
M 443 959 10654110 895 . 467 935
M 935 1394 10654111 985 . 959 1370
"""

def count_samples(samplescsv):
  count = 0
  with open(samplescsv) as f:
    for line in f:
      if line != "\n":
        count += 1
  return count

def main():
  samplescsv = sys.argv[1]
  date = sys.argv[2]
  description = sys.argv[3]
  projecturl = sys.argv[4]

  cohortsize = count_samples(samplescsv)
  print(text.format(date, cohortsize, description, projecturl))

if __name__ == '__main__':
  main()
