package notifier

import (
	"context"
	"log"
)

func (n *Notifier) Run(ctx context.Context) func() {

	return func() {

		for {

			select {
			case <-ctx.Done():
				log.Println("Shutting Notifier ...")
				for _, v := range n.outputs {
					_ = v.Sink.Close()
				}
				return
			case alert := <-n.AlertChan:
				for _, v := range n.outputs {
					data, err := v.Formatter.FormatAlert(alert)
					if err != nil {
						log.Println(err)
						continue
					}
					if err = v.Sink.Send(data); err != nil {
						log.Println(err)
					}

				}
			case info := <-n.InfoChan:
				for _, v := range n.outputs {
					data, err := v.Formatter.FormatInfo(info)
					if err != nil {
						log.Println(err)
						continue
					}
					if err = v.Sink.Send(data); err != nil {
						log.Println(err)
					}

				}
			}
		}
	}
}
