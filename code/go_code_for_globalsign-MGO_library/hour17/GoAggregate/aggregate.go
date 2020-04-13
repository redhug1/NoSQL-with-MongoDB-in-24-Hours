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

func displayAggregate(iter *mgo.Iter) {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		displayDoc(doc)
	}
	err := iter.Close()
	check(err)
}

func (m *Mongo) largeSmallVowels() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var vowels = []string{"a", "e", "i", "o", "u"}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"first": bson.M{"$in": vowels}}},
		bson.M{"$group": bson.M{"_id": "$first",
			"largest":  bson.M{"$max": "$size"},
			"smallest": bson.M{"$min": "$size"},
			"total":    bson.M{"$sum": 1}}},
		bson.M{"$sort": bson.M{"first": 1}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nLargest and smallest word sizes for word begining with a vowel:\n")
	displayAggregate(iter)
}

func (m *Mongo) top5AverageWordFirst() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	pipeline := []bson.M{
		bson.M{"$group": bson.M{"_id": "$first",
			"average": bson.M{"$avg": "$size"}}},
		bson.M{"$sort": bson.M{"average": -1}},
		bson.M{"$limit": 5},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nFirst letter of top 5 largest average word size:\n")
	displayAggregate(iter)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.largeSmallVowels()
	mongodb.top5AverageWordFirst()
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
