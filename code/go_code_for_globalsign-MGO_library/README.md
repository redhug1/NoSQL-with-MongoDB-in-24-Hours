Docker file setup files and code for book:

## This folder is for 'go' ONLY and works with the "globalsign/mgo" mongodb library and mongodb version 3.6.17


See the comments in 'Dockerfile' for creating a Docker image that creates a Ubuntu 18.04 Container with `mongo 3.6.17` installed + `go 1.14.1`

See also the instructions in the 'Dockerfile' for starting the Container.

When you run the container for the first time, you will need to create the initial database by running in `/hour05` the `generate_words.js` with command:

`mongo generate_words.js`

You can then go down into the /hourXX folders and run the various scripts as for example:

`go run connect.go` from the /hour16/GoConnect folder

-=-=-

The folder `mongodb-linux-x86_64-ubuntu1404-3.6.17` came from unzipping the file `linux/mongodb-linux-x86_64-ubuntu1404-3.6.17.tgz`, which came from site:
https://www.mongodb.org/dl/linux/x86_64

-=-=-

## Additional notes for linux:
## Page 81
- Use `nano` to edit: `mongod_config_auth.txt` in folder `/usr/local/bin`
- You should be able to copy the contents of said named file from hour04 folder into nano editor in the container.
- Use `ps -alt` to see the mongod process number and then use `kill -9 nn` to stop its process
- To restart the MongoDB Server: `mongod --config mongo_config_auth.txt &`
## Page 102
NOTE: The code in this repo' for `hour05/generate_words.js` has been modified to read in 10,000 words from a file to create the words collection.

To save typing, copy:

`mongo words --eval "db.word_stats.find().count()"`
## Page 121
Some extra words were added to `google-10000-english.txt` in folder `hour05` to enable the last two find's in `find_specific.js` to show some output.
## Page 350
You may need to run `generate_words.js` in hour05 if you wish to run PythonAdd.py more than once.