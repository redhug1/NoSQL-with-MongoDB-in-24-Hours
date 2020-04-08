from pymongo import MongoClient	# does a little of what is in file: find_count.js (hour07)
def countWords(collection):
    cursor = collection.find()
    print("Total words in the collection:")
    print(cursor.count())
    query = {'first': 'a'}
    cursor = collection.find(query)
    print("\nTotal words starting with A:")
    print(cursor.count())
if __name__=="__main__":
    mongo = MongoClient('mongodb://localhost:27017')
    db = mongo['words']
    collection = db['word_stats']
    countWords(collection)
