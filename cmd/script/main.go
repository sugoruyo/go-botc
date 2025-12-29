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
	scriptPath := os.Args[1]
	scriptData, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("failed to read %s: %s", scriptPath, err)
	}
	var s botc.Script
	err = json.Unmarshal(scriptData, &s)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %s", err)
	}
	rosterPath := os.Args[2]
	rosterData, err := os.ReadFile(rosterPath)
	if err != nil {
		log.Fatalf("failed to read %s: %s", rosterPath, err)
	}
	var r botc.Roster
	err = json.Unmarshal(rosterData, &r)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %s", err)
	}

	fmt.Printf("Script: %s by %s\n", s.Meta.Name, s.Author())
	if s.Meta.Almanac != "" {
		fmt.Printf("Learn more at %s\n", s.Meta.Almanac)
	}
	missing := s.PopulateIndex(r)
	if len(missing) > 0 {
		fmt.Printf("Missing originals: %s", strings.Join(missing, ", "))
	}
	fmt.Println("First Night Order:")
	for i, j := range s.FirstNight() {
		fmt.Printf("%02d. %s\n", i+1, j.GetName())
	}
	fmt.Println("Other Night Order:")
	for i, j := range s.OtherNights() {
		fmt.Printf("%02d. %s\n", i+1, j.GetName())
	}
	fmt.Println("Characters")
	if len(s.OriginalCharacterIds) > 0 {
		fmt.Printf("Original: %s\n", strings.Join(s.OriginalCharacterIds, ", "))
	}
	if len(s.CustomCharacters) > 0 {
		customs := make([]string, len(s.CustomCharacters))
		for i, c := range s.CustomCharacters {
			customs[i] = c.Name
		}
		fmt.Printf("Custom: %s\n", strings.Join(customs, ", "))
	}
}
