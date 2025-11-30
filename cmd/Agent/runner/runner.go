package runner

import (
	"gold/cmd/Agent/config"
	"gold/cmd/Agent/fetcher"
	"gold/cmd/Agent/sender"
	"log"
	"time"
)

type Runner struct {
	Fetcher fetcher.Fetcher
	Sender  sender.Sender
	Period  time.Duration
	config  config.AgentConfig
}

func NewRunner(Fetcher fetcher.Fetcher, Sender sender.Sender, Period time.Duration, config config.AgentConfig) *Runner {
	return &Runner{
		Fetcher: Fetcher,
		Sender:  Sender,
		Period:  Period,
		config:  config,
	}
}

func (r *Runner) Start() {
	ticker := time.NewTicker(r.Period)
	for range ticker.C {
		snap, err := r.Fetcher.Collect()
		if err != nil {
			log.Println("collect error:", err)
			continue
		}
		if err := r.Sender.Send(snap); err != nil {
			log.Println("send error:", err)
		}
	}
}
