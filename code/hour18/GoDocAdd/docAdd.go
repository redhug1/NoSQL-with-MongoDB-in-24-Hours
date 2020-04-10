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

func findSpecificWords(collection *mgo.Collection) {
	var abc = []string{"tweet", "gogle", "selfie"}
	query := bson.M{"word": bson.M{"$in": abc}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func showNewDocs(collection *mgo.Collection) {
	var doc bson.M
	query := bson.M{"category": "New"} // NOTE: the case of the letters does matter
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		displayDoc(doc)
	}
	findSpecificWords(collection) // added to just show the word of interest
}

func addSelfie(collection *mgo.Collection) {
	var letters = []string{"s", "e", "l", "f", "i"}
	var constChars = []string{"s", "l", "f"}
	var vowelChars = []string{"e", "i"}
	selfie := bson.M{"word": "selfie",
		"first":    "s",
		"last":     "e",
		"size":     6,
		"category": "New",
		"stats":    []bson.M{bson.M{"vowels": 3}, bson.M{"consonants": 3}},
		"letters":  letters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": constChars},
			bson.M{"type": "vowels", "chars": vowelChars},
		}}
	fmt.Printf("About to insert ...\n")
	err := collection.Insert(selfie)
	check(err)
	fmt.Printf("After Inserting One:\n")
	showNewDocs(collection)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	fmt.Printf("\nBefore Inserting:\n")
	showNewDocs(collection)
	addSelfie(collection)
	//	addGoogleAndTweet(collection)
}
