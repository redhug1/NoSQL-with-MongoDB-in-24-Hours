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

func (m *Mongo) over12() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\n\nWords with more than 12 characters:\n")
	cursor := collection.Find(bson.M{"size": bson.M{"$gt": 12}})
	displayCursor(cursor)
}

func (m *Mongo) startingABC() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nWords starting with A, B, C:\n")
	var abc = []string{"a", "b", "c"}
	query := bson.M{"first": bson.M{"$in": abc}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func (m *Mongo) startEndVowels() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nWords starting and ending with a vowel:\n")
	var vowels = []string{"a", "e", "i", "o", "u"}
	query := bson.M{"$and": []bson.M{
		bson.M{"first": bson.M{"$in": vowels}},
		bson.M{"last": bson.M{"$in": vowels}},
	}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func (m *Mongo) over6Vowels() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nWords with more than 5 vowels:\n")
	query := bson.M{"stats.vowels": bson.M{"$gt": 5}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func (m *Mongo) nonAlphaCharacters() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	print("\nWords with 1 non-alphabet character:\n")
	query := bson.M{"charsets": bson.M{"$elemMatch": bson.M{"$and": []bson.M{
		bson.M{"type": "other"},
		bson.M{"chars": bson.M{"$size": 1}}}}}}
	cursor := collection.Find(query)
	displayCursor(cursor)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.over12()
	mongodb.startingABC()
	mongodb.startEndVowels()
	mongodb.over6Vowels()
	mongodb.nonAlphaCharacters()
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
