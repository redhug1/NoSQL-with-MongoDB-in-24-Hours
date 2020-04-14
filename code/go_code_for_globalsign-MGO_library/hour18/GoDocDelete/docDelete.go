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

func findSpecificWords(collection *mgo.Collection) error {
	var abc = []string{"tweet", "gogle", "selfie", "jimmmy"}
	query := bson.M{"word": bson.M{"$in": abc}}
	cursor := collection.Find(query)
	return displayCursor(cursor)
}

func showNewDocs(collection *mgo.Collection) error {
	var doc bson.M
	query := bson.M{"category": "New"} // NOTE: the case of the letters does matter
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return err
		}
	}
	err := iter.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Showing structure of document for word 'the' written by javascript to check that the ones written by this go program are the same ...\n")
	fmt.Printf("You need to do a visual check / comparison !\n")
	query = bson.M{"word": "the"} // NOTE: the case of the letters does matter
	cursor = collection.Find(query)
	// Show all the doc's found ...
	iter = cursor.Iter()
	for iter.Next(&doc) {
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return err
		}
	}
	err = iter.Close()
	if err != nil {
		return err
	}

	return findSpecificWords(collection) // added to just show the word of interest
}

func (m *Mongo) removeNewDocs() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\n\nBefore Deleting:\n")
	err := showNewDocs(collection)
	if err != nil {
		return err
	}

	query := bson.M{"category": "New"} // NOTE: the case of the letters does matter
	changeInfo, err := collection.RemoveAll(query)
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d \n", changeInfo.Removed, changeInfo.Updated)
	}
	if err != nil {
		return err
	}
	fmt.Printf("\nAfter Deleting:\n")
	return showNewDocs(collection)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	err = mongodb.removeNewDocs()
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
