#!/bin/bash

export SHELL='/bin/bash'

lib=$1
threshold=$2

ls $lib/*.sglf.gz | parallel "basename {} .sglf.gz | tr -d '\n' && printf $'\t' && gzip -dc {} | wc -c" | sort > sglfsize.tsv

awk '$2 > a { print $1 }' a=$(($threshold*1024**2)) sglfsize.tsv > skippaths.txt

