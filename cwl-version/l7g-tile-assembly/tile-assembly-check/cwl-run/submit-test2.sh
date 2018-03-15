#!/bin/bash

arvados-cwl-runner \
  cwl/tile-assembly-sanity-check.cwl \
  yml/tile-assembly-check_test2.yml
