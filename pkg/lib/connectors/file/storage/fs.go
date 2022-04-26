package storage

import (
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func getFSStorage(config map[string]interface{}) (types.Storager, error) {
	return services.NewStoragerFromString("fs://")
}
