Docker file setup files and code for book:

## NoSQL with MongoDB in 24 Hours,
( book printed: September 2014 )

## This book uses mongo version 2.4.8
## Added folders to hours 16 to 18 that contain versions of the similarly named python files in `go`

`NOT` all of the examples in the book will work for later versions of mongo.

See the comments in 'Dockerfile' for creating a Docker image that creates a Ubuntu 18.04 Container with `mongo 2.4.8` installed + `python 3.6.9` + `go 1.14.1`

See also the instructions in the 'Dockerfile' for starting the Container.

You can then go down into the /studies/code/hourXX folders and run the various scripts as for example:

`mongo shell_script.js` from the /studies/code/hour02 folder

-=-=-

The folder `mongodb-linux-x86_64-2.4.8` came from unzipping the file `mongodb-linux-x86_64-2.4.8.tgz`, which came from site:
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
