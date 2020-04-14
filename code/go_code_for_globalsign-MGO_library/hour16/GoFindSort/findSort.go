package main

import (
	"errors"
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

func displayCursor(cursor *mgo.Query) error {
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
	if err != nil {
		return err
	}
	if len(words) > 65 {
		words = words[:65] + "..."
	}
	fmt.Println(words)
	return nil
}

func (m *Mongo) sortWordsAscending() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("word")
	fmt.Printf("\nW words ordered ascending: ")
	return displayCursor(cursor)
}

func (m *Mongo) sortWordsDescending() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "w"}
	cursor := collection.Find(query)
	cursor.Sort("-word")
	fmt.Printf("\nW words ordered descending: ")
	return displayCursor(cursor)
}

func (m *Mongo) sortWordsAscAndSize() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"first": "q"}
	cursor := collection.Find(query)
	cursor.Sort("last", "-size")
	fmt.Printf("\nQ words ordered first by last letter and then by size: ")
	return displayCursor(cursor)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	err = mongodb.sortWordsAscending()
	if err != nil {
		log.Println(err)
		return
	}

	err = mongodb.sortWordsDescending()
	if err != nil {
		log.Println(err)
		return
	}

	err = mongodb.sortWordsAscAndSize()
	if err != nil {
		log.Println(err)
	}
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

	names, err := session.DB(mongodb.Database).CollectionNames()
	if err != nil {
		log.Printf("Failed to get collection names: %v", err)
		return nil, err
	}

	// look for required 'collection name' in slice ...
	var found bool = false
	for _, name := range names {
		if name == mongodb.Collection {
			found = true
			break
		}
	}
	if found == false {
		log.Printf("Can NOT find collection: %v, in Database: %v", mongodb.Collection, mongodb.Database)
		return nil, errors.New("Collection missing")
	}

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
