package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func (m *Mongo) addUpsert() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nBefore Upserting:\n")
	showWord(collection)

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

func (m *Mongo) updateUpsert() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

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

func (m *Mongo) removeRighty() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("Removing 'righty' ...\n")
	query := bson.M{"word": "righty"} // NOTE: the case of the letters does matter
	changeInfo, err := collection.RemoveAll(query)
	check(err)
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d, Matched: %d, Upserted ID: %v\n", changeInfo.Removed, changeInfo.Updated, changeInfo.Matched, changeInfo.UpsertedId)
	}
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer func() {
		fmt.Printf("Closing mongodb session\n")
		mongodb.Session.Close()
	}()

	mongodb.addUpsert()
	mongodb.updateUpsert()

	mongodb.removeRighty()
}

// The following code is suitable for putting in its own file ...
// (if i had placed this file in the github path)

const (
	mongoURI string = "127.0.0.1"
)

type Mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
}

func GetMongoDB() (*Mongo, error) {
	mongodb := &Mongo{
		Collection: "word_stats",
		Database:   "words",
		URI:        mongoURI,
	}

	session, err := mongodb.init()
	if err != nil {
		log.Printf("failed to initialise mongo %v", err)
		return nil, err
	}
	mongodb.Session = session

	return mongodb, nil
}

func (m *Mongo) init() (session *mgo.Session, err error) {
	if session, err = mgo.Dial(m.URI); err != nil {
		return nil, err
	}

	//	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	//	session.SetMode(mgo.Strong, true)
	return session, nil
}
