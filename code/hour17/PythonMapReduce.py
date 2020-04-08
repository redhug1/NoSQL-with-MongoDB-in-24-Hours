from pymongo import MongoClient # this does same as file: doc_map_reduce.js (hour09)
from bson.code import Code
from bson.son import SON

def displayMapReduce(cursor):
    for doc in cursor:
        print(doc)

def totalVowelBeginningCertainLetter(collection):
    map = Code("""
                function(){
                    emit(this.first, this.stats.vowels);
                }
                """)
    reduce = Code("""
                function(key, values){
                    return Array.sum(values);
                }
                """)
    out = "results_collection"
    results = collection.map_reduce(map, reduce, out)
    print("\n\nTotal vowel count in words beginning with a certain letter:")
    col = db[out]
    cursor = col.find()
    displayMapReduce(cursor)

def moreComplexMapReduce(collection):
    map = Code("""
                function() { emit(this.first,
                            { count: 0,
                            vowels: this.stats.vowels,
                            consonants: this.stats.consonants})
                }
                """)
    reduce = Code("""
                function(key, values){
                    result = {count: values.length,
                            vowels: 0, consonants: 0};
                    for(var i=0; i<values.length; i++){
                        if (values[i].vowels)
                            result.vowels += values[i].vowels;
                        if (values[i].consonants)
                            result.consonants += values[i].consonants;
                    }
                    return result;
                }
                """)
    out = "results_collection2"
    # the following finalize seems to have no effect
    finalize = Code("""
            finalize: function(key, obj) {
                obj.characters = obj.vowels + obj.consonants;
                return obj;
            }
            """)

    results = collection.map_reduce(map, reduce, out, finalize, query={'last': {'$in':['a','e','i','o','u']}})
    print("\n\nTotal words, vowels, consonants and characters in words beginning with a certain letter that ends with a vowel:")
    col = db["results_collection2"]
    # now apply the the equivalent of the 'finalize' that did not do anything

    # It would be nice if something like the following would 'work' ...
##    query = {'value.count': {'$gte': 0}}
##    update = {
##        '$set' : {'value.characters': {'$sum': ['$value.vowels', '$value.consonants']}}
##    }
##    results = col.update(query, update, upsert= False, multi= True)

    # BUT, i can't get it to, so the following brute force works:

    cursor = col.find()
    for doc in cursor:
        sum = doc['value']['vowels'] + doc['value']['consonants']   # total the 'vowels' and 'consonants' numbers from document
        doc_id = doc['_id']                                         # get doocument's _id
        query_str = {'_id':doc_id}                                  # construct the query
        sum_str = str(sum)
        update_str = {'$set':{'characters':sum_str}}                # construct the update
        print(query_str, update_str)                                # show what is about to be done 'for sanity check'
        result = col.update(query_str, update_str, upsert= False, multi= False)

    cursor = col.find()
    displayMapReduce(cursor)

if __name__=="__main__":
    mongo = MongoClient('mongodb://localhost:27017')
    db = mongo['words']
    collection = db['word_stats']
    totalVowelBeginningCertainLetter(collection)
    moreComplexMapReduce(collection)
