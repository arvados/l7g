#!/bin/bash

arvados-cwl-runner \
  --disable-reuse \
  cwl/sglf-sanity-check-p_wf.cwl \
  yml/sglf-sanity-check-p_wf.yml

exit

arvados-cwl-runner \
  --submit \
  --disable-reuse \
  --api=containers \
  --no-wait \
  cwl/sglf-sanity-check-p_wf.cwl \
  yml/sglf-sanity-check-p_wf.yml

