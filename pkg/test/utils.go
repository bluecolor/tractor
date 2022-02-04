package test

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bluecolor/tractor/pkg/utils"
)

func GenEnvFile() (string, string, error) {
	content := utils.Dedent(`
	LOG_LEVEL=info

	DB_HOST=localhost
	DB_PORT=3306
	DB_DATABASE=db
	DB_USERNAME=user
	DB_PASSWORD=password


	S3_ENDPOINT=localhost:9000
	S3_ACCESS_KEY_ID=NKB928W953OFDO1Z6IVG
	S3_SECRET_ACCESS_KEY=lPLZ7GagvvRYrq23RkwB043xuwSAz2cOGUIGu3AP
	S3_USE_SSL=false


	APP_SECRET=yOPODmbdwCbzPYhaSD4U1+CDchLNGyHRwzxMfd9VfcQgClmZ79Gmd0yP32VKS8kEWh5nRuqiyR/57o/PTy8st7rrMmmwx9cENKxcVtwwC+E6rstAWJD+yWlDE9EJ/mfdkZJKJ36EqtDU8xuuzD4L53IuxORvsTn9E9Prem+0JcLRqWNtL2Fj2f5sod0PLk6wCTICn5VIiJhtIvPcGnrzJ/UNkk8KsLl67NmTEyl1dqobwgqZpOHPiRrLi/JQ5qFwmSqD8MRd5GxEONfcK43dRNvnGvMJh3Rw3yHe965h42ygRXHzAOqSdhaWJlbeWHXR1Ge8hm7LvtmBoGK7+OGnGw==
	APP_SEED_PATH=./assets/seed
	`)

	dirname, err := ioutil.TempDir("/tmp", "tractor-test")
	if err != nil {
		return "", "", err
	}
	filename := filepath.Join(dirname, ".env")
	err = ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return "", "", err
	}
	return filename, dirname, nil
}
