#!/bin/bash

set -eo pipefail

#echo {0..862} | xargs -n1 -P 15 -I{} bash -c " time ./run_cglf_single.sh {}

#printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} time ./run_cglf_single.sh {}
printf '%04x\n' {0..862} | xargs -n1 -P 15 -I{} time ./run_cglf_single.sh {}
#printf '%04x\n' {0..10} | xargs -n1 -P 15 -I{} echo ">>" {}
