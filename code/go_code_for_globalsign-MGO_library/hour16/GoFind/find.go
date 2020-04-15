package main

import (
	"fmt"
	"os"
	"time"

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

func (m *Mongo) getOne() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var doc bson.M
	err := collection.Find(bson.M{}).One(&doc)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println("\nSingle Document:")
	fmt.Println(doc)
	return nil
}

func (m *Mongo) getManyFor() error {
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
	return errors.Wrap(err, "")
}

func (m *Mongo) getManySlice() error {
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
	if err != nil {
		return errors.Wrap(err, "")
	}

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
	return nil
}

func main() {
	mongodb, err := GetMongoDB()
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		log.Fatal().Err(errors.New(errString)).Str("", "").Msgf("Database problem")
		// log.Fatal() above exits the program
	}

	defer mongodb.Session.Close()

	err = mongodb.getOne()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.getManyFor()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.getManySlice()
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
