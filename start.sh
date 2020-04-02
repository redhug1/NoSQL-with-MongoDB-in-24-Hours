#!/bin/bash
mongod &
sleep 7
echo " "
echo "Checking mongo service is listening:"
lsof -i -P -n | grep LISTEN
echo " "
cd studies/code
ls
bash
