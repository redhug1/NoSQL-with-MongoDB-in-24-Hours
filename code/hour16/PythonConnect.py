from pymongo import MongoClient
import pymongo
mongo = MongoClient('mongodb://localhost:27017/')
db = mongo['words']
collection = db['word_stats']
print("Number of Documents: ")
print(collection.find().count())