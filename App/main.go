package main

import (
	"fmt"
	"log"

	news "github.com/Com1Software/go-clevelandnews"
)

func main() {
	items, err := news.GetLatest()
	if err != nil {
		log.Fatal(err)
	}

	for i, it := range items {
		if i >= 10 { // just show first 10
			break
		}
		fmt.Printf("%s\n%s\n\n", it.Title, it.Link)
	}
}
