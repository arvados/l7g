#!/bin/bash
#
# wrapper script to run the './run_cglf_single-cgi-69.sh' script on each
# of the tile paths.
#

set -eo pipefail
printf '%04x\n' {0..862} | xargs -n1 -P 15 -I{} time ./run_cglf_single-cgi-69.sh {}
