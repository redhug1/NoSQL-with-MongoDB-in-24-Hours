#!/bin/bash
mongod --config /usr/local/bin/mongod_config.txt &
sleep 7
echo " "
echo "Checking mongo service is listening:"
lsof -i -P -n | grep LISTEN
echo " "
echo "Checking mongo version:"
mongo --version
echo " "
echo "'studies/code' folder contains:"
echo " "
cd studies/code
ls
echo " "
echo "alias python=python3" >> ~/.bashrc
bash
