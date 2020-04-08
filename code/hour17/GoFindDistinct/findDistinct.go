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

func sizesOfAllWords(collection *mgo.Collection) {
	var result []int
	err := collection.Find(nil).Distinct("size", &result)
	check(err)
	fmt.Printf("\nDistinct Sizes of words: %v\n", result)
}

func sizesOfQWords(collection *mgo.Collection) {
	var result []int
	query := bson.M{"first": "q"}
	cursor := collection.Find(query)
	err := cursor.Distinct("size", &result)
	check(err)
	fmt.Printf("\nDistinct Sizes of words starting with Q: %v\n", result)
}

func firstLetterOfLongWords(collection *mgo.Collection) {
	var result []string
	query := bson.M{"size": bson.M{"$gt": 12}}
	cursor := collection.Find(query)
	err := cursor.Distinct("first", &result)
	check(err)
	fmt.Printf("\nDistinct first letters of words longer than 12 characters: %v\n", result)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	sizesOfAllWords(collection)
	sizesOfQWords(collection)
	firstLetterOfLongWords(collection)
}
