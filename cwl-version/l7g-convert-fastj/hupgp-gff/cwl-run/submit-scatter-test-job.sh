#!/bin/bash

arvados-cwl-runner --submit --no-wait scatter_hupgp-gff-to-fastj.cwl cwl-run/scatter_hupgp-gff-to-fastj_Arvados-test.yml
#arvados-cwl-runner --disable-reuse --local scatter_hupgp-gff-to-fastj.cwl cwl-run/scatter_hupgp-gff-to-fastj_Arvados-test.yml
#arvados-cwl-runner --local scatter_hupgp-gff-to-fastj.cwl cwl-run/scatter_hupgp-gff-to-fastj_Arvados-test.yml
