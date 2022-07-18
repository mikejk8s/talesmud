package characters

import (
	"github.com/mikejk8s/talesmud/pkg/entities/items"
)

//Inventory data
type Inventory struct {
	Size  int32         `json:"size"`
	Items []*items.Item `json:"items"`
}
