package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func displayDoc(doc bson.M) {
	fmt.Printf("%v\n", doc)
	jsonString, err := json.MarshalIndent(doc, "", " ")
	check(err)
	fmt.Println("\nResult as JSON:")

	var out bytes.Buffer
	err = json.Indent(&out, jsonString, "", "  ")
	check(err)

	var st string = out.String()
	fmt.Printf("%v\n", st)
}

func includeFields(collection *mgo.Collection, fields []string) {
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
	check(err)
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"first": "p"}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	check(err)

	fmt.Printf("\nIncluding %v fields:\n", fields)
	displayDoc(doc)
}

func excludeFields(collection *mgo.Collection, fields []string) {
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
	check(err)
	//fmt.Printf("bson %v\n", fieldObj)

	query := bson.M{"first": "p"}
	err = collection.Find(query).Select(fieldObj).One(&doc)
	check(err)

	fmt.Printf("\nExcluding %v fields:\n", fields)
	displayDoc(doc)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	excludeFields(collection, []string{}) // show all fields
	includeFields(collection, []string{"word", "size"})
	includeFields(collection, []string{"word", "letters"})
	excludeFields(collection, []string{"letters", "stats", "charsets"})
}
