#!/bin/bash
mongod -vvv --config /usr/local/bin/mongod_config.txt &
sleep 2
echo " "
echo "Checking mongo service is listening:"
lsof -i -P -n | grep LISTEN
echo " "
cd studies
ls
#echo "alias python=python3" >> ~/.bashrc
bash
