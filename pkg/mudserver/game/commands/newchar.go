package commands

import (
	"github.com/mikejk8s/talesmud/pkg/mudserver/game/def"
	"github.com/mikejk8s/talesmud/pkg/mudserver/game/messages"
)

// NewCharacterCommand ...
type NewCharacterCommand struct {
}

// Key ...
func (command *NewCharacterCommand) Key() CommandKey { return &ExactCommandKey{} }

// Execute ... executes the command
func (command *NewCharacterCommand) Execute(game def.GameCtrl, message *messages.Message) bool {

	// just send a message to the client to start character creation
	game.SendMessage() <- messages.NewCreateCharacterMessage(message.FromUser.ID)
	return true
}
