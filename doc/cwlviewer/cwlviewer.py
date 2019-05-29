#!/usr/bin/env python

import requests
import time
import os
import argparse

# Code largely based off of example from
# https://view.commonwl.org/apidocs


def cwlviewer():
    ''' Takes input cwl workflow and creates a diagram to represent
        it using cwl-viewer.
        Run using the following:  python cwlviewer URLOFCWLWORKFLOW'''

    # Setting up inputs
    parser = argparse.ArgumentParser()
    parser.add_argument('workflowURL', metavar='WORKFLOWURL',
                        help='URL of workflow')
    args = parser.parse_args()
    workflowURL = args.workflowURL

    # Whole workflow URL on github
    BASE_URL = 'https://view.commonwl.org'

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
        qLocation = addResponse.headers['Location']
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


        if 'Location' in qResponse.headers:
            # Success, get the workflow
            workflowResponse = requests.get(
                BASE_URL + qResponse.headers['location'], headers=HEADERS)
            if (workflowResponse.status_code == requests.codes.ok):
                workflowJson = workflowResponse.json()
                print(BASE_URL + workflowJson['visualisationSvg'])
                print('Verified with cwltool version ' +
                      workflowJson['cwltoolVersion'])
            else:
                print('Error: Could not get returned workflow')
        elif qResponse.json()['cwltoolStatus'] == 'ERROR':
            # Cwltool failed to run here
#            print(qResponse.json()['message'])
            print('Error: Cwltool failed to verify')
        elif maxAttempts == 0:
            print('Timeout: Cwltool did not finish')

    # current hack to get around bug that is returning workflow
    elif addResponse.status_code == requests.codes.ok:
        workflowResponse = addResponse
        workflowJson = workflowResponse.json()
        print(BASE_URL + workflowJson['visualisationSvg'])
        print('Verified with cwltool version ' +
              workflowJson['cwltoolVersion'])

    else:
        print('Error adding workflow')


if __name__ == '__main__':
    cwlviewer()

