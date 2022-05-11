package notionterm

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/ariary/go-utils/pkg/check"
	"github.com/jomei/notionapi"
)

const buttonPageTpl = `
<html>
    <body>
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

func ActivateNotionTerm(client *notionapi.Client, pageid string, play chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📶 notionterm activated")
		play <- struct{}{}
	})
}

func DeactivateNotionTerm(pause chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📴 notionterm deactivated")
		pause <- struct{}{}
	})
}

func SetupRoutes(client *notionapi.Client, pageid string, play chan struct{}, pause chan struct{}) {
	http.Handle("/button", buttonPage())
	http.Handle("/activate", ActivateNotionTerm(client, pageid, play))
	http.Handle("/deactivate", DeactivateNotionTerm(pause))
}
