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

	var resRoom domain.RoomModel

	resRoom.Name = room.Name
	resRoom.OwnerID = room.OwnerID
	resRoom.ID = room.ID
	resRoom.ParticipantsCount = len(model.ParticipantIds)

	return &resRoom, nil
}

func (r *roomRepo) CreateMessage(ctx context.Context, model *domain.CreateMessageModel) error {
	querier := postgresql.New(r.pg.GetDB())

	_, err := querier.CreateMessage(ctx, postgresql.CreateMessageParams{
		RoomID:  model.RoomID,
		UserID:  model.UserID,
		Content: model.Content,
	})
	if err != nil {
		return errors.Wrap(err, "querier.CreateMessage")
	}

	return nil
}

func (r *roomRepo) GetRoomsByUserID(ctx context.Context, id uuid.UUID) (_ []*domain.RoomModel, _ error) {
	querier := postgresql.New(r.pg.GetDB())

	rooms, err := querier.GetRoomsByUserId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetRoomsByUserID")
	}

	var roomsModel []*domain.RoomModel

	for _, room := range rooms {
		roomsModel = append(roomsModel, &domain.RoomModel{
			Name:              room.Name,
			OwnerID:           room.OwnerID,
			MovieID:           room.MovieID.UUID,
			Timecode:          0,
			ID:                room.ID,
			ParticipantsCount: int(room.Count),
		})
	}

	return roomsModel, nil
}

func (r *roomRepo) GetParticipantsByRoomID(ctx context.Context, id uuid.UUID) (uuid.UUIDs,error) {
	querier := postgresql.New(r.pg.GetDB())

	participants, err := querier.GetParticipantsByRoomId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetParticipantsByRoomID")
	}

	ids := make(uuid.UUIDs, 0)
	for _, parparticipant := range participants {
		ids = append(ids, parparticipant.ID)
	}


	return ids, nil
}

func (r *roomRepo) GetMessagesByRoomID(ctx context.Context, model *domain.MessagesByRoomIDModel) ([]*domain.MessageModel, error) {
	querier := postgresql.New(r.pg.GetDB())

	messages, err := querier.GetMessagesByRoomId(ctx, postgresql.GetMessagesByRoomIdParams{
		RoomID: model.RoomID,
		Limit:  int32(model.Limit),
		Offset: int32(model.Offset),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetMessagesByRoomID")
	}

	var messagesModel []*domain.MessageModel

	for _, message := range messages {
		messagesModel = append(messagesModel, &domain.MessageModel{
			Content:   message.Content,
			CreatedAt: int(message.CreatedAt.UnixNano()),
			UserID:    message.UserID,
			RoomID:    message.RoomID,
			MessageID: message.ID,
		})
	}

	return messagesModel, nil
}

func (roomRepo) DeleteRoom(_ context.Context, _ *domain.DeleteRoomModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (r *roomRepo) AddParticipants(ctx context.Context, model *domain.AddParticipantsModel) error {
	db := r.pg.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "db.Begin")
	}
	defer tx.Rollback()

	querier := postgresql.New(db)

	qtx := querier.WithTx(tx)

	st := status.New(codes.InvalidArgument, "users not found")
	br := &errdetails.BadRequest{}

	for _, p := range model.ParticipantIds {
		err = qtx.AddParticipant(ctx, postgresql.AddParticipantParams{
			RoomID: model.RoomID,
			UserID: p,
		})
		if err != nil {
			v := &errdetails.BadRequest_FieldViolation{
				Field:       "userID",
				Description: "userID " + p.String() + " does not exist",
			}
			br.FieldViolations = append(br.FieldViolations, v)
		}
	}
	if len(br.FieldViolations) > 0 {
		st, err = st.WithDetails(br)
		if err != nil {
			slog.Error("Caught error", "trace", errors.Wrap(err, "st.WithDetails"))
			return sharedkernel.ErrServer
		}
		return st.Err()
	}

	return tx.Commit()
}

func (roomRepo) RemoveParticipant(_ context.Context, _ *domain.RemoveParticipantModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (roomRepo) UpdateRoom(_ context.Context, _ *domain.UpdateRoomModel) (_ *domain.RoomModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (r *roomRepo) GetRoomOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	querier := postgresql.New(r.pg.GetDB())

	ownerID, err := querier.GetRoomOwner(ctx, id)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "querier.GetRoomOwner")
	}

	return ownerID, nil
}
