package rooms

import (
	"context"

	"github.com/google/uuid"

	"github.com/Watch2Gather/server/internal/room/domain"
)

type (
	RoomRepo interface {
		CreateRoom(context.Context, *domain.CreateRoomModel) (*domain.RoomModel, error)
		CreateMessage(context.Context, *domain.CreateMessageModel) (*domain.SendMessageModel, error)
		GetRoomsByUserID(context.Context, uuid.UUID) ([]*domain.GetRoomsModel, error)
		GetParticipantsByRoomID(context.Context, uuid.UUID) ([]*domain.ParticipantModel, error)
		GetMessagesByRoomID(context.Context, uuid.UUID) ([]*domain.MessageModel, error)
		UpdateRoom(context.Context, *domain.UpdateRoomModel) error
		DeleteRoom(context.Context, *domain.DeleteRoomModel) error
		AddParticipants(context.Context, *domain.AddParticipantsModel) error
		RemoveParticipant(context.Context, *domain.RemoveParticipantModel) error
		GetRoomOwner(context.Context, uuid.UUID) (uuid.UUID, error)
		AddMovieToRoom(context.Context, *domain.AddMovieModel) error
	}
	UseCase interface {
		CreateRoom(context.Context, *domain.CreateRoomModel) (*domain.RoomModel, error)
		DeleteRoom(context.Context, uuid.UUID) error
		LeaveRoom(context.Context, uuid.UUID) error
		SendMessage(context.Context, *domain.CreateMessageModel) (*domain.SendMessageModel, error)
		GetRoomsByUser(context.Context) ([]*domain.GetRoomsModel, error)
		GetParticipantsByRoomID(context.Context, uuid.UUID) ([]*domain.ParticipantModel, error)
		InviteToRoom(context.Context, *domain.AddParticipantsModel) (*domain.RoomModel, error)
		GetRoomMessages(context.Context, uuid.UUID) ([]*domain.MessageModel, error)
		UpdateRoom(context.Context, *domain.UpdateRoomModel) error
		GetUserInfo(context.Context, uuid.UUID) (*domain.UserModel, error)
		AddMovieToRoom(context.Context, *domain.AddMovieModel) error
	}
)
