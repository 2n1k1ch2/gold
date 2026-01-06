package runner

import (
	"context"
	snapshot "gold/api"
	"gold/config"
	"gold/internal/analyzer"
	"gold/internal/clusterizer"
	"gold/internal/parser"
	"gold/internal/receiver"
	"sync"
)

type Runner struct {
	receiver    *receiver.Receiver
	parser      *parser.Parser
	clusterizer *clusterizer.Clusterizer
	analyzer    *analyzer.Analyzer

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}
type Component interface {
	Start(ctx context.Context) func()
}

func NewRunner(cfg config.Config) *Runner {
	var r *Runner

	cancelCtx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	r.ctx = cancelCtx
	r.cancel = cancel
	r.wg = wg

	Snaps := make(chan *snapshot.SnapShot)
	Parsed := make(chan *parser.ParsedSnapShot)
	Normal := make(chan *clusterizer.Cluster)
	Alerts := make(chan *analyzer.Alert)
	Infos := make(chan *analyzer.Info)
	Receiver := receiver.NewReceiver(cfg, Snaps)
	Parser := parser.NewParser(Snaps, Parsed)
	Clusterizer := clusterizer.NewClusterizer(Parsed, Normal)
	Analyzer := analyzer.NewAnalyzer(Normal, Infos, Alerts)
	r.receiver = Receiver
	r.parser = Parser
	r.clusterizer = Clusterizer
	r.analyzer = Analyzer

	return r
}
func (r *Runner) runStage(fn func()) {
	defer r.wg.Done()
	fn()
}
func (r *Runner) Run() {
	r.wg.Add(4)

	go r.runStage(r.receiver.Start(r.ctx))
	go r.runStage(r.parser.Start(r.ctx))
	go r.runStage(r.clusterizer.Start(r.ctx))
	//	go r.runStage(r.analyzer.Start(r.ctx))
}
func (r *Runner) Stop() {
	r.cancel()
	r.wg.Wait()
}
