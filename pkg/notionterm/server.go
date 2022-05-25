package notionterm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

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
func ButtonPage() http.Handler {
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

const waitPageTpl = `
<html>
    <body>
	<p>⌛ Waiting for notionterm set up</p> 
    </body>
</html>
`

//WaitPageIDHandler: handler waiting for pageid of notion page in parameter
func WaitPageIDHandler(config *Config, urlCh chan string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlCh <- r.URL.Query().Get("url")
		t, err := template.New("interactive").Parse(waitPageTpl)
		if err != nil {
			fmt.Println(err, "failed loading waiting page template")
		}
		data := struct {
		}{}
		check.Check(t.Execute(w, data), "failed writing in waiting template")

	})
}

//ActivateNotionTerm: endpoint to activate notionterm
func ActivateNotionTerm(client *notionapi.Client, pageid string, play chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📶 notionterm activated")
		play <- struct{}{}
		w.Write([]byte("OK"))
	})
}

//DeactivateNotionTerm: endpoint to deactivate notionterm
func DeactivateNotionTerm(pause chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📴 notionterm deactivated")
		w.Write([]byte("OK"))
		pause <- struct{}{}
	})
}

//LaunchNotiontermServer: Launch server with appropriate endpoint: /button which return a button widget, /activate to activate button
//, /deactivate for the reverse, and if
func LaunchNotiontermServer(config *Config, isServer bool, play chan struct{}, pause chan struct{}) (stopServer chan bool) {
	stopServer = make(chan bool)
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + config.Port, Handler: m}

	// define handlers
	m.Handle("/button", ButtonPage())
	m.Handle("/activate", ActivateNotionTerm(config.Client, config.PageID, play))
	m.Handle("/deactivate", DeactivateNotionTerm(pause))

	// launch
	go func() {
		fmt.Printf("🌀 Start server on %s:%s", config.ExternalIP, config.Port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return stopServer
}

//LaunchUrlWaitingServer: launch the server waiting for url parameter, kill it when task is finished
func LaunchUrlWaitingServer(port string) (pageid string) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + port, Handler: m}
	urlCh := make(chan string)
	resp := `
<html>
<body>
	<p>⌛ Waiting for notionterm set up</p> 
</body>
</html>
`
	m.HandleFunc("/notionterm", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
		// Send query parameter to the channel
		urlCh <- r.URL.Query().Get("url")
	})
	go func() {
		fmt.Println("🌀 Start server on", port, ".. waiting for notion page url")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	pageid = listenAndWaitPageId(&s, urlCh)

	return pageid
}

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
