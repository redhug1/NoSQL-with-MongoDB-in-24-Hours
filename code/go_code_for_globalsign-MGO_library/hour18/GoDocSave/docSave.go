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

func includeFieldsForWord(collection *mgo.Collection, word string, fields []string) {
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
	check(err)
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"word": word}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	check(err)

	fmt.Printf("\nIncluding %v fields:\n", fields)
	displayDoc(doc)
}

func showWord(collection *mgo.Collection) {
	includeFieldsForWord(collection, "the", []string{"word", "the", "category"})
}

func (m *Mongo) saveBlueDoc() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("Before Saving:")
	showWord(collection)
	query := bson.M{"word": "the"}
	update := bson.M{
		"$set": bson.M{"category": "blue"}, // 'Add' new field
	}
	err := collection.Update(query, update) // NOTE: the equivalent to the python 'save' is to do an 'update'
	check(err)
	print("\nAfter Saving Doc:\n")
	showWord(collection)
}

func (m *Mongo) resetDoc() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	query := bson.M{"word": "the"}
	update := bson.M{
		"$unset": bson.M{"category": nil}, // 'Remove' new field
	}
	err := collection.Update(query, update) // NOTE: the equivalent to the python 'save' is to do an 'update'
	check(err)
	fmt.Printf("\nAfter Resetting Doc:\n")
	showWord(collection)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.saveBlueDoc()
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
