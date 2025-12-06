package receiver

import (
	"encoding/json"
	snapshot "gold/api"
	"gold/config"
	"log"
	"net"
)

type Receiver struct {
	Listener     net.Listener
	SnapshotChan chan<- snapshot.SnapShot
	Addr         string
}

func NewReceiver(cfg config.Config, snapshotChan chan<- snapshot.SnapShot) *Receiver {
	return &Receiver{
		SnapshotChan: snapshotChan,
		Addr:         cfg.ReceiverAddr,
	}
}

func (r *Receiver) Start() error {
	listener, err := net.Listen("tcp", r.Addr)
	if err != nil {
		return err
	}
	r.Listener = listener

	log.Printf("Receiver listening on %s", r.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go r.handleConnection(conn)
	}
}

func (r *Receiver) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	decoder := json.NewDecoder(conn)
	for {
		var snap snapshot.SnapShot
		if err := decoder.Decode(&snap); err != nil {
			log.Printf("Decode error: %v", err)
			return
		}

		r.SnapshotChan <- snap
	}
}

func (r *Receiver) Stop() {
	if r.Listener != nil {
		err := r.Listener.Close()
		if err != nil {
			log.Printf("Close error: %v", err)
		}
	}
}
