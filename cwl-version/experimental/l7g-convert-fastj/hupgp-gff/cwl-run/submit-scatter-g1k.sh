#!/bin/bash

arvados-cwl-runner \
  --submit \
  --disable-reuse \
  --project-uuid=su92l-j7d0g-sndpr7v3au985dt \
  --api=containers \
  --no-wait \
  cwl/scatter_hupgp-gff-to-fastj.cwl \
  yml/scatter_hupgp-gff-to-fastj_human_g1k_v37.yml
