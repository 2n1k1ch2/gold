package sender

import (
	"bytes"
	"fmt"
	snapshot "gold/api"
	"gold/cmd/Agent/fetcher"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

type DefaultSender struct {
	url string
}

func NewDefaultSender(url string) *DefaultSender {
	return &DefaultSender{url: url}
}
func (s *DefaultSender) Send(snapshot *fetcher.RuntimeSnapshot) error {
	protoSnapShot := convertToProto(snapshot)
	data, err := proto.Marshal(protoSnapShot)
	if err != nil {
		return fmt.Errorf("marshal proto: %w", err)
	}

	resp, err := http.Post(s.url, "application/x-protobuf", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func convertToProto(s *fetcher.RuntimeSnapshot) *snapshot.SnapShot {
	return &snapshot.SnapShot{
		ServiceName:    s.ServiceName,
		Timestamp:      s.Timestamp,
		GoroutineDump:  s.GoroutineDump,
		BlockProfile:   s.BlockProfile,
		MutexProfile:   s.MutexProfile,
		RuntimeMetrics: s.RuntimeMetrics,
		Version:        s.Version,
	}
}
