package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"

	elastic "gopkg.in/olivere/elastic.v3"
)

func init() {

	// Read in nginx_json_template
	buf, err := ioutil.ReadFile("/nginx_data/nginx_json_template.json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Fatal(err)
	}

	exists, err := client.IndexExists("nginx_json_elastic_stack_example").Do()
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("nginx_json_elastic_stack_example").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Fatalf("expected index creation to be ack'd; got: %v", createIndex.Acknowledged)
		}
	}
	fmt.Println(string(buf))

	putres, err := client.IndexPutTemplate("nginx_json_elastic_stack_example").
		BodyString(string(buf)).
		Do()
	if err != nil {
		log.Fatalf("expected no error; got: %v", err)
	}
	if putres == nil {
		log.Fatalf("expected response; got: %v", putres)
	}
	if !putres.Acknowledged {
		log.Fatalf("expected index template to be ack'd; got: %v", putres.Acknowledged)
	}
}

func index(data string) {
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Fatal(err)
	}

	newScan, err := client.Index().
		Index("nginx_json_elastic_stack_example").
		Type("logs").
		OpType("index").
		// Id("1").
		BodyJson(data).
		Do()
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"id":    newScan.Id,
		"index": newScan.Index,
		"type":  newScan.Type,
	}).Debug("Indexed sample.")
}

func main() {

	file, err := os.Open("/nginx_data/nginx_json_logs")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		// Index into Elasticsearch
		index(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// json.Unmarshal(file, &jsontype)
	// fmt.Printf("Results: %v\n", jsontype)
}
