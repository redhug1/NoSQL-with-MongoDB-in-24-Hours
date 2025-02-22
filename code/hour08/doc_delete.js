mongo = new Mongo("localhost");
wordsDB = mongo.getDB("words");
wordsDB.runCommand({getLastError: 1, w: 1, j: true, wtimeout: 1000});
wordsColl = wordsDB.getCollection("word_stats");
print("Before Delete One: ");
cursor = wordsColl.find({category: 'New'}, {word:1});
printjson(cursor.toArray());
wordsColl.remove({category: 'New'}, true);
cursor = wordsColl.find({category: 'New'}, {word:1});
print("After Delete One: ")
printjson(cursor.toArray());
wordsColl.remove({category: 'New'});
cursor = wordsColl.find({category: 'New'}, {word:1});
print("After Delete All: ")
printjson(cursor.toArray());