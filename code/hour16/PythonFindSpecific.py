from pymongo import MongoClient	# does similar to what is in file find_specific.js (hour06)
def displayCursor(cursor):
    words = ''
    for doc in cursor:
        words += doc["word"] + ","
    if len(words) > 65:
        words = words[:65] + "..."
    print(words)
def over12(collection):
    print("\n\nWords with more than 12 characters:")
    query = {'size': {'$gt': 12}}
    cursor = collection.find(query)
    displayCursor(cursor)
def startingABC(collection):
    print("\nWords starting with A, B, C:")
    query = {'first': {'$in': ["a","b","c"]}}
    cursor = collection.find(query)
    displayCursor(cursor)
def startEndVowels(collection):
    print("\nWords starting and ending with a vowel:")
    query = {'$and': [
        {'first': {'$in': ["a","e","i","o","u"]}},
        {'last': {'$in': ["a","e","i","o","u"]}}]}
    cursor = collection.find(query)
    displayCursor(cursor)
def over6Vowels(collection):
    print("\nWords with more than 5 vowels:")
    query = {'stats.vowels': {'$gt': 5}}
    cursor = collection.find(query)
    displayCursor(cursor)
def nonAlphaCharacters(collection):
    print("\nWords with 1 non-alphabet character:")
    query = {'charsets':
        {'$elemMatch':
            {'$and': [
                {'type': 'other'},
                {'chars': {'$size': 1}}]}}}
    cursor = collection.find(query)
    displayCursor(cursor)
if __name__=="__main__":
    mongo = MongoClient('mongodb://localhost:27017')
    db = mongo['words']
    collection = db['word_stats']
    over12(collection)
    startingABC(collection)
    startEndVowels(collection)
    over6Vowels(collection)
    nonAlphaCharacters(collection)
