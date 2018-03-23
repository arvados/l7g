#!/bin/bash

arvados-cwl-runner --local \
  cwl/sglf-create-wf.cwl \
  yml/sglf-create-test.yml

arvados-cwl-runner --local \
  cwl/sglf-create-wf.cwl \
  yml/sglf-create-test2.yml
