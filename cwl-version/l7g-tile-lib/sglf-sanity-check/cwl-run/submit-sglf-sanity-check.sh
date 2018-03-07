#!/bin/bash

arvados-cwl-runner \
  --submit \
  --disable-reuse \
  --api=containers \
  --no-wait \
  cwl/sglf-sanity-check.cwl \
  yml/sglf-sanity-check.yml

