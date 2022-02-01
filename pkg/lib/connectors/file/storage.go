package file

import (
	_ "go.beyondstorage.io/services/fs/v4"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func getFS(storageConfig StorageConfig) (types.Storager, error) {
	url := storageConfig.GetURL()
	return services.NewStoragerFromString(url)
}

func getStorage(storageType string, storageConfig StorageConfig) (types.Storager, error) {

	switch storageType {
	case "fs":
		return getFS(storageConfig)
	}
	return nil, nil
}
