from pymongo import MongoClient # this does what is in file: doc_add.js (hour08)
def displayCursor(cursor):
    words = ''
    for doc in cursor:
        words += doc["word"] + ","
    if len(words) > 65:
        words = words[:65] + "..."
    print(words)
def findSpecificWords(collection):
    query = {'word': {'$in': ["tweet","gogle","selfie"]}}
    cursor = collection.find(query)
    displayCursor(cursor)

def showNewDocs(collection):
    query = {'category': 'New'}
    cursor = collection.find(query)
    for doc in cursor:
        print(doc)
    findSpecificWords(collection) # added to just show the word of interest
def addSelfie(collection):
    selfie = {
        'word': 'selfie', 'first': 's', 'last': 'e',
        'size': 6, 'category': 'New',
        'stats': {'vowels': 3, 'consonants': 3},
        'letters': ["s","e","l","f","i"],
        'charsets': [
             {'type': 'consonants', 'chars': ["s","l","f"]},
             {'type': 'vowels', 'chars': ["e","i"]}]}
    results = collection.insert(selfie)
    print("\nInserting One Results:")
    print(str(results))
    print("After Inserting One:")
    showNewDocs(collection)
def addGoogleAndTweet(collection):
    # deliberate mis-spelling as google is already in the 100K list of words
    gogle = {
        'word': 'gogle', 'first': 'g', 'last': 'e',
        'size': 6, 'category': 'New',
        'stats': {'vowels': 2, 'consonants': 3},
        'letters': ["g","o","l","e"],
        'charsets': [
             {'type': 'consonants', 'chars': ["g","l"]},
             {'type': 'vowels', 'chars': ["o","e"]}]}
    tweet = {
        'word': 'tweet', 'first': 't', 'last': 't',
        'size': 5, 'category': 'New',
        'stats': {'vowels': 2, 'consonants': 3},
        'letters': ["t","w","e"],
        'charsets': [
             {'type': 'consonants', 'chars': ["t","w"]},
             {'type': 'vowels', 'chars': ["e"]}]}
    results = collection.insert([gogle, tweet])
    print("\nInserting Multiple Results:")
    print(str(results))
    print("After Inserting Multiple:")
    showNewDocs(collection)
if __name__=="__main__":
    mongo = MongoClient('mongodb://localhost:27017')
    mongo.write_concern = {'w': 1, 'j': True}
    db = mongo['words']
    collection = db['word_stats']
    print("Before Inserting:")
    showNewDocs(collection)
    addSelfie(collection)
    addGoogleAndTweet(collection)