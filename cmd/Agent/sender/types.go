package sender

import (
	"gold/cmd/Agent/fetcher"
)

type Sender interface {
	Send(snapshot *fetcher.RuntimeSnapshot) error
}
