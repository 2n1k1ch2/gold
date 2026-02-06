package notifier

import (
	"gold/internal/analyzer"
)

type Notifier struct {
	// NOTE: delivery is synchronous; async workers may be added later
	AlertChan <-chan *analyzer.Alert
	InfoChan  <-chan *analyzer.Info
	outputs   []Output
}

type Sink interface {
	Send(payload []byte) error
	Close() error
}

type Formatter interface {
	FormatAlert(alert *analyzer.Alert) ([]byte, error)
	FormatInfo(info *analyzer.Info) ([]byte, error)
}

type Output struct {
	Sink      Sink
	Formatter Formatter
}
