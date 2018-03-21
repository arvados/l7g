#!/bin/bash

arvados-cwl-runner \
  --submit \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --no-wait \
  --api=containers \
  --submit-runner-ram 80480 \
  cwl/scatter_hupgp-gff-to-fastj.cwl \
  yml/scatter_hupgp-gff-to-fastj_hg19_0.yml

arvados-cwl-runner \
  --submit \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --no-wait \
  --api=containers \
  --submit-runner-ram 80480 \
  cwl/scatter_hupgp-gff-to-fastj.cwl \
  yml/scatter_hupgp-gff-to-fastj_hg19_1.yml

