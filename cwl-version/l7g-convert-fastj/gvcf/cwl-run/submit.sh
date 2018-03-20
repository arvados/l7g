#!/bin/bash

arvados-cwl-runner \
  --disable-reuse \
  cwl/gvcf-to-fastj_wf.cwl \
  yml/gvcf-to-fastj.yml

