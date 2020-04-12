package main

import (
	"bytes"
	"encoding/json"
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

func displayDoc(doc bson.M) {
	fmt.Printf("%v\n", doc)
	jsonString, err := json.MarshalIndent(doc, "", " ")
	check(err)
	fmt.Println("\nResult as JSON:")

	var out bytes.Buffer
	err = json.Indent(&out, jsonString, "", "  ")
	check(err)

	var st string = out.String()
	fmt.Printf("%v\n", st)
}

func displayAggregate(iter *mgo.Iter) {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		displayDoc(doc)
	}
	err := iter.Close()
	check(err)
}

func largeSmallVowels(collection *mgo.Collection) {
	var vowels = []string{"a", "e", "i", "o", "u"}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"first": bson.M{"$in": vowels}}},
		bson.M{"$group": bson.M{"_id": "$first",
			"largest":  bson.M{"$max": "$size"},
			"smallest": bson.M{"$min": "$size"},
			"total":    bson.M{"$sum": 1}}},
		bson.M{"$sort": bson.M{"first": 1}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nLargest and smallest word sizes for word begining with a vowel:\n")
	displayAggregate(iter)
}

func top5AverageWordFirst(collection *mgo.Collection) {
	pipeline := []bson.M{
		bson.M{"$group": bson.M{"_id": "$first",
			"average": bson.M{"$avg": "$size"}}},
		bson.M{"$sort": bson.M{"average": -1}},
		bson.M{"$limit": 5},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nFirst letter of top 5 largest average word size:\n")
	displayAggregate(iter)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	largeSmallVowels(collection)
	top5AverageWordFirst(collection)
}
