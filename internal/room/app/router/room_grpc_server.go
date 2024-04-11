package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/lib/pq"
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

func (g *roomGRPCServer) CreateRoom(ctx context.Context, req *gen.CreateRoomRequest) (*gen.CreateRoomResponse, error) {
	slog.Info("POST: CreateRoom")

	if len(req.RoomName) == 0 {
		req.RoomName = "testMovie" // TODO add movie names
	}

	var ids uuid.UUIDs
	partErr := status.New(codes.InvalidArgument, "users not found")
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, domain.ErrUserAlreadyExists
			}
		}
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

	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) InviteToRoom(_ context.Context, _ *gen.InviteToRoomRequest) (_ *gen.InviteToRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) EnterRoom(req *gen.EnterRoomRequest, srv gen.RoomService_EnterRoomServer) error {
	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) UpdateRoom(_ context.Context, _ *gen.UpdateRoomRequest) (_ *gen.UpdateRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) DeleteRoom(_ context.Context, _ *gen.DeleteRoomRequest) (_ *gen.DeleteRoomResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomGRPCServer) SendMessage(_ context.Context, _ *gen.Message) (_ *gen.SendMessageResponse, _ error) {
	panic("not implemented") // TODO: Implement
}
