#!/usr/bin/env python

import fnmatch
import requests
import time
import os
import glob

# You can alternatively define these in travis.yml as env vars or arguments
BASE_URL = 'https://view.commonwl.org/workflows'

#get the cwl in l7g/cwl-version
matches = []
for root, dirnames, filenames in os.walk('cwl-version'):
    for filename in fnmatch.filter(filenames, '*.cwl'):
        matches.append(os.path.join(root, filename))

print matches

REPO_SLUG = 'curoverse/l7g/blob/master/'

# Headers
HEADERS = {
'user-agent': 'my-app/0.0.1',
'accept': 'application/json'
}

#Testing WORKFLOW_PATH
#WORKFLOW_PATH = 'cwl-version/clean/cwl/tiling_clean_gvcf.cwl'

#This will loop through matches, need to indent everything after to make work
for WORKFLOW_PATH in matches:
# Whole workflow URL on github
    workflowURL = 'https://github.com/' + REPO_SLUG + WORKFLOW_PATH
    print '\n',workflowURL,'\n'

    # Add new workflow with the specific commit ID of this build
    addResponse = requests.post(BASE_URL,
                                data={'url': workflowURL},
                                headers = HEADERS)

    print BASE_URL,'\n',workflowURL,'\n\n'

    print(addResponse)
    print(addResponse.encoding)
    print(addResponse.content)
    print(addResponse.url)
    print(addResponse.request)
    print(addResponse.raw)
    print(addResponse.headers)

    print('\n\n End Sarah\'s code \n\n')
    print('Sleep 1 second\n\n')
    time.sleep(1)
