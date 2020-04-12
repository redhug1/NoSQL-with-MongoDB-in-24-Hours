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

func pageResults(collection *mgo.Collection, skip int) {
	query := bson.M{"first": "y"}
	cursor := collection.Find(query).Limit(10).Skip(skip)
	res_count, err := cursor.Count()
	check(err)
	if res_count > 0 {
		fmt.Printf("\nPage %v to %v :\n", skip+1, skip+res_count)
		displayCursor(cursor)
		if res_count == 10 {
			pageResults(collection, skip+10) // recurse
		}
	}
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	pageResults(collection, 0)
}
