package ws

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
)

var movieDir = os.Getenv("MOVIE_PATH_PREFIX")

func PlayerControl(
	id uuid.UUID,
	c *websocket.Conn,
	ctx context.Context,
	streams map[uuid.UUID]*sharedkernel.Broker[*sharedkernel.Message],
) {
	broker := streams[id]
	sender := uuid.New()

	quit := make(chan struct{}, 2)
	go func() {
		for {
			select {
			case <-ctx.Done():
				quit <- struct{}{}
				return
			default:
				slog.Debug("output")
				msg := &sharedkernel.Message{
					Instruction: &sharedkernel.Instruction{},
					Sender:      sender.String(),
				}
				err := c.ReadJSON(msg.Instruction)
				if err != nil {
					slog.Error("Read error", "err", err)
					quit <- struct{}{}
					return
				}
				broker.Publish(msg)
			}
		}
	}()

	go func() {
		msgCh := broker.Subscribe()
		for {
			select {
			case msg := <-msgCh:
				if msg.Sender == sender.String() {
					continue
				}
				c.WriteJSON(msg.Instruction)
			case <-ctx.Done():
				quit <- struct{}{}
				broker.Unsubscribe(msgCh)
				return
			}
		}
	}()

	<-quit
	c.Close()
}
