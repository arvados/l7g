import arvados
import re
import sys
import argparse

# List collections with owner_uuid=owner_uuid
collectnames = []
pdh = []

owner_uuid = 'su92l-j7d0g-tkhhwr3tec6pjyq'

# Get list of collections owned by uuid
call = arvados.util.list_all(arvados.api().collections().list, filters=[["owner_uuid","=",owner_uuid]])

n = len(call)

for i in xrange(0,n):
  collectnames.append("%s\n" % (call[i]['name']))
  pdh.append("%s\n" % (call[i]['portable_data_hash']))

# Write out collection names and phds to output files
f1 = open('collectnames.txt', 'w')
f2 = open('pdhs.txt','w')

for l in collectnames: f1.write(l)
for l in pdh: f2.write(l)

