package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sugoruyo/go-botc"
)

func main() {
	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %s", path, err)
	}
	var s botc.Script
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %s", err)
	}
	fmt.Printf("Script: %s by %s\n", s.Meta.Name, s.Author())
	if s.Meta.Almanac != "" {
		fmt.Printf("Learn more at %s\n", s.Meta.Almanac)
	}
	fmt.Println("First Night Order:")
	for i, j := range s.Meta.FirstNight {
		fmt.Printf("%02d. %s\n", i+1, j)
	}
	fmt.Println("Other Night Order:")
	for i, j := range s.Meta.OtherNight {
		fmt.Printf("%02d. %s\n", i+1, j)
	}
	fmt.Println("Characters")
	fmt.Printf("Original: %s\n", strings.Join(s.OriginalCharacterIds, ", "))
	customs := make([]string, len(s.CustomCharacters))
	for i, c := range s.CustomCharacters {
		customs[i] = c.Name
	}
	fmt.Printf("Custom: %s\n", strings.Join(customs, ", "))
}
