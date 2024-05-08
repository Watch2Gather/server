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

var RoomGRPCServerSet = wire.NewSet(NewGRPCRoomServer)

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

var roomMessages map[uuid.UUID][]chan *gen.Message = make(map[uuid.UUID][]chan *gen.Message)

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
			Id:                room.ID.String(),
			OwnerId:           room.OwnerID.String(),
			ParticipantsCount: int32(room.ParticipantsCount),
			FilmTitle:         "",
			FilmPoster:        "",
		})
	}

	return res, nil
}

func (g *roomGRPCServer) GetUserIDsByRoom(ctx context.Context, req *gen.GetUserIDsByRoomRequest) (*gen.GetUserIDsByRoomResponse, error) {
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

	res := &gen.GetUserIDsByRoomResponse{
		ParticipantIds: participants.Strings(),
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

func (g *roomGRPCServer) EnterRoom(req *gen.EnterRoomRequest, srv gen.RoomService_EnterRoomServer) error {
	slog.Info("GET: EnterRoom")

	ctx := srv.Context()
	done := ctx.Done()

	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return sharedkernel.ErrServer
	}

	roomChan := make(chan *gen.Message)
	roomMessages[roomID] = append(roomMessages[roomID], roomChan)

	messages, err := g.uc.EnterRoom(ctx, roomID)
	// slog.Debug("Getting Messages", "messages", messages)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return sharedkernel.ErrServer
	}

	if len(messages) == 0 {
		slog.Debug("Empty messages")
		srv.Send(nil)
	}

	for _, message := range messages {
		srv.Send(&gen.Message{
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

	go func() {
		var message *gen.Message
		for {
			select {
			case message = <-roomChan:
				// slog.Debug("received Messsage", "message", message)
				srv.Send(message)
			case <-done:
				return
			}
		}
	}()

	<-done
	slog.Debug("Chat leaved")
	return nil
}

func (roomGRPCServer) UpdateRoom(_ context.Context, _ *gen.UpdateRoomRequest) (_ *gen.UpdateRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) DeleteRoom(_ context.Context, _ *gen.DeleteRoomRequest) (_ *gen.DeleteRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (g *roomGRPCServer) SendMessage(ctx context.Context, req *gen.SendMessageRequest) (*gen.SendMessageResponse, error) {
	slog.Info("POST: SendMessage")

	id, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}
	err = g.uc.SendMessage(ctx, &domain.CreateMessageModel{
		Content: req.GetText(),
		RoomID:  id,
	})
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.CreateRoom"))
		return nil, sharedkernel.ErrServer
	}

	go func() {
		for _, rChan := range roomMessages[id] {
			rChan <- &gen.Message{
				Id:  req.RoomId,
				Text: req.Text,
			}
		}
	}()

	return &gen.SendMessageResponse{}, nil
}
