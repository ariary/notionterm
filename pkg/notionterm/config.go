package notionterm

import (
	"time"

	"github.com/jomei/notionapi"
)

type CaptionBlock struct {
	Type notionapi.BlockType
	Id   notionapi.BlockID
}

type Config struct {
	Delay           time.Duration
	Client          *notionapi.Client
	ExternalUrl     string
	Port            string
	PageID          string
	Path            string
	PS1             string
	Shell           string
	ExternalIP      string
	CaptionBlock    CaptionBlock
	TerminalBlockId notionapi.BlockID
}
