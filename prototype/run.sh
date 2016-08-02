#!/bin/bash

LOC_DATA_DIR="/scratch/l7g"
img="$1"

if [[ "$img" == "" ]] ; then
  echo "provide image"
  exit 1
fi

docker run -v $LOC_DATA_DIR:/data -t -i -p 8081:8081 -p 8082:8082 -p 8083:8083 -p 8084:8084 -p 8085:8085 -e LD_LIBRARY_PATH=/usr/local/lib $img /bin/bash
