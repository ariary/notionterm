package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ariary/go-utils/pkg/host"
	"github.com/ariary/notionterm/pkg/notionterm"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

var PageUrl, Token string

func main() {

	var delay int
	var isServerMode, isConfigFromPage bool
	var config notionterm.Config
	config.PS1 = "$ "

	//CMD ROOT
	var rootCmd = &cobra.Command{Use: "notionterm", Run: func(cmd *cobra.Command, args []string) {

		// Initialization
		Init(&config, isServerMode)

		if config.ExternalIP == "" {
			fmt.Println("❌ Failed to get external URL/IP from flags and detection")
			os.Exit(92)
		}

		if config.ExternalUrl == "" {
			config.ExternalUrl = "https://" + config.ExternalIP + ":" + config.Port
		} else {
			config.ExternalUrl = "https://" + config.ExternalUrl
		}

		if isServerMode {
			config.PageID = notionterm.LaunchUrlWaitingServer(&config)
			notionterm.DeleteEmbed(config)
		}

		var play = make(chan struct{})
		var pause = make(chan struct{})
		stopNotion := notionterm.LaunchNotiontermServer(&config, isServerMode, play, pause)

		// creation
		config.CaptionBlock.Id = notionterm.CreateButtonBlock(config)
		config.CaptionBlock.Type = notionapi.BlockTypeEmbed
		if _, err := notionterm.UpdateCaptionById(config.Client, config.PageID, config.CaptionBlock, config.Path); err != nil { //add caption
			fmt.Println("Failed setting button caption:", err)
		}
		config.TerminalBlockId = notionterm.CreateTerminalBlock(config)

		// run
		go notionterm.NotiontermRun(&config, play, pause)
		pause <- struct{}{}
		<-stopNotion
	}}

	//CMD OUTGOING/LIGHT
	var cmdOutgoing = &cobra.Command{ //only outgoing HTTP traffic
		Use:   "light",
		Short: "only grab information from notion page. No HTTP ingoing traffic is used to work.",
		Long:  `only grab information from notion page by performing outgoing HTTP request. No HTTP ingoing traffic is allowed/required so the notionterm is not used as a server.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Initialization
			Init(&config, isServerMode)

			var play = make(chan struct{})
			var pause = make(chan struct{})

			stopNotion := make(chan bool)

			// creation
			config.TerminalBlockId = notionterm.CreateTerminalBlock(config)
			config.CaptionBlock.Id = config.TerminalBlockId
			config.CaptionBlock.Type = notionapi.BlockTypeCode
			go notionterm.NotiontermRun(&config, play, pause)
			<-stopNotion
		},
	}

	// FLAGS
	//token,pageid,p,external,shell,delay,server,config-from-page
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", Token, "notion integration key")
	rootCmd.PersistentFlags().StringVarP(&PageUrl, "page-url", "u", PageUrl, "notion page url")
	rootCmd.PersistentFlags().StringVarP(&config.Port, "port", "p", "9292", "notionterm HTTP listening port")
	rootCmd.PersistentFlags().StringVarP(&config.Shell, "shell", "r", "sh", "shell runtime to execute command with notionterm (sh, bash, and cmd.exe)")
	rootCmd.PersistentFlags().IntVarP(&delay, "delay", "d", 500, "delay between each request to the notion page by notionterm")

	rootCmd.Flags().StringVarP(&config.ExternalIP, "external", "e", "", "external URL/IP of the machine where notionterm runs")
	rootCmd.Flags().StringVarP(&config.ExternalUrl, "override-url", "o", "", "override external url (useful if coupled with ngrok)")
	rootCmd.Flags().BoolVarP(&isServerMode, "server", "s", false, "retrieve notion page id from ingoing HTTP request")
	rootCmd.Flags().BoolVarP(&isConfigFromPage, "config-from-page", "c", false, "retrieve notionterm configuration from page")

	rootCmd.AddCommand(cmdOutgoing)
	rootCmd.Execute()
}

//Init: intialize additional variables
func Init(config *notionterm.Config, isServer bool) {
	//integration token
	if Token == "" {
		Token = os.Getenv("NOTION_TOKEN")
		if Token == "" {
			fmt.Println("❌ Please set NOTION_TOKEN envvar with your integration token before launching notionterm or use --token flag")
			os.Exit(92)
		}
	}

	//pageid
	if !isServer {
		config.PageID = getPageId(PageUrl)
	}

	// external ip/url
	if config.ExternalIP == "" {
		//try to find it
		var ext string
		ext, err := host.GetExternalIP()
		config.ExternalIP = ext
		if err != nil {
			fmt.Println("Failed to detect external ip (dig):", err)
		} else if config.ExternalIP == "" {
			config.ExternalIP, err = host.GetHostIP()
			if err != nil {
				fmt.Println("Failed to detect external ip (hostname):", err)
			}
		}
	}

	//client
	config.Client = notionapi.NewClient(notionapi.Token(Token))

	//get current path
	if path, err := os.Getwd(); err != nil {
		fmt.Println(err)
		os.Exit(92)
	} else {
		config.Path = path
	}
}

//getPAgeId: return notion page id from a notion page url
func getPageId(pageurl string) (pageid string) {
	if pageurl == "" {
		pageurl = os.Getenv("NOTION_TERM_PAGE_URL")
		if pageurl == "" {
			fmt.Println("❌ Please set NOTION_TERM_PAGE_URL envvar with your page id before launching notionterm (CTRL+L on desktop app), or use --page flag")
			os.Exit(92)
		}
	}

	pageid = pageurl[strings.LastIndex(pageurl, "-")+1:]
	if pageid == pageurl {
		fmt.Println("❌ PAGEID was not found in NOTION_TERM_PAGE_URL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
		os.Exit(92)
	}

	return pageid
}
