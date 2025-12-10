package clusterizer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gold/internal/parser"
	"log"
	"strings"
)

const (
	RUNNING      string = "running"
	RUNNABLE     string = "runnable"
	SLEEP        string = "sleep"
	CHAN_SEND    string = "chan_send"
	CHAN_RECEIVE string = "chan_receive"
	SELECT       string = "select"
	IO_WAIT      string = "io_wait"
	SYSTEM_CALL  string = "system_call"
	GC_SWEEP     string = "gc_sweep"
	DEAD         string = "dead"
)

func (c *Clusterizer) Clustering(snap *parser.ParsedSnapShot) Cluster {
	cluster := NewCluster()
	cluster.Goroutines = clusteringGoroutines(snap.Goroutines)
	cluster.Block = clusteringBlock(snap.BlockProfile)
	cluster.Mutex = clusteringMutex(snap.MutexProfile)
	return *cluster

}
func clusteringGoroutines(gors []parser.Goroutine) map[string]GoroutineObject {
	gorObj := map[string]GoroutineObject{}
	for _, g := range gors {
		status, err := findStatus(&g)
		if err != nil {
			log.Println(err)
			continue
		}
		name := giveName(&g)
		hsh, err := hashGoroutine(&g)
		if err != nil {
			log.Println(err)
			continue
		}

		if existing, exists := gorObj[hsh]; exists {
			existing.Count += 1
			existing.Ids = append(existing.Ids, g.Id)
			gorObj[hsh] = existing
		} else {
			gorObj[hsh] = GoroutineObject{
				Hash:   hsh,
				Status: status,
				Name:   name,
				Frames: g.Data[1:],
				Count:  1,
				Ids:    []uint64{g.Id},
			}
		}
	}
	return gorObj
}

func clusteringBlock(blocks []parser.Block) map[string]BlockObject {
	out := map[string]BlockObject{}
	for _, b := range blocks {
		hash := hashFrames(b.Frames)
		if existing, ok := out[hash]; ok {
			existing.Count += b.Count
			existing.Cycles += uint64(b.Cycles)
			out[hash] = existing
		} else {
			out[hash] = BlockObject{
				Hash:   hash,
				Frames: b.Frames,
				Count:  b.Count,
				Cycles: uint64(b.Cycles),
			}
		}
	}
	return out
}

func clusteringMutex(mutexs []parser.Mutex) map[string]MutexObject {
	out := map[string]MutexObject{}
	for _, m := range mutexs {
		hash := hashFrames(m.Frames)
		if existing, ok := out[hash]; ok {
			existing.Count += m.Count
			existing.Cycles += uint64(m.Cycles)
			out[hash] = existing
		} else {
			out[hash] = MutexObject{
				Hash:   hash,
				Frames: m.Frames,
				Count:  m.Count,
				Cycles: uint64(m.Cycles),
			}
		}
	}
	return out
}
func findStatus(g *parser.Goroutine) (string, error) {
	if strings.Contains(g.Data[0], "running") {
		return RUNNING, nil
	}
	if strings.Contains(g.Data[0], "runnable") {
		return RUNNABLE, nil
	}
	if strings.Contains(g.Data[0], "sleep") {
		return SLEEP, nil
	}
	if strings.Contains(g.Data[0], "chan send") {
		return CHAN_SEND, nil
	}
	if strings.Contains(g.Data[0], "chan receive") {
		return CHAN_RECEIVE, nil
	}
	if strings.Contains(g.Data[0], "select") {
		return SELECT, nil
	}
	if strings.Contains(g.Data[0], "io wait") {
		return IO_WAIT, nil
	}
	if strings.Contains(g.Data[0], "system call") {
		return SYSTEM_CALL, nil
	}
	if strings.Contains(g.Data[0], "gc sweep") {
		return GC_SWEEP, nil
	}
	if strings.Contains(g.Data[0], "dead") {
		return DEAD, nil
	}

	return "", errors.New("cant get status")

}

func giveName(g *parser.Goroutine) string {
	str := strings.Split(g.Data[0], ":")
	return str[1]
}

func hashGoroutine(g *parser.Goroutine) (string, error) {
	if len(g.Data) < 2 {
		return "", errors.New("cant get hash")
	}

	// make frames
	combined := strings.Join(g.Data[1:], "")

	// then get hash from frames
	hasher := sha256.New()
	hasher.Write([]byte(combined))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func hashFrames(frames []string) string {

	hasher := sha256.New()
	for _, frame := range frames {
		hasher.Write([]byte(frame))
	}
	return hex.EncodeToString(hasher.Sum(nil))

}
