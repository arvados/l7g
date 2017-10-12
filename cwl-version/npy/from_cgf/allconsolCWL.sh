#!/bin/bash

#export indir='npy'
#export outdir='outdir'
#export outprefix='all'
#export npyconsolfile='./npy-consolidate'

indir=$1
outdir=$2
outprefix=$3
npyconsolfile=$4

mkdir -p outdir

filearray=($(ls -f $indir/[0-9]*))
export filearray
$npyconsolfile ${filearray[@]}  $outdir'/'$outprefix

