package rooms

import (
	"context"

	"github.com/google/uuid"

	"github.com/Watch2Gather/server/internal/room/domain"
)

type (
	RoomRepo interface {
		CreateRoom(context.Context, *domain.CreateRoomModel) (*domain.RoomModel, error)
		CreateMessage(context.Context, *domain.CreateMessageModel) error
		GetRoomsByUserID(context.Context, uuid.UUID) ([]*domain.RoomModel, error)
		GetParticipantsByRoomID(context.Context, uuid.UUID) ([]*domain.ParticipantModel, error)
		GetMessagesByRoomID(context.Context, *domain.MessagesByRoomIDModel) ([]*domain.MessageModel, error)
		UpdateRoom(context.Context, *domain.UpdateRoomModel) (*domain.RoomModel, error)
		DeleteRoom(context.Context, *domain.DeleteRoomModel) error
		AddParticipants(context.Context, *domain.AddParticipantsModel) error
		RemoveParticipant(context.Context, *domain.RemoveParticipantModel) error
		GetRoomOwner(context.Context, uuid.UUID) (uuid.UUID, error)
	}
	UseCase interface {
		CreateRoom(context.Context, *domain.CreateRoomModel) (*domain.RoomModel, error)
		DeleteRoom(context.Context, uuid.UUID) error
		LeaveRoom(context.Context, uuid.UUID) error
		SendMessage(context.Context, *domain.MessageModel) error
		GetRoomsByUser(context.Context) ([]*domain.RoomModel, error)
		InviteToRoom(context.Context, *domain.AddParticipantsModel) (*domain.RoomModel, error)
		EnterRoom(context.Context, uuid.UUID) (any, error)
		UpdateRoom(context.Context, *domain.RoomModel) (*domain.RoomModel, error)
	}
)
