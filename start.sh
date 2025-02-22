#!/bin/bash
mongod --config /usr/local/bin/mongod_config.txt &
pid_mongo=$!
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
echo " "
kill -SIGTERM $pid_mongo
if ps -p $pid_mongo > /dev/null; then
    echo "waiting for 'mongod' to finish ..."
fi

while ps -p $pid_mongo > /dev/null; do sleep 1; done
echo "'mongod' finished."