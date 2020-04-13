package main

import (
	"fmt"
	"log"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func (m *Mongo) sortWordsAscending() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("word")
	fmt.Printf("\nW words ordered ascending: ")
	displayCursor(cursor)
}

func (m *Mongo) sortWordsDescending() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("-word")
	fmt.Printf("\nW words ordered descending: ")
	displayCursor(cursor)
}

func (m *Mongo) sortWordsAscAndSize() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "q"}
	cursor := collection.Find(query)
	cursor.Sort("last", "-size")
	fmt.Printf("\nQ words ordered first by last letter and then by size: ")
	displayCursor(cursor)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.sortWordsAscending()
	mongodb.sortWordsDescending()
	mongodb.sortWordsAscAndSize()
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
