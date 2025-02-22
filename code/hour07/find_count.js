mongo = new Mongo("localhost");
wordsDB = mongo.getDB("words");
wordsColl = wordsDB.getCollection("word_stats");
cursor = wordsColl.find({first: {$in: ['a', 'b', 'c']}});
print("words starting with a, b or c: ", cursor.count());
cursor = wordsColl.find({size:{$gt: 12}});
print("Words longer than 12 characters: ", cursor.count());
cursor = wordsColl.find({size:{$mod: [2,0]}});
print("Words with even Lengths: ", cursor.count());
cursor = wordsColl.find({letters:{$size: 12}});
print("Words with 12 distinct characters: ", cursor.count());
cursor = wordsColl.find({$and:
    [{first:{
        $in: ['a', 'e', 'i', 'o', 'u']}},
     {last:{
        $in: ['a', 'e', 'i', 'o', 'u']}}]});
print("Words that start and end with a vowel: ", cursor.count());
cursor = wordsColl.find({"stats.vowels":{$gt: 6}});
print("Words containing 7 or more vowel: ", cursor.count());
cursor = wordsColl.find({letters:{$all: ['a', 'e', 'i', 'o', 'u']}});
print("Words with all 5 vowels: ", cursor.count());
cursor = wordsColl.find({otherChars: {$exists: true}});
print("Words with non-alphabet characters: ", cursor.count());
cursor = wordsColl.find({charsets:{
    $elemMatch:{
        $and:[{type: 'other'},
            {chars: {$size: 2}}]}}});
print("Words with 2 non-alphabet characters: ", cursor.count());