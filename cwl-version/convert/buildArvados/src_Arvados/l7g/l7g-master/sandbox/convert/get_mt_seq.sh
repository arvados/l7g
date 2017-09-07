#!/bin/bash

okg_ref="/data-sde/data/ref/human_g1k_v37.fasta.gz"
gff="/data-sde/data/pgp-gff/hu089792-GS02269-DNA_B02.gff.gz"

pasta -r <( samtools faidx $okg_ref MT | egrep -v '^>' | tr '[:upper:]' '[:lower:]' ) \
  -i <( tabix /data-sde/data/pgp-gff/hu089792-GS02269-DNA_B02.gff.gz chrM ) \
  -a gff-pasta

