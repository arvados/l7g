#!/usr/bin/env python

import fnmatch
import requests
import time
import os
import glob

# You can alternatively define these in travis.yml as env vars or arguments
BASE_URL = 'https://view.commonwl.org'
WORKFLOW_PATH = '/workflows/workflow.cwl'

#get the cwl in l7g/cwl-version
matches = []
for root, dirnames, filenames in os.walk('cwl-version'):
    for filename in fnmatch.filter(filenames, '*.cwl'):
        matches.append(os.path.join(root, filename))

print matches

REPO_SLUG = 'curoverse/l7g/blob/master/'
#Testing WORKFLOW_PATH
WORKFLOW_PATH = 'cwl-version/filter/cwl/tiling_filtergvcf19.cwl'

#This will loop through matches, need to indent everything after to make work
#for WORKFLOW_PATH in matches:
# Whole workflow URL on github
workflowURL = 'https://github.com/' + REPO_SLUG + WORKFLOW_PATH
print workflowURL,'\n'

# Headers
HEADERS = {
'user-agent': 'my-app/0.0.1',
'accept': 'application/json'
}

# Add new workflow with the specific commit ID of this build
addResponse = requests.post(BASE_URL + '/workflows',
data={'url': workflowURL},
headers=HEADERS)

if addResponse.status_code == requests.codes.accepted:
    qLocation = addResponse.headers['location']

    # Get the queue item until success
    qResponse = requests.get(BASE_URL + qLocation,
    headers=HEADERS,
    allow_redirects=False)
    maxAttempts = 5
    while qResponse.status_code == requests.codes.ok and qResponse.json()['cwltoolStatus'] == 'RUNNING' and maxAttempts > 0:
        time.sleep(5)
        qResponse = requests.get(BASE_URL + qLocation,
        headers=HEADERS,
        allow_redirects=False)
        maxAttempts -= 1

        if qResponse.headers['location']:
            # Success, get the workflow
            workflowResponse = requests.get(BASE_URL + qResponse.headers['location'], headers=HEADERS)
            if (workflowResponse.status_code == requests.codes.ok):
                workflowJson = workflowResponse.json()
                # Do what you want with the workflow JSON
                # Include details in documentation files etc
                print(BASE_URL + workflowJson['visualisationSvg'])
                print('Verified with cwltool version ' + workflowJson['cwltoolVersion'])
                # etc...
            else:
                print('Could not get returned workflow')
        elif qResponse.json()['cwltoolStatus'] == 'ERROR':
            # Cwltool failed to run here
            print(qResponse.json()['message'])
        elif maxAttempts == 0:
            print('Timeout: Cwltool did not finish')

else:
    print('Error adding workflow')
