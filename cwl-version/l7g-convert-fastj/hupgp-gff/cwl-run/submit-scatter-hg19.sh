#!/bin/bash

## local
##
arvados-cwl-runner \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --api=containers \
  --eval-timeout 20000 \
  --local \
  scatter_hupgp-gff-to-fastj.cwl \
  cwl-run/scatter_hupgp-gff-to-fastj_hg19.yml

exit

## remote
##
arvados-cwl-runner \
  --submit \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --no-wait \
  --api=containers \
  --submit-runner-ram 80480 \
  scatter_hupgp-gff-to-fastj.cwl \
  cwl-run/scatter_hupgp-gff-to-fastj_hg19.yml

exit

arvados-cwl-runner \
  --submit \
  --disable-reuse \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --no-wait \
  --submit-runner-ram 40480 \
  scatter_hupgp-gff-to-fastj.cwl \
  cwl-run/scatter_hupgp-gff-to-fastj_hg19.yml
