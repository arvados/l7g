#!/bin/bash

set -aeo pipefail

export ref="$1"
export sample="$2"
export offset="$3"
export prefix="$4"

lightning diff-fasta -timeout=100ms -offset $offset -sequence $prefix <(printf $ref) <(printf $sample)
