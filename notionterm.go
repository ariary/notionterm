package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ariary/notionterm/pkg/notionterm"
)

func main() {
	config, buttonID, buttonUrl := notionterm.Init()

	var play = make(chan struct{})
	var pause = make(chan struct{})
	go notionterm.NotionTerm(config, play, pause)
	pause <- struct{}{}
	notionterm.SetupRoutes(config.Client, config.PageID, play, pause)

	//WORKAROUND to make the button loading (find a way to launch server in background) => Launch server THEN update embed link
	go func() {
		time.Sleep(2 * time.Second)
		if buttonUrl != "" {
			if _, err := notionterm.UpdateButtonUrl(config.Client, buttonID, buttonUrl); err != nil {
				fmt.Println("Failed updating button url:", err)
				os.Exit(92)
			}
		}
	}()

	fmt.Printf("🖥️ Launch notionterm on port %s !\n\n", config.Port)
	log.Println(http.ListenAndServe(":"+config.Port, nil))
	//Try to update button after starting server
}
