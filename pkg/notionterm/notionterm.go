package notionterm

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ariary/go-utils/pkg/host"
	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

//Init: init notionterm: param, envar etc
func Init() (config Config) {
	var buttonUrlOverride, port string
	flag.StringVar(&buttonUrlOverride, "button-url", "", "override button url (useful if notionterm service is behind a proxy)")
	flag.StringVar(&port, "p", "", "specify target listening port (HTTP traffic)")
	flag.Parse()

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

	config.Pageid = pageurl[strings.LastIndex(pageurl, "-")+1:]
	if config.Pageid == pageurl {
		fmt.Println("❌ PAGEID was not found in NOTION_TERM_PAGE_URL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
	}

	// CHECK PAGE CONTENT
	config.Client = notionapi.NewClient(notionapi.Token(token))

	children, err := notionion.RequestProxyPageChildren(config.Client, config.Pageid)
	if err != nil {
		fmt.Println("Failed retrieving page children blocks:", err)
		os.Exit(92)
	}
	// target  config
	//targetUrl: find target reachable url (neither in args or in page otherwise try to find it)
	var targetUrl string
	if len(flag.Args()) > 0 { //in args
		targetUrl = flag.Arg(0)
	} else {
		//in page
		targetUrlTmp, _ := RequestTargetUrlFromConfig(config.Client, config.Pageid)
		// if err != nil {
		// 	fmt.Println("Failed to retrieve target URL from notion page:", err)
		// }
		if targetUrlTmp == "" {
			//try to find it
			targetUrlTmp, err = host.GetExternalIP()
			if err != nil {
				fmt.Println("Failed to detect external ip (dig):", err)
			} else if targetUrlTmp == "" {
				targetUrlTmp, err = host.GetHostIP()
				if err != nil {
					fmt.Println("Failed to detect external ip (hostname):", err)
				}
			}
		}
		targetUrl = targetUrlTmp
	}

	// port config
	if port == "" {
		if port, _ = RequestPortFromConfig(config.Client, config.Pageid); port == "" {
			port = "9292"
		}
	}
	config.Port = port

	// embed button section checks
	var buttonUrl string
	if targetUrl == "" {
		fmt.Println("❌ Failed to get target URL/IP")
		os.Exit(92)
	} else if buttonUrlOverride == "" {
		fmt.Println("📡 Target:", targetUrl)
		buttonUrl = "https://" + targetUrl + ":" + port + "/button"
	} else {
		fmt.Println("📡 Target button url:", buttonUrlOverride)
		buttonUrl = buttonUrlOverride
	}
	if button, err := GetButtonBlock(children); err != nil {
		fmt.Println("❌ button not found in the notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("🕹️ button widget found")
		if buttonUrl != "" {
			if _, err := UpdateButtonUrl(config.Client, button.ID, buttonUrl); err != nil {
				fmt.Println("Failed updating button url:", err)
				os.Exit(92)
			}
		}
		//get current path & update Caption accordingly
		config.Path, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(92)
		}
		UpdateButtonCaption(config.Client, button, config.Path)
	}

	// code/terminal section check
	if code, err := GetTerminalBlock(children); err != nil {
		fmt.Println("❌ terminal section not found in notionterm page")
		os.Exit(92)
	} else {
		fmt.Println("👨‍💻 terminal block found")
		UpdateCodeContent(config.Client, code.ID, "$ ")
	}

	// for i := 0; i < len(children); i++ {
	// 	fmt.Printf("%+v", children[i])
	// }

	config.PS1 = "$ "
	config.Delay = 500 * time.Millisecond

	return config

}

//NotionTerm: "Infinite loop" to read the content of terminal code block and execute it if it is a command, then returning stdout
func NotionTerm(config Config, play chan struct{}, pause chan struct{}) {
	for {
		time.Sleep(config.Delay)
		select {
		case <-pause:
			//fmt.Println("pause")
			select {
			case <-play:
				//fmt.Println("play")
			}
		default:
			termBlock, err := RequestTerminalBlock(config.Client, config.Pageid)
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
						cmd = cmdSplit[1]
					}
					cmd = strings.Replace(cmd, "\n", "", -1)
					if !handleSpecialCommand(&config, termBlock, cmd) {
						//Execute it and print
						ExecAndPrint(config.Client, termBlock, config.Path, cmd)
					}

					//refresh+add new terminal line ($)
					termBlock, err = RequestTerminalBlock(config.Client, config.Pageid)
					if err != nil {
						fmt.Println(err)
						continue
					}
					AddTermLine(config.Client, termBlock)
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

func handleSpecialCommand(config *Config, termBlock notionapi.CodeBlock, cmd string) (isSpecial bool) {

	if strings.HasPrefix(cmd, "cd ") { //TODO handle if beginnign with .. to mak ethe path absolute for caption
		//change path
		cmdSplit := strings.Split(cmd, " ")
		if len(cmdSplit) > 1 {
			path := cmdSplit[1]
			if button, err := RequestButtonBlock(config.Client, config.Pageid); err != nil {
				fmt.Println(err)
			} else {
				UpdateButtonCaption(config.Client, button, path)
				config.Path = path
				fmt.Println("📁 Change directory:", path)
			}
			return true
		} else {
			fmt.Println("Failed retrieving directory in 'cd' command:", cmd)
		}
	} else if strings.HasPrefix(cmd, "clear") {
		fmt.Println("🦆 Clear terminal")
		UpdateCodeContent(config.Client, termBlock.ID, "")
		return true
	}
	return false
}

//ExecAndPrint: execute command and print the result in code block
func ExecAndPrint(client *notionapi.Client, termBlock notionapi.CodeBlock, path string, cmd string) {
	fmt.Println("📟", cmd)
	commandExec := exec.Command("sh", "-c", cmd)
	commandExec.Dir = path
	stdout, err := commandExec.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Print the output
	//fmt.Println(string(stdout))
	if _, err := AddRichText(client, termBlock, string(stdout)); err != nil {
		fmt.Println("failed to add rich text in terminal code block:", err)
	}
}
