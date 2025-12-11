package parser

import (
	"errors"
	snapshot "gold/api"
	"log"
	"sync"
)

var (
	ErrBufEmpty = errors.New("buffer is empty")
)

type Parser struct {
	logger             log.Logger
	RawSnapShotChan    <-chan snapshot.SnapShot
	ParsedSnapShotChan chan<- ParsedSnapShot
	wg                 sync.WaitGroup
}

type ParsedSnapShot struct {
	Service   string
	Timestamp int64

	Goroutines     []Goroutine // debug=2
	BlockProfile   []Block     // pprof binary
	MutexProfile   []Mutex     // pprof binary
	RuntimeMetrics map[string]float64
	Version        string
}

type Goroutine struct {
	Data []string
	Id   uint64
}
type Block struct {
	Frames []string // стеки блокировки (фреймы)
	Count  uint64   // количество блокировок
	Cycles int64    // длительность ожиданий
}
type Mutex struct {
	Frames []string // stacktrace места ожидания locka
	Count  uint64   // количество ожиданий
	Cycles int64    // общее время ожидания
}
