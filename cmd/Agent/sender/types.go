package sender

import (
	"golang.org/x/net/context"
	"gold/cmd/Agent/fetcher"
)

type Sender interface {
	Send(ctx context.Context, snapshot *fetcher.RuntimeSnapshot) error
}
type Encoder interface {
	Encode(snapshot *fetcher.RuntimeSnapshot) ([]byte, error)
	ContentType() string
}
type Transport interface {
	Post(ctx context.Context, url string, contentType string, body []byte) error
}
