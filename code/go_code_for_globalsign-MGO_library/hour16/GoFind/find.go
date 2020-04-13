package main

import (
	"fmt"
	"log"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Client struct {
	session *mgo.Session
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (m *mongo) getOne() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var doc bson.M
	err := collection.Find(bson.M{}).One(&doc)
	check(err)
	fmt.Println("\nSingle Document:")
	fmt.Println(doc)
}

func (m *mongo) getManyFor() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Println("\nMany Using 'iter' Loop:")
	iter := collection.Find(bson.M{}).Iter()
	var doc bson.M
	var i = 0
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		fmt.Println(doc)
		i++
		if i >= 8 {
			break
		}
	}
	err := iter.Close()
	check(err)
}

func (m *mongo) getManySlice() {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Print("\nMany Using Skip & Limit + Loop:  ")
	var docs []bson.M
	start := time.Now()
	// set the cursor to skip the first 4, then span the next 4
	cursor := collection.Find(bson.M{}).Skip(4).Limit(4)
	// then get up to 4 documents at the cursor
	err := cursor.All(&docs)
	check(err)
	fmt.Printf("Search took %s\n", time.Since(start))
	var words []string
	for i := 0; i < 4; i++ {
		doc := docs[i]
		valueWord := doc["word"]
		switch v := valueWord.(type) { // do "type assertion" for field
		case string:
			//fmt.Printf("\nExtracted word is: '%v'\n", v)
			words = append(words, v)
		default:
			fmt.Println("word error")
		}
	}
	fmt.Print("Words: ")
	fmt.Println(words)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	mongodb.getOne()
	mongodb.getManyFor()
	mongodb.getManySlice()
}

const (
	mongoURI string = "127.0.0.1"
)

type mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
}

func GetMongoDB() (*mongo, error) {
	mongodb := &mongo{
		Collection: "word_stats",
		Database:   "words",
		URI:        mongoURI,
	}

	session, err := mongodb.Init()
	if err != nil {
		log.Printf("failed to initialise mongo %v", err)
		return nil, err
	}
	mongodb.Session = session

	return mongodb, nil
}

func (m *mongo) Init() (session *mgo.Session, err error) {
	if session, err = mgo.Dial(m.URI); err != nil {
		return nil, err
	}

	//	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	//	session.SetMode(mgo.Strong, true)
	return session, nil
}
