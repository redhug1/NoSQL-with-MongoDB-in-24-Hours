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
		log.Fatal(err)
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

func displayGroup(iter *mgo.Iter) {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		displayDoc(doc)
	}
	err := iter.Close()
	check(err)
}

func updateTotalAndDisplay(iter *mgo.Iter) {
	var doc bson.M

	defer func() {
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for iter.Next(&doc) {
		jsonString, err := json.MarshalIndent(doc, "", " ")
		check(err)

		var fields map[string]interface{}
		err = json.Unmarshal([]byte(jsonString), &fields)
		check(err)
		fmt.Printf("Before adding 'total', fields: %v\n", fields)

		vowelCount := fields["vowels"].(float64)
		consonantCount := fields["consonants"].(float64)
		total := vowelCount + consonantCount
		fields["total"] = total

		displayDoc(bson.M(fields))
	}
}

func (m *Mongo) firstIsALastIsVowel() {
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
	displayGroup(iter)
}

func (m *Mongo) firstLetterTotals() {
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
	updateTotalAndDisplay(iter)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.firstIsALastIsVowel()
	mongodb.firstLetterTotals()
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
