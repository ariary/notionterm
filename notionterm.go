package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/ariary/notionterm/pkg/notionterm"
	"github.com/jomei/notionapi"
)

func main() {
	port := "9292"
	var buttonUrl string
	flag.StringVar(&buttonUrl, "button", "", "button url")
	flag.Parse()
	if len(flag.Args()) > 0 {
		port = flag.Arg(0)
	}
	// integration token
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("❌ Please set NOTION_TOKEN envvar with your integration token before launching notionion")
		os.Exit(92)
	}
	// page id
	pageurl := os.Getenv("NOTION_TERM_PAGE_URL")
	if pageurl == "" {
		fmt.Println("❌ Please set NOTION_TERM_PAGE_URL envvar with your page id before launching notionion (CTRL+L on desktop app)")
		os.Exit(92)
	}

	pageid := pageurl[strings.LastIndex(pageurl, "-")+1:]
	if pageid == pageurl {
		fmt.Println("❌ PAGEID was not found in NOTION_TERM_PAGE_URL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
	}

	// CHECK PAGE CONTENT
	client := notionapi.NewClient(notionapi.Token(token))

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		fmt.Println("Failed retrieving page children blocks:", err)
		os.Exit(92)
	}

	// embed button section checks
	if button, err := notionterm.GetButtonBlock(children); err != nil {
		fmt.Println("❌ button not found in the notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("🕹️ button widget found")
		if buttonUrl != "" {
			notionterm.UpdateButtonUrl(client, button.ID, buttonUrl)
		}
	}

	// code/terminal section check
	if code, err := notionterm.GetTerminalBlock(children); err != nil {
		fmt.Println("❌ terminal section not found in notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("👨‍💻 terminal block found")
		notionterm.UpdateCodeContent(client, code.ID, "$ ")
	}

	// for i := 0; i < len(children); i++ {
	// 	fmt.Printf("%+v", children[i])
	// }

	var play = make(chan struct{})
	var pause = make(chan struct{})
	go notionterm.NotionTerm(client, pageid, play, pause)
	pause <- struct{}{}

	notionterm.SetupRoutes(client, pageid, play, pause)
	fmt.Printf("🖥️ Launch notionterm on port %s !\n\n", port)
	log.Println(http.ListenAndServe(":"+port, nil))

}
