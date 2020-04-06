package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func displayCursor(cursor *mgo.Query) {
	var doc bson.M
	var words string
	iter := cursor.Iter()
	for iter.Next(&doc) {
		valueWord := doc["word"]
		switch v := valueWord.(type) { // do "type assertion" for field
		case string:
			if len(words) > 0 {
				words = words + "," + v
			} else {
				words = v
			}
		default:
			fmt.Println("word error")
		}
	}
	err := iter.Close()
	check(err)
	if len(words) > 65 {
		words = words[:65] + "..."
	}
	fmt.Println(words)
}

func over12(collection *mgo.Collection) {
	fmt.Printf("\n\nWords with more than 12 characters:\n")
	cursor := collection.Find(bson.M{"size": bson.M{"$gt": 12}})
	displayCursor(cursor)
}

func startingABC(collection *mgo.Collection) {
	fmt.Printf("\nWords starting with A, B, C:\n")
	var abc = []string{"a", "b", "c"}
	query := bson.M{"first": bson.M{"$in": abc}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func startEndVowels(collection *mgo.Collection) {
	fmt.Printf("\nWords starting and ending with a vowel:\n")
	var vowels = []string{"a", "e", "i", "o", "u"}
	query := bson.M{"$and": []bson.M{
		bson.M{"first": bson.M{"$in": vowels}},
		bson.M{"last": bson.M{"$in": vowels}},
	}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func over6Vowels(collection *mgo.Collection) {
	fmt.Printf("\nWords with more than 5 vowels:\n")
	query := bson.M{"stats.vowels": bson.M{"$gt": 5}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func nonAlphaCharacters(collection *mgo.Collection) {
	print("\nWords with 1 non-alphabet character:\n")
	query := bson.M{"charsets": bson.M{"$elemMatch": bson.M{"$and": []bson.M{
		bson.M{"type": "other"},
		bson.M{"chars": bson.M{"$size": 1}}}}}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	over12(collection)
	startingABC(collection)
	startEndVowels(collection)
	over6Vowels(collection)
	nonAlphaCharacters(collection)
}
