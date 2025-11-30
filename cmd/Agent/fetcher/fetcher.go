package fetcher

import (
	"bytes"
	"encoding/gob"
	"runtime"
	"runtime/metrics"
	"time"
)

type DefaultFetcher struct {
	IncludeBlockProfile bool
	IncludeMutexProfile bool
	IncludeMetrics      bool
}

func NewRuntimeFetcher(BlockProfile, MutexProfile, Metrics bool) *DefaultFetcher {
	return &DefaultFetcher{
		IncludeBlockProfile: BlockProfile,
		IncludeMutexProfile: MutexProfile,
		IncludeMetrics:      Metrics,
	}
}

func (f *DefaultFetcher) Collect() (*RuntimeSnapshot, error) {
	// ---------------------
	// 1) Goroutine dump
	// ---------------------
	buf := make([]byte, 2<<20) // 2MB
	n := runtime.Stack(buf, true)
	gorDump := buf[:n]

	// ---------------------
	// 2) BlockProfile
	// ---------------------
	var blockRecords []runtime.BlockProfileRecord
	if f.IncludeBlockProfile {
		n, _ := runtime.BlockProfile(nil)
		blockRecords = make([]runtime.BlockProfileRecord, n)
		runtime.BlockProfile(blockRecords)
	}

	// ---------------------
	// 3) MutexProfile
	// ---------------------
	var mutexRecords []runtime.BlockProfileRecord
	if f.IncludeMutexProfile {
		n, _ := runtime.MutexProfile(nil)
		mutexRecords = make([]runtime.BlockProfileRecord, n)
		runtime.MutexProfile(mutexRecords)
	}

	// ---------------------
	// 4) runtime/metrics
	// ---------------------
	metricsMap := make(map[string]float64)
	if f.IncludeMetrics {

		descriptions := metrics.All()
		samples := make([]metrics.Sample, len(descriptions))
		for i, desc := range descriptions {
			samples[i] = metrics.Sample{
				Name: desc.Name,
			}
		}
		metrics.Read(samples)

		for _, s := range samples {
			switch s.Value.Kind() {
			case metrics.KindUint64:
				metricsMap[s.Name] = float64(s.Value.Uint64())
			case metrics.KindFloat64:
				metricsMap[s.Name] = s.Value.Float64()
			case metrics.KindBad:
				continue
			default:
				// Для других типов можно добавить обработку
				continue
			}
		}
	}

	// ---------------------
	// 5) Convert block/mutex to []byte
	// ---------------------
	blockBuf, _ := convertWithGob(blockRecords)
	mutexBuf, _ := convertWithGob(mutexRecords)

	// ---------------------
	// 6) Build snapshot
	// ---------------------
	snap := &RuntimeSnapshot{
		Service:        "",
		Timestamp:      time.Now().Unix(),
		GoroutineDump:  gorDump,
		BlockProfile:   blockBuf,
		MutexProfile:   mutexBuf,
		RuntimeMetrics: metricsMap,
		Version:        runtime.Version(),
	}

	return snap, nil
}

func convertWithGob(records []runtime.BlockProfileRecord) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	gob.Register([]runtime.BlockProfileRecord{})

	err := encoder.Encode(records)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
