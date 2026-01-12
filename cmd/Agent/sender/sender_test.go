package sender

import (
	"fmt"
	"gold/cmd/Agent/config"
	"gold/cmd/Agent/fetcher"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		defer request.Body.Close()
		fmt.Printf("Recieve request: %s %s\n", request.Method, string(body))
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("OK"))
	})
	go func() {
		fmt.Println("Start server in 8080 port")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			fmt.Printf("Server couldn't not start %v\n", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)

	cfg := config.AgentConfig{ServiceName: "TestService",
		Version: "1.0.0",
		Url:     "http://localhost:8080",
	}

	sender := NewDefaultSender(cfg.Url)
	Fetcher := fetcher.NewRuntimeFetcher(false, false, false, cfg)
	snap, err := Fetcher.Collect()
	if err != nil {
		t.Error(err)
	}
	if err = sender.Send(snap); err != nil {
		t.Error(err)
	}

}
