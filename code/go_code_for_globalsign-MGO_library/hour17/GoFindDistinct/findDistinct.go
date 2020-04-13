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

func (m *Mongo) sizesOfAllWords() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var result []int
	err := collection.Find(nil).Distinct("size", &result)
	check(err)
	fmt.Printf("\nDistinct Sizes of words: %v\n", result)
}

func (m *Mongo) sizesOfQWords() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var result []int
	query := bson.M{"first": "q"}
	cursor := collection.Find(query)
	err := cursor.Distinct("size", &result)
	check(err)
	fmt.Printf("\nDistinct Sizes of words starting with Q: %v\n", result)
}

func (m *Mongo) firstLetterOfLongWords() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var result []string
	query := bson.M{"size": bson.M{"$gt": 12}}
	cursor := collection.Find(query)
	err := cursor.Distinct("first", &result)
	check(err)
	fmt.Printf("\nDistinct first letters of words longer than 12 characters: %v\n", result)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.sizesOfAllWords()
	mongodb.sizesOfQWords()
	mongodb.firstLetterOfLongWords()
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
