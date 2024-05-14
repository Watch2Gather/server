package rooms

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/room/domain"
)

type usecase struct {
	roomRepo          RoomRepo
	userInfoDomainSvc domain.UserInfoDomainService
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	roomRepo RoomRepo,
	userInfoDomainSvc domain.UserInfoDomainService,
) UseCase {
	return &usecase{
		roomRepo:          roomRepo,
		userInfoDomainSvc: userInfoDomainSvc,
	}
}

var roomMessages map[uuid.UUID]sharedkernel.Broker[struct{}]

func (u *usecase) CreateRoom(ctx context.Context, model *domain.CreateRoomModel) (*domain.RoomModel, error) {
	if len(model.Name) == 0 {
		model.Name = "testMovie" // TODO add movie names
	}

	id, err := sharedkernel.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetUserIDFromContext")
	}

	model.OwnerID = id
	model.ParticipantIds = append(model.ParticipantIds, model.OwnerID)

	room, err := u.roomRepo.CreateRoom(ctx, model)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, err
		}
		return nil, errors.Wrap(err, "roomRepo.CreateRoom")
	}

	return room, nil
}

func (usecase) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}

func (usecase) LeaveRoom(ctx context.Context, _ uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}

func (u *usecase) SendMessage(ctx context.Context, model *domain.CreateMessageModel) (*domain.SendMessageModel, error) {
	id, err := sharedkernel.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetUserIDFromContext")
	}

	model.UserID = id

	msg, err := u.roomRepo.CreateMessage(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "roomRepo.CreateMessage")
	}

	return msg, nil
}

func (u *usecase) GetRoomsByUser(ctx context.Context) ([]*domain.GetRoomsModel, error) {
	id, err := sharedkernel.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetUserIDFromContext")
	}

	rooms, err := u.roomRepo.GetRoomsByUserID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "roomRepo.GetRoomsByUserID")
	}

	return rooms, nil
}

func (u *usecase) GetParticipantsByRoomID(ctx context.Context, id uuid.UUID) (uuid.UUIDs, error) {
	participants, err := u.roomRepo.GetParticipantsByRoomID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "roomRepo.GetRoomsByUserID")
	}

	return participants, nil
}

func (usecase) InviteToRoom(ctx context.Context, model *domain.AddParticipantsModel) (*domain.RoomModel, error) {
	panic("not implemented") // TODO: Implement
}

func (u *usecase) GetRoomMessages(ctx context.Context, id uuid.UUID) ([]*domain.MessageModel, error) {
	messages, err := u.roomRepo.GetMessagesByRoomID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "roomRepo.CreateMessage")
	}
	return messages, nil
}

func (u *usecase) UpdateRoom(ctx context.Context, model *domain.UpdateRoomModel) error {
	err := u.roomRepo.UpdateRoom(ctx, &domain.UpdateRoomModel{
		Name:    model.Name,
		MovieID: model.MovieID,
		RoomID:  model.RoomID,
	})
	if err != nil {
		return errors.Wrap(err, "roomRepo.CreateMessage")
	}
	return nil
}

func (u *usecase) GetUserInfo(ctx context.Context, id uuid.UUID) (*domain.UserModel, error) {
	user, err := u.userInfoDomainSvc.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "userInfo.GetUserInfo")
	}
	return user, nil
}

func (u *usecase) AddMovieToRoom(ctx context.Context, model *domain.AddMovieModel) error {
	err := u.roomRepo.AddMovieToRoom(ctx, model)
	if err != nil {
		return errors.Wrap(err, "roomRepo.AddMovieToRoom")
	}
	return nil
}
