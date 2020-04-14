package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

func displayCursor(cursor *mgo.Query) error {
	var doc bson.M
	var words string
	iter := cursor.Iter()
	for iter.Next(&doc) {
		valueWord := doc["word"]
		switch v := valueWord.(type) { // do "type assertion" for field
		case string:
			if len(words) > 0 {
				words = words + "," + v
			} else {
				words = v
			}
		default:
			fmt.Println("word error")
		}
	}
	err := iter.Close()
	if err != nil {
		return err
	}
	if len(words) > 65 {
		words = words[:65] + "..."
	}
	fmt.Println(words)
	return nil
}

func displayDoc(doc bson.M) error {
	fmt.Printf("%v\n", doc)
	jsonString, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		return err
	}
	fmt.Println("\nResult as JSON:")

	var out bytes.Buffer
	err = json.Indent(&out, jsonString, "", "  ")
	if err != nil {
		return err
	}

	var st string = out.String()
	fmt.Printf("%v\n", st)
	return nil
}

func findSpecificWords(collection *mgo.Collection) error {
	var abc = []string{"tweet", "gogle", "selfie", "jimmmy"}
	query := bson.M{"word": bson.M{"$in": abc}}
	cursor := collection.Find(query)
	return displayCursor(cursor)
}

func showNewDocs(collection *mgo.Collection) error {
	var doc bson.M
	query := bson.M{"category": "New"} // NOTE: the case of the letters does matter
	cursor := collection.Find(query)
	iter := cursor.Iter()
	for iter.Next(&doc) {
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return err
		}
	}
	err := iter.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Showing structure of document for word 'the' written by javascript to check that the ones written by this go program are the same ...\n")
	fmt.Printf("You need to do a visual check / comparison !\n")
	query = bson.M{"word": "the"} // NOTE: the case of the letters does matter
	cursor = collection.Find(query)
	// Show all the doc's found ...
	iter = cursor.Iter()
	for iter.Next(&doc) {
		err := displayDoc(doc)
		if err != nil {
			iter.Close()
			return err
		}
	}
	err = iter.Close()
	if err != nil {
		return err
	}

	return findSpecificWords(collection) // added to just show the word of interest
}

func (m *Mongo) addSelfie() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	fmt.Printf("\nBefore Inserting:\n")
	err := showNewDocs(collection)
	if err != nil {
		return err
	}

	var letters = []string{"s", "e", "l", "f", "i"}
	var constChars = []string{"s", "l", "f"}
	var vowelChars = []string{"e", "i"}
	selfie := bson.M{"word": "selfie",
		"first":    "s",
		"last":     "e",
		"size":     6,
		"category": "New",
		"stats":    bson.M{"vowels": 3, "consonants": 3},
		"letters":  letters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": constChars},
			bson.M{"type": "vowels", "chars": vowelChars},
		}}
	fmt.Printf("About to insert ...\n")
	err = collection.Insert(selfie)
	if err != nil {
		return err
	}
	fmt.Printf("After Inserting One:\n")
	return showNewDocs(collection)
}

func (m *Mongo) addGoogleAndTweet() error {
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	// deliberate mis-spelling as google is already in the 100K list of words
	var gLetters = []string{"g", "o", "l", "e"}
	var gConstChars = []string{"g", "l"}
	var gVowelChars = []string{"o", "e"}
	gogle := bson.M{
		"word":     "gogle",
		"first":    "g",
		"last":     "e",
		"size":     6,
		"category": "New",
		"stats":    bson.M{"vowels": 2, "consonants": 3},
		"letters":  gLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": gConstChars},
			bson.M{"type": "vowels", "chars": gVowelChars},
		}}
	var tLetters = []string{"t", "w", "e"}
	var tConstChars = []string{"t", "w"}
	var tVowelChars = []string{"e"}
	tweet := bson.M{
		"word":     "tweet",
		"first":    "t",
		"last":     "t",
		"size":     5,
		"category": "New",
		"stats":    bson.M{"vowels": 2, "consonants": 3},
		"letters":  tLetters,
		"charsets": []bson.M{
			bson.M{"type": "consonants", "chars": tConstChars},
			bson.M{"type": "vowels", "chars": tVowelChars},
		}}
	fmt.Printf("About to insert multiple ...\n")
	// Add multiple documents 'ONE' at a time as 'mgo' lib only does one at a time ...
	err := collection.Insert(gogle)
	if err != nil {
		return err
	}
	err = collection.Insert(tweet)
	if err != nil {
		return err
	}
	fmt.Printf("After Inserting Multiple:\n")
	return showNewDocs(collection)
}

func (m *Mongo) addJimmmyViaStruct() error { // thats 3 m's in Jimmmy, to ensure word is not already in the 10K list
	s := m.Session.Copy()
	defer s.Close()

	collection := s.DB(m.Database).C(m.Collection)

	type StatsType struct {
		Vowels     float64 `bson:"vowels"`
		Consonants float64 `bson:"consonants"`
	}
	type CharsetsType struct {
		Type  string   `bson:"type"`
		Chars []string `bson:"chars"`
	}
	type DocStruct struct {
		Word     string         `bson:"word"`
		First    string         `bson:"first"`
		Last     string         `bson:"last"`
		Size     float64        `bson:"size"`
		Category string         `bson:"category"`
		Stats    StatsType      `bson:"stats"`
		Letters  []string       `bson:"letters"`
		Charsets []CharsetsType `bson:"charsets"`
	}
	jimmmy := DocStruct{}
	jimmmy.Word = "jimmmy"
	jimmmy.First = "j"
	jimmmy.Last = "y"
	jimmmy.Size = 6
	jimmmy.Category = "New"
	jimmmy.Stats.Vowels = 1
	jimmmy.Stats.Consonants = 5
	jimmmy.Letters = []string{"j", "i", "m", "y"}
	jimmmy.Charsets = []CharsetsType{
		{Type: "consonants", Chars: []string{"j", "m", "y"}},
		{Type: "vowels", Chars: []string{"i"}},
	}
	fmt.Printf("About to insert 'jimmmy' via 'go' structure ...\n")
	err := collection.Insert(jimmmy)
	if err != nil {
		return err
	}
	fmt.Printf("After Inserting 'jimmmy':\n")
	return showNewDocs(collection)
}

func main() {
	mongodb, err := GetMongoDB()
	check(err)
	defer mongodb.Session.Close()

	err = mongodb.addSelfie()
	if err != nil {
		log.Println(err)
		return
	}

	err = mongodb.addGoogleAndTweet()
	if err != nil {
		log.Println(err)
		return
	}

	err = mongodb.addJimmmyViaStruct()
	if err != nil {
		log.Println(err)
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
		log.Printf("failed to initialise mongo %v", err)
		return nil, err
	}
	mongodb.Session = session

	names, err := session.DB(mongodb.Database).CollectionNames()
	if err != nil {
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
