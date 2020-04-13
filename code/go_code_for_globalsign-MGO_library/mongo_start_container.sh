#!/bin/bash
clear
echo "Starting mongo 3.6.17 Container ... (allow 3 seconds for # prompt) ..."
echo " "
docker run --rm -hostname{mongo_3-6} -it -v $PWD:/studies -v /data/db36:/data/db36 mong:3.6.17
