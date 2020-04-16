package main

import (
	"fmt"
	"os"
	"sync"
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

var (
	mutex        sync.Mutex
	pingInFlight bool = false
)

// Ping the mongodb database
func (m *Mongo) Ping() error {
	mutex.Lock()
	if pingInFlight == true {
		mutex.Unlock()
		return nil // reject re-entrant calls (should this function get called from different go routines)
	}
	pingInFlight = true
	mutex.Unlock()

	s := m.Session.Copy()
	defer func() {
		s.Close()
		mutex.Lock()
		pingInFlight = false
		mutex.Unlock()
	}()

	pingDoneChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("db ping")
		start := time.Now()
		// NOTE: if at this point the mongodb stops / stops responding ...
		// the following Ping will timeout after ~50 seconds and
		// return "no reachable servers" as the error string.
		err := s.Ping()
		log.Printf("Ping took : %s", time.Since(start))
		if err != nil {
			log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("Ping mongo")
		} else {
			log.Printf("ping OK")
		}
		pingDoneChan <- err
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(pingDoneChan)
	}()

	var result error
	select {
	case err := <-pingDoneChan:
		result = err
	}
	return result
}

func main() {
	mongodb, err := GetMongoDB()
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		log.Fatal().Err(errors.New(errString)).Str("", "").Msgf("Database problem")
		// log.Fatal() above exits the program
	}

	defer mongodb.Session.Close()

	collection := mongodb.Session.DB(mongodb.Database).C(mongodb.Collection)

	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	fmt.Println("Number of Documents:", count)
	mongodb.Ping()
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
