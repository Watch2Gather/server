package repo

import (
	"context"
	"database/sql"
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

func (r *roomRepo) CreateMessage(ctx context.Context, model *domain.CreateMessageModel) (*domain.SendMessageModel, error) {
	querier := postgresql.New(r.pg.GetDB())

	msg, err := querier.CreateMessage(ctx, postgresql.CreateMessageParams{
		RoomID:  model.RoomID,
		UserID:  model.UserID,
		Content: model.Content,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.CreateMessage")
	}

	return &domain.SendMessageModel{
		Text:      msg.Content,
		CreatedAt: msg.CreatedAt.Unix(),
		ID:        msg.ID,
	}, nil
}

func (r *roomRepo) GetRoomsByUserID(ctx context.Context, id uuid.UUID) ([]*domain.GetRoomsModel, error) {
	querier := postgresql.New(r.pg.GetDB())

	rooms, err := querier.GetRoomsByUserId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetRoomsByUserID")
	}

	var roomsModel []*domain.GetRoomsModel

	for _, room := range rooms {
		var mID, poster string
		if room.PosterPath.Valid {
			poster = room.PosterPath.String
		}
		if room.MovieID.Valid {
			mID = room.MovieID.UUID.String()
		}
		roomsModel = append(roomsModel, &domain.GetRoomsModel{
			Name:              room.Name,
			FilmTitle:         room.Title.String,
			PosterPath:        poster,
			ParticipantsCount: int(room.Count),
			Timecode:          room.Timecode.Time.Second(),
			OwnerID:           room.OwnerID,
			ID:                room.ID,
			MovieID:           mID,
		})
	}

	return roomsModel, nil
}

func (r *roomRepo) GetParticipantsByRoomID(ctx context.Context, id uuid.UUID) ([]*domain.ParticipantModel, error) {
	querier := postgresql.New(r.pg.GetDB())

	participants, err := querier.GetParticipantsByRoomId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetParticipantsByRoomID")
	}

	res := []*domain.ParticipantModel{}
	for _, parparticipant := range participants {
		res = append(res, &domain.ParticipantModel{
			Username: parparticipant.Username,
			Avatar:   parparticipant.Avatar.String,
			UserID:   parparticipant.ID,
		})
	}

	return res, nil
}

func (r *roomRepo) GetMessagesByRoomID(ctx context.Context, id uuid.UUID) ([]*domain.MessageModel, error) {
	querier := postgresql.New(r.pg.GetDB())

	messages, err := querier.GetMessagesByRoomId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetMessagesByRoomID")
	}

	var messagesModel []*domain.MessageModel

	for _, message := range messages {
		messagesModel = append(messagesModel, &domain.MessageModel{
			Text: message.Content,
			User: domain.User{
				Name:   message.Username,
				Avatar: message.Avatar.String,
				ID:     message.UID,
			},
			CreatedAt: message.CreatedAt.Unix(),
			ID:        message.MID,
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

func (r *roomRepo) UpdateRoom(ctx context.Context, model *domain.UpdateRoomModel) error {
	querier := postgresql.New(r.pg.GetDB())

	var mID uuid.NullUUID
	if model.MovieID != uuid.Nil {
		mID.UUID = model.MovieID
		mID.Valid = true
	}
	err := querier.UpdateRoom(ctx, postgresql.UpdateRoomParams{
		Name: sql.NullString{
			String: model.Name,
			Valid:  true,
		},
		MovieID: mID,
		ID: uuid.NullUUID{
			UUID:  model.RoomID,
			Valid: true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "querier.GetRoomOwner")
	}

	return nil
}

func (r *roomRepo) GetRoomOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	querier := postgresql.New(r.pg.GetDB())

	ownerID, err := querier.GetRoomOwner(ctx, id)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "querier.GetRoomOwner")
	}

	return ownerID, nil
}

func (r *roomRepo) AddMovieToRoom(ctx context.Context, model *domain.AddMovieModel) error {
	querier := postgresql.New(r.pg.GetDB())

	err := querier.AddMovieToRoom(ctx, postgresql.AddMovieToRoomParams{
		ID: model.RoomID,
		MovieID: uuid.NullUUID{
			UUID:  model.MovieID,
			Valid: true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "querier.AddMovieToRoom")
	}

	return nil
}
