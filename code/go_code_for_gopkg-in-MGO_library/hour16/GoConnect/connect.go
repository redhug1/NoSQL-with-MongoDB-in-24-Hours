package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of Documents:", count)
}
