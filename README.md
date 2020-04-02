Docker file setup files and code for book:

## NoSQL with MongoDB in 24 Hours,
( book printed: September 2014 )

## This book uses mongo version 2.4.8

`NOT` all of the examples in the book will work for later versions of mongo.

See the comments in 'Dockerfile' for creating a Docker image that creates a Ubuntu 18.04 Container with mongo 2.4.8 installed.

See also the instructions in the 'Dockerfile' for starting the Container.

You can then go down into the /studies/code/hourXX folders and run the various scripts as for example:

`mongo shell_script.js` from the /studies/code/hour02 folder

-=-=-

The folder `mongodb-linux-x86_64-2.4.8` came from unzipping the file `mongodb-linux-x86_64-2.4.8.tgz`, which came from site:
https://www.mongodb.org/dl/linux/x86_64


-=-=-

## Additional notes for linux:
## Page 81
- Use `nano` to create/edit: `mongod_config_auth.txt` in folder `/usr/local/bin`
- You should be able to copy the contents of said named file from hour04 folder into nano editor in the container.
- Use `ps -alt` to see the mongod process number and then use `kill -9 nn` to stop its process
- To restart the MongoDB Server: `mongod --config mongo_config_auth.txt &`
