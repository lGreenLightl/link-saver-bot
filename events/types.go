package events

import "github.com/lGreenLightl/link-saver-bot/lib/consts"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Event struct {
	Type consts.Type
	Text string
	Meta any
}
