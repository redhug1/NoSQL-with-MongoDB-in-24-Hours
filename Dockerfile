# to build this dockerfile: docker build . -t mong:2.4.8
# then to run: docker run --rm -hostname{mongo_2-4} -it -v $PWD:/studies mong:2.4.8
# then use mongo: mongo, etc
#
# Start from specific image version of ubuntu
FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install lsof nano

COPY mongodb-linux-x86_64-2.4.8/bin/. /usr/local/bin
COPY code/hour02/mongod_config.txt /usr/local/bin

# Create the default data directory for mongo
RUN mkdir -p /data/db

COPY start.sh .
RUN chmod +x start.sh

# run commands in running container, to finalise container setup for use
CMD ./start.sh
