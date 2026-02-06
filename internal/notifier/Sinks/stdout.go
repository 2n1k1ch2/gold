package Sinks

import (
	"io"
)

type StdOutSink struct {
	io.Writer
}

func (s *StdOutSink) Send(payload []byte) error {
	_, err := s.Write(payload)
	return err
}
func (s *StdOutSink) Close() error {
	return nil
}
