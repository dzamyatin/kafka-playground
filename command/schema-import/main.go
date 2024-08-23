package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type importType string

const (
	INSERT importType = "insert"
	IMPORT importType = "import"
)

var IMPORT_TYPES = []importType{INSERT, IMPORT}

func main() {
	log.Println("starting schema-import")

	var distSchemaRegistryUrl, importerStr string
	var importerName importType

	flag.StringVar(
		&distSchemaRegistryUrl,
		"url",
		"",
		"destination schema registry url, example: http://localhost:8081",
	)

	var importNames []string
	for _, importName := range IMPORT_TYPES {
		importNames = append(importNames, string(importName))
	}

	flag.StringVar(
		&importerStr,
		"importer",
		string(INSERT),
		"type of import "+strings.Join(importNames, ", "),
	)
	flag.Parse()

	importerName = importType(importerStr)

	if distSchemaRegistryUrl == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Println("Reading schema from stdin")
	b, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	var schemas []importSchema

	err = json.Unmarshal(b, &schemas)

	if err != nil {
		panic(err)
	}

	log.Println("Starting import")

	importer := getImporter(distSchemaRegistryUrl, importerName)
	for _, schema := range schemas {
		err = importer.importSchema(schema)
		if err != nil {
			log.Println(err)
		}
	}
}

func getImporter(distSchemaRegistryUrl string, importer importType) schemaImporter {
	switch importer {
	case IMPORT:
		return newImportImporter(distSchemaRegistryUrl)
	case INSERT:
		return newStandardImporter(distSchemaRegistryUrl)
	default:
		panic("unknown importer")
	}
}

type schemaImporter interface {
	importSchema(schema importSchema) error
}

type importImporter struct {
	insertImporter
}

func (r *importImporter) importSchema(schema importSchema) error {
	err := r.doSwithModeRequest(schema.Subject, IMPORT_MODE_IMPORT)
	if err != nil {
		return err
	}
	schemaRequest := importerImportSchemaRequest{
		SchemaType: schema.SchemaType,
		Schema:     schema.Schema,
		Version:    schema.Version,
		Id:         schema.Id,
	}

	err = r.doImportRequest(schemaRequest, schema.Subject)
	if err != nil {
		return err
	}

	err = r.doSwithModeRequest(schema.Subject, IMPORT_MODE_READWRITE)

	return err
}

type importMode string

const (
	IMPORT_MODE_IMPORT    importMode = "IMPORT"
	IMPORT_MODE_READWRITE importMode = "READWRITE"
)

func (r *importImporter) doSwithModeRequest(subject string, mode importMode) error {
	var res *http.Response
	body := bytes.NewReader([]byte(`{"mode": "` + mode + `"}`))
	req, err := http.NewRequest(
		"PUT",
		r.distSchemaRegistryUrl+"/mode/"+subject+"?force=true",
		body,
	)
	res, err = r.client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(res.Body)
		errText := fmt.Sprintf("Import mode change fail %s by %s", res.Status, respBody)
		log.Println(errText)
		return errors.New(errText)
	}

	return nil
}

func newImportImporter(distSchemaRegistryUrl string) *importImporter {
	return &importImporter{
		insertImporter{
			client:                http.Client{},
			distSchemaRegistryUrl: distSchemaRegistryUrl,
		},
	}
}

type insertImporter struct {
	client                http.Client
	distSchemaRegistryUrl string
}

func newStandardImporter(distSchemaRegistryUrl string) *insertImporter {
	return &insertImporter{
		client:                http.Client{},
		distSchemaRegistryUrl: distSchemaRegistryUrl,
	}
}

func (r *insertImporter) doImportRequest(body any, subject string) error {
	var err error
	var res *http.Response

	res, err = r.client.Do(r.createImportRequest(body, subject))

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		respBody, _ := io.ReadAll(res.Body)
		log.Printf("unexpected status code %d from %s with message %s", res.StatusCode, r.distSchemaRegistryUrl, respBody)

		return errors.New("unexpected status code " + strconv.Itoa(res.StatusCode))
	}

	log.Println("Successfully imported schema")

	return nil
}

func (r *insertImporter) createImportRequest(body any, subject string) *http.Request {
	var err error
	var b []byte
	var req *http.Request
	header := http.Header{
		"Content-Type": {"application/vnd.schemaregistry.v1+json"},
	}

	b, err = json.Marshal(&body)
	if err != nil {
		panic(err)
	}

	req, err = http.NewRequest(
		"POST",
		r.distSchemaRegistryUrl+"/subjects/"+subject+"/versions",
		bytes.NewReader(b),
	)
	req.Header = header

	if err != nil {
		panic(err)
	}

	return req
}

func (r *insertImporter) importSchema(schema importSchema) error {
	schemaRequest := importSchemaRequest{
		SchemaType: schema.SchemaType,
		Schema:     schema.Schema,
	}

	return r.doImportRequest(schemaRequest, schema.Subject)
}

type importSchema struct {
	Subject    string
	Version    int
	Id         int
	SchemaType string
	Schema     string
}

type importSchemaRequest struct {
	SchemaType string `json:"schemaType,omitempty"`
	Schema     string `json:"schema"`
}

type importerImportSchemaRequest struct {
	SchemaType string `json:"schemaType,omitempty"`
	Schema     string `json:"schema"`
	Version    int    `json:"version"`
	Id         int    `json:"id"`
}
