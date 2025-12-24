package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sugoruyo/go-botc"
)

func main() {
	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %s", path, err)
	}
	var r botc.Roster
	err = json.Unmarshal(data, &r)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %s", err)
	}
	fmt.Printf("Name: %s\n", r.Name)
	fmt.Printf("Author: %s\n", r.Author)
	fmt.Printf("Almanac: %s\n", r.Almanac)
	for _, c := range r.Characters {
		fmt.Printf("%s: %s\n", c.Name, c.Ability)
	}
}
