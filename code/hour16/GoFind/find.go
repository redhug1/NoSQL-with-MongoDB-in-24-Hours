package main

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getOne(collection *mgo.Collection) {
	var doc bson.M
	err := collection.Find(bson.M{}).One(&doc)
	if err != nil {
		log.Fatal(err)
	}
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
}

func getManySlice(collection *mgo.Collection) {
	fmt.Print("\nMany Using Skip & Limit + Loop:  ")
	start := time.Now()
	cursor := collection.Find(bson.M{}).Skip(4).Limit(4)
	var docs []bson.M
	if err := cursor.All(&docs); err != nil {
		log.Fatal(err)
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
	fmt.Println("Words:")
	fmt.Println(words)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	getOne(collection)
	getManyFor(collection)
	getManySlice(collection)
}
