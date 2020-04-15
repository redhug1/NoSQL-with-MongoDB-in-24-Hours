package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

func displayAggregate(iter *mgo.Iter) error {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		err := displayDoc(doc)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}
	err := iter.Close()
	return errors.Wrap(err, "")
}

func (m *Mongo) largeSmallVowels() error {
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
	return displayAggregate(iter)
}

func (m *Mongo) top5AverageWordFirst() error {
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
	return displayAggregate(iter)
}

func main() {
	mongodb, err := GetMongoDB()
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		log.Fatal().Err(errors.New(errString)).Str("", "").Msgf("Database problem")
		// log.Fatal() above exits the program
	}

	defer mongodb.Session.Close()

	err = mongodb.largeSmallVowels()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}
	err = mongodb.top5AverageWordFirst()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
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
		// no session to close
		log.Printf("failed to initialise mongo %v", err)
		return nil, err
	}
	mongodb.Session = session

	names, err := session.DB(mongodb.Database).CollectionNames()
	if err != nil {
		session.Close()
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
		session.Close()
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
