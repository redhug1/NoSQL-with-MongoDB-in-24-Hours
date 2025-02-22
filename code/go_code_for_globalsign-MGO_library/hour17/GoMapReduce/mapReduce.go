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

func displayGroup(iter *mgo.Iter) error {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}
	}
	err := iter.Close()
	return errors.Wrap(err, "")
}

func updateTotalAndDisplay(iter *mgo.Iter) error {
	var doc bson.M

	for iter.Next(&doc) {
		jsonString, err := json.MarshalIndent(doc, "", " ")
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}

		var fields map[string]interface{}
		err = json.Unmarshal([]byte(jsonString), &fields)
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}
		fmt.Printf("Before adding 'total', fields: %v\n", fields)

		vowelCount := fields["vowels"].(float64)
		consonantCount := fields["consonants"].(float64)
		total := vowelCount + consonantCount
		// Insert the new key and its value
		fields["characters"] = total

		err = displayDoc(bson.M(fields))
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}
	}
	err := iter.Close()
	return errors.Wrap(err, "")
}

func (m *Mongo) totalVowelBeginningCertainLetter() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	// ======================================================================================================
	// NOTE: this code the uses MapReduce, as far as any documentation goes, should work ... but does NOT ...

	var result []struct {
		Id    int "_id" // NOTE: this does not get filled properly !
		Value float64
	}

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.first, this.stats.vowels); }",
		Reduce: "function(key, values) { return Array.sum(values); }",
		//Out:    "results_collection",
	}

	query := bson.M{"first": bson.M{"$in": alphabet}}

	_, err := collection.Find(query).MapReduce(job, &result)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Printf("\n\nTotal vowel count in words beginning with a certain letter:\n")
	for i, item := range result {
		fmt.Printf("%T, %v", item, item)
		fmt.Printf("    %s : %v\n", alphabet[i], item.Value)
	}

	// ======================================
	// NOTE: So, we do it a different way ...

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"first": bson.M{"$in": alphabet}}},
		bson.M{"$group": bson.M{"_id": "$first",
			"value": bson.M{"$sum": "$stats.vowels"},
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nTotal vowel count in words beginning with a certain letter:\n")
	return displayGroup(iter)
}

func (m *Mongo) moreComplexMapReduce(session *mgo.Session) error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	var vowels = []string{"a", "e", "i", "o", "u"}

	// ======================================================================================================
	// NOTE: this code the uses MapReduce, as far as any documentation goes, should work ... but does NOT ...

	var result []struct {
		Id    int "_id"
		Value float64
	}

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.first, { vowels: this.stats.vowels, consonants: this.stats.consonants} ); }",
		Reduce: "function(key, values) { result = {count: values.length, vowels: 0, consonants: 0}; for(var i=0; i<values.length; i++){ if (values[i].vowels) result.vowels += values[i].vowels; if (values[i].consonants) result.consonants += values[i].consonants; } return result; }",
		//		Out:    "results_collection",
		Finalize: "function(key, obj) { obj.characters = obj.vowels + obj.consonants; return obj; }",
	}

	query := bson.M{"last": bson.M{"$in": vowels}}

	_, err := collection.Find(query).MapReduce(job, &result)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Printf("\nTotal words, vowels, consonants and characters in words beginning with a certain letter that ends with a vowel:\n")
	for _, item := range result {
		fmt.Printf("%T, %v, %v\n", item, item, item.Value)
	}

	// ======================================
	// NOTE: So, we do it a different way ...

	query = bson.M{"$and": []bson.M{
		bson.M{"last": bson.M{"$in": vowels}},
		bson.M{"first": bson.M{"$in": alphabet}},
	}}

	cursor := collection.Find(query)

	// Clear the collection from any previous run
	session.DB("words").C("temp_stats").RemoveAll(bson.M{})

	// Save result of query into temporary collection
	newCollection := session.DB("words").C("temp_stats")
	var doc bson.M
	iter := cursor.Iter()
	for iter.Next(&doc) {
		err := newCollection.Insert(&doc)
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "")
		}
	}
	err = iter.Close()
	if err != nil {
		return errors.Wrap(err, "")
	}

	// Do the next [art] of the 'map reduce' using the temporary collection
	collection = session.DB("words").C("temp_stats")

	pipeline := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$group": bson.M{"_id": "$first",
			// NOTE: this count value produced for the 10K word list is correct
			// and the value in the output for the Javascript and Python Code is wrong ... hmmm
			// For example, in the 10K word list, there are 9 words beginning with 'z' and ending in a vowel, thus:
			// ["zone","zero","zoo","zimbabwe","zoophilia","zope","zambia","za","zu"]
			// and similarly, there are 12 words beginning with 'y' and ending in a vowel, thus:
			// ["you","yahoo","ya","yoga","yorkshire","ye","yo","yamaha","yu","yea","yale","yugoslavia"]
			"count":      bson.M{"$sum": 1},
			"vowels":     bson.M{"$sum": "$stats.vowels"},
			"consonants": bson.M{"$sum": "$stats.consonants"},
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
	}
	iter = collection.Pipe(pipeline).Iter()
	fmt.Printf("\nTotal words, vowels, consonants and characters in words beginning with a certain letter that ends with a vowel:\n")

	// The final part of the 'map reduce' ...
	// The original python code created the 'characters' field in this function, but
	// i can't see how to achieve that with the 'mgo' mongodb go library ...
	// So, the following function call digs into each document and achieves
	// the desired result in a programatic manner.
	return updateTotalAndDisplay(iter)
}

func main() {
	mongodb, err := GetMongoDB()
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		log.Fatal().Err(errors.New(errString)).Str("", "").Msgf("Database problem")
		// log.Fatal() above exits the program
	}

	defer mongodb.Session.Close()

	err = mongodb.totalVowelBeginningCertainLetter()
	if err != nil {
		log.Printf("%+v", err)
		log.Error().Err(errors.New(fmt.Sprintf("%+v", err))).Msgf("")
		return // do this so that 'defer' gets done
	}

	err = mongodb.moreComplexMapReduce(mongodb.Session)
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
