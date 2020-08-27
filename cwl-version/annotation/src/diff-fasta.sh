#!/bin/bash

set -aeo pipefail

export ref="$1"
export sample="$2"
export offset="$3"

lightning diff-fasta -offset $offset <(printf $ref) <(printf $sample)
