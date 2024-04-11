package rooms

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"

	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/room/domain"
)

type usecase struct {
	roomRepo RoomRepo
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	roomRepo RoomRepo,
) UseCase {
	return &usecase{
		roomRepo: roomRepo,
	}
}

func (u *usecase) CreateRoom(ctx context.Context, model *domain.CreateRoomModel) (*domain.RoomModel, error) {
	if len(model.Name) == 0 {
		model.Name = "testMovie" // TODO add movie names
	}

	id, err := sharedkernel.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetUserIDFromContext")
	}

	model.OwnerID = id

	room, err := u.roomRepo.CreateRoom(ctx, model)
	if err != nil {
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

func (usecase) SendMessage(ctx context.Context, model *domain.MessageModel) error {
	panic("not implemented") // TODO: Implement
}

func (usecase) GetRoomsByUser(ctx context.Context, id uuid.UUID) ([]*domain.RoomModel, error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) InviteToRoom(ctx context.Context, ids []uuid.UUID) (*domain.RoomModel, error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) EnterRoom(ctx context.Context, id uuid.UUID) (any, error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) UpdateRoom(ctx context.Context, model *domain.RoomModel) (*domain.RoomModel, error) {
	panic("not implemented") // TODO: Implement
}
