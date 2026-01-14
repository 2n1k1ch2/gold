package sender

import (
	"bytes"
	"context"
	"fmt"
	snapshot "gold/api"
	"gold/cmd/Agent/fetcher"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net/http"
)

type DefaultSender struct {
	URL string
}

func NewDefaultSender(url string) *DefaultSender {
	return &DefaultSender{URL: url}
}

type ProtoEncoder struct{}

func (p *ProtoEncoder) Encode(snapshot *fetcher.RuntimeSnapshot) ([]byte, error) {
	protoSnapShot := convertToProto(snapshot)
	data, err := proto.Marshal(protoSnapShot)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (p *ProtoEncoder) ContentType() string {
	return "application/x-protobuf"
}

type HTTPTransport struct {
}

func (t *HTTPTransport) Post(ctx context.Context, url string, contentType string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http error %d: %s", resp.StatusCode, b)
	}
	return nil
}
func (s *DefaultSender) Send(ctx context.Context, snapshot *fetcher.RuntimeSnapshot) error {
	encoder := ProtoEncoder{}
	data, err := encoder.Encode(snapshot)
	if err != nil {
		return err
	}
	transport := HTTPTransport{}
	return transport.Post(ctx, s.URL, "application/x-protobuf", data)
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

// --------FOR TESTS--------------------------------

type TestSender struct {
	Transport *TestTransport
}

func (t *TestSender) Send(ctx context.Context, snapshot *fetcher.RuntimeSnapshot) error {
	encoder := ProtoEncoder{}
	data, err := encoder.Encode(snapshot)
	if err != nil {
		return err
	}

	return t.Transport.Post(ctx, "", "application/x-protobuf", data)
}

type TestTransport struct {
	LastBody []byte
}

func (t *TestTransport) Post(ctx context.Context, url string, contentType string, body []byte) error {
	t.LastBody = body
	return nil
}
