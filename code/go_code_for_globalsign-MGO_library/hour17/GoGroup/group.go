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

func displayGroup(iter *mgo.Iter) error {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return err
		}
	}
	err := iter.Close()
	return err
}

func updateTotalAndDisplay(iter *mgo.Iter) error {
	var doc bson.M

	for iter.Next(&doc) {
		jsonString, err := json.MarshalIndent(doc, "", " ")
		if err != nil {
			iter.Close()
			return err
		}

		var fields map[string]interface{}
		err = json.Unmarshal([]byte(jsonString), &fields)
		if err != nil {
			iter.Close()
			return err
		}
		fmt.Printf("Before adding 'total', fields: %v\n", fields)

		vowelCount := fields["vowels"].(float64)
		consonantCount := fields["consonants"].(float64)
		total := vowelCount + consonantCount
		fields["total"] = total

		err = displayDoc(bson.M(fields))
		if err != nil {
			iter.Close()
			return err
		}
	}
	err := iter.Close()
	return err
}

func (m *Mongo) firstIsALastIsVowel() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var vowels = []string{"a", "e", "i", "o", "u"}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"$and": []bson.M{
			bson.M{"first": "a"},
			bson.M{"last": bson.M{"$in": vowels}}},
		},
		},
		bson.M{"$group": bson.M{"_id": "$last", // select the last letter to produce the count from
			"first": bson.M{"$first": "$first"},
			"last":  bson.M{"$first": "$last"},
			"count": bson.M{"$sum": 1}}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\n'A' words grouped by first and last letter that ends with a vowel:\n")
	return displayGroup(iter)
}

func (m *Mongo) firstLetterTotals() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"first": bson.M{"$in": alphabet}}},
		bson.M{"$group": bson.M{"_id": "$first",
			"vowels":     bson.M{"$sum": "$stats.vowels"},
			"consonants": bson.M{"$sum": "$stats.consonants"},
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nWords grouped by first letter with totals:\n")

	// the original python code created the 'total' field in this function, but
	// i can't see how to achieve that with the 'mgo' mongodb go library ...
	// So, the following function call digs into each document and achieves
	// the desired result in a programatic manner.
	return updateTotalAndDisplay(iter)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.firstIsALastIsVowel()
	if err != nil {
		log.Println(err)
		return
	}

	mongodb.firstLetterTotals()
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
