package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	username string
	password string
)

func init() {
	username = os.Getenv("ES_USERNAME")
	password = os.Getenv("ES_PASSWORD")
	fmt.Printf("Using username: %s and password: %s\n", username, password)

	fmt.Println("Creating index: nginx_json_elastic_stack_example")
	createIndex()

	unzip("/nginx_data/data.zip", "/nginx_data")

	fmt.Println("Adding template: nginx_json_elastic_stack_example")
	putTemplate()

	fmt.Println("Adding ingest pipeline: nginx-pipeline")
	putPipeline()
}

func createIndex() {

	var err error
	var client *elastic.Client

	if username != "" && password != "" {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetBasicAuth(username, password),
		)
	} else {
		client, err = elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	}

	if err != nil {
		log.Fatal(err)
	}

	exists, err := client.IndexExists("nginx_json_elastic_stack_example").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.
			CreateIndex("nginx_json_elastic_stack_example").
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Fatalf("expected index creation to be ack'd; got: %v", createIndex.Acknowledged)
		}
	}
}

func putTemplate() {

	var err error
	var client *elastic.Client

	// Read in nginx_json_template
	buf, err := ioutil.ReadFile("/nginx_data/nginx_json_template.json")
	if err != nil {
		log.Fatal(err)
	}

	if username != "" && password != "" {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetBasicAuth(username, password),
		)
	} else {
		client, err = elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	}

	putres, err := client.IndexPutTemplate("nginx_json_elastic_stack_example").
		BodyString(string(buf)).
		Do(context.Background())
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

func putPipeline() {
	var err error
	var client *elastic.Client

	// Read in nginx_json_template
	buf, err := ioutil.ReadFile("/nginx_data/nginx-ingest-pipeline.json")
	if err != nil {
		log.Fatal(err)
	}

	if username != "" && password != "" {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetBasicAuth(username, password),
		)
	} else {
		client, err = elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	}

	putres, err := client.IngestPutPipeline("nginx-pipeline").
		BodyString(string(buf)).
		Do(context.Background())

	if err != nil {
		log.Fatalf("expected no error; got: %v", err)
	}
	if putres == nil {
		log.Fatalf("expected response; got: %v", putres)
	}
	if !putres.Acknowledged {
		log.Fatalf("expected ingest pipeline to be ack'd; got: %v", putres.Acknowledged)
	}
}

func unzip(archive, target string) error {

	fmt.Println("Unzip archive ", target)

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		filePath := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	var err error
	var client *elastic.Client

	file, err := os.Open("/nginx_data/nginx_json_logs")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if username != "" && password != "" {
		client, err = elastic.NewClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetBasicAuth(username, password),
		)
	} else {
		client, err = elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	}

	bulkRequest := client.Bulk()
	bulkRequest.Pipeline("nginx-pipeline")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bulkRequest = bulkRequest.Add(elastic.NewBulkIndexRequest().
			Index("nginx_json_elastic_stack_example").
			Type("logs").
			Doc(scanner.Text()))
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("This Bulk Requests is %s.\n", humanize.Bytes(uint64(bulkRequest.EstimatedSizeInBytes())))

	bulkResponse, err := bulkRequest.Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	if bulkResponse == nil {
		log.Errorf("expected bulkResponse to be != nil; got nil")
	}

	if bulkRequest.NumberOfActions() != 0 {
		log.Errorf("expected bulkRequest.NumberOfActions %d; got %d", 0, bulkRequest.NumberOfActions())
	}
}
