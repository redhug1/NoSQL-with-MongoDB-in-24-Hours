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
	query := bson.M{"word": "righty"}
	err := collection.Find(query).One(&doc)
	if err != nil {
		log.Println(err) // continue on as the word being searched for will not initially be in the list ...
	}
	displayDoc(doc)
}

func addUpsert(collection *mgo.Collection) {
	var rLetters = []string{"r", "i", "g", "h"}
	var rConstChars = []string{"r", "g", "h"}
	var rVowelChars = []string{"i"}

	righty := bson.M{
		"word":     "righty",
		"first":    "r",
		"last":     "y",
		"size":     4,
		"category": "New",
		"stats":    bson.M{"vowels": 1, "consonants": 4},
		"letters":  rLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": rConstChars},
			bson.M{"type": "vowels", "chars": rVowelChars},
		}}
	err := collection.Insert(righty)
	check(err)

	print("\nAfter Upsert as insert:\n")
	showWord(collection)
}

func updateUpsert(collection *mgo.Collection) {
	var rLetters = []string{"r", "i", "g", "h", "t", "y"}
	var rConstChars = []string{"r", "g", "h", "t", "y"}
	var rVowelChars = []string{"i"}

	query := bson.M{"word": "righty"}
	update := bson.M{
		"$set": bson.M{
			"word":     "righty",
			"first":    "r",
			"last":     "y",
			"size":     6,
			"category": "Updated",
			"stats":    bson.M{"vowels": 1, "consonants": 5},
			"letters":  rLetters,
			"charsets": []bson.M{
				bson.M{"type": "consonants", "chars": rConstChars},
				bson.M{"type": "vowels", "chars": rVowelChars},
			}},
	}
	changeInfo, err := collection.Upsert(query, update)
	check(err)
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d, Matched: %d, Upserted ID: %v\n", changeInfo.Removed, changeInfo.Updated, changeInfo.Matched, changeInfo.UpsertedId)
	}

	print("\nAfter Upsert as update:\n")
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

	fmt.Printf("\nBefore Upserting:\n")
	showWord(collection)
	addUpsert(collection)
	updateUpsert(collection)

	fmt.Printf("Removing 'righty' ...\n")
	query := bson.M{"word": "righty"} // NOTE: the case of the letters does matter
	changeInfo, err := collection.RemoveAll(query)
	check(err)
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d, Matched: %d, Upserted ID: %v\n", changeInfo.Removed, changeInfo.Updated, changeInfo.Matched, changeInfo.UpsertedId)
	}
}
