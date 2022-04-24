package storage

import (
	"errors"

	_ "go.beyondstorage.io/services/fs/v4"
	_ "go.beyondstorage.io/services/s3/v3"
	"go.beyondstorage.io/v5/types"
)

func GetStorage(provider map[string]interface{}) (types.Storager, error) {
	code := provider["code"].(string)
	config := provider["config"].(map[string]interface{})
	switch code {
	case "fs":
		return getFSStorage(config)
	case "s3":
		return getS3Storage(config)
	default:
		return nil, errors.New("unsupported storage provider " + code)
	}
}
