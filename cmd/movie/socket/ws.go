package socket

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/Watch2Gather/server/internal/movie/infras/ws"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
)

var streamChans = make(map[uuid.UUID]*sharedkernel.Broker[*sharedkernel.Message])

type webSocketHandler struct {
	handler func(
		id uuid.UUID,
		c *websocket.Conn,
		ctx context.Context,
		streams map[uuid.UUID]*sharedkernel.Broker[*sharedkernel.Message],
	)

	upgrader websocket.Upgrader
}

var moviePath = os.Getenv("MOVIE_PATH_PREFIX")

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	id, err := uuid.Parse(r.PathValue("roomID"))
	if err != nil {
		slog.Info("missing uuid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("401 - wrong id"))
		cancel()
		return
	}

	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error %s when upgrading connection to websocket", err)
		cancel()
		return
	}
	defer c.Close()

	slog.Info("start responding to client...")

	_, ok := streamChans[id]
	if !ok {
		broker := sharedkernel.NewBroker[*sharedkernel.Message]()
		streamChans[id] = broker
		go broker.Start()
	}

	wsh.handler(id, c, ctx, streamChans)
	slog.Debug("11")

	cancel()
	msg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "haha")
	w.Write(msg)
}

func ServeMovie(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("movieID"))
	if err != nil {
		slog.Info("missing uuid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("401 - wrong id"))
		return
	}

	http.ServeFile(w, r, moviePath+id.String()+".mp4")
}

func NewWsHandler(mux *http.ServeMux) {
	ctx, cancel := context.WithCancel(context.Background())
	wsHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
		handler:  ws.PlayerControl,
	}

	mux.HandleFunc(
		"/api/v1/stream/{roomID}",
		wsHandler.ServeHTTP,
	)
	mux.HandleFunc("/api/v1/stream/movie/{movieID}", ServeMovie)

	<-ctx.Done()
	cancel()
}
