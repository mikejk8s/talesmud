package game

import (
	"time"

	"github.com/mikejk8s/talesmud/pkg/entities"
	log "github.com/sirupsen/logrus"

	c "github.com/mikejk8s/talesmud/pkg/mudserver/game/commands"
	m "github.com/mikejk8s/talesmud/pkg/mudserver/game/messages"
	"github.com/mikejk8s/talesmud/pkg/service"
)

//Game ... contains live game state
type Game struct {
	id    string
	title string

	// access to repository data
	Facade service.Facade

	SystemUser *entities.User

	// messages
	onMessageReceived chan interface{}
	sendMessage       chan interface{}

	OnUserJoined chan *m.UserJoined
	OnUserQuit   chan *m.UserQuit

	//OnAvatarJoinedRoom chan *AvatarJoinedRoom
	//OnAvatarLeftRoom   chan *AvatarLeftRoom

	Receivers []Receiver

	CommandProcessor *c.CommandProcessor
	RoomProcessor    *c.RoomProcessor

	Avatars map[string]*Avatar

	//world *World
}

// New creates a new game instance
func New(facade service.Facade) *Game {
	return &Game{

		title: "Lair of the Dragon",

		CommandProcessor: c.NewCommandProcessor(),
		RoomProcessor:    c.NewRoomProcessor(),

		// event channels
		onMessageReceived: make(chan interface{}, 20),
		sendMessage:       make(chan interface{}, 20),
		OnUserJoined:      make(chan *m.UserJoined, 20),
		OnUserQuit:        make(chan *m.UserQuit, 20),

		// game update listeners
		//	Receivers: make([]Receiver, 0, 10),

		Avatars: make(map[string]*Avatar),

		Facade: facade,
	}
}

// Subscribe ... sub
//func (g *Game) Subscribe(receiver Receiver) {
//	g.Receivers = append(g.Receivers, receiver)
//}

// Receiver ... rec
type Receiver interface {
	OnMessage(message interface{})
}

//Unsubscribe ...
func (g *Game) Unsubscribe(receiver *Receiver) {
	//TODO:
	//game.Receivers = delete(game.Receivers, receiver)
}

/*
// SendMessage ...
func (g *Game) SendMessage(msg interface{}) {
	// broeadcast message to all game listeners (currently only websocket server)
	for _, receiver := range g.Receivers {
		receiver.OnMessage(msg)
	}
}*/

// SendMessage ...
func (g *Game) SendMessage() chan interface{} {
	return g.sendMessage
}

//OnMessageReceived returns onMessageReceived channel
func (g *Game) OnMessageReceived() chan interface{} {
	return g.onMessageReceived
}

//GetFacade ...
func (g *Game) GetFacade() service.Facade {
	return g.Facade
}

const roomUpdateInterval = 10
const npcUpdateInterval = 10

func (g *Game) handleGameUpdates() {

	roomTicker := time.NewTicker(roomUpdateInterval * time.Second)
	npcTicker := time.NewTicker(npcUpdateInterval * time.Second)

	for {
		select {
		case <-roomTicker.C:
			g.handleRoomUpdates()
		case <-npcTicker.C:
			//server.handleUserPings()
		}
	}
}

//Run main game loop
func (g *Game) Run() {

	go g.handleGameUpdates()

	go func() {
		for {
			select {
			case userJoined := <-g.OnUserJoined:
				log.Info("Received UserJoinged message")
				g.handleUserJoined(userJoined.User)

			case userQuit := <-g.OnUserQuit:
				log.WithField("user", userQuit.User).Info("Received UserQuit message")
				g.handleUserQuit(userQuit.User)

			case msg := <-g.onMessageReceived:
				switch message := msg.(type) {
				case *m.Message:
					// attach current character if a user is set
					g.attachCharacterToMessage(message)

					// only broadcast if global commandprocessor didnt process it
					if !g.CommandProcessor.Process(g, message) {
						// check room commands
						if !g.RoomProcessor.Process(g, message) {
							// generic messages will be converted to plain OutgoingMessages (type message)
							// and send to the room audience including the origin nickname or charactername
							g.handleDefaultMessage(message)
						}
					}
				}
			}
		}
	}()
}
