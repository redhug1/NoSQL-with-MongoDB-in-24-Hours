#!/bin/bash
clear
echo "Starting mongo 2.4.8 Container ... (allow 7 seconds for # prompt) ..."
echo " "
docker run --rm -hostname{mongo_2-4} -it -v $PWD:/studies -v /data/db24:/data/db24 mong:2.4.8
