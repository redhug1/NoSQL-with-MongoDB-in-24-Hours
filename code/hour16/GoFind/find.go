package main

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getOne(collection *mgo.Collection) {
	var doc bson.M
	err := collection.Find(bson.M{}).One(&doc)
	check(err)
	fmt.Println("\nSingle Document:")
	fmt.Println(doc)
}

func getManyFor(collection *mgo.Collection) {
	fmt.Println("\nMany Using 'iter' Loop:")
	iter := collection.Find(bson.M{}).Iter()
	//	var words []string
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

func getManySlice(collection *mgo.Collection) {
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
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	getOne(collection)
	getManyFor(collection)
	getManySlice(collection)
}
