#!/bin/bash

arvados-cwl-runner \
  --submit \
  --disable-reuse \
  --api=containers \
  --no-wait \
  sglf-sanity-check.cwl \
  cwl-run/sglf-sanity-check.yml

