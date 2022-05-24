package notionterm

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/ariary/go-utils/pkg/check"
	"github.com/jomei/notionapi"
)

const buttonPageTpl = `
<html>
    <body onload="document.getElementById('switch').innerText= 'ON'">
        <script>
function Activate() {
    var switchButton = document.getElementById("switch");
    url =""
    if (switchButton.innerText== "ON"){
        url = "/activate"
        switchButton.innerText= "OFF"
    }else{
        url = "/deactivate"
        switchButton.innerText= "ON"
    }
    fetch(url);
    }
        </script>
        <button id="switch" onclick="Activate()">ON</button> 
    </body>
</html>
`

//buttonPage: button widget handler
func buttonPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🕹️ Load button page")
		t, err := template.New("interactive").Parse(buttonPageTpl)
		if err != nil {
			fmt.Println(err, "failed loading button widget template")
		}
		data := struct {
		}{}
		check.Check(t.Execute(w, data), "failed writing button template in page")
	})
}

//ActivateNotionTerm: endpoint to activate notionterm
func ActivateNotionTerm(client *notionapi.Client, pageid string, play chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📶 notionterm activated")
		play <- struct{}{}
	})
}

//DeactivateNotionTerm: endpoint to deactivate notionterm
func DeactivateNotionTerm(pause chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📴 notionterm deactivated")
		pause <- struct{}{}
	})
}

//SetupRoutes: set up notionterm routes
func SetupRoutes(client *notionapi.Client, pageid string, play chan struct{}, pause chan struct{}) {
	http.Handle("/button", buttonPage())
	http.Handle("/activate", ActivateNotionTerm(client, pageid, play))
	http.Handle("/deactivate", DeactivateNotionTerm(pause))
}

//Listen on /notionterm URL, and wait for page contained in "url" parameter, exit when found and stop server
func listenAndWaitPageId(s *http.Server, urlCh chan string) string {
	for {
		select {
		case url := <-urlCh:

			pageId := url[strings.LastIndex(url, "-")+1:]
			if pageId == url {
				fmt.Println("Page ID was not found in url provided:", url, ". Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
			} else {
				// Post process after shutdown here (postponed a bit the shutdown cause notion make 2 requests)
				s.Shutdown(context.Background())
				fmt.Println("🎫 Got Page ID from request:", pageId)
				return pageId
			}
		}
	}
}

//createNotionTermBlock: udate embed block url to the button URL + create terminal block (code)
func createNotionTermBlock(config *Config, children notionapi.Blocks, url string) {
	//CREATE BUTTON BLOCK
	//TO FIX: Updating the last embed does not work very well => delete embed + create one
	time.Sleep(1 * time.Second) //wait the embed widget to be loaded
	// embed, err := GetButtonBlock(children)
	// if err != nil {
	// 	fmt.Println("Failed to create notion block:", err)
	// 	os.Exit(92)
	// }

	// if _, err := UpdateButtonUrl(config.Client, embed.ID, url); err != nil {
	// 	fmt.Println("Failed tranform embed block to button (update URL):", err)
	// 	os.Exit(92)
	// }

	//Check if last embed is well loaded
	embed := children[len(children)-1]
	if embed.GetType() == notionapi.BlockTypeBookmark {
		//delete bookmark + create button (embed)
		if _, err := config.Client.Block.Delete(context.Background(), embed.GetID()); err != nil {
			fmt.Println("Failed deleting bookmark:", err)
			os.Exit(92)
		}
		if err := AppendEmbedBlock(config.Client, config.PageID, url); err != nil {
			fmt.Println("Failed creating embed block", err)
			os.Exit(92)
		}

	} else {
		if _, err := UpdateButtonUrl(config.Client, embed.GetID(), url); err != nil {
			fmt.Println("Failed tranform embed block to button (update URL):", err)
			os.Exit(92)
		}
	}
	//CREATE TERMINAL BLOCK
	if err := AppendCodeBlock(config.Client, config.PageID, ""); err != nil {
		fmt.Println("Failed creating terminal block", err)
		os.Exit(92)
	}

	time.Sleep(1500 * time.Millisecond)
}
