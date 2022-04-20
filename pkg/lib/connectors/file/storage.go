package file

import (
	_ "go.beyondstorage.io/services/fs/v4"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func getStorage(provider map[string]interface{}) (types.Storager, error) {
	code := provider["code"].(string)
	return services.NewStoragerFromString(code + "://")

}
