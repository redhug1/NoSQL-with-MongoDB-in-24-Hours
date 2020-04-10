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

func displayGroup(iter *mgo.Iter) {
	var doc bson.M
	for iter.Next(&doc) {
		fmt.Println("Document is:")
		displayDoc(doc)
	}
	err := iter.Close()
	check(err)
}

func updateTotalAndDisplay(iter *mgo.Iter) {
	var doc bson.M

	defer func() {
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for iter.Next(&doc) {
		jsonString, err := json.MarshalIndent(doc, "", " ")
		check(err)

		var fields map[string]interface{}
		err = json.Unmarshal([]byte(jsonString), &fields)
		check(err)
		fmt.Printf("Before adding 'total', fields: %v\n", fields)

		vowelCount := fields["vowels"].(float64)
		consonantCount := fields["consonants"].(float64)
		total := vowelCount + consonantCount
		// Insert the new key and its value
		fields["characters"] = total

		displayDoc(bson.M(fields))
	}
}

func totalVowelBeginningCertainLetter(collection *mgo.Collection) {
	/* NOTE: this commented out code, as far as any documentation goes, should work ... but does NOT ...
	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	var result []struct {
		Id    int "_id"
		Value float64
	}

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.first, this.stats.vowels); }",
		Reduce: "function(key, values) { return Array.sum(values); }",
		Out:    "results_collection",
	}

	query := bson.M{"first": bson.M{"$in": alphabet}}

	_, err := collection.Find(query).MapReduce(job, &result)	// !!! this fails
	check(err)
	print("\n\nTotal vowel count in words beginning with a certain letter:")
	for _, item := range result {
		fmt.Println(item.Value)
	}
	*/
	/* NOTE: So, we do it a different way ... */
	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"first": bson.M{"$in": alphabet}}},
		bson.M{"$group": bson.M{"_id": "$first",
			"value": bson.M{"$sum": "$stats.vowels"},
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
	}
	iter := collection.Pipe(pipeline).Iter()
	fmt.Printf("\nTotal vowel count in words beginning with a certain letter:\n")
	displayGroup(iter)
}

func moreComplexMapReduce(collection *mgo.Collection, session *mgo.Session) {
	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	var vowels = []string{"a", "e", "i", "o", "u"}

	query := bson.M{"$and": []bson.M{
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
		check(err)
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
	updateTotalAndDisplay(iter)
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	check(err)
	defer session.Close()

	collection := session.DB("words").C("word_stats")

	totalVowelBeginningCertainLetter(collection)
	moreComplexMapReduce(collection, session)
}
