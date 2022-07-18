package messages

import (
	"github.com/mikejk8s/talesmud/pkg/entities"
	e "github.com/mikejk8s/talesmud/pkg/entities"
	"github.com/mikejk8s/talesmud/pkg/entities/characters"
	"github.com/mikejk8s/talesmud/pkg/entities/rooms"
	"github.com/mikejk8s/talesmud/pkg/mudserver/game/def"
	"github.com/mikejk8s/talesmud/pkg/mudserver/game/util"
)

// CharacterJoinedRoom ... asdd
type CharacterJoinedRoom struct {
	MessageResponse
}

// CharacterLeftRoom ... asdd
type CharacterLeftRoom struct {
	MessageResponse
}

// CharacterSelected ...
type CharacterSelected struct {
	MessageResponse
	Character *characters.Character `json:"character"`
}

// NewUserQuit ... creates a new User Joined event
func NewUserQuit(user *e.User) *UserQuit {
	return &UserQuit{
		User: user,
	}
}

// NewUserJoined ... creates a new User Joined event
func NewUserJoined(user *e.User) *UserJoined {
	return &UserJoined{
		User: user,
	}
}

//AudienceType type
type AudienceType int

const (
	MessageAudienceOrigin = iota + 1
	MessageAudienceUser
	MessageAudienceRoom
	MessageAudienceRoomWithoutOrigin
	MessageAudienceGlobal
	MessageAudienceSystem
)

// MessageResponse ... Define our message object
type MessageResponse struct {
	Audience   AudienceType `json:"-"`
	AudienceID string       `json:"-"`
	OriginID   string       `json:"-"`

	Type     MessageType `json:"type"`
	Username string      `json:"username"`
	Message  string      `json:"message"`
}

//GetAudience ,,,
func (m MessageResponse) GetAudience() AudienceType {
	return m.Audience
}

//GetAudienceID ,,,
func (m MessageResponse) GetAudienceID() string {
	return m.AudienceID
}

//GetOriginID ,,,
func (m MessageResponse) GetOriginID() string {
	return m.OriginID
}

//GetMessage ,,,
func (m MessageResponse) GetMessage() string {
	return m.Message
}

//MessageResponder ...
type MessageResponder interface {
	GetAudience() AudienceType
	GetAudienceID() string
	GetOriginID() string
	GetMessage() string
}

// MultiResponse ...
type MultiResponse struct {
	Responses []MessageResponse
}

//NewMultiResponse ...
func NewMultiResponse(responses ...MessageResponse) MultiResponse {
	mr := MultiResponse{
		Responses: []MessageResponse{},
	}
	for _, rsp := range responses {
		mr.Responses = append(mr.Responses, rsp)
	}
	return mr
}

// EnterRoomMessage ... Define our message object
type EnterRoomMessage struct {
	MessageResponse
	Room rooms.Room `json:"room"`
}

//NewEnterRoomMessage ...
func NewEnterRoomMessage(room *rooms.Room, user *entities.User, game def.GameCtrl) *EnterRoomMessage {
	return &EnterRoomMessage{
		MessageResponse: MessageResponse{
			Audience: MessageAudienceOrigin,
			Type:     MessageTypeEnterRoom,
			Message:  util.CreateRoomDescription(room, user, game),
		},
		Room: *room,
	}
}

// NewRoomBasedMessage ... creates a new Websocket message
func NewRoomBasedMessage(user string, message string) MessageResponse {
	return MessageResponse{
		// default
		Audience: MessageAudienceRoom,
		Type:     MessageTypeDefault,
		Message:  message,
		Username: user,
	}
}

// Reply ... creates a reply message
func Reply(userID string, message string) MessageResponse {
	return MessageResponse{
		// default
		Audience:   MessageAudienceOrigin,
		AudienceID: userID,
		Type:       MessageTypeDefault,
		Message:    message,
	}
}

// NewCreateCharacterMessage ...
func NewCreateCharacterMessage(user string) MessageResponse {
	return MessageResponse{
		Type:       MessageTypeCreateCharacter,
		Message:    "User has no characters created.",
		Audience:   MessageAudienceOrigin,
		AudienceID: user,
	}
}
