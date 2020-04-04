blogy = {   // adjust spelling as 'blog' already in 100K list of words
    word: 'blogy', first: 'b', last: 'y',
    size: 5, letters: ['b','l','o','g','y'],
    stats: {vowels:1, consonants: 4},
    charsets: [{type: 'consonants', chars: ['b','l','g','y']},
                {type: 'vowels', chars: ['o']}],
    category: 'New' };
mongo = new Mongo("localhost");
wordsDB = mongo.getDB("words");
wordsDB.runCommand({getLastError: 1, w: 1, j: true, wtimeout: 1000});
wordsColl = wordsDB.getCollection("word_stats");
cursor = wordsColl.find({category:"blue"}, {word: 1, category:1});
print("Before Existing Save: ");
printjson(cursor.toArray());
word = wordsColl.findOne({word:"ocean"});
word.category="blue";
wordsColl.save(word);
word = wordsColl.findOne({word:"sky"});
word.category="blue";
wordsColl.save(word);
cursor = wordsColl.find({category:"blue"}, {word: 1, category:1});
print("After Existing Save: ");
printjson(cursor.toArray());
word = wordsColl.findOne({word:"blogy"});// adjust spelling as 'blog' already in 100K list of words
print("Before New Document Save: ");
printjson(word);
wordsColl.save(blogy);
word = wordsColl.findOne({word:"blogy"}, {word: 1, category:1});
print("After New Document Save: ");
printjson(word);