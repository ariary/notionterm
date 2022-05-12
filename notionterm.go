package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ariary/notionterm/pkg/notionterm"
)

func main() {
	port, pageid, client := notionterm.Init()

	var play = make(chan struct{})
	var pause = make(chan struct{})
	go notionterm.NotionTerm(client, pageid, play, pause)
	pause <- struct{}{}
	notionterm.SetupRoutes(client, pageid, play, pause)
	fmt.Printf("🖥️ Launch notionterm on port %s !\n\n", port)
	log.Println(http.ListenAndServe(":"+port, nil))

}
