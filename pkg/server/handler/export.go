package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/mikejk8s/talesmud/pkg/entities"
	"github.com/mikejk8s/talesmud/pkg/entities/characters"
	"github.com/mikejk8s/talesmud/pkg/entities/items"
	"github.com/mikejk8s/talesmud/pkg/entities/rooms"
	"github.com/mikejk8s/talesmud/pkg/repository"
	"github.com/mikejk8s/talesmud/pkg/scripts"
	"github.com/mikejk8s/talesmud/pkg/service"
)

//ExportHandler ...
type ExportHandler struct {
	RoomsService      service.RoomsService
	CharactersService service.CharactersService
	UserService       service.UsersService
	ItemsService      service.ItemsService
	ScriptService     service.ScriptsService
}

type exportStructure struct {
	Rooms         []*rooms.Room           `json:"rooms"`
	Items         []*items.Item           `json:"items"`
	ItemTemplates []*items.ItemTemplate   `json:"itemTemplates"`
	Characters    []*characters.Character `json:"characters"`
	Scripts       []*scripts.Script       `json:"scripts"`
	Users         []*e.User               `json:"users"`
}

//Export Exports all data structures as JSON
func (handler *ExportHandler) Export(c *gin.Context) {

	d := exportStructure{}

	d.Rooms, _ = handler.RoomsService.FindAll()
	d.Characters, _ = handler.CharactersService.FindAll()
	d.Users, _ = handler.UserService.FindAll()
	d.ItemTemplates, _ = handler.ItemsService.ItemTemplates().FindAll(repository.ItemsQuery{})
	d.Items, _ = handler.ItemsService.Items().FindAll(repository.ItemsQuery{})
	d.Scripts, _ = handler.ScriptService.FindAll()

	//c.JSON(http.StatusOK, d)
	c.IndentedJSON(http.StatusOK, d)
}

//Import Imports all data structures
func (handler *ExportHandler) Import(c *gin.Context) {

	// drop all collections before importing
	handler.RoomsService.Drop()
	handler.CharactersService.Drop()
	handler.UserService.Drop()
	handler.ItemsService.Items().Drop()
	handler.ItemsService.ItemTemplates().Drop()
	handler.ScriptService.Drop()

	var data exportStructure
	//if err := c.ShouldBindYAML(&data); err != nil {
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, room := range data.Rooms {
		handler.RoomsService.Import(room)
	}

	for _, character := range data.Characters {
		handler.CharactersService.Import(character)
	}

	for _, user := range data.Users {
		handler.UserService.Import(user)
	}

	for _, item := range data.Items {
		handler.ItemsService.Items().Import(item)
	}
	for _, itemTemplate := range data.ItemTemplates {
		handler.ItemsService.ItemTemplates().Import(itemTemplate)
	}

	for _, script := range data.Scripts {
		handler.ScriptService.Import(script)
	}

	c.JSON(http.StatusOK, gin.H{"status": "Import successful"})
}
