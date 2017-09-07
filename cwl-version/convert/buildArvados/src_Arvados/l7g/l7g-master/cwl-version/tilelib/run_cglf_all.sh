#!/bin/bash

set -eo pipefail

#echo {0..862} | xargs -n1 -P 15 -I{} bash -c " time ./run_cglf_single.sh {}

printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} ./tilelibCWL.sh {} lib /data-sdd/scirpts/tilelib/cgi-69-data /data-sdd/cwl_tiling/tilelib/fastj2cgflib /data-sdd/cwl_tiling/tilelib/data /data-sdd/cwl_tiling/tilelib/verbose_tagset /data-sdd/data/l7g/tagset.fa/tagset.fa.gz 
#printf '%04x\n' {0..862} | xargs -n1 -P 15 -I{} time ./run_cglf_single.sh {}
#printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} echo ">>" {}
