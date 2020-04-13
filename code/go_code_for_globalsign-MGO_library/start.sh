#!/bin/bash
mongod -vvv --config /usr/local/bin/mongod_config.txt &
pid_mongo=$!
sleep 2
echo " "
echo "Checking mongo service is listening:"
lsof -i -P -n | grep LISTEN
echo " "
echo "Checking mongo version:"
mongo --version
echo " "
echo "'studies' folder (for go_code_for_globalsign-MGO_library) contains:"
echo " "
cd studies
ls
echo " "
#echo "alias python=python3" >> ~/.bashrc
bash
echo " "
kill -SIGTERM $pid_mongo
if ps -p $pid_mongo > /dev/null; then
    echo "waiting for 'mongod' to finish ..."
fi

while ps -p $pid_mongo > /dev/null; do sleep 1; done
echo "'mongod' finished."