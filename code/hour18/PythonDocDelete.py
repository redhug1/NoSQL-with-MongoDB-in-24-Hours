from pymongo import MongoClient # this does similar to what is in file: doc_delete.js (hour08)
def showNewDocs(collection):
    query = {'category': 'New'}
    cursor = collection.find(query)
    for doc in cursor:
        print(doc)
def removeNewDocs(collection):
    query = {'category': 'New'}
    results = collection.remove(query)
    print("\nDelete Docs Result:")
    print(str(results))
    print("\nAfter Deleting Docs:")
    showNewDocs(collection)
if __name__=="__main__":
    mongo = MongoClient('mongodb://localhost:27017')
    mongo.write_concern = {'w': 1, 'j': True}
    db = mongo['words']
    collection = db['word_stats']
    print("Before Deleting:")
    showNewDocs(collection)
    removeNewDocs(collection)
