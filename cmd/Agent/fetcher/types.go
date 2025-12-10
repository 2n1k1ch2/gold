package fetcher

type Fetcher interface {
	Collect() (*RuntimeSnapshot, error)
}

type RuntimeSnapshot struct {
	ServiceName string
	Timestamp   int64

	GoroutineDump  []byte // debug=2
	BlockProfile   []byte // pprof binary
	MutexProfile   []byte // pprof binary
	RuntimeMetrics map[string]float64
	Version        string
}
