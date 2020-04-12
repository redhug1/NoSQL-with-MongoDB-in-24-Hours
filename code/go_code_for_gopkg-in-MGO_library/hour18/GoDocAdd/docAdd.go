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

func addSelfie(collection *mgo.Collection) {
	var letters = []string{"s", "e", "l", "f", "i"}
	var constChars = []string{"s", "l", "f"}
	var vowelChars = []string{"e", "i"}
	selfie := bson.M{"word": "selfie",
		"first":    "s",
		"last":     "e",
		"size":     6,
		"category": "New",
		"stats":    bson.M{"vowels": 3, "consonants": 3},
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

func addGoogleAndTweet(collection *mgo.Collection) {
	// deliberate mis-spelling as google is already in the 100K list of words
	var gLetters = []string{"g", "o", "l", "e"}
	var gConstChars = []string{"g", "l"}
	var gVowelChars = []string{"o", "e"}
	gogle := bson.M{
		"word":     "gogle",
		"first":    "g",
		"last":     "e",
		"size":     6,
		"category": "New",
		"stats":    bson.M{"vowels": 2, "consonants": 3},
		"letters":  gLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": gConstChars},
			bson.M{"type": "vowels", "chars": gVowelChars},
		}}
	var tLetters = []string{"t", "w", "e"}
	var tConstChars = []string{"t", "w"}
	var tVowelChars = []string{"e"}
	tweet := bson.M{
		"word":     "tweet",
		"first":    "t",
		"last":     "t",
		"size":     5,
		"category": "New",
		"stats":    bson.M{"vowels": 2, "consonants": 3},
		"letters":  tLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": tConstChars},
			bson.M{"type": "vowels", "chars": tVowelChars},
		}}
	fmt.Printf("About to insert multiple ...\n")
	// Add multiple documents 'ONE' at a time as 'mgo' lib only does one at a time ...
	err := collection.Insert(gogle)
	check(err)
	err = collection.Insert(tweet)
	check(err)
	fmt.Printf("After Inserting Multiple:\n")
	showNewDocs(collection)
}

func addJimmmyViaStruct(collection *mgo.Collection) { // thats 3 m's in Jimmmy, to ensure word is not already in the 10K list
	type StatsType struct {
		Vowels     float64 `bson:"vowels"`
		Consonants float64 `bson:"consonants"`
	}
	type CharsetsType struct {
		Type  string   `bson:"type"`
		Chars []string `bson:"chars"`
	}
	type DocStruct struct {
		Word     string         `bson:"word"`
		First    string         `bson:"first"`
		Last     string         `bson:"last"`
		Size     float64        `bson:"size"`
		Category string         `bson:"category"`
		Stats    StatsType      `bson:"stats"`
		Letters  []string       `bson:"letters"`
		Charsets []CharsetsType `bson:"charsets"`
	}
	jimmmy := DocStruct{}
	jimmmy.Word = "jimmmy"
	jimmmy.First = "j"
	jimmmy.Last = "y"
	jimmmy.Size = 6
	jimmmy.Category = "New"
	jimmmy.Stats.Vowels = 1
	jimmmy.Stats.Consonants = 5
	jimmmy.Letters = []string{"j", "i", "m", "y"}
	jimmmy.Charsets = []CharsetsType{
		{Type: "consonants", Chars: []string{"j", "m", "y"}},
		{Type: "vowels", Chars: []string{"i"}},
	}
	fmt.Printf("About to insert 'jimmmy' via 'go' structure ...\n")
	err := collection.Insert(jimmmy)
	check(err)
	fmt.Printf("After Inserting 'jimmmy':\n")
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

	fmt.Printf("\nBefore Inserting:\n")
	showNewDocs(collection)
	addSelfie(collection)
	addGoogleAndTweet(collection)
	addJimmmyViaStruct(collection)
}
