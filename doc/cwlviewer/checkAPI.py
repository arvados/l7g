#!/usr/bin/env python

import fnmatch
import requests
import time
import os
import glob

BASE_URL = 'https://view.commonwl.org/workflows'

workflowurl = 'https://github.com/curoverse/l7g/blob/master/cwl-version/convert2fastj/gvcf_version/cwl/tiling_convert2fastj_gvcf.cwl'

HEADERS = {
    'user-agent': 'my-app/0.0.1',
    'accept': 'application/json'
}

addResponse = requests.post(BASE_URL,
                            data={'url': workflowurl}, headers=HEADERS)


print(addResponse.status_code)
