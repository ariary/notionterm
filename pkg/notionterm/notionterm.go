package notionterm

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

//Init: init notionterm: param, envar etc
func Init() (port string, pageid string, client *notionapi.Client) {
	port = "9292"
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

	pageid = pageurl[strings.LastIndex(pageurl, "-")+1:]
	if pageid == pageurl {
		fmt.Println("❌ PAGEID was not found in NOTION_TERM_PAGE_URL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
	}

	// CHECK PAGE CONTENT
	client = notionapi.NewClient(notionapi.Token(token))

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		fmt.Println("Failed retrieving page children blocks:", err)
		os.Exit(92)
	}

	// embed button section checks
	if button, err := GetButtonBlock(children); err != nil {
		fmt.Println("❌ button not found in the notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("🕹️ button widget found")
		if buttonUrl != "" {
			UpdateButtonUrl(client, button.ID, buttonUrl)
		}
	}

	// code/terminal section check
	if code, err := GetTerminalBlock(children); err != nil {
		fmt.Println("❌ terminal section not found in notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("👨‍💻 terminal block found")
		UpdateCodeContent(client, code.ID, "$ ")
	}

	// for i := 0; i < len(children); i++ {
	// 	fmt.Printf("%+v", children[i])
	// }

	return port, pageid, client

}

//NotionTerm: "Infinite loop" to read the content of terminal code block and execute it if it is a command, then returning stdout
func NotionTerm(client *notionapi.Client, pageid string, play chan struct{}, pause chan struct{}) {
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-pause:
			//fmt.Println("pause")
			select {
			case <-play:
				//fmt.Println("play")
			}
		default:
			termBlock, err := RequestTerminalBlock(client, pageid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cmd, err := GetTerminalLastRichText(termBlock)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println("last:", cmd)
			if strings.Contains(cmd, "\n") && strings.HasPrefix(cmd, "$ ") {
				if isCommand(cmd) {
					cmdSplit := strings.Split(cmd, "$ ")
					if len(cmdSplit) > 1 {
						cmd = cmdSplit[1] //todo check len
					}
					//execute it and print
					fmt.Println(cmd)
					cmmandExec := exec.Command("sh", "-c", cmd)
					stdout, err := cmmandExec.Output()

					if err != nil {
						fmt.Println(err.Error())
						return
					}
					// Print the output
					//fmt.Println(string(stdout))
					if _, err := AddRichText(client, termBlock, string(stdout)); err != nil {
						fmt.Println(err)
					}

					//refresh+add new terminal line ($)
					termBlock, err = RequestTerminalBlock(client, pageid)
					if err != nil {
						fmt.Println(err)
						continue
					}
					AddTermLine(client, termBlock)
				}
			}
		}
	}
}

//check if a command really is
func isCommand(command string) bool {
	if command[len(command)-2] == '\\' {
		return false
	}
	return true
}
