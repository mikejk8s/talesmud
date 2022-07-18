package commands

import (
	"fmt"

	"github.com/mikejk8s/talesmud/pkg/mudserver/game/def"
	"github.com/mikejk8s/talesmud/pkg/mudserver/game/messages"
)

// ListCharactersCommand ... foo
type ListCharactersCommand struct {
}

// Key ...
func (command *ListCharactersCommand) Key() CommandKey { return &ExactCommandKey{} }

// Execute ... executes scream command
func (command *ListCharactersCommand) Execute(game def.GameCtrl, message *messages.Message) bool {

	if characters, err := game.GetFacade().CharactersService().FindAllForUser(message.FromUser.ID); err == nil {

		result := "Your Characters:\n"

		for _, character := range characters {
			result += fmt.Sprintf("- %v [LVL %v %v %vxp] %v - %v\n", character.Name, character.Level, character.Class.Name, character.XP, character.Race.Name, character.Description)
		}
		result += "To select character use: sc [charactername]"

		game.SendMessage() <- messages.Reply(message.FromUser.ID, result)

	} else {
		game.SendMessage() <- messages.Reply(message.FromUser.ID, "You have no characters.")
	}

	return true
}
