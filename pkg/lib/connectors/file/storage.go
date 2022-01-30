package file

import "go.beyondstorage.io/v5/types"

func getFS(storageConfig StorageConfig) (types.Storager, error) {
	return nil, nil
}

func getStorage(storageType string, storageConfig StorageConfig) (types.Storager, error) {
	switch storageType {
	case "fs":
		return getFS(storageConfig)
	}
	return nil, nil
}
