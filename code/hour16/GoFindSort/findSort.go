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

func sortWordsAscending(collection *mgo.Collection) {
	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("word")
	fmt.Printf("\nW words ordered ascending: ")
	displayCursor(cursor)
}

func sortWordsDescending(collection *mgo.Collection) {
	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("-word")
	fmt.Printf("\nW words ordered descending: ")
	displayCursor(cursor)
}

func sortWordsAscAndSize(collection *mgo.Collection) {
	query := bson.M{"first": "q"}
	cursor := collection.Find(query)
	cursor.Sort("last", "-size")
	fmt.Printf("\nQ words ordered first by last letter and then by size: ")
	displayCursor(cursor)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	sortWordsAscending(collection)
	sortWordsDescending(collection)
	sortWordsAscAndSize(collection)
}
