package receiver

import (
	"context"
	"encoding/json"
	snapshot "gold/api"
	"gold/config"
	"log"
	"net"
)

type Receiver struct {
	Listener     net.Listener
	SnapshotChan chan<- *snapshot.SnapShot
	Addr         string
}

func NewReceiver(cfg config.Config, Snaps chan<- *snapshot.SnapShot) *Receiver {
	return &Receiver{
		SnapshotChan: Snaps,
		Addr:         cfg.ReceiverAddr,
	}
}

func (r *Receiver) Start(ctx context.Context) func() {
	return func() {
		listener, err := net.Listen("tcp", r.Addr)
		if err != nil {
			log.Fatal(err)
		}
		r.Listener = listener

		log.Printf("Receiver listening on %s", r.Addr)
		go func() {
			<-ctx.Done()
			_ = listener.Close()
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					log.Printf("Accept error: %v", err)
					continue
				}
			}
			go r.handleConnection(ctx, conn)
		}
	}
}
func (r *Receiver) handleConnection(ctx context.Context, conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("close connection: %v", err)
		}
	}(conn)

	decoder := json.NewDecoder(conn)
	for {
		var snap snapshot.SnapShot
		if err := decoder.Decode(&snap); err != nil {
			log.Printf("Decode error from %s: %v", conn.RemoteAddr(), err)
			return
		}

		select {
		case r.SnapshotChan <- &snap:
		case <-ctx.Done():
			return
		}
	}
}

func (r *Receiver) Stop() {

}
