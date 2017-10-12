#!/bin/bash

set -eo pipefail

export bashscript="$1"
export fastj2cgflib="$2"
export datadir="$3"
export verbose_tagset="$4"
export tagset="$5"
export numthreads="$6"

#export bashscript="/data-sdd/cwl_tiling/tilelib/tilelibCWL.sh"
#export fastj2cgflib="/data-sdd/cwl_tiling/tilelib/fastj2cgflib"
#export datadir="/data-sdd/cwl_tiling/tilelib/data"
#export verbose_tagset="/data-sdd/cwl_tiling/tilelib/verbose_tagset"
#export tagset="/data-sdd/data/l7g/tagset.fa/tagset.fa.gz"
#export numthreads="15"

#echo {0..862} | xargs -n1 -P 15 -I{} bash -c " time ./run_cglf_single.sh {}

printf '%04x\n' {0..862} | xargs -n1 -P $numthreads -I{} $bashscript {} lib $fastj2cgflib $datadir $verbose_tagset $tagset 

#printf '%04x\n' {0..862} | xargs -n1 -P 15 -I{} $bashscript {} lib $fastj2cgflib $datadir $verbose_tagset $tagset

#printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} ./tilelibCWL.sh {} lib /data-sdd/cwl_tiling/tilelib/fastj2cgflib /data-sdd/cwl_tiling/tilelib/data /data-sdd/cwl_tiling/tilelib/verbose_tagset /data-sdd/data/l7g/tagset.fa/tagset.fa.gz 

#printf '%04x\n' {0..862} | xargs -n1 -P 15 -I{} time ./run_cglf_single.sh {}
#printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} echo ">>" {}
