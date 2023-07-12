package main

import (
	cyoa "cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with CYOA story")
	port := flag.String("port", "8050", "port, we will be running the server")
	flag.Parse()
	fmt.Printf("The story is %s\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Could not open the file, %v", err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		log.Fatalf("Could not parse the JSON, %v", err)
	}
	
	http.Handle("/", cyoa.NewHandler(story))

	fmt.Printf("Running in port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
