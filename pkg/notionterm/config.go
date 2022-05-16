package notionterm

import (
	"time"

	"github.com/jomei/notionapi"
)

type Config struct {
	Delay  time.Duration
	Client *notionapi.Client
	Port   string
	PageID string
	Path   string
	PS1    string
	Shell  string
}
