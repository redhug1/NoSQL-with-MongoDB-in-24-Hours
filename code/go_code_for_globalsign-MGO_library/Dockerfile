# to build this dockerfile: docker build . -t mong:3.6.17
# then to run:
#
#   docker run --rm -hostname{mongo_3-6} -it -v $PWD:/studies mong:3.6.17
#
# OR to retain the database between seperate runs of the container ... use:
#
#   docker run --rm -hostname{mongo_3-6} -it -v $PWD:/studies -v /data/db36:/data/36 mong:3.6.17
#
# The above is in a script named: mongo_start_container.sh (make sure it is executable)
#
# then use mongo: mongo, etc
#
# Start from specific image version of ubuntu
FROM ubuntu:18.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y -q dialog apt-utils tasksel

RUN apt-get install lsof nano

# comment out the following 2 lines if you do not need python3 and it's mongo library that works with mongodb 2.4.8 !!! NOTE: needs updating for mongo 3.6.17 !!!
#RUN apt-get install -y python3 python3-pip libsnappy-dev libkrb5-dev
#RUN python3 -m pip install pymongo==2.7.2

COPY mongodb-linux-x86_64-ubuntu1404-3.6.17/bin/. /usr/local/bin

# Create the default data directory for mongo
# NOTE: this is a different folder, as mongo 2.4.8 database files are not compatible with mongo 3.6.17
RUN mkdir -p /data/db36

# Configure Go
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV GOBIN $GOPATH/bin

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN apt-get install -y wget

RUN wget -O go.tgz https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz 
RUN tar -C /usr/local -xzf go.tgz
RUN rm go.tgz

RUN apt-get install -y git

# Install go mongo library that will work with mongo 3.6.17
RUN go get github.com/globalsign/mgo

COPY mongod_config.txt /usr/local/bin
RUN mkdir /var/log/mongodb

COPY start.sh .
RUN chmod +x start.sh

EXPOSE 27017

# run commands in running container, to finalise container setup for use
CMD ./start.sh
