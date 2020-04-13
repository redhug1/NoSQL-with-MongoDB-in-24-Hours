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
	var words = []string{"left", "lefty"}
	query := bson.M{"word": bson.M{"$in": words}}
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		displayDoc(doc)
	}
}

func (m *Mongo) updateDoc() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nBefore Updating:\n")
	showWord(collection)

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

func (m *Mongo) resetDoc() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

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
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.updateDoc()
	mongodb.resetDoc()
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
