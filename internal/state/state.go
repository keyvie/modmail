package state

import (
	"github.com/keyvie/modmail/internal/database"
)

type Public struct {
	Data *database.Data
}

func New(data *database.Data) *Public {
	return &Public{
		Data: data,
	}
}
