#!/bin/bash

arvados-cwl-runner --local \
  cwl/sglf-create.cwl \
  yml/sglf-create-test.yml
