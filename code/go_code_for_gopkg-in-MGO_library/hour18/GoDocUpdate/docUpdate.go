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
		log.Printf("Go application has failed, here's why:\n")
		log.Fatal(err)
		// NOTE: a real application needs to do a lot more with error handling than just stop here
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

func showWord(collection *mgo.Collection) {
	var doc bson.M
	var words = []string{"left", "lefty"}
	query := bson.M{"word": bson.M{"$in": words}}
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		displayDoc(doc)
	}
}

func updateDoc(collection *mgo.Collection) {
	query := bson.M{"word": "left"}
	update := bson.M{
		"$set":  bson.M{"word": "lefty"},
		"$inc":  bson.M{"size": 1, "stats.consonants": 1},
		"$push": bson.M{"letters": "y"},
	}
	err := collection.Update(query, update)
	check(err)
	print("\nAfter Updating Doc:\n")
	showWord(collection)
}

func resetDoc(collection *mgo.Collection) {
	query := bson.M{"word": "lefty"}
	update := bson.M{
		"$set": bson.M{"word": "left"},
		"$inc": bson.M{"size": -1, "stats.consonants": -1},
		"$pop": bson.M{"letters": "y"},
	}
	err := collection.Update(query, update)
	check(err)
	fmt.Printf("\nAfter Resetting Doc:\n")
	showWord(collection)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer func() {
		fmt.Printf("Closing mongodb session\n")
		session.Close()
	}()

	collection := session.DB("words").C("word_stats")

	fmt.Printf("\nBefore Updating:\n")
	showWord(collection)
	updateDoc(collection)
	resetDoc(collection)
}
