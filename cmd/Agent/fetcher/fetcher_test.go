package fetcher

import (
	"gold/cmd/Agent/config"
	"testing"
	"time"
)

func TestFetcher(t *testing.T) {

	for i := 1; i < 5; i++ {
		go func() {
			time.Sleep(5 * time.Second)
		}()
	}
	time.Sleep(3 * time.Second)

	cfg := config.AgentConfig{ServiceName: "TestService",
		Version: "1.0.0",
		Url:     "localhost",
	}
	fetcher_ := NewRuntimeFetcher(false, false, false, cfg)
	snap, err := fetcher_.Collect()
	if err != nil {
		t.Errorf("Error on Collect,%s", err)
	}
	if snap.ServiceName != cfg.ServiceName {
		t.Errorf("Error on Collect,snap.ServiceName=%s,%s", snap.ServiceName, cfg.ServiceName)
	}
	if snap.Version != cfg.Version {
		t.Errorf("Error on Collect,snap.Version=%s,%s", snap.Version, cfg.Version)
	}
	if len(snap.GoroutineDump) == 0 {
		t.Errorf("Error on Collect,snap.GoroutineDump=%#v", snap.GoroutineDump)
	}
	
}
