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

func (m *Mongo) includeFields(fields []string) error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

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
		return errors.Wrap(err, "")
	}
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"first": "p"}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	if err != nil {
		return errors.Wrap(err, "")
	}

	fmt.Printf("\nIncluding %v fields:\n", fields)
	return displayDoc(doc)
}

func (m *Mongo) excludeFields(fields []string) error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var doc bson.M

	var fieldObj bson.M

	var sel string
	if len(fields) > 0 {
		// place variable number of 'fields' into JSON style string
		sel = `{"`
		for _, field := range fields {
			sel += field + `": 0,"` // Add 'key' and a value of '0' to exclude this 'key' from the results
		}
		sel = sel[:len(sel)-2] + `}`
	} else {
		sel = `{}`
	}
	fmt.Printf("\nselect: %v\n", sel)

	// convert variable length JSON search string into format required by mongodb
	err := bson.UnmarshalJSON([]byte(sel), &fieldObj)
	if err != nil {
		return errors.Wrap(err, "")
	}
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"first": "p"}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	if err != nil {
		return errors.Wrap(err, "")
	}

	fmt.Printf("\nExcluding %v fields:\n", fields)
	return displayDoc(doc)
}

func main() {
	mongodb, err := GetMongoDB()
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		log.Fatal().Err(errors.New(errString)).Str("", "").Msgf("Database problem")
		// log.Fatal() above exits the program
	}

	defer mongodb.Session.Close()

	err = mongodb.excludeFields([]string{}) // show all fields
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.includeFields([]string{"word", "size"})
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.includeFields([]string{"word", "letters"})
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.excludeFields([]string{"letters", "stats", "charsets"})
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
