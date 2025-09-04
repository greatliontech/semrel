package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/greatliontech/semrel/pkg/semrel"
	"github.com/swaggest/jsonschema-go"
)

func main() {
	reflector := jsonschema.Reflector{}

	schema, err := reflector.Reflect(semrel.ConfigFile{})
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.MarshalIndent(schema, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))
}
