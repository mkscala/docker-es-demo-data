package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

func init() {
	fmt.Println("Creating index: nginx_json_elastic_stack_example")
	createIndex()
	unzip("/nginx_data/data.zip", "/nginx_data")
	fmt.Println("Adding template: nginx_json_elastic_stack_example")
	putTemplate()
	fmt.Println("Adding ingest pipeline: nginx-pipeline")
	putPipeline()
}

func createIndex() {
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
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
	// Read in nginx_json_template
	buf, err := ioutil.ReadFile("/nginx_data/nginx_json_template.json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Fatal(err)
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
	// Read in nginx_json_template
	pipelineJSON, err := ioutil.ReadFile("/nginx_data/nginx-ingest-pipeline.json")
	if err != nil {
		log.Fatal(err)
	}

	url := "http://elasticsearch:9200/_ingest/pipeline/nginx-pipeline"

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(pipelineJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
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

	file, err := os.Open("/nginx_data/nginx_json_logs")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Fatal(err)
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
