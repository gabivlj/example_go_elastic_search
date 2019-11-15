package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"users":{
				"type":"keyword"
			},
			"user":{
				"type":"keyword"
			},
			"textContent":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"date":{
				"type":"date"
			},
			"indexPosition": {
				"type":"integer"
			}
		}
	}
}`

func main() {

	ctx := context.Background()
	usersArray := []string{"gabivlj", "gabivlj2"}

	var (
		wg sync.WaitGroup
	)

	client, err := elastic.NewClient()

	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("chat").Do(ctx)
	if err != nil {
		panic(err)
	}

	if exists == false {
		createIndex, err := client.CreateIndex("chat").BodyString(mapping).Do(ctx)
		// createIndex, err := client.CreateIndex("chat").Do(ctx)
		fmt.Println(mapping)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	// req, err := http.NewRequest("PUT", "http://127.0.0.1:9200/chat", bytes.NewBuffer([]byte(mapping)))
	// req.Header.Add("Content-Type", "application/json")
	// c := &http.Client{}
	// resp, err := c.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	defer resp.Body.Close()
	// 	content, _ := ioutil.ReadAll(resp.Body)
	// 	fmt.Print(string(content))
	// }
	r := RandomStrings(10)
	for i, message := range r {
		(func(message string, i int) {
			wg.Add(1)
			defer wg.Done()
			chitchat := Message{Users: usersArray, User: ChooseString(usersArray), Text: message, IndexPosition: int64(i + 10), Date: time.Now()}
			if err != nil {
				panic(err)
			}
			put, err :=
				client.Index().Index("chat").Type("_doc").BodyJson(chitchat).Do(ctx)
			if err != nil {
				panic(err)
			}
			println("New message from: %s in %s", put.Id, put.Index, put.Type, put.Version)
		})(message, i)
	}

	wg.Wait()

	// My own way..
	js, err := json.Marshal(map[string]interface{}{"query": map[string]interface{}{"match": map[string]interface{}{"user": "gabivlj"}}})
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println(string(js))
	req, err := http.NewRequest("GET", "http://127.0.0.1:9200/chat/_search", bytes.NewBuffer(js))
	req.Header.Add("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		fmt.Print(string(content))
	}

	// Elastic library...

	termQuery := elastic.NewTermQuery("user", "gabivlj")
	fmt.Print(termQuery)
	fmt.Println()
	termQuery.QueryName("match")
	searchResult, err := client.Search().
		Index("chat").
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	var chit Message
	fmt.Print(searchResult)
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	for _, item := range searchResult.Each(reflect.TypeOf(chit)) {
		if t, ok := item.(Message); ok {
			fmt.Printf("Message by %s: %s\n", t.User, t.Text)
		}
	}
}

// 	// Initialize a client with the default settings.
// 	//
// 	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
// 	//
// 	es, err := elasticsearch.NewDefaultClient()
// 	es, err :=  v

// 	es.DeleteIndex(users).Do(context.Background())

// 	if err != nil {
// 		log.Fatalf("Error creating the client: %s", err)
// 	}

// 	// 1. Get cluster info
// 	//
// 	res, err := es.Info()
// 	if err != nil {
// 		log.Fatalf("Error getting response: %s", err)
// 	}
// 	// Check response status
// 	if res.IsError() {
// 		log.Fatalf("Error: %s", res.String())
// 	}
// 	// Deserialize the response into a map.
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Fatalf("Error parsing the response body: %s", err)
// 	}
// 	// Print client and server version numbers.
// 	log.Printf("Client: %s", elasticsearch.Version)
// 	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
// 	log.Println(strings.Repeat("~", 37))

// 	// 2. Index documents concurrently
// 	//
// 	for i, title := range RandomStrings(5) {
// 		wg.Add(1)

// 		go func(i int, title string) {
// 			defer wg.Done()

// 			// Build the request body.
// 			var b strings.Builder
// 			b.WriteString(`{"message" : "`)
// 			b.WriteString(title)
// 			b.WriteString(`", `)
// 			b.WriteString(`"user" : "`)
// 			b.WriteString(ChooseString(usersArray))
// 			b.WriteString(`"}`)

// 			fmt.Println(b.String())

// 			// Set up the request object.
// 			req := esapi.IndexRequest{
// 				Index:      users,
// 				DocumentID: strconv.Itoa(i + 1),
// 				Body:       strings.NewReader(b.String()),
// 				Refresh:    "true",
// 			}

// 			// Perform the request with the client.
// 			res, err := req.Do(context.Background(), es)
// 			if err != nil {
// 				log.Fatalf("Error getting response: %s", err)
// 			}
// 			defer res.Body.Close()

// 			if res.IsError() {
// 				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
// 			} else {
// 				// Deserialize the response into a map.
// 				var r map[string]interface{}
// 				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 					log.Printf("Error parsing the response body: %s", err)
// 				} else {
// 					// Print the response status and indexed document version.
// 					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
// 				}
// 			}
// 		}(i, title)
// 	}
// 	wg.Wait()

// 	log.Println(strings.Repeat("-", 37))

// 	// 3. Search for the indexed documents
// 	//
// 	// Build the request body.
// 	var buf bytes.Buffer
// 	query := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"match": map[string]interface{}{
// 				"user": "gabivlj2",
// 			},
// 		},
// 	}
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		log.Fatalf("Error encoding query: %s", err)
// 	}

// 	// Perform the search request.

// 	res, err = es.Search(
// 		es.Search.WithContext(context.Background()),
// 		es.Search.WithIndex(users),
// 		es.Search.WithBody(&buf),
// 		es.Search.WithTrackTotalHits(true),
// 		es.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		log.Fatalf("Error getting response: %s", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var e map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
// 			log.Fatalf("Error parsing the response body: %s", err)
// 		} else {
// 			// Print the response status and error information.
// 			log.Fatalf("[%s] %s: %s",
// 				res.Status(),
// 				e["error"].(map[string]interface{})["type"],
// 				e["error"].(map[string]interface{})["reason"],
// 			)
// 		}
// 	}

// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Fatalf("Error parsing the response body: %s", err)
// 	}
// 	// Print the response status, number of results, and request duration.
// 	log.Printf(
// 		"[%s] %d hits; took: %dms",
// 		res.Status(),
// 		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
// 		int(r["took"].(float64)),
// 	)
// 	// Print the ID and document source for each hit.
// 	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
// 		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
// 	}

// 	log.Println(strings.Repeat("=", 37))
// }
