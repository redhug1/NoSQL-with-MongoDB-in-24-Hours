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

func countWords(collection *mgo.Collection) {
	count, err := collection.Find(bson.M{}).Count()
	check(err)
	fmt.Printf("\nTotal words in the collection: %d\n", count)
	query := bson.M{"first": "a"}
	count, err = collection.Find(query).Count()
	check(err)
	fmt.Printf("\nTotal words starting with A: %d\n", count)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	countWords(collection)
}
