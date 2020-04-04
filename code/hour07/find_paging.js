function displayWords(skip, cursor){
    print("Page: " + parseInt(skip+1) + " to " +
        parseInt(skip+cursor.size()));
    words = cursor.map(function(word){
        return word.word;
    });
    wordStr = JSON.stringify(words);
    if (wordStr.length > 80){
        wordStr = wordStr.slice(0, 80) + "...";
    }
    print(wordStr);
}
mongo = new Mongo("localhost");
wordsDB = mongo.getDB("words");
wordsColl = wordsDB.getCollection("word_stats");
cursor = wordsColl.find({first:'y'});
count = cursor.size();
skip = 0;
for (i = 0; i < count; i += 10){
    cursor = wordsColl.find({first:'y'});
    cursor.skip(skip);
    cursor.limit(10);
    displayWords(skip, cursor);
    skip += 10;
}