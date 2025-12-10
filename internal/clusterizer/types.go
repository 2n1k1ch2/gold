package clusterizer

var ServiceClusters = map[string]Cluster{}

type Cluster struct {
	Goroutines map[string]GoroutineObject
	Block      map[string]BlockObject
	Mutex      map[string]MutexObject
}

type GoroutineObject struct {
	Hash   string   `json:"hash"`
	Status string   `json:"status"`
	Name   string   `json:"name"`
	Frames []string `json:"-"`
	Count  uint64   `json:"count"`
	Ids    []uint64 `json:"ids"`
}

type BlockObject struct {
	Hash   string
	Frames []string
	Count  uint64
	Cycles uint64
}

type MutexObject struct {
	Hash   string
	Frames []string
	Count  uint64
	Cycles uint64
}

func NewCluster() *Cluster {
	return &Cluster{
		Goroutines: make(map[string]GoroutineObject),
		Block:      make(map[string]BlockObject),
		Mutex:      make(map[string]MutexObject),
	}
}
