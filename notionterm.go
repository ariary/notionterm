package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ariary/notionterm/pkg/notionterm"
)

func main() {
	config := notionterm.Init()

	var play = make(chan struct{})
	var pause = make(chan struct{})
	go notionterm.NotionTerm(config, play, pause)
	pause <- struct{}{}
	notionterm.SetupRoutes(config.Client, config.Pageid, play, pause)
	fmt.Printf("🖥️ Launch notionterm on port %s !\n\n", config.Port)
	log.Println(http.ListenAndServe(":"+config.Port, nil))

}
