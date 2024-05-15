package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/Watch2Gather/server/cmd/room/config"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/room/domain"
	"github.com/Watch2Gather/server/internal/room/usecases/rooms"
	"github.com/Watch2Gather/server/proto/gen"
)

type roomGRPCServer struct {
	cfg *config.Config
	uc  rooms.UseCase
}

var _ gen.RoomServiceServer = (*roomGRPCServer)(nil)

var (
	RoomGRPCServerSet = wire.NewSet(NewGRPCRoomServer)
	roomMessages      = make(map[uuid.UUID]*sharedkernel.Broker[*gen.Message])
)

func NewGRPCRoomServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc rooms.UseCase,
) gen.RoomServiceServer {
	svc := roomGRPCServer{
		cfg: cfg,
		uc:  uc,
	}

	gen.RegisterRoomServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (g *roomGRPCServer) CreateRoom(ctx context.Context, req *gen.CreateRoomRequest) (*gen.CreateRoomResponse, error) {
	slog.Info("POST: CreateRoom")

	if len(req.RoomName) == 0 {
		req.RoomName = "testMovie" // TODO add movie names
	}

	var ids uuid.UUIDs
	for _, id := range req.GetParticipantIds() {
		uid, err := uuid.Parse(id)
		if err != nil {
			slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
			return nil, domain.ErrInvalidID
		}
		ids = append(ids, uid)
	}

	room, err := g.uc.CreateRoom(ctx, &domain.CreateRoomModel{
		Name:           req.RoomName,
		ParticipantIds: ids,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, err
		}

		// if pqErr, ok := err.(*pq.Error); ok {
		// 	switch pqErr.Code.Name() {
		// 	case "unique_violation":
		// 		return nil, domain.ErrUserAlreadyExists
		// 	}
		// }
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	res := gen.CreateRoomResponse{
		OwnerId: room.OwnerID.String(),
		RoomId:  room.ID.String(),
	}

	return &res, nil
}

func (g *roomGRPCServer) GetRoomsByUser(ctx context.Context, req *gen.GetRoomsByUserRequest) (*gen.GetRoomsByUserResponse, error) {
	slog.Info("GET: GetRoomsByUser")

	rooms, err := g.uc.GetRoomsByUser(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	res := &gen.GetRoomsByUserResponse{
		Rooms: []*gen.Room{},
	}

	for _, room := range rooms {
		res.Rooms = append(res.Rooms, &gen.Room{
			Name:              room.Name,
			Id:                room.ID.String(),
			OwnerId:           room.OwnerID.String(),
			ParticipantsCount: int32(room.ParticipantsCount),
			FilmTitle:         room.FilmTitle,
			PosterPath:        room.PosterPath,
			MovieId:           room.MovieID,
		})
	}

	return res, nil
}

func (g *roomGRPCServer) GetUsersByRoom(ctx context.Context, req *gen.GetUsersByRoomRequest) (*gen.GetUsersByRoomResponse, error) {
	slog.Info("GET: GetRoomsByUser")

	id, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	participants, err := g.uc.GetParticipantsByRoomID(ctx, id)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetParticipantsByRoomID"))
		return nil, sharedkernel.ErrServer
	}

	res := &gen.GetUsersByRoomResponse{}

	for _, participant := range participants {
    res.Participants = append(res.Participants, &gen.Participant{
    	Id:     participant.UserID.String(),
    	Name:   participant.Username,
    	Avatar: participant.Avatar,
    })
  }

	return res, nil
}

func (g *roomGRPCServer) InviteToRoom(ctx context.Context, req *gen.InviteToRoomRequest) (*gen.InviteToRoomResponse, error) {
	slog.Info("POST: InviteToRoom")

	_, err := g.uc.InviteToRoom(ctx, &domain.AddParticipantsModel{
		RoomID: [16]byte{},
	})
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	res := &gen.InviteToRoomResponse{}

	return res, nil
}

func (g *roomGRPCServer) GetMessagesByRoom(ctx context.Context, req *gen.GetMessagesByRoomRequest) (*gen.GetMessagesByRoomResponse, error) {
	slog.Info("GET: GetMessagesByRoom")
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	messages, err := g.uc.GetRoomMessages(ctx, roomID)
	// slog.Debug("Getting Messages", "messages", messages)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	res := &gen.GetMessagesByRoomResponse{}

	for _, message := range messages {
		res.Messages = append(res.Messages, &gen.Message{
			Id:        message.ID.String(),
			Text:      message.Text,
			CreatedAt: message.CreatedAt,
			User: &gen.Sender{
				Id:     message.User.ID.String(),
				Name:   message.User.Name,
				Avatar: message.User.Avatar,
			},
		})
	}

	return res, nil
}

func (g *roomGRPCServer) EnterRoom(req *gen.EnterRoomRequest, srv gen.RoomService_EnterRoomServer) error {
	slog.Info("GET: EnterRoom")

	ctx := srv.Context()

	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return sharedkernel.ErrServer
	}

	broker, ok := roomMessages[roomID]
	if !ok {
		broker = sharedkernel.NewBroker[*gen.Message]()
		roomMessages[roomID] = broker
		go broker.Start()
	}

	go func() {
		msgCh := broker.Subscribe()
		defer broker.Unsubscribe(msgCh)

	out:
		for {
			select {
			case msg := <-msgCh:
				srv.Send(msg)
			case <-ctx.Done():
				break out
			}
		}
	}()

	<-ctx.Done()
	slog.Debug("Chat leaved")
	return nil
}

func (g *roomGRPCServer) SendMessage(ctx context.Context, req *gen.SendMessageRequest) (*gen.SendMessageResponse, error) {
	slog.Info("POST: SendMessage")

	rID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	uID, err := sharedkernel.GetUserIDFromContext(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	user, err := g.uc.GetUserInfo(ctx, uID)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetUserInfo"))
		return nil, sharedkernel.ErrServer
	}

	broker, ok := roomMessages[rID]
	if !ok {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, domain.ErrNoRoomOpen
	}

	msg, err := g.uc.SendMessage(ctx, &domain.CreateMessageModel{
		Content: req.GetText(),
		RoomID:  rID,
	})
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	broker.Publish(&gen.Message{
		Id:        msg.ID.String(),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
		User: &gen.Sender{
			Id:     uID.String(),
			Name:   user.Name,
			Avatar: user.Avatar,
		},
	})

	return &gen.SendMessageResponse{}, nil
}

func (g *roomGRPCServer) UpdateRoom(ctx context.Context, req *gen.UpdateRoomRequest) (*gen.UpdateRoomResponse, error) {
	slog.Info("POST: UpdateRoom")

	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		slog.Debug("uuid1", "roomID", req.GetRoomId())
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	var movieID uuid.UUID
	if req.GetFilmId() != "" {
		movieID, err = uuid.Parse(req.GetFilmId())
		if err != nil {
			slog.Debug("uuid2", "movieID", req.GetFilmId())
			slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
			return nil, sharedkernel.ErrServer
		}
	} else {
		movieID = uuid.Nil
	}

	err = g.uc.UpdateRoom(ctx, &domain.UpdateRoomModel{
		Name:    req.GetRoomName(),
		MovieID: movieID,
		RoomID:  roomID,
	})
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	return &gen.UpdateRoomResponse{}, nil
}

func (roomGRPCServer) DeleteRoom(_ context.Context, _ *gen.DeleteRoomRequest) (_ *gen.DeleteRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}
