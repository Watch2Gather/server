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

type RemoveParticipantModel struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

type UserModel struct {
	Name   string
	Avatar string
	ID     uuid.UUID
}

type AddMovieModel struct {
	RoomID  uuid.UUID
	MovieID uuid.UUID
}

type RoomModel struct {
	Name              string
	ParticipantsCount int
	Timecode          int
	OwnerID           uuid.UUID
	MovieID           uuid.UUID
	ID                uuid.UUID
}

type GetRoomsModel struct {
	Name              string
	FilmTitle         string
	PosterPath        string
	ParticipantsCount int
	Timecode          int
	OwnerID           uuid.UUID
	ID                uuid.UUID
	MovieID           string
}

type ParticipantModel struct {
	Username string
	Avatar   string
	UserID   uuid.UUID
}

type RoomsByUserModel struct {
	ID     uuid.UUID
	Limit  int
	Offset int
}

type MessageModel struct {
	Text      string
	User      User
	CreatedAt int64
	ID        uuid.UUID
}

type User struct {
	Name   string
	Avatar string
	ID     uuid.UUID
}

type AddParticipantsModel struct {
	ParticipantIds uuid.UUIDs
	RoomID         uuid.UUID
	OwnerID        uuid.UUID
}

type SendMessageModel struct {
	Text      string
	CreatedAt int64
	ID        uuid.UUID
}
