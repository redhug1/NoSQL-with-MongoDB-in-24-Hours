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
	var abc = []string{"tweet", "gogle", "selfie", "jimmmy"}
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
	fmt.Printf("Showing structure of document for word 'the' written by javascript to check that the ones written by this go program are the same ...\n")
	fmt.Printf("You need to do a visual check / comparison !\n")
	query = bson.M{"word": "the"} // NOTE: the case of the letters does matter
	cursor = collection.Find(query)
	// Show all the doc's found ...
	iter = cursor.Iter()
	for iter.Next(&doc) {
		displayDoc(doc)
	}

	findSpecificWords(collection) // added to just show the word of interest
}

func removeNewDocs(collection *mgo.Collection) {
	query := bson.M{"category": "New"} // NOTE: the case of the letters does matter
	changeInfo, err := collection.RemoveAll(query)
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d \n", changeInfo.Removed, changeInfo.Updated)
	}
	check(err)
	fmt.Printf("\nAfter Deleting:\n")
	showNewDocs(collection)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer func() {
		fmt.Printf("Closing mongodb session\n")
		session.Close()
	}()

	collection := session.DB("words").C("word_stats")

	fmt.Printf("\n\nBefore Deleting:\n")
	showNewDocs(collection)
	removeNewDocs(collection)
}
