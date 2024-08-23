package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
)

func main() {
	var importPath, distSchemaRegistryUrl string

	flag.StringVar(&importPath, "path", "", "path to import file")
	flag.StringVar(
		&distSchemaRegistryUrl,
		"destination",
		"",
		"destination schema registry url, example: http://localhost:8081",
	)
	flag.Parse()

	if importPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.OpenFile(importPath, os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	b, err := io.ReadAll(f)

	if err != nil {
		panic(err)
	}

	var schemas []importSchema

	err = json.Unmarshal(b, &schemas)

	if err != nil {
		panic(err)
	}

	var header http.Header = http.Header{
		"Content-Type": {"application/vnd.schemaregistry.v1+json"},
	}
	var req *http.Request

	for _, schema := range schemas {
		req, err = http.NewRequest(
			"POST",
			distSchemaRegistryUrl+"/subjects/"+schema.subject+"/versions",
			bytes.NewReader([]byte(schema.schema)),
		)
		req.Header = header

		if err != nil {
			panic(err)
		}
	}
}

type importSchema struct {
	subject    string
	version    int
	id         int
	schemaType string
	schema     string
}
