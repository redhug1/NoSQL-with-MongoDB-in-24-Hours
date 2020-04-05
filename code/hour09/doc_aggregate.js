mongo = new Mongo('localhost');
wordsDB = mongo.getDB('words');
wordsColl = wordsDB.getCollection('word_stats');
results = wordsColl.aggregate(
    { $match: {first:{$in:['a','e','i','o','u']}}},
    { $group: { _id:"$first",
                largest:{$max:"$size"},
                smallest:{$min:"$size"},
                total:{$sum:1}}},
    { $sort: {_id:1}}
);
print("Largest and smallest word sizes for " +
    "words beginning with a vowel: ");
results.result.forEach(function(item){
    print(JSON.stringify(item));
});
results = wordsColl.aggregate(
    {$match: {size:4}},
    {$limit: 5},
    {$project: {_id:"$word", stats:1}}
);
print("Stats for 5 four leter words: ");
results.result.forEach(function(item){
    print(JSON.stringify(item));
});
results = wordsColl.aggregate(
    {$group: {_id:"$first", average:{$avg:"$size"}}},
    {$sort: {average:-1}},
    {$limit: 5}
);
print("First letter of top 6 largest average size: ");
results.result.forEach(function(item){
    print(JSON.stringify(item));
});