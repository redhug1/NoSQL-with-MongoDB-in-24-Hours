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

func showWord(collection *mgo.Collection) error {
	var doc bson.M
	query := bson.M{"word": "righty"}
	err := collection.Find(query).One(&doc)
	if err != nil {
		log.Printf("%v\n", err) // continue on as the word being searched for will not initially be in the list ...
	}
	return displayDoc(doc)
}

func (m *Mongo) addUpsert() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nBefore Upserting:\n")
	showWord(collection)

	var rLetters = []string{"r", "i", "g", "h"}
	var rConstChars = []string{"r", "g", "h"}
	var rVowelChars = []string{"i"}

	righty := bson.M{
		"word":     "righty",
		"first":    "r",
		"last":     "y",
		"size":     4,
		"category": "New",
		"stats":    bson.M{"vowels": 1, "consonants": 4},
		"letters":  rLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": rConstChars},
			bson.M{"type": "vowels", "chars": rVowelChars},
		}}
	err := collection.Insert(righty)
	if err != nil {
		return errors.Wrap(err, "")
	}

	print("\nAfter Upsert as insert:\n")
	return showWord(collection)
}

func (m *Mongo) updateUpsert() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var rLetters = []string{"r", "i", "g", "h", "t", "y"}
	var rConstChars = []string{"r", "g", "h", "t", "y"}
	var rVowelChars = []string{"i"}

	query := bson.M{"word": "righty"}
	update := bson.M{
		"$set": bson.M{
			"word":     "righty",
			"first":    "r",
			"last":     "y",
			"size":     6,
			"category": "Updated",
			"stats":    bson.M{"vowels": 1, "consonants": 5},
			"letters":  rLetters,
			"charsets": []bson.M{
				bson.M{"type": "consonants", "chars": rConstChars},
				bson.M{"type": "vowels", "chars": rVowelChars},
			}},
	}
	changeInfo, err := collection.Upsert(query, update)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d, Matched: %d, Upserted ID: %v\n", changeInfo.Removed, changeInfo.Updated, changeInfo.Matched, changeInfo.UpsertedId)
	}

	print("\nAfter Upsert as update:\n")
	return showWord(collection)
}

func (m *Mongo) removeRighty() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("Removing 'righty' ...\n")
	query := bson.M{"word": "righty"} // NOTE: the case of the letters does matter
	changeInfo, err := collection.RemoveAll(query)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if changeInfo != nil {
		fmt.Printf("Removed: %d, Updated: %d, Matched: %d, Upserted ID: %v\n", changeInfo.Removed, changeInfo.Updated, changeInfo.Matched, changeInfo.UpsertedId)
	}
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

	err = mongodb.addUpsert()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.updateUpsert()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.removeRighty()
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
