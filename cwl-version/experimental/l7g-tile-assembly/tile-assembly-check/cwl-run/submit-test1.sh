#!/bin/bash

arvados-cwl-runner \
  cwl/tile-assembly-sanity-check.cwl \
  yml/tile-assembly-check_test1.yml
