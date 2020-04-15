package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/pkg/errors"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func displayDoc(doc bson.M) error {
	fmt.Printf("%v\n", doc)
	jsonString, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println("\nResult as JSON:")

	var out bytes.Buffer
	err = json.Indent(&out, jsonString, "", "  ")
	if err != nil {
		return errors.Wrap(err, "")
	}

	var st string = out.String()
	fmt.Printf("%v\n", st)
	return nil
}

func showWord(collection *mgo.Collection) error {
	var doc bson.M
	var words = []string{"left", "lefty"}
	query := bson.M{"word": bson.M{"$in": words}}
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}
	}
	err := iter.Close()
	return err
}

func (m *Mongo) updateDoc() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nBefore Updating:\n")
	err := showWord(collection)
	if err != nil {
		return errors.Wrap(err, "")
	}

	query := bson.M{"word": "left"}
	update := bson.M{
		"$set":  bson.M{"word": "lefty"},
		"$inc":  bson.M{"size": 1, "stats.consonants": 1},
		"$push": bson.M{"letters": "y"},
	}
	err = collection.Update(query, update)
	if err != nil {
		return errors.Wrap(err, "")
	}
	print("\nAfter Updating Doc:\n")
	return showWord(collection)
}

func (m *Mongo) resetDoc() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"word": "lefty"}
	update := bson.M{
		"$set": bson.M{"word": "left"},
		"$inc": bson.M{"size": -1, "stats.consonants": -1},
		"$pop": bson.M{"letters": 1},
	}
	err := collection.Update(query, update)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Printf("\nAfter Resetting Doc:\n")
	return showWord(collection)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	err = mongodb.updateDoc()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	err = mongodb.resetDoc()
	if err != nil {
		fmt.Printf("%+v\n", err)
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
