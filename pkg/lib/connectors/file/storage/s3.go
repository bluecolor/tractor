package storage

import (
	"go.beyondstorage.io/services/s3/v3"
	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func getS3Storage(config map[string]interface{}) (types.Storager, error) {
	params := []types.Pair{}

	credential := "hmac:" + config["accessKeyId"].(string) + ":" + config["secretAccessKeyId"].(string)
	params = append(params, pairs.WithCredential(credential))

	if config["endpoint"] != nil {
		params = append(params, pairs.WithEndpoint(config["endpoint"].(string)))
	}
	if config["location"] != nil {
		params = append(params, pairs.WithLocation(config["location"].(string)))
	}
	if config["bucket"] != nil {
		params = append(params, pairs.WithName(config["bucket"].(string)))
	}

	return s3.NewStorager(params...)
}
