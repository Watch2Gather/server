package repo

import (
	"context"
	"log/slog"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/room/domain"
	"github.com/Watch2Gather/server/internal/room/infras/postgresql"
	"github.com/Watch2Gather/server/internal/room/usecases/rooms"
	"github.com/Watch2Gather/server/pkg/postgres"
)

const _defaultEntityCap = 64

type roomRepo struct {
	pg postgres.DBEngine
}

var _ rooms.RoomRepo = (*roomRepo)(nil)

var RepositorySet = wire.NewSet(NewRoomRepo)

func NewRoomRepo(
	pg postgres.DBEngine,
) rooms.RoomRepo {
	return &roomRepo{pg: pg}
}

func (r roomRepo) CreateRoom(ctx context.Context, model *domain.CreateRoomModel) (*domain.RoomModel, error) {
	resRoom := &domain.RoomModel{}
	db := r.pg.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "db.Begin")
	}
	defer tx.Rollback()

	querier := postgresql.New(db)

	qtx := querier.WithTx(tx)
	room, err := qtx.CreateRoom(ctx, postgresql.CreateRoomParams{
		Name:    model.Name,
		OwnerID: model.OwnerID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.CreateRoom")
	}
	resRoom.Name = room.Name
	resRoom.OwnerID = room.OwnerID
	resRoom.ID = room.ID

	st := status.New(codes.InvalidArgument, "users not found")
	br := &errdetails.BadRequest{}

	for _, p := range model.ParticipantIds {
		err = qtx.AddParticipant(ctx, postgresql.AddParticipantParams{
			RoomID: room.ID,
			UserID: p,
		})
		if err != nil {
			v := &errdetails.BadRequest_FieldViolation{
				Field:       "userID",
				Description: "userID " + p.String() + " does not exist",
			}
			br.FieldViolations = append(br.FieldViolations, v)
		}
		resRoom.ParticipantIds = append(resRoom.ParticipantIds, p)
	}
	if len(br.FieldViolations) > 0 {
		st, err = st.WithDetails(br)
		if err != nil {
			slog.Error("Caught error", "trace", errors.Wrap(err, "st.WithDetails"))
			return nil, sharedkernel.ErrServer
		}
		return nil, st.Err()
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "tx.Commit")
	}
	return resRoom, nil
}

func (roomRepo) CreateMessage(_ context.Context, _ *domain.CreateMessageModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) GetRoomsByUserID(_ context.Context, _ uuid.UUID) (_ []*domain.RoomModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) GetParticipantsByRoomID(_ context.Context, _ uuid.UUID) (_ []*domain.ParticipantModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) GetMessagesByRoomID(_ context.Context, _ *domain.MessagesByRoomIDModel) (_ []*domain.MessageModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) DeleteRoom(_ context.Context, _ *domain.DeleteRoomModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) AddParticipant(_ context.Context, _ *domain.AddParticipantModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) RemoveParticipant(_ context.Context, _ *domain.RemoveParticipantModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) UpdateRoom(_ context.Context, _ *domain.UpdateRoomModel) (_ *domain.RoomModel, _ error) {
	panic("not implemented") // TODO: Implement
}
