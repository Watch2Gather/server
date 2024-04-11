package domain

import "github.com/google/uuid"

type CreateRoomModel struct {
	Name           string
	ParticipantIds uuid.UUIDs
	OwnerID        uuid.UUID
}

type DeleteRoomModel struct {
	RoomID  uuid.UUID
	OwnerID uuid.UUID
}

type UpdateRoomModel struct {
	Name    string
	MovieID uuid.UUID
	RoomID  uuid.UUID
}

type CreateMessageModel struct {
	Content string
	RoomID  uuid.UUID
	UserID  uuid.UUID
}

type MessagesByRoomIDModel struct {
	RoomID uuid.UUID
	Limit  int
	Offset int
}

type AddParticipantModel struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

type RemoveParticipantModel struct {
	roomID uuid.UUID
	UserID uuid.UUID
}

type RoomModel struct {
	Name           string
	OwnerID        uuid.UUID
	MovieID        uuid.UUID
	Timecode       int
	ID             uuid.UUID
	ParticipantIds uuid.UUIDs
}

type ParticipantModel struct {
	Username string
	Avatar   string
	UserID   uuid.UUID
}

type MessageModel struct {
	Content   string
	CreatedAt int
	UserID    uuid.UUID
	RoomID    uuid.UUID
	MessageID uuid.UUID
}
