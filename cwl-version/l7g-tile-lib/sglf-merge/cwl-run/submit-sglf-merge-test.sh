#!/bin/bash

arvados-cwl-runner --local \
  cwl/sglf-merge.cwl \
  yml/sglf-merge-test.yml

