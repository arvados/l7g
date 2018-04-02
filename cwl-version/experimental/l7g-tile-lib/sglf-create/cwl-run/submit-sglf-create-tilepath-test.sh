#!/bin/bash

arvados-cwl-runner --local \
  cwl/sglf-create-tilepath.cwl \
  yml/sglf-create-tileepath-test.yml
