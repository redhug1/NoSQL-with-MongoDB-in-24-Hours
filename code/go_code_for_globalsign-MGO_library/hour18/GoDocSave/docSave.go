package main

import (
	"bytes"
	"encoding/json"
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

func displayDoc(doc bson.M) error {
	fmt.Printf("%v\n", doc)
	jsonString, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		return err
	}
	fmt.Println("\nResult as JSON:")

	var out bytes.Buffer
	err = json.Indent(&out, jsonString, "", "  ")
	if err != nil {
		return err
	}

	var st string = out.String()
	fmt.Printf("%v\n", st)
	return nil
}

func includeFieldsForWord(collection *mgo.Collection, word string, fields []string) error {
	var doc bson.M

	var fieldObj bson.M

	var sel string
	if len(fields) > 0 {
		// place variable number of 'fields' into JSON style string
		sel = `{"`
		for _, field := range fields {
			sel += field + `": 1,"` // Add 'key' and a value of '1' to include this 'key' in the results
		}
		sel = sel[:len(sel)-2] + `}`
	} else {
		sel = `{}`
	}
	fmt.Printf("\nselect: %v\n", sel)

	// convert variable length JSON search string into format required by mongodb
	err := bson.UnmarshalJSON([]byte(sel), &fieldObj)
	if err != nil {
		return err
	}
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"word": word}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	if err != nil {
		return err
	}

	fmt.Printf("\nIncluding %v fields:\n", fields)
	return displayDoc(doc)
}

func showWord(collection *mgo.Collection) error {
	return includeFieldsForWord(collection, "the", []string{"word", "the", "category"})
}

func (m *Mongo) saveBlueDoc() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("Before Saving:")
	err := showWord(collection)
	if err != nil {
		return err
	}
	query := bson.M{"word": "the"}
	update := bson.M{
		"$set": bson.M{"category": "blue"}, // 'Add' new field
	}
	err = collection.Update(query, update) // NOTE: the equivalent to the python 'save' is to do an 'update'
	if err != nil {
		return err
	}
	print("\nAfter Saving Doc:\n")
	return showWord(collection)
}

func (m *Mongo) resetDoc() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"word": "the"}
	update := bson.M{
		"$unset": bson.M{"category": nil}, // 'Remove' new field
	}
	err := collection.Update(query, update) // NOTE: the equivalent to the python 'save' is to do an 'update'
	if err != nil {
		return err
	}
	fmt.Printf("\nAfter Resetting Doc:\n")
	return showWord(collection)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	err = mongodb.saveBlueDoc()
	if err != nil {
		log.Println(err)
		return
	}

	err = mongodb.resetDoc()
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
