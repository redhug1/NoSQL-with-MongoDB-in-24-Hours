package main

import (
	"context"
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
func (m *Mongo) Ping(ctx context.Context) (time.Time, error) {
	mutex.Lock()
	if (pingInFlight == true) || (time.Since(m.lastPingTime) < 1*time.Second) {
		if pingInFlight == true {
			fmt.Printf("reject, as Ping is in Flight\n")
		}
		// reject re-entrant calls (should this function get called from different go routines)
		lpt := m.lastPingTime // protect from race
		lres := m.lastPingResult
		mutex.Unlock()
		return lpt, lres
	}
	pingInFlight = true
	mutex.Unlock()

	s := m.Session.Copy()
	defer func() {
		s.Close()
		pingInFlight = false
		mutex.Unlock()
		fmt.Printf("ping defer\n")
	}()

	mutex.Lock()
	m.lastPingTime = time.Now()
	mutex.Unlock()

	pingDoneChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Trace().Msg("db ping")
		start := time.Now()
		// NOTE: if at this point the mongodb stops / stops responding ...
		// (which is entirely possible when mongo db being accessed is
		//  on another server)
		// the following Ping will timeout after ~60 seconds and
		// return "no reachable servers" as the error string.
		err := s.Ping()
		log.Trace().Msgf("Ping took : %s", time.Since(start))
		if err != nil {
			log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("Ping mongo")
		} else {
			log.Trace().Msg("ping OK")
		}
		pingDoneChan <- err
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(pingDoneChan)
	}()
	fmt.Printf("Ping go returned\n")

	select {
	case err := <-pingDoneChan:
		mutex.Lock()
		m.lastPingResult = err
	case <-ctx.Done():
		mutex.Lock()
		m.lastPingResult = ctx.Err()
	}
	return m.lastPingTime, m.lastPingResult
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
	ctx, cancelHealthChecks := context.WithCancel(context.Background())

	for i := 0; i < 10000; i++ {
		go func(num int) {
			/*ptime, err := */ mongodb.Ping(ctx)

			//sfmt.Printf("%v : Last Ping'd at: %v, %v  :: current time %v\n", num, ptime, err, time.Now())
		}(i)
		time.Sleep(150 * time.Microsecond) // NOTE: Adjust this delay to suite your system speed
		// to small a value and you may not see any Ping results due to all 'go routines' completing within
		// ONE second ...
	}

	time.Sleep(3 * time.Second)

	log.Trace().Msg("Canceling healthchecks")
	cancelHealthChecks()
	<-ctx.Done()
	log.Trace().Msg("Canceled healthchecks")
}

// The following code is suitable for putting in its own file ...
// (if i had placed this file in the github path)

const (
	mongoURI string = "127.0.0.1"
)

type Mongo struct {
	Collection     string
	Database       string
	Session        *mgo.Session
	URI            string
	lastPingTime   time.Time
	lastPingResult error
}

func GetMongoDB() (*Mongo, error) {
	mongodb := &Mongo{
		Collection: "word_stats",
		Database:   "words",
		URI:        mongoURI,
	}

	mongodb.lastPingTime = time.Now()
	mongodb.lastPingResult = nil

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
